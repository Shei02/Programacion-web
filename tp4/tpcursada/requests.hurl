GET http://localhost:8080/api/usuarios
HTTP/1.1 200

# Agregar usuarios 
POST http://localhost:8080/api/usuarios
Content-Type: application/json

{"nomusu":"Sheila","contrasenia":"1234","email":"sheila@test.com","fechanac":"1990-05-12T00:00:00Z"}

HTTP/1.1 201

POST http://localhost:8080/api/usuarios
Content-Type: application/json

{"nomusu":"Esteban","contrasenia":"abcd","email":"esteban@test.com","fechanac":"1985-11-23T00:00:00Z"}

HTTP/1.1 201

POST http://localhost:8080/api/usuarios
Content-Type: application/json

{"nomusu":"Belen","contrasenia":"1234","email":"belen@test.com","fechanac":"1958-05-12T00:00:00Z"}

HTTP/1.1 201

POST http://localhost:8080/api/usuarios
Content-Type: application/json

{"nomusu":"Eden","contrasenia":"abcd","email":"eden@test.com","fechanac":"2002-11-23T00:00:00Z"}

HTTP/1.1 201

# Películas
GET http://localhost:8080/api/peliculas
HTTP/1.1 200

POST http://localhost:8080/api/peliculas
Content-Type: application/json

{"titulo":"Matrix","duracion":136,"director":"Lana Wachowski","actores":"Keanu Reeves, Laurence Fishburne","edadmin":13,"sinopsis":"Un hacker descubre la realidad simulada en la que vive.","anioestr":1999}

HTTP/1.1 201

POST http://localhost:8080/api/peliculas
Content-Type: application/json

{"titulo":"Inception","duracion":148,"director":"Christopher Nolan","actores":"Leonardo DiCaprio, Joseph Gordon-Levitt","edadmin":13,"sinopsis":"Un ladrón especializado en sueños roba secretos de la mente.","anioestr":2010}

HTTP/1.1 201

POST http://localhost:8080/api/peliculas
Content-Type: application/json

{"titulo":"Interstellar","duracion":169,"director":"Christopher Nolan","actores":"Matthew McConaughey, Anne Hathaway","edadmin":13,"sinopsis":"Un grupo de exploradores viaja a través de un agujero de gusano en busca de un nuevo hogar para la humanidad.","anioestr":2014}

HTTP/1.1 201

POST http://localhost:8080/api/peliculas
Content-Type: application/json

{"titulo":"The Godfather","duracion":175,"director":"Francis Ford Coppola","actores":"Marlon Brando, Al Pacino","edadmin":18,"sinopsis":"La historia de la familia mafiosa Corleone.","anioestr":1972}

HTTP/1.1 201

# Miras (calificaciones)
GET http://localhost:8080/api/miras?idu=1
HTTP/1.1 200

POST http://localhost:8080/api/miras
Content-Type: application/json

{"idp":1,"idu":1,"gustoono":"sí","calif":5}

HTTP/1.1 201

POST http://localhost:8080/api/miras
Content-Type: application/json

{"idp":2,"idu":2,"gustoono":"no","calif":1}

HTTP/1.1 201

POST http://localhost:8080/api/miras
Content-Type: application/json

{"idp":3,"idu":1,"gustoono":"sí","calif":5}

HTTP/1.1 201

POST http://localhost:8080/api/miras
Content-Type: application/json

{"idp":4,"idu":2,"gustoono":"no","calif":3}

HTTP/1.1 201
