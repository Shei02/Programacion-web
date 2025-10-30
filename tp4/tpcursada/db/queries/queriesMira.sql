-- name: GetMira :many
SELECT * FROM mira WHERE IDU = $1; 

-- name: ListMira :many
SELECT * FROM mira ORDER BY calif;

-- name: CreateMira :exec
INSERT INTO mira (IDP, IDU, gustoOno, calif) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateMira :exec
UPDATE mira SET gustoOno = $3, calif = $4 WHERE IDP = $1 AND IDU = $2;

-- name: DeleteMira :exec
DELETE FROM mira WHERE IDU = $1;
