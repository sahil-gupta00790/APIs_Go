-- migrate:up
CREATE TABLE schools(
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    level text NOT NULL,
    contact text NOT NULL,
    phone text NOT NULL,
    email text NOT NULL,
    website text NOT NULL,
    address text NOT NULL,
    mode text[] NOT NULL,
    version integer NOT NULL DEFAULT 1
);

-- migrate:down
DROP TABLE IF EXISTS schools;

