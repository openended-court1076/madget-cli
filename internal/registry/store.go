package registry

import (
	"errors"
)

type Store interface {
	Close() error
	IsValidToken(raw string) (bool, error)
	UpsertPackage(name string) (int64, error)
	InsertPackageVersion(packageID int64, version, checksum, tarballPath, manifestXML, metadataJSON string) error
	ListPackageVersions(name string) ([]PackageVersion, error)
	FindPackageVersion(name, version string) (PackageVersion, error)
	IsNoRows(err error) bool
}

func NewStore(driver, dsn string) (Store, error) {
	switch driver {
	case "", "postgres":
		if dsn == "" {
			return nil, errors.New("MADGET_DATABASE_URL is required when MADGET_DB_DRIVER is postgres (or set MADGET_DB_DRIVER=sqlite for local ./dev.db)")
		}
		return newPostgresStore(dsn)
	case "sqlite":
		if dsn == "" {
			dsn = "./dev.db"
		}
		return newSQLiteStore(dsn)
	default:
		return nil, errors.New("MADGET_DB_DRIVER must be postgres or sqlite")
	}
}
