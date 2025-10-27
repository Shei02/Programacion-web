#!/bin/bash

echo "=== Crear Película ==="
curl -X POST -H "Content-Type: application/json" \
    -d '{"Titulo":"Tenet","Duracion":"150 min","Director":"Christopher Nolan","Actores":"John David Washington","Edadmin":13,"Sinopsis":"Inversión temporal","Anioestr":2020}' \
    http://localhost:8080/peliculas
echo

echo "=== Listar Películas ==="
curl http://localhost:8080/peliculas
echo

echo "=== Crear Usuario ==="
curl -X POST -H "Content-Type: application/json" \
    -d '{"Nomusu":"sheii","Contrasenia":"123456","Email":"sheila@example.com","Fechanac":"2002-05-06T00:00:00Z"}' \
    http://localhost:8080/usuarios
echo

echo "=== Listar Usuarios ==="
curl http://localhost:8080/usuarios
echo

echo "=== Crear Mira (usuario 1, película 1) ==="
curl -X POST -H "Content-Type: application/json" \
    -d '{"Idp":1,"Idu":1,"Gustoono":"Me gustó","Calif":5}' \
    http://localhost:8080/miras
echo

echo "=== Ver Miras del Usuario 1 ==="
curl http://localhost:8080/miras/1
echo