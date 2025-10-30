#!/usr/bin/env bash
# requests.bash

set -euo pipefail
BASE="http://localhost:8080"
TMP=/tmp/req_resp_$$

show_resp() {
  local f=$1
  if [ -s "$f" ]; then
    # determine first non-whitespace character to guess JSON vs plain text
    first=$(sed -n '1s/^[[:space:]]*//;p' "$f" | head -c1 || true)
    if echo "$first" | grep -q "^[\{\[]"; then
      if command -v jq >/dev/null 2>&1; then
        jq . "$f" || cat "$f"
      elif command -v python3 >/dev/null 2>&1; then
        python3 -m json.tool "$f" || cat "$f"
      else
        cat "$f"
      fi
    else
      cat "$f"
    fi
  else
    echo "<empty body>"
  fi
}

print_status() {
  local code=$1
  local label
  case "$code" in
    200) label="200 OK" ;; 
    201) label="201 Created" ;; 
    204) label="204 No Content" ;; 
    400) label="400 Bad Request" ;; 
    404) label="404 Not Found" ;; 
    409) label="409 Conflict" ;; 
    500) label="500 Internal Server Error" ;; 
    *)   label="$code" ;;
  esac
  echo "=> $label"
}

echo "=== POST /peliculas ==="
cat > $TMP <<'JSON'
{"titulo":"Tenet","duracion":150,"director":"Christopher Nolan","actores":"John David Washington","edadmin":13,"sinopsis":"Inversión temporal","anioestr":2020}
JSON
status=$(curl -sS -o $TMP -w '%{http_code}' -X POST -H "Content-Type: application/json" -d @$TMP "$BASE/peliculas")
print_status $status
show_resp $TMP
echo

echo "=== GET /peliculas ==="
status=$(curl -sS -o $TMP -w '%{http_code}' "$BASE/peliculas")
print_status $status
show_resp $TMP
echo

echo "=== GET /peliculas/1 ==="
status=$(curl -sS -o $TMP -w '%{http_code}' "$BASE/peliculas/1")
print_status $status
show_resp $TMP
echo

echo "=== PUT /peliculas/1 ==="
cat > $TMP <<'JSON'
{"titulo":"Tenet (editada)","duracion":152,"director":"Christopher Nolan","actores":"John David Washington","edadmin":13,"sinopsis":"Inversión temporal (versión editada)","anioestr":2020}
JSON
status=$(curl -sS -o $TMP -w '%{http_code}' -X PUT -H "Content-Type: application/json" -d @$TMP "$BASE/peliculas/1")
print_status $status
show_resp $TMP
echo

echo "=== DELETE /peliculas/1 ==="
# capture body (if any) and code
status=$(curl -sS -o $TMP -w '%{http_code}' -X DELETE "$BASE/peliculas/1")
print_status $status
show_resp $TMP
echo

# ---------- Usuarios ----------

echo "=== POST /usuarios ==="
cat > $TMP <<'JSON'
{"nomusu":"sheila","contrasenia":"123456","email":"sheila@example.com","fechanac":"1990-05-12T00:00:00Z"}
JSON
status=$(curl -sS -o $TMP -w '%{http_code}' -X POST -H "Content-Type: application/json" -d @$TMP "$BASE/usuarios")
print_status $status
show_resp $TMP
echo

echo "=== GET /usuarios ==="
status=$(curl -sS -o $TMP -w '%{http_code}' "$BASE/usuarios")
print_status $status
show_resp $TMP
echo

echo "=== GET /usuarios/1 ==="
status=$(curl -sS -o $TMP -w '%{http_code}' "$BASE/usuarios/1")
print_status $status
show_resp $TMP
echo

echo "=== PUT /usuarios/1 ==="
cat > $TMP <<'JSON'
{"nomusu":"sheila_mod","contrasenia":"654321","email":"sheila+mod@example.com","fechanac":"1990-05-12T00:00:00Z"}
JSON
status=$(curl -sS -o $TMP -w '%{http_code}' -X PUT -H "Content-Type: application/json" -d @$TMP "$BASE/usuarios/1")
print_status $status
show_resp $TMP
echo

echo "=== DELETE /usuarios/1 ==="
status=$(curl -sS -o $TMP -w '%{http_code}' -X DELETE "$BASE/usuarios/1")
print_status $status
show_resp $TMP
echo

# ---------- Miras ----------

# Create a mira (assumes pelicula 1 and usuario 1 exist)
echo "=== POST /miras ==="
cat > $TMP <<'JSON'
{"idp":1,"idu":1,"gustoono":"sí","calif":5}
JSON
status=$(curl -sS -o $TMP -w '%{http_code}' -X POST -H "Content-Type: application/json" -d @$TMP "$BASE/miras")
print_status $status
show_resp $TMP
echo

echo "=== GET /miras ==="
status=$(curl -sS -o $TMP -w '%{http_code}' "$BASE/miras")
print_status $status
show_resp $TMP
echo

echo "=== GET /miras/1 ==="
status=$(curl -sS -o $TMP -w '%{http_code}' "$BASE/miras/1")
print_status $status
show_resp $TMP
echo

echo "=== DELETE /miras/1/1 ==="
status=$(curl -sS -o $TMP -w '%{http_code}' -X DELETE "$BASE/miras/1/1")
print_status $status
show_resp $TMP

echo "Done."
rm -f $TMP
