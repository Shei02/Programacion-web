package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"          // driver PostgreSQL
	"tpe.com/tpcursada/api"        // handlers CRUD
	db "tpe.com/tpcursada/db/sqlc" // módulo SQLC
)

func main() {
	// --- Conexión a la base de datos ---
	connStr := "user=postgres password=postgres dbname=platpelis sslmode=disable"
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error al abrir conexión:", err)
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		log.Fatal("No se pudo conectar a la BD:", err)
	}
	fmt.Println("✅ Conexión exitosa a PostgreSQL")

	q := db.New(conn)
	server := api.NewServer(q)

	// --- Ruta raíz ---
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		file, err := os.ReadFile("index.html")
		if err != nil {
			http.Error(w, "Archivo no encontrado", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(file)
	})

	// --- Rutas PELÍCULAS ---
	http.HandleFunc("/peliculas", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			server.ListPeliculas(w, r)
		case "POST":
			server.CreatePelicula(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/peliculas/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			server.GetPelicula(w, r)
		case "PUT":
			server.UpdatePelicula(w, r)
		case "DELETE":
			server.DeletePelicula(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	// --- Rutas USUARIOS ---
	http.HandleFunc("/usuarios", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			server.ListUsuarios(w, r)
		case "POST":
			server.CreateUsuario(w, r) // dejar que la función haga el Decode
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/usuarios/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			server.GetUsuario(w, r)
		case "PUT":
			server.UpdateUsuario(w, r)
		case "DELETE":
			server.DeleteUsuario(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	// --- Rutas MIRAS ---
	http.HandleFunc("/miras", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			server.ListMiras(w, r)
		case "POST":
			server.CreateMira(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/miras/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			server.GetMira(w, r)
		case "DELETE":
			if r.URL.Path == "/miras/" || r.URL.Path == "/miras" {
				server.DeleteMira(w, r)
			} else {
				server.DeleteMiraByPath(w, r)
			}
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	// --- Iniciar servidor ---
	port := 8080
	fmt.Printf("🚀 Servidor corriendo en http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
