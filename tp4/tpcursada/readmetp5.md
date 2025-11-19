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
```markdown

# Proyecto Películas

Este repositorio contiene una pequeña API en Go que gestiona películas, usuarios y "miras" (calificaciones). Las pruebas de integración usan Docker (PostgreSQL) y archivos de `requests.hurl` para validar los endpoints.

**Requisitos (mínimos)**
- `git`
- `go` (1.21+)
- `docker` y `docker compose`
- `make`
- `sqlc` (opcional si vas a regenerar el código SQLC, se puede usar - -  `make sqlc`)
- `curl` (recomendado)
- `hurl` (opcional: sólo necesario si quieres ejecutar las pruebas automáticas `requests.hurl`)

Resumen rápido de comandos (rápido, para ejecutar ahora)
```bash
git clone https://github.com/Shei02/Programacion-web.git
cd Programacion-web
git fetch origin feature/tp5-templates
git checkout feature/tp5-templates
cd tp4/tpcursada
make all
```

1) Preparar el entorno Go
```bash
go mod download
go mod tidy
```
2) Generar código desde SQL (si haces cambios en `db/queries`)
```bash
make sqlc
```
3) Ejecutar todo (recomendado)
```bash
make all
```

La API quedará disponible en: `http://localhost:8080`

Comandos útiles
- Parar y eliminar contenedores y volúmenes (borra datos):
```bash
make db-down
```
- Ejecutar sólo las pruebas Hurl (si ya instalaste `hurl`):
```bash
hurl requests.hurl
```

Notas y detalles importantes
- Puerto: la API escucha en `http://localhost:8080`.
- Base de datos: se crea un volumen Docker para PostgreSQL. Ejecutar `make db-down` elimina ese volumen (datos borrados).

Rutas útiles (ejemplos)
- UI (si el proyecto sirve páginas):
  - `http://localhost:8080/`
  - `http://localhost:8080/peliculas`
  - `http://localhost:8080/usuarios`
  - `http://localhost:8080/miras`
- API (JSON):
  - `GET /api/peliculas`
  - `GET /api/peliculas/1`
  - `GET /api/usuarios/1`
  - `GET /api/miras?idu=2`

Integrantes del grupo:

- Acevedo Belen
- Artaza Sheila

```
