-- name: AddPerson :exec
INSERT INTO people (
  name
) VALUES (
  $1
);