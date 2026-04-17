# Run Madget registry API with SQLite (same defaults as `make dev-sqlite`).
$ErrorActionPreference = "Stop"
$env:MADGET_DB_DRIVER = "sqlite"
$env:MADGET_DATABASE_URL = "./dev.db"
$env:MADGET_STORAGE_ROOT = "./storage"
go run ./apps/registry-api
