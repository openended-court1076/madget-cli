-- Mevcut Postgres kurulumlarında bir kez çalıştırın (sütun yoksa).
ALTER TABLE package_versions ADD COLUMN IF NOT EXISTS manifest_xml TEXT;
ALTER TABLE package_versions ADD COLUMN IF NOT EXISTS metadata_json TEXT;
