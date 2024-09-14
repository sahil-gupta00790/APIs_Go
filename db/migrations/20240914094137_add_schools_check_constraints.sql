-- migrate:up
ALTER TABLE schools ADD CONSTRAINT mode_length_check CHECK (array_length(mode,1) BETWEEN 1 AND 5);

-- migrate:down
ALTER TABLE schools DROP CONSTRAINT IF EXIST mode_length_check;
