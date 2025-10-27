# Nombre del ejecutable
BINARY=server

# Variables de Go
GO=go
GOBUILD=$(GO) build
GORUN=$(GO) run
GOTEST=$(GO) test
GOMOD=$(GO) mod

# Docker Compose
DC=docker compose
DB_SERVICE=db
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=platpelis
DB_PORT=5432

# --- Targets ---
.PHONY: all build run start clean test fmt tidy db-up db-down db-wait sqlc seed \
	get-usuarios get-peliculas get-miras add-usuarios add-peliculas add-miras

# --- ALL: todo junto ---
all: fmt tidy db-up db-wait sqlc build start

# --- Base de datos ---
db-up:
	@echo "🐳 Levantando PostgreSQL con Docker..."
	@$(DC) up -d

db-down:
	@echo "🛑 Deteniendo PostgreSQL..."
	@$(DC) down -v --remove-orphans

db-wait:
	@echo "⏳ Esperando a que la base de datos esté lista (dentro del contenedor)..."
	@docker exec postgres_peliculas bash -c 'until pg_isready -U postgres > /dev/null 2>&1; do sleep 1; done'
	@echo "✅ Base de datos lista dentro del contenedor."


# --- SQLC ---
sqlc:
	@echo "📝 Generando código SQLC..."
	@sqlc generate

# --- Go ---
build: 
	@echo "🔨 Compilando $(BINARY)..."
	@$(GOBUILD) -o $(BINARY) main.go
	@echo "✅ Compilación finalizada."

run:
	@echo "🚀 Ejecutando API REST..."
	@$(GORUN) main.go

start: build
	@echo "🚀 Ejecutando $(BINARY)..."
	@./$(BINARY)

fmt:
	@echo "🎨 Formateando código..."
	@$(GO) fmt ./...

tidy:
	@echo "📦 Actualizando dependencias..."
	@$(GOMOD) tidy

test:
	@echo "🧪 Corriendo tests..."
	@$(GOTEST) ./...

clean:
	@echo "🧹 Limpiando proyecto..."
	@rm -f $(BINARY)

# --- Agregar y listar datos usando la API ---

# Listar recursos
get-usuarios:
	@echo "📋 Listando usuarios..."
	@curl -s http://localhost:8080/usuarios | sh -c 'command -v jq >/dev/null 2>&1 && jq . || (command -v python3 >/dev/null 2>&1 && python3 -m json.tool || cat)'
	@echo

get-peliculas:
	@echo "📋 Listando películas..."
	@curl -s http://localhost:8080/peliculas | sh -c 'command -v jq >/dev/null 2>&1 && jq . || (command -v python3 >/dev/null 2>&1 && python3 -m json.tool || cat)'
	@echo

get-miras:
	@echo "📋 Listando miras del usuario 1..."
	@curl -s "http://localhost:8080/miras?idu=1" | sh -c 'command -v jq >/dev/null 2>&1 && jq . || (command -v python3 >/dev/null 2>&1 && python3 -m json.tool || cat)'

	@echo "📋 Listando miras del usuario 2..."
	@curl -s "http://localhost:8080/miras?idu=2" | sh -c 'command -v jq >/dev/null 2>&1 && jq . || (command -v python3 >/dev/null 2>&1 && python3 -m json.tool || cat)'

	@echo "📋 Listando miras del usuario 3..."
	@curl -s "http://localhost:8080/miras?idu=3" | sh -c 'command -v jq >/dev/null 2>&1 && jq . || (command -v python3 >/dev/null 2>&1 && python3 -m json.tool || cat)'

# Agregar datos de ejemplo

add-usuarios:
	@echo "🌱 Agregando usuarios de ejemplo..."
	@curl -sS -X POST http://localhost:8080/usuarios \
	-H "Content-Type: application/json" \
	-d '{"nomusu":"Sheila","contrasenia":"1234","email":"sheila@test.com","fechanac":"1990-05-12T00:00:00Z"}' | sh -c 'command -v jq >/dev/null 2>&1 && jq . || (command -v python3 >/dev/null 2>&1 && python3 -m json.tool || cat)'

	@curl -sS -X POST http://localhost:8080/usuarios \
	-H "Content-Type: application/json" \
	-d '{"nomusu":"Esteban","contrasenia":"abcd","email":"esteban@test.com","fechanac":"1985-11-23T00:00:00Z"}' | sh -c 'command -v jq >/dev/null 2>&1 && jq . || (command -v python3 >/dev/null 2>&1 && python3 -m json.tool || cat)'

	@curl -sS -X POST http://localhost:8080/usuarios \
	-H "Content-Type: application/json" \
	-d '{"nomusu":"Belen","contrasenia":"1234","email":"belen@test.com","fechanac":"1958-05-12T00:00:00Z"}' | sh -c 'command -v jq >/dev/null 2>&1 && jq . || (command -v python3 >/dev/null 2>&1 && python3 -m json.tool || cat)'

	@curl -sS -X POST http://localhost:8080/usuarios \
	-H "Content-Type: application/json" \
	-d '{"nomusu":"Eden","contrasenia":"abcd","email":"eden@test.com","fechanac":"2002-11-23T00:00:00Z"}' | sh -c 'command -v jq >/dev/null 2>&1 && jq . || (command -v python3 >/dev/null 2>&1 && python3 -m json.tool || cat)'


add-peliculas:
	@echo "🌱 Agregando películas de ejemplo..."
	@curl -sS -X POST http://localhost:8080/peliculas \
	-H "Content-Type: application/json" \
	-d '{"titulo":"Matrix","duracion":"136min","director":"Lana Wachowski","actores":"Keanu Reeves, Laurence Fishburne","edadmin":13,"sinopsis":"Un hacker descubre la realidad simulada en la que vive.","anioestr":1999}' | sh -c 'command -v jq >/dev/null 2>&1 && jq . || (command -v python3 >/dev/null 2>&1 && python3 -m json.tool || cat)'

	@curl -sS -X POST http://localhost:8080/peliculas \
	-H "Content-Type: application/json" \
	-d '{"titulo":"Inception","duracion":"148min","director":"Christopher Nolan","actores":"Leonardo DiCaprio, Joseph Gordon-Levitt","edadmin":13,"sinopsis":"Un ladrón especializado en sueños roba secretos de la mente.","anioestr":2010}' | sh -c 'command -v jq >/dev/null 2>&1 && jq . || (command -v python3 >/dev/null 2>&1 && python3 -m json.tool || cat)'

	@curl -sS -X POST http://localhost:8080/peliculas \
	-H "Content-Type: application/json" \
	-d '{"titulo":"Interstellar","duracion":"169min","director":"Christopher Nolan","actores":"Matthew McConaughey, Anne Hathaway","edadmin":13,"sinopsis":"Un grupo de exploradores viaja a través de un agujero de gusano en busca de un nuevo hogar para la humanidad.","anioestr":2014}' | sh -c 'command -v jq >/dev/null 2>&1 && jq . || (command -v python3 >/dev/null 2>&1 && python3 -m json.tool || cat)'

	@curl -sS -X POST http://localhost:8080/peliculas \
	-H "Content-Type: application/json" \
	-d '{"titulo":"The Godfather","duracion":"175min","director":"Francis Ford Coppola","actores":"Marlon Brando, Al Pacino","edadmin":18,"sinopsis":"La historia de la familia mafiosa Corleone.","anioestr":1972}' | sh -c 'command -v jq >/dev/null 2>&1 && jq . || (command -v python3 >/dev/null 2>&1 && python3 -m json.tool || cat)'


add-miras:
	@echo "🌱 Agregando miras de ejemplo..."
	@curl -sS -X POST http://localhost:8080/miras \
	-H "Content-Type: application/json" \
	-d '{"idp":1,"idu":1,"gustoono":"sí","calif":5}' | sh -c 'command -v jq >/dev/null 2>&1 && jq . || (command -v python3 >/dev/null 2>&1 && python3 -m json.tool || cat)'

	@curl -sS -X POST http://localhost:8080/miras \
	-H "Content-Type: application/json" \
	-d '{"idp":2,"idu":2,"gustoono":"no","calif":1}' | sh -c 'command -v jq >/dev/null 2>&1 && jq . || (command -v python3 >/dev/null 2>&1 && python3 -m json.tool || cat)'

	@curl -sS -X POST http://localhost:8080/miras \
	-H "Content-Type: application/json" \
	-d '{"idp":3,"idu":1,"gustoono":"sí","calif":5}' | sh -c 'command -v jq >/dev/null 2>&1 && jq . || (command -v python3 >/dev/null 2>&1 && python3 -m json.tool || cat)'

	@curl -sS -X POST http://localhost:8080/miras \
	-H "Content-Type: application/json" \
	-d '{"idp":4,"idu":2,"gustoono":"no","calif":3}' | sh -c 'command -v jq >/dev/null 2>&1 && jq . || (command -v python3 >/dev/null 2>&1 && python3 -m json.tool || cat)'

# Target que agrega todos los datos de ejemplo
seed: add-usuarios add-peliculas add-miras
	@echo "✅ Datos de ejemplo agregados."

# --- Pruebas CRUD (curl) ---
test-peliculas-get:
	@echo "🔎 GET /peliculas/1"
	@curl -sS http://localhost:8080/peliculas/1 | sh -c 'command -v jq >/dev/null 2>&1 && jq . || (command -v python3 >/dev/null 2>&1 && python3 -m json.tool || cat)'
	@echo

test-peliculas-put:
	@echo "✏️  PUT /peliculas/1"
	@curl -sS -o /tmp/make_resp.peliculas -w 'HTTP_CODE:%{http_code}\n' -X PUT http://localhost:8080/peliculas/1 \
	-H "Content-Type: application/json" \
	-d '{"titulo":"Titulo modificado","duracion":"120min","director":"Mod","actores":"A","edadmin":13,"sinopsis":"mod","anioestr":2025}'
	@if [ -s /tmp/make_resp.peliculas ]; then \
		first=$$(head -c1 /tmp/make_resp.peliculas 2>/dev/null || echo); \
		if echo "$$first" | grep -q "^[\{\[]"; then \
			command -v jq >/dev/null 2>&1 && jq . /tmp/make_resp.peliculas || (command -v python3 >/dev/null 2>&1 && python3 -m json.tool /tmp/make_resp.peliculas || cat /tmp/make_resp.peliculas); \
		else \
			cat /tmp/make_resp.peliculas; \
		fi; \
	fi
	@rm -f /tmp/make_resp.peliculas
	@echo

test-peliculas-delete:
	@echo "🗑️  DELETE /peliculas/1"
	@curl -sS -o /tmp/make_resp.peliculas -w 'HTTP_CODE:%{http_code}\n' -X DELETE http://localhost:8080/peliculas/1
	@if [ -s /tmp/make_resp.peliculas ]; then \
		first=$$(head -c1 /tmp/make_resp.peliculas 2>/dev/null || echo); \
		if echo "$$first" | grep -q "^[\{\[]"; then \
			command -v jq >/dev/null 2>&1 && jq . /tmp/make_resp.peliculas || (command -v python3 >/dev/null 2>&1 && python3 -m json.tool /tmp/make_resp.peliculas || cat /tmp/make_resp.peliculas); \
		else \
			cat /tmp/make_resp.peliculas; \
		fi; \
	fi
	@rm -f /tmp/make_resp.peliculas
	@echo

test-peliculas-crud: test-peliculas-get test-peliculas-put test-peliculas-delete

test-usuarios-get:
	@echo "🔎 GET /usuarios/1"
	@curl -sS http://localhost:8080/usuarios/1 | sh -c 'command -v jq >/dev/null 2>&1 && jq . || (command -v python3 >/dev/null 2>&1 && python3 -m json.tool || cat)'
	@echo

test-usuarios-put:
	@echo "✏️  PUT /usuarios/1"
	@curl -sS -o /tmp/make_resp.usuarios -w 'HTTP_CODE:%{http_code}\n' -X PUT http://localhost:8080/usuarios/1 \
	-H "Content-Type: application/json" \
	-d '{"nomusu":"UsuarioMod","contrasenia":"x","email":"mod@test.com","fechanac":"1990-01-01T00:00:00Z"}'
	@if [ -s /tmp/make_resp.usuarios ]; then \
		first=$$(head -c1 /tmp/make_resp.usuarios 2>/dev/null || echo); \
		if echo "$$first" | grep -q "^[\{\[]"; then \
			command -v jq >/dev/null 2>&1 && jq . /tmp/make_resp.usuarios || (command -v python3 >/dev/null 2>&1 && python3 -m json.tool /tmp/make_resp.usuarios || cat /tmp/make_resp.usuarios); \
		else \
			cat /tmp/make_resp.usuarios; \
		fi; \
	fi
	@rm -f /tmp/make_resp.usuarios
	@echo

test-usuarios-delete:
	@echo "🗑️  DELETE /usuarios/1"
	@curl -sS -o /tmp/make_resp.usuarios -w 'HTTP_CODE:%{http_code}\n' -X DELETE http://localhost:8080/usuarios/1
	@if [ -s /tmp/make_resp.usuarios ]; then \
		first=$$(head -c1 /tmp/make_resp.usuarios 2>/dev/null || echo); \
		if echo "$$first" | grep -q "^[\{\[]"; then \
			command -v jq >/dev/null 2>&1 && jq . /tmp/make_resp.usuarios || (command -v python3 >/dev/null 2>&1 && python3 -m json.tool /tmp/make_resp.usuarios || cat /tmp/make_resp.usuarios); \
		else \
			cat /tmp/make_resp.usuarios; \
		fi; \
	fi
	@rm -f /tmp/make_resp.usuarios
	@echo

test-usuarios-crud: test-usuarios-get test-usuarios-put test-usuarios-delete

test-miras-get:
	@echo "🔎 GET /miras/1/1"
	@curl -sS http://localhost:8080/miras/1/1 | sh -c 'command -v jq >/dev/null 2>&1 && jq . || (command -v python3 >/dev/null 2>&1 && python3 -m json.tool || cat)'
	@echo

test-miras-put:
	@echo "✏️  PUT /miras/1/1 (si no existe, mostrará código)"
	@curl -sS -o /tmp/make_resp.miras -w 'HTTP_CODE:%{http_code}\n' -X PUT http://localhost:8080/miras/1/1 \
	-H "Content-Type: application/json" \
	-d '{"idp":1,"idu":1,"gustoono":"sí","calif":4}'
	@if [ -s /tmp/make_resp.miras ]; then \
		first=$$(head -c1 /tmp/make_resp.miras 2>/dev/null || echo); \
		if echo "$$first" | grep -q "^[\{\[]"; then \
			command -v jq >/dev/null 2>&1 && jq . /tmp/make_resp.miras || (command -v python3 >/dev/null 2>&1 && python3 -m json.tool /tmp/make_resp.miras || cat /tmp/make_resp.miras); \
		else \
			cat /tmp/make_resp.miras; \
		fi; \
	fi
	@rm -f /tmp/make_resp.miras
	@echo

test-miras-delete:
	@echo "🗑️  DELETE /miras/1/1"
	@curl -sS -o /tmp/make_resp.miras -w 'HTTP_CODE:%{http_code}\n' -X DELETE http://localhost:8080/miras/1/1
	@if [ -s /tmp/make_resp.miras ]; then \
		first=$$(head -c1 /tmp/make_resp.miras 2>/dev/null || echo); \
		if echo "$$first" | grep -q "^[\{\[]"; then \
			command -v jq >/dev/null 2>&1 && jq . /tmp/make_resp.miras || (command -v python3 >/dev/null 2>&1 && python3 -m json.tool /tmp/make_resp.miras || cat /tmp/make_resp.miras); \
		else \
			cat /tmp/make_resp.miras; \
		fi; \
	fi
	@rm -f /tmp/make_resp.miras
	@echo

test-miras-crud: test-miras-get test-miras-put test-miras-delete

test-crud: test-peliculas-crud test-usuarios-crud test-miras-crud
