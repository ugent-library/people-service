-- +goose Up

CREATE TABLE organizations (
  id BIGSERIAL PRIMARY KEY,
  parent_id BIGINT REFERENCES organizations ON DELETE SET NULL CHECK (parent_id <> id),
  name TEXT NOT NULL CHECK (name <> ''),
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE organization_identifiers (
  organization_id BIGINT NOT NULL REFERENCES organizations ON DELETE CASCADE,
  type TEXT CHECK (type <> ''),
  value TEXT CHECK (value <> ''),
  PRIMARY KEY (type, value)
);

CREATE INDEX organization_identifiers_person_id_fkey ON organization_identifiers (organization_id);

-- +goose Down

DROP TABLE organizations CASCADE;
DROP TABLE organization_identifiers CASCADE;
