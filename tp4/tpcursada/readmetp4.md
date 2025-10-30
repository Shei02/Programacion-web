# Proyecto Películas 

Este proyecto es un ejemplo de aplicación en Go que permite **crear, listar, actualizar y borrar películas** en una base de datos PostgreSQL. Se utiliza **SQLC** para generar código Go a partir de consultas SQL y **Docker-Compose** para levantar la base de datos fácilmente.

## Requisitos
Antes de comenzar, asegúrate de tener instalados:
- [Go 1.21+]
- [Docker]
- [Docker Compose]
- [SQLC]
- [make] ##recomendado para seed y pruebas
- [curl] ##recomendado para seed y pruebas
- [jq] ##recomendado para seed y pruebas

## Estructura del proyecto
tp4
├─ tpcursada
|   ├─ api
|   | ├─ mira.go
|   | ├─ peliculas.go
|   | ├─ usuarios.go
|   ├─ db/
|   │ ├─ schema/ # Archivos SQL para crear tablas
|   │ │ ├─ schema.sql
|   │ ├─ queries/ # Consultas SQL usadas por SQLC
|   │ │ ├─ queriesMira.sql
|   │ │ ├─ queriesPelicula.sql
|   │ │ ├─ queriesUsuario.sql
|   │ ├─ sqcl #aparece luego de hacer el sqlc generate 
|   ├─ static/
|   | ├─ app.js
|   | ├─ miras.html
|   | ├─ miras.js
|   | ├─ peliculas.html
|   | ├─ peliculas.js
|   | ├─ usuarios.html
|   | ├─ usuarios.js
|   | ├─ styles.css
|   ├─ main.go # Archivo principal de la aplicación
|   ├─ go.mod
|   ├─ go.sum
|   ├─ sqlc.yaml
|   ├─ index.html
|   ├─ Makefile
|   ├─ readmetp4.md
|   ├─ requests.bash
|   ├─ server.log
|   ├─ test_api.sh
|   └─ docker-compose.yml

## Pasos para ejecutar el proyecto

## 0. Clonar el repo desde git
    git clone <link>.git
    cd <nombre_del_repo>

## 1. Inicializar/preparar modulos GO
    go mod download
    go mod init <nombre_del_modulo>
    go mod tidy

## 2. Generar el codigo Go desde SQLC
   make sqcl

## 3. Construir todo
    make all

## 4. Cargar datos de ejemplo 
    make seed

## 5. Ejecutar la aplicacion 
    make run 
## Otra opcion de ejecucion 
    go run main.go
    
## 6. Abrir la app desde el navegador poner en la barra de direcciones
    http://localhost:8080

## Rutas útiles para pruebas (ejemplos)
    UI:
        http://localhost:8080/
        http://localhost:8080/peliculas
        http://localhost:8080/usuarios
        http://localhost:8080/miras
    API (JSON):
        GET /api/peliculas
        GET /api/peliculas/1
        GET /api/usuarios/1
        GET /api/miras?idu=2

## 7. Para finalizar y que no quede nada en ejecucion, eliminando contenedores y volumenes 
    make db-down

## Integrantes del grupo:
## Acevedo Belen
## Artaza Sheila
