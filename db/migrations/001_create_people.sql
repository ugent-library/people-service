CREATE TABLE people (
  id BIGSERIAL PRIMARY KEY,
  active BOOLEAN DEFAULT false NOT NULL,
  name TEXT NOT NULL CHECK (name <> ''),
  preferred_name TEXT CHECK (name <> ''),
  given_name TEXT CHECK (name <> ''),
  family_name TEXT CHECK (name <> ''),
  preferred_given_name TEXT CHECK (name <> ''),
  preferred_family_name TEXT CHECK (name <> ''),
  honorific_prefix TEXT CHECK (name <> ''),
  email TEXT CHECK (name <> ''),
  attributes JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE people_identifiers (
  person_id BIGINT NOT NULL REFERENCES people ON DELETE CASCADE,
  type TEXT CHECK (type <> ''),
  value TEXT CHECK (value <> ''),
  PRIMARY KEY (type, value)
);

---- create above / drop below ----

DROP TABLE people CASCADE;
DROP TABLE people_identifiers CASCADE;
