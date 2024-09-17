-- migrate:up
CREATE INDEX IF NOT EXISTS schools_name_idx ON schools USING GIN(to_tsvector('simple',name));
CREATE INDEX IF NOT EXISTS schools_level_idx ON schools USING GIN(to_tsvector('simple',level));
CREATE INDEX IF NOT EXISTS schools_mode_idx ON schools USING GIN(mode);

-- migrate:down

DROP INDEX IF EXISTS schools_name_idx;
DROP INDEX IF EXISTS schools_level_idx;
DROP INDEX IF EXISTS schools_mode_idx;
