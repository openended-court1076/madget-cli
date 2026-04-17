package registry

import (
	"database/sql"
	_ "embed"
	"fmt"
	"strings"
)

//go:embed schema_sqlite.sql
var embeddedSQLiteSchema string

func applySQLiteMigrations(db *sql.DB) error {
	for _, stmt := range splitSQLStatements(embeddedSQLiteSchema) {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}
	return ensurePackageVersionMetadataColumns(db)
}

func ensurePackageVersionMetadataColumns(db *sql.DB) error {
	for _, col := range []string{"manifest_xml", "metadata_json"} {
		var n int
		if err := db.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('package_versions') WHERE name = ?`, col).Scan(&n); err != nil {
			return err
		}
		if n == 0 {
			q := fmt.Sprintf("ALTER TABLE package_versions ADD COLUMN %s TEXT", col)
			if _, err := db.Exec(q); err != nil {
				return err
			}
		}
	}
	return nil
}

func splitSQLStatements(sql string) []string {
	lines := strings.Split(sql, "\n")
	var out []string
	var b strings.Builder
	for _, line := range lines {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "--") {
			continue
		}
		b.WriteString(line)
		b.WriteByte('\n')
	}
	raw := strings.TrimSpace(b.String())
	for _, part := range strings.Split(raw, ";") {
		stmt := strings.TrimSpace(part)
		if stmt != "" {
			out = append(out, stmt)
		}
	}
	return out
}
