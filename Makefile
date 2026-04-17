# Madget — local development
# SQLite: schema is applied automatically on startup (embedded SQL + dev token).

.PHONY: dev-sqlite dev-postgres example-tgz help

help:
	@echo "Targets:"
	@echo "  make dev-sqlite    - Run registry API with SQLite (./dev.db, default port 8080)"
	@echo "  make dev-postgres  - Run registry API with Postgres (set URL or use deployments/docker-compose.yml)"
	@echo "  make example-tgz   - Build example/package.tgz from example/payload/"

dev-sqlite:
	MADGET_DB_DRIVER=sqlite MADGET_DATABASE_URL=./dev.db MADGET_STORAGE_ROOT=./storage go run ./apps/registry-api

dev-postgres:
	MADGET_DB_DRIVER=postgres MADGET_DATABASE_URL=postgres://madget:madget@localhost:5432/madget?sslmode=disable MADGET_STORAGE_ROOT=./storage go run ./apps/registry-api

example-tgz:
	tar -czf example/package.tgz -C example/payload .
