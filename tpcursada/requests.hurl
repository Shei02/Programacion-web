# requests.hurl
# Hurl examples + simple assertions for the API (peliculas, usuarios, miras).
# Run with: hurl requests.hurl

### Create pelicula
POST http://localhost:8080/peliculas
Content-Type: application/json

{
  "titulo": "Tenet",
  "duracion": "150 min",
  "director": "Christopher Nolan",
  "actores": "John David Washington",
  "edadmin": 13,
  "sinopsis": "Inversión temporal",
  "anioestr": 2020
}

HTTP/1.1 201
header Content-Type: application/json
body-contains "titulo"

### List peliculas
GET http://localhost:8080/peliculas

HTTP/1.1 200
header Content-Type: application/json

### Get pelicula 1
GET http://localhost:8080/peliculas/1

HTTP/1.1 200
header Content-Type: application/json
body-contains "idp"

### Update pelicula 1
PUT http://localhost:8080/peliculas/1
Content-Type: application/json

{
  "titulo": "Tenet (editada)",
  "duracion": "152 min",
  "director": "Christopher Nolan",
  "actores": "John David Washington",
  "edadmin": 13,
  "sinopsis": "Inversión temporal (versión editada)",
  "anioestr": 2020
}

HTTP/1.1 200
header Content-Type: application/json
body-contains "titulo"

### Delete pelicula 1
DELETE http://localhost:8080/peliculas/1

# If there are miras referencing the pelicula, expect 409 or 204 if deleted
HTTP/1.1 200|204|409

### Create usuario
POST http://localhost:8080/usuarios
Content-Type: application/json

{
  "nomusu": "sheila",
  "contrasenia": "123456",
  "email": "sheila@example.com",
  "fechanac": "1990-05-12T00:00:00Z"
}

HTTP/1.1 201
header Content-Type: application/json
body-contains "idu"

### Get usuario 1
GET http://localhost:8080/usuarios/1

HTTP/1.1 200
header Content-Type: application/json
body-contains "idu"

### Update usuario 1
PUT http://localhost:8080/usuarios/1
Content-Type: application/json

{
  "nomusu": "sheila_mod",
  "contrasenia": "654321",
  "email": "sheila+mod@example.com",
  "fechanac": "1990-05-12T00:00:00Z"
}

HTTP/1.1 200
header Content-Type: application/json
body-contains "nomusu"

### Delete usuario 1
DELETE http://localhost:8080/usuarios/1

# If there are miras referencing the usuario, expect 409 or 204 if deleted
HTTP/1.1 200|204|409

### Create mira
POST http://localhost:8080/miras
Content-Type: application/json

{
  "idp": 1,
  "idu": 1,
  "gustoono": "sí",
  "calif": 5
}

# Could be 201 (created) or 409 (FK violation or duplicate)
HTTP/1.1 201|409
header Content-Type: application/json|text/plain

### Get miras for user
GET http://localhost:8080/miras?idu=1

HTTP/1.1 200
header Content-Type: application/json

### Get mira (idp=1,idu=1)
GET http://localhost:8080/miras/1/1

HTTP/1.1 200|404
header Content-Type: application/json

### Delete mira (1,1)
DELETE http://localhost:8080/miras/1/1

HTTP/1.1 204|404
