CREATE TABLE people (
  id BIGSERIAL PRIMARY KEY,
  active BOOLEAN DEFAULT false NOT NULL,
  name TEXT NOT NULL CHECK (name <> ''),
  preferred_name TEXT,
  given_name TEXT,
  family_name TEXT,
  preferred_given_name TEXT,
  preferred_family_name TEXT,
  honorific_prefix TEXT,
  email TEXT,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE people_identifiers (
  person_id BIGINT NOT NULL REFERENCES people,
  type TEXT CHECK (type <> ''),
  value TEXT CHECK (value <> ''),
  PRIMARY KEY (type, value)
);

---- create above / drop below ----

DROP TABLE IF EXISTS people CASCADE;
