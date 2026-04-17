CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE NOT NULL,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tokens (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER REFERENCES users(id),
    token TEXT UNIQUE NOT NULL,
    scopes TEXT NOT NULL DEFAULT '[]',
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    revoked_at TEXT
);

CREATE TABLE IF NOT EXISTS packages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS package_versions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    package_id INTEGER NOT NULL REFERENCES packages(id) ON DELETE CASCADE,
    version TEXT NOT NULL,
    checksum TEXT NOT NULL,
    tarball_path TEXT NOT NULL,
    manifest_xml TEXT,
    metadata_json TEXT,
    published_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(package_id, version)
);

-- Local dev: use Authorization: Bearer dev-token
INSERT OR IGNORE INTO users (id, email) VALUES (1, 'dev@local');
INSERT OR IGNORE INTO tokens (id, user_id, token, scopes) VALUES (1, 1, 'dev-token', '[]');
