# Madget v2 - Registry + API + CLI

Madget v2 is a package management ecosystem built in Go with three core parts:

- CLI (`madget`) for publish/install workflows
- Registry API for package metadata, version resolution, and tarball delivery
- Metadata storage: PostgreSQL or SQLite (local dev defaults to SQLite)

## Project Layout

- `apps/cli`: Cobra-based CLI commands (`init`, `login`, `publish`, `install`)
- `apps/registry-api`: HTTP API service
- `internal/resolver`: semver range resolution (`^`, `~`, exact)
- Publish: tam **MadGet.xml** gövdesi `manifest_xml` olarak saklanır; ayrıştırılmış **`metadata`** JSON olarak `package_versions` satırında tutulur (`GET /v1/packages/{name}/versions`, resolve yanıtında döner)
- `internal/integrity`: SHA256 checksum helpers
- `migrations`: SQL schema bootstrap
- `deployments`: local Docker compose
- `example`: örnek `MadGet.xml` + payload; `make example-tgz` ile `package.tgz` üretilir

## Quick Start

### Registry API (SQLite, no Docker)

From the repo root, with no env vars, the API defaults to **SQLite** at `./dev.db` and applies the embedded schema on startup (including a dev publisher token `dev-token`).

```bash
go run ./apps/registry-api
```

Or one command:

```bash
make dev-sqlite
```

On Windows (PowerShell):

```powershell
.\scripts\dev-sqlite.ps1
```

### Registry API (PostgreSQL)

1. Start Postgres:

```bash
docker compose -f deployments/docker-compose.yml up -d
```

2. Run the API:

```bash
set MADGET_DB_DRIVER=postgres
set MADGET_DATABASE_URL=postgres://madget:madget@localhost:5432/madget?sslmode=disable
set MADGET_STORAGE_ROOT=./storage
go run ./apps/registry-api
```

Or:

```bash
make dev-postgres
```

### CLI (against http://localhost:8080)

Örnek paket (`example/`): önce tarball üret, sonra publish + install.

```bash
make example-tgz
go run . init --registry http://localhost:8080
go run . login --token dev-token
go run . publish ./example/MadGet.xml ./example/package.tgz
go run . install example-pkg@^1.0.0
```

Windows’ta tarball için: `.\scripts\build-example-package.ps1`

Kurulum çıktısı `vendor/example-pkg/1.0.0/` altında oluşur (ör. `README.txt`).

## Paket manifesti: MadGet.xml

Tüm paketler **`MadGet.xml`** kullanır (artık `package.json` yok). Registry adı için **`package_name`** önceliklidir; yoksa **`name`** kullanılır. **`version`** ve **`description`** zorunludur.

Örnek: `example/MadGet.xml` veya repo kökündeki `MadGet.xml`.

```xml
<?xml version="1.0" encoding="UTF-8"?>
<application>
    <info name="ExamplePkg"
          version="1.0.0"
          package_name="example-pkg"
          license="MIT"
          categories="Example">
        <description>Kısa açıklama</description>
        <author id="mehmetalidsy" />
    </info>
</application>
```

## Next Steps

- Add dependency graph persistence (`dependencies` table + resolver integration)
- Add token issuance endpoint and hashed-token storage
- Add lockfile generation and update command
- Add E2E tests for publish/resolve/install flow
