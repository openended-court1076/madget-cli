package registry

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type PostgresStore struct {
	db *sql.DB
}

func newPostgresStore(dsn string) (Store, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Close() error {
	return s.db.Close()
}

func (s *PostgresStore) IsNoRows(err error) bool {
	return err == sql.ErrNoRows
}

func (s *PostgresStore) IsValidToken(raw string) (bool, error) {
	var exists bool
	err := s.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM tokens WHERE token = $1 AND revoked_at IS NULL)`, raw).Scan(&exists)
	return exists, err
}

func (s *PostgresStore) UpsertPackage(name string) (int64, error) {
	var packageID int64
	err := s.db.QueryRow(`
		INSERT INTO packages (name)
		VALUES ($1)
		ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`, name).Scan(&packageID)
	return packageID, err
}

func (s *PostgresStore) InsertPackageVersion(packageID int64, version, checksum, tarballPath, manifestXML, metadataJSON string) error {
	_, err := s.db.Exec(`
		INSERT INTO package_versions (package_id, version, checksum, tarball_path, manifest_xml, metadata_json)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, packageID, version, checksum, tarballPath, manifestXML, metadataJSON)
	return err
}

func (s *PostgresStore) ListPackageVersions(name string) ([]PackageVersion, error) {
	rows, err := s.db.Query(`
		SELECT p.name, pv.version, pv.checksum, pv.tarball_path, pv.published_at::text,
		       COALESCE(pv.manifest_xml, ''), COALESCE(pv.metadata_json, '')
		FROM packages p
		JOIN package_versions pv ON pv.package_id = p.id
		WHERE p.name = $1
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

func (s *PostgresStore) FindPackageVersion(name, version string) (PackageVersion, error) {
	var v PackageVersion
	err := s.db.QueryRow(`
		SELECT p.name, pv.version, pv.checksum, pv.tarball_path, pv.published_at::text,
		       COALESCE(pv.manifest_xml, ''), COALESCE(pv.metadata_json, '')
		FROM packages p
		JOIN package_versions pv ON pv.package_id = p.id
		WHERE p.name = $1 AND pv.version = $2
		LIMIT 1
	`, name, version).Scan(&v.Name, &v.Version, &v.Checksum, &v.Tarball, &v.Published, &v.ManifestXML, &v.MetadataJSON)
	return v, err
}
