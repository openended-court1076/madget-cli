package registry

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type SQLiteStore struct {
	db *sql.DB
}

func newSQLiteStore(dsn string) (Store, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	if err := applySQLiteMigrations(db); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

func (s *SQLiteStore) IsNoRows(err error) bool {
	return err == sql.ErrNoRows
}

func (s *SQLiteStore) IsValidToken(raw string) (bool, error) {
	var count int
	err := s.db.QueryRow(`SELECT COUNT(1) FROM tokens WHERE token = ? AND revoked_at IS NULL`, raw).Scan(&count)
	return count > 0, err
}

func (s *SQLiteStore) UpsertPackage(name string) (int64, error) {
	if _, err := s.db.Exec(`
		INSERT INTO packages (name)
		VALUES (?)
		ON CONFLICT(name) DO UPDATE SET name = excluded.name
	`, name); err != nil {
		return 0, err
	}
	var packageID int64
	if err := s.db.QueryRow(`SELECT id FROM packages WHERE name = ? LIMIT 1`, name).Scan(&packageID); err != nil {
		return 0, err
	}
	return packageID, nil
}

func (s *SQLiteStore) InsertPackageVersion(packageID int64, version, checksum, tarballPath, manifestXML, metadataJSON string) error {
	_, err := s.db.Exec(`
		INSERT INTO package_versions (package_id, version, checksum, tarball_path, manifest_xml, metadata_json)
		VALUES (?, ?, ?, ?, ?, ?)
	`, packageID, version, checksum, tarballPath, manifestXML, metadataJSON)
	return err
}

func (s *SQLiteStore) ListPackageVersions(name string) ([]PackageVersion, error) {
	rows, err := s.db.Query(`
		SELECT p.name, pv.version, pv.checksum, pv.tarball_path, CAST(pv.published_at AS TEXT),
		       IFNULL(pv.manifest_xml, ''), IFNULL(pv.metadata_json, '')
		FROM packages p
		JOIN package_versions pv ON pv.package_id = p.id
		WHERE p.name = ?
	`, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]PackageVersion, 0)
	for rows.Next() {
		var v PackageVersion
		if err := rows.Scan(&v.Name, &v.Version, &v.Checksum, &v.Tarball, &v.Published, &v.ManifestXML, &v.MetadataJSON); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) FindPackageVersion(name, version string) (PackageVersion, error) {
	var v PackageVersion
	err := s.db.QueryRow(`
		SELECT p.name, pv.version, pv.checksum, pv.tarball_path, CAST(pv.published_at AS TEXT),
		       IFNULL(pv.manifest_xml, ''), IFNULL(pv.metadata_json, '')
		FROM packages p
		JOIN package_versions pv ON pv.package_id = p.id
		WHERE p.name = ? AND pv.version = ?
		LIMIT 1
	`, name, version).Scan(&v.Name, &v.Version, &v.Checksum, &v.Tarball, &v.Published, &v.ManifestXML, &v.MetadataJSON)
	return v, err
}
