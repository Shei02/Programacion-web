# Proyecto Películas 

Este proyecto es un ejemplo de aplicación en Go que permite **crear, listar, actualizar y borrar películas** en una base de datos PostgreSQL. Se utiliza **SQLC** para generar código Go a partir de consultas SQL y **Docker-Compose** para levantar la base de datos fácilmente.

---

## Requisitos

Antes de comenzar, asegúrate de tener instalados:

- [Go 1.21+]
- [Docker]
- [Docker Compose]
- [SQLC]

---

## Estructura del proyecto
my_sql_project/
│
├─ db/
│ ├─ schema/ # Archivos SQL para crear tablas
│ │ ├─ schema.sql
│ ├─ queries/ # Consultas SQL usadas por SQLC
│ │ ├─ queriesMira.sql
│ │ ├─ queriesPelicula.sql
│ │ ├─ queriesUsuario.sql
│ 
├─ main.go # Archivo principal de la aplicación
├─ go.mod
├─ go.sum
├─ sqlc.yaml
└─ docker-compose.yml

## Pasos para ejecutar el proyecto
## 1. Levantar POSTGRESQL con docker 
    Ejecuta en la terminal:
    docker compose up -d 

## 2. Generar el codigo Go desde SQLC
   Ejecuta en la terminal: 
    sqcl generate 

## 3. Verificar que el contenedor docker este corriendo 
    Ejecuta en la terminal: 
    docker ps

## 4. Ejecutar el proyecto
   Ejecuta en la terminal:
    go run main.go

## Si todo esta correcto lo que deberias observar es:
    Creación de películas de ejemplo.
    Listado de todas las películas.
    Actualización de una película.
    Búsqueda de una película por ID.
    Borrado de una película.
    
## 5. Para finalizar y que no quede nada en ejecucion, eliminando contenedores y volumenes 
    Ejecutar en la terminal:
    docker compose down -v 

## Integrantes del grupo:
## Acevedo Belen
## Artaza Sheila
