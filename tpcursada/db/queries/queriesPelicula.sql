-- name: GetPeli :one 
SELECT * FROM pelicula WHERE IDP = $1;

-- name: ListPelis :many
SELECT * FROM pelicula ORDER BY titulo;

-- name: CreatePeli :one
INSERT INTO pelicula (titulo, duracion, director, actores, edadMin, sinopsis, anioEstr) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: UpdatePeli :exec
UPDATE pelicula SET titulo = $2, duracion = $3, director = $4, actores = $5, edadMin = $6, sinopsis =$7, anioEstr =$8 WHERE IDP = $1;

-- name: DeletePeli :exec
DELETE FROM pelicula WHERE IDP = $1;
