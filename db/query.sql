-- name: GetPerson :one
SELECT p.*
FROM people p, people_identifiers pi
WHERE p.id = pi.person_id AND pi.type = $1 AND pi.value = $2;

-- name: CreatePerson :one
INSERT INTO people (
  name,
  preferred_name,
  given_name,
  family_name,
  preferred_given_name,
  preferred_family_name,
  honorific_prefix,
  email,
  active,
  username,
  attributes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING id;

-- name: UpdatePerson :exec
UPDATE people SET (
  name,
  preferred_name,
  given_name,
  family_name,
  preferred_given_name,
  preferred_family_name,
  honorific_prefix,
  email,
  active,
  username,
  attributes,
  updated_at
) = ($2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, CURRENT_TIMESTAMP)
WHERE id = $1;

-- name: DeletePerson :exec
DELETE FROM people
WHERE id = $1;

-- name: GetPersonIdentifiers :many
SELECT *
FROM people_identifiers
WHERE person_id = $1;

-- name: CreatePersonIdentifier :exec
INSERT INTO people_identifiers (
  person_id,
  type,
  value
) VALUES ($1, $2, $3);

-- name: TransferPersonIdentifier :exec
UPDATE people_identifiers SET person_id = ($3)
WHERE type = $1 AND value = $2;

-- name: DeletePersonIdentifier :exec
DELETE FROM people_identifiers
WHERE type = $1 AND value = $2;