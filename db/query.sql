-- name: GetPerson :one
WITH i AS (
  SELECT i1.*
  FROM people_identifiers i1
  LEFT JOIN  people_identifiers i2 ON i1.person_id = i2.person_id
  WHERE i2.type = $1 AND i2.value = $2	
)
SELECT p.*, json_agg(json_build_object('type', i.type, 'value', i.value)) AS identifiers
FROM people p, i WHERE p.id = i.person_id
GROUP BY p.id;

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

-- name: DeactivatePeople :exec
UPDATE people SET active = FALSE
WHERE updated_at < $1;

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