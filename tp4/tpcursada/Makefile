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
.PHONY: all build run start clean test fmt tidy db-up db-down db-wait sqlc \

.PHONY: test-hurl

# --- ALL: todo junto ---
all: fmt tidy db-up db-wait sqlc app-wait test-hurl 
	@echo "🚀 Todo listo. El servidor está corriendo dentro del contenedor en http://localhost:8080"

# --- Esperar a que la API web esté lista antes de ejecutar seed ---
.PHONY: app-wait
app-wait:
	@echo "⏳ Esperando a que la API esté lista en http://localhost:8080..."
	@i=0; \
	while true; do \
		code=$$(curl -s -o /dev/null -w '%{http_code}' http://localhost:8080/api/usuarios 2>/dev/null || echo 000); \
		if [ "$$code" = "200" ]; then echo "✅ API lista"; break; fi; \
		i=$$((i+1)); \
		if [ $$i -ge 60 ]; then echo "❌ Timeout esperando la API"; exit 1; fi; \
		sleep 1; \
	done

# --- Templ (compilador de plantillas a Go) ---
.PHONY: templ-install templ-generate

templ-install:
	@echo "🔧 Instalando 'templ' (si no está instalado)..."
	@go install github.com/a-h/templ/cmd/templ@latest

templ-generate: templ-install
	@echo "📦 Generando plantillas compiladas con 'templ'..."
	@templ generate ./views -o views/templates_gen.go -p views || \
	(echo "Advertencia: 'templ generate' falló o 'templ' no está en PATH; asegurate de tener la herramienta instalada"; exit 1)
	@echo "✅ Plantillas generadas en views/templates_gen.go"

# --- Base de datos ---
db-up:
	@echo "🐳 Levantando servicios con Docker..."
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

# --- Pruebas CRUD (curl) ---
test-peliculas-get:
	@echo "🔎 GET /peliculas/1"
	@curl -sS http://localhost:8080/api/peliculas/1 | sh -c 'command -v jq >/dev/null 2>&1 && jq . || (command -v python3 >/dev/null 2>&1 && python3 -m json.tool || cat)'
	@echo

test-peliculas-put:
	@echo "✏️  PUT /peliculas/1"
	@curl -sS -o /tmp/make_resp.peliculas -w 'HTTP_CODE:%{http_code}\n' -X PUT http://localhost:8080/peliculas/1 \
	@curl -sS -o /tmp/make_resp.peliculas -w 'HTTP_CODE:%{http_code}\n' -X PUT http://localhost:8080/api/peliculas/1 \
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
	@curl -sS -o /tmp/make_resp.peliculas -w 'HTTP_CODE:%{http_code}\n' -X DELETE http://localhost:8080/api/peliculas/1
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
	@curl -sS http://localhost:8080/api/usuarios/1 | sh -c 'command -v jq >/dev/null 2>&1 && jq . || (command -v python3 >/dev/null 2>&1 && python3 -m json.tool || cat)'
	@echo

test-usuarios-put:
	@echo "✏️  PUT /usuarios/1"
	@curl -sS -o /tmp/make_resp.usuarios -w 'HTTP_CODE:%{http_code}\n' -X PUT http://localhost:8080/usuarios/1 \
	@curl -sS -o /tmp/make_resp.usuarios -w 'HTTP_CODE:%{http_code}\n' -X PUT http://localhost:8080/api/usuarios/1 \
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
	@curl -sS -o /tmp/make_resp.usuarios -w 'HTTP_CODE:%{http_code}\n' -X DELETE http://localhost:8080/api/usuarios/1
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
	@curl -sS http://localhost:8080/api/miras/1/1 | sh -c 'command -v jq >/dev/null 2>&1 && jq . || (command -v python3 >/dev/null 2>&1 && python3 -m json.tool || cat)'
	@echo

test-miras-put:
	@echo "✏️  PUT /miras/1/1 (si no existe, mostrará código)"
	@curl -sS -o /tmp/make_resp.miras -w 'HTTP_CODE:%{http_code}\n' -X PUT http://localhost:8080/miras/1/1 \
	@curl -sS -o /tmp/make_resp.miras -w 'HTTP_CODE:%{http_code}\n' -X PUT http://localhost:8080/api/miras/1/1 \
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
	@curl -sS -o /tmp/make_resp.miras -w 'HTTP_CODE:%{http_code}\n' -X DELETE http://localhost:8080/api/miras/1/1
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

test-hurl:
	@echo "🔬 Ejecutando pruebas Hurl (requests.hurl)..."
	@command -v hurl >/dev/null 2>&1 || (echo "Error: 'hurl' no está instalado. Instalalo desde https://github.com/Orange-OpenSource/hurl"; exit 1)
	@hurl requests.hurl
