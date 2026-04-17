package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/mehmetalidsy/madget-cli/internal/config"
	"github.com/mehmetalidsy/madget-cli/internal/manifest"
)

func publishPackage(cfg config.LocalConfig, manifestPath, tarballPath string) error {
	if cfg.Token == "" {
		return fmt.Errorf("publisher token is missing, run login first")
	}

	xmlBytes, err := os.ReadFile(manifestPath)
	if err != nil {
		return err
	}
	meta, err := manifest.LoadFromBytes(xmlBytes)
	if err != nil {
		return err
	}

	file, err := os.Open(tarballPath)
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("name", meta.RegistryName())
	_ = writer.WriteField("version", meta.Version)
	_ = writer.WriteField("description", meta.Description)
	_ = writer.WriteField("manifest_xml", string(xmlBytes))

	part, err := writer.CreateFormFile("tarball", filepath.Base(tarballPath))
	if err != nil {
		return err
	}
	if _, err := io.Copy(part, file); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, cfg.RegistryURL+"/v1/packages", body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+cfg.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("publish failed with status %d", resp.StatusCode)
	}

	var out map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return err
	}
	fmt.Printf("Published %v@%v\n", out["name"], out["version"])
	return nil
}

func parseInstallSpec(spec string) (name, versionRange string, err error) {
	parts := strings.Split(spec, "@")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid spec %q, expected name@range", spec)
	}
	return parts[0], parts[1], nil
}
