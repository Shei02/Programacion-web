# Proyecto Películas

Este repositorio contiene una API REST escrita en Go para gestionar películas, usuarios y "miras" (relaciones usuario↔película). El proyecto usa:

- Go
- PostgreSQL (contenerizado con Docker Compose)
- SQLC para generar el acceso a la base de datos

El entregable incluye el código fuente en este repositorio y un conjunto de scripts/targets (`Makefile` y `requests.bash`) para levantar la base, ejecutar el servidor y probar la API de forma reproducible.

## Requisitos

Asegúrate de tener instalados en tu máquina:

- Go 1.21+ (o la versión indicada en `go.mod`)
- Docker
- Docker Compose (v2 `docker compose` o equivalente)
- SQLC (para regenerar el código opcionalmente)

Opcional pero recomendado:
- `jq` o `python3` para formatear JSON en las pruebas (el script `requests.bash` detecta ambos).

## Clonar el repositorio

Clona este repositorio y sitúate en su carpeta antes de ejecutar los pasos siguientes. 

```bash
# usando HTTPS
git clone https://github.com/Shei02/Programacion-web.git

cd <repo>
```

Si ya lo clonaste o estás trabajando en la copia local, simplemente sitúate en la raíz del proyecto:

```bash
cd /ruta/al/proyecto
```

## Atajos con Makefile (recomendado)

El repositorio incluye un `Makefile` con targets que simplifican los pasos comunes. A continuación los comandos principales que deberías ejecutar en este orden desde la raíz del proyecto:

0) Crear el modulo de go
```bash
go mod init tp.com/tpcursada
```
   
1) Hacer el sqlc generate (si no no funciona), no pudimos hacer que se genere al hacer make all (automatico):
```bash
make sqlc
```

1) Levantar la base de datos en un contenedor:
```bash
make db-up
```

2) Esperar a que PostgreSQL esté listo (se revisa dentro del contenedor):
```bash
make db-wait
```

3) (Opcional) Generar código a partir de las consultas SQL con sqlc:

```bash
make sqlc
```

4) Compilar y ejecutar el servidor (o usar `make all` para ejecutar la secuencia completa):
```bash
# compilar
make build

# iniciar el binario en primer plano
make start

# O bien, todo en uno (formatea, tidy, levanta BD, espera, sqlc, compila y corre):
make all
```

5) (Opcional) Poblar datos de ejemplo en la API:
```bash
make seed
```

6) Parar y eliminar contenedores/volúmenes cuando termines:
```bash
make db-down
```

## Probar la API (scripts reproducibles)

Hay un script `requests.bash` que ejecuta las peticiones `curl` sobre los endpoints principales (películas, usuarios, miras).
Ejecutalo desde la raíz del repo:

```bash
bash ./requests.bash
```

También hay objetivos individuales en el `Makefile` para probar casos puntuales (útiles en correcciones rápidas):

- `make test-peliculas-put` — prueba solo el PUT /peliculas/1
- `make test-peliculas-crud` — GET/PUT/DELETE para películas
- `make test-usuarios-crud` — análogo para usuarios
- `make test-miras-crud` — análogo para miras

Ejemplo para probar solo el PUT desde el Makefile:

```bash
make test-usuarios-put
```

## Estructura principal del repo

- `main.go` — router y arranque del servidor
- `api/` — handlers HTTP (peliculas, usuarios, miras)
- `db/` — SQL, código generado por sqlc (`db/sqlc` contiene los archivos generados)
- `requests.bash` — script de pruebas con `curl`
- `Makefile` — objetivos para facilitar ejecución, pruebas y datos de ejemplo
- `docker-compose.yml` — definición del servicio PostgreSQL

## Entrega

El código fuente completo está en este repositorio. Para la entrega, incluí:

- Este `readme.md` con instrucciones de ejecución.
- El `Makefile` con los atajos para levantar DB, compilar y correr las pruebas.
- El script `requests.bash` con ejemplos reproducibles (usa `curl`).

## Integrantes

- Acevedo Belén
- Artaza Sheila

