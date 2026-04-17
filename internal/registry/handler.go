package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mehmetalidsy/madget-cli/internal/integrity"
	"github.com/mehmetalidsy/madget-cli/internal/manifest"
	"github.com/mehmetalidsy/madget-cli/internal/resolver"
)

type Handler struct {
	store       Store
	storageRoot string
}

func NewHandler(store Store, storageRoot string) *Handler {
	return &Handler{store: store, storageRoot: storageRoot}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	r.Post("/v1/packages", h.publish)
	r.Get("/v1/packages/{name}/versions", h.listVersions)
	r.Get("/v1/packages/{name}/resolve", h.resolve)
	r.Get("/v1/tarballs/{name}/{version}", h.tarball)
}

func (h *Handler) publish(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	ok, err := h.store.IsValidToken(strings.TrimSpace(token))
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !ok {
		respondError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	if err := r.ParseMultipartForm(20 << 20); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	manifestXML := strings.TrimSpace(r.FormValue("manifest_xml"))
	if manifestXML == "" {
		respondError(w, http.StatusBadRequest, "manifest_xml is required (full MadGet.xml body)")
		return
	}
	app, err := manifest.UnmarshalApplication([]byte(manifestXML))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid MadGet.xml: "+err.Error())
		return
	}
	pm, err := manifest.LoadFromBytes([]byte(manifestXML))
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	name := pm.RegistryName()
	version := strings.TrimSpace(pm.Version)
	if name == "" || version == "" {
		respondError(w, http.StatusBadRequest, "MadGet.xml must include package_name or name and version")
		return
	}
	formName := strings.TrimSpace(r.FormValue("name"))
	formVer := strings.TrimSpace(r.FormValue("version"))
	if formName != "" && formName != name {
		respondError(w, http.StatusBadRequest, "name form field must match MadGet.xml registry name")
		return
	}
	if formVer != "" && formVer != version {
		respondError(w, http.StatusBadRequest, "version form field must match MadGet.xml version")
		return
	}
	metaBytes, err := app.MetadataJSON()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	metadataJSON := string(metaBytes)

	file, header, err := r.FormFile("tarball")
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".tgz"
	}
	objectPath := filepath.Join(name, version, uuid.NewString()+ext)
	absolutePath := filepath.Join(h.storageRoot, objectPath)
	if err := os.MkdirAll(filepath.Dir(absolutePath), 0o755); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	out, err := os.Create(absolutePath)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if _, err := io.Copy(out, file); err != nil {
		out.Close()
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	out.Close()

	checksum, err := integrity.FileSHA256(absolutePath)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	packageID, err := h.store.UpsertPackage(name)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := h.store.InsertPackageVersion(packageID, version, checksum, objectPath, manifestXML, metadataJSON); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	created := map[string]any{
		"name":     name,
		"version":  version,
		"checksum": checksum,
	}
	var metaObj any
	if err := json.Unmarshal([]byte(metadataJSON), &metaObj); err == nil {
		created["metadata"] = metaObj
	}
	respondJSON(w, http.StatusCreated, created)
}

func (h *Handler) listVersions(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	versions, err := h.store.ListPackageVersions(name)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	out := make([]map[string]any, 0, len(versions))
	for _, v := range versions {
		out = append(out, packageVersionToMap(v))
	}
	respondJSON(w, http.StatusOK, out)
}

func (h *Handler) resolve(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	rng := r.URL.Query().Get("range")
	if rng == "" {
		respondError(w, http.StatusBadRequest, "range query is required")
		return
	}

	versions, err := h.store.ListPackageVersions(name)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	candidates := make([]string, 0, len(versions))
	for _, v := range versions {
		candidates = append(candidates, v.Version)
	}

	resolvedVersion, err := resolver.Resolve(candidates, rng)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	match, err := h.store.FindPackageVersion(name, resolvedVersion)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resolved := map[string]any{
		"package":     match.Name,
		"version":     match.Version,
		"checksum":    match.Checksum,
		"tarball_url": fmt.Sprintf("/v1/tarballs/%s/%s", match.Name, match.Version),
	}
	for k, v := range packageVersionToMap(match) {
		if k == "name" || k == "version" || k == "checksum" || k == "tarball_path" || k == "published_at" {
			continue
		}
		resolved[k] = v
	}
	respondJSON(w, http.StatusOK, resolved)
}

func (h *Handler) tarball(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	version := chi.URLParam(r, "version")
	match, err := h.store.FindPackageVersion(name, version)
	if err != nil {
		if h.store.IsNoRows(err) {
			respondError(w, http.StatusNotFound, "package version not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	http.ServeFile(w, r, filepath.Join(h.storageRoot, match.Tarball))
}

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

func packageVersionToMap(v PackageVersion) map[string]any {
	out := map[string]any{
		"name":         v.Name,
		"version":      v.Version,
		"checksum":     v.Checksum,
		"tarball_path": v.Tarball,
		"published_at": v.Published,
	}
	if v.ManifestXML != "" {
		out["manifest_xml"] = v.ManifestXML
	}
	if v.MetadataJSON != "" {
		var meta any
		if err := json.Unmarshal([]byte(v.MetadataJSON), &meta); err == nil {
			out["metadata"] = meta
		}
	}
	return out
}
