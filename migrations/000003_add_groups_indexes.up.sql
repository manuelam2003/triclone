CREATE INDEX IF NOT EXISTS groups_name_idx ON groups USING GIN (to_tsvector('simple', name));
