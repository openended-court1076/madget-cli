package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/mehmetalidsy/madget-cli/internal/registry"
)

func main() {
	driver, dsn := resolveDatabaseEnv()
	storageRoot := os.Getenv("MADGET_STORAGE_ROOT")
	if storageRoot == "" {
		storageRoot = "./storage"
	}

	store, err := registry.NewStore(driver, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	handler := registry.NewHandler(store, storageRoot)
	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	addr := ":8080"
	log.Printf("madget registry listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

func resolveDatabaseEnv() (driver, dsn string) {
	driver = os.Getenv("MADGET_DB_DRIVER")
	dsn = os.Getenv("MADGET_DATABASE_URL")

	if driver == "" {
		if dsn == "" {
			return "sqlite", "./dev.db"
		}
		if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
			return "postgres", dsn
		}
		return "sqlite", dsn
	}
	if dsn == "" && driver == "sqlite" {
		return driver, "./dev.db"
	}
	return driver, dsn
}
