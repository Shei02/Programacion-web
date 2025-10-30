-- name: GetUsuario :one 
SELECT * FROM usuario WHERE IDU = $1;

-- name: ListUsuario :many
SELECT * FROM usuario ORDER BY nomUsu;

-- name: CreateUsuario :one
INSERT INTO usuario (nomUsu, contrasenia, email, fechaNac) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateUsuario :exec
UPDATE usuario SET nomUsu = $2, contrasenia = $3, email = $4, fechaNac = $5  WHERE IDU = $1;

-- name: DeleteUsuario :exec
DELETE FROM usuario WHERE IDU = $1;
