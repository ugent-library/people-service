CREATE TABLE people (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL CHECK (name <> ''),
  preferred_name TEXT CHECK (preferred_name <> ''),
  given_name TEXT CHECK (given_name <> ''),
  family_name TEXT CHECK (family_name <> ''),
  preferred_given_name TEXT CHECK (preferred_given_name <> ''),
  preferred_family_name TEXT CHECK (preferred_family_name <> ''),
  honorific_prefix TEXT CHECK (honorific_prefix <> ''),
  email TEXT CHECK (email <> ''),
  active BOOLEAN DEFAULT false NOT NULL,
  username TEXT CHECK (username <> ''),
  attributes JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- CREATE UNIQUE INDEX "people_email_active_key" ON "people" ("email") WHERE (active=true);
-- CREATE UNIQUE INDEX "people_username_active_key" ON "people" ("username") WHERE (active=true);

CREATE TABLE people_identifiers (
  person_id BIGINT NOT NULL REFERENCES people ON DELETE CASCADE,
  type TEXT CHECK (type <> ''),
  value TEXT CHECK (value <> ''),
  PRIMARY KEY (type, value)
);

---- create above / drop below ----

DROP TABLE people CASCADE;
DROP TABLE people_identifiers CASCADE;
