-- name: GetPersonByIdentifier :one
SELECT p.*
FROM people p
INNER JOIN people_identifiers pi
  ON pi.person_id = p.id
  WHERE pi.type = $1 AND pi.value = $2;

-- name: CreatePerson :one
INSERT INTO people (
  active,
  name,
  preferred_name,
  given_name,
  family_name,
  preferred_given_name,
  preferred_family_name,
  honorific_prefix,
  email
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id;

-- name: CreatePersonIdentifier :exec
INSERT INTO people_identifiers (
  person_id,
  type,
  value
) VALUES ($1, $2, $3);