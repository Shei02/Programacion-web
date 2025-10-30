package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

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
	// Servir la página HTML de películas en /peliculas
	http.HandleFunc("/peliculas", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		// Servimos la página estática que mostrará las películas en HTML
		http.ServeFile(w, r, "static/peliculas.html")
	})

	// Endpoints JSON para uso por AJAX/Fetch en /api/peliculas
	http.HandleFunc("/api/peliculas", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			server.ListPeliculas(w, r)
		case "POST":
			server.CreatePelicula(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/peliculas/", func(w http.ResponseWriter, r *http.Request) {
		// clone request and trim /api prefix so handlers that expect paths like
		// /peliculas/{id} continue working unchanged
		r2 := r.Clone(r.Context())
		r2.URL.Path = strings.TrimPrefix(r.URL.Path, "/api")
		switch r.Method {
		case "GET":
			server.GetPelicula(w, r2)
		case "PUT":
			server.UpdatePelicula(w, r2)
		case "DELETE":
			server.DeletePelicula(w, r2)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	// Serve the same static pelicula page for both list and detail paths so the
	// client-side JS can detect an ID in the URL and fetch a single resource.
	http.HandleFunc("/peliculas/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		http.ServeFile(w, r, "static/peliculas.html")
	})

	http.HandleFunc("/pelicula/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		http.ServeFile(w, r, "static/peliculas.html")
	})

	// --- Rutas USUARIOS ---

	// Servir la página HTML de usuarios en /usuarios
	http.HandleFunc("/usuarios", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		http.ServeFile(w, r, "static/usuarios.html")
	})

	// Endpoints JSON para usuarios en /api/usuarios
	http.HandleFunc("/api/usuarios", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			server.ListUsuarios(w, r)
		case "POST":
			server.CreateUsuario(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/usuarios/", func(w http.ResponseWriter, r *http.Request) {
		r2 := r.Clone(r.Context())
		r2.URL.Path = strings.TrimPrefix(r.URL.Path, "/api")
		switch r.Method {
		case "GET":
			server.GetUsuario(w, r2)
		case "PUT":
			server.UpdateUsuario(w, r2)
		case "DELETE":
			server.DeleteUsuario(w, r2)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	// Serve the usuarios page for list and detail so client JS can fetch single usuario
	http.HandleFunc("/usuarios/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		http.ServeFile(w, r, "static/usuarios.html")
	})

	http.HandleFunc("/usuario/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		http.ServeFile(w, r, "static/usuarios.html")
	})

	// --- Rutas MIRAS ---
	// Servir la página HTML de miras en /miras
	http.HandleFunc("/miras", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		http.ServeFile(w, r, "static/miras.html")
	})

	// Endpoints JSON para miras en /api/miras
	http.HandleFunc("/api/miras", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			server.ListMiras(w, r)
		case "POST":
			server.CreateMira(w, r)
		case "DELETE":
			// soportar eliminación por query ?idu=&idp=
			server.DeleteMira(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/miras/", func(w http.ResponseWriter, r *http.Request) {
		r2 := r.Clone(r.Context())
		r2.URL.Path = strings.TrimPrefix(r.URL.Path, "/api")
		switch r.Method {
		case "GET":
			server.GetMira(w, r2)
		case "PUT":
			server.UpdateMira(w, r2)
		case "DELETE":
			// delegar a DeleteMiraByPath si se requiere
			server.DeleteMiraByPath(w, r2)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	// Serve the miras page for list and detail so client JS can fetch miras for a user id
	http.HandleFunc("/miras/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		http.ServeFile(w, r, "static/miras.html")
	})

	http.HandleFunc("/mira/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		http.ServeFile(w, r, "static/miras.html")
	})

	//-------------------conexion archivos de carpeta static---------------
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Servir la página de películas
	//http.HandleFunc("/peliculas", func(w http.ResponseWriter, r *http.Request) {
	//http.ServeFile(w, r, "static/peliculas.html")
	//})

	// Servir la página de miras
	//http.HandleFunc("/miras", func(w http.ResponseWriter, r *http.Request) {
	//   http.ServeFile(w, r, "static/miras.html")
	//})

	// --- Iniciar servidor ---
	port := 8080
	fmt.Printf("🚀 Servidor corriendo en http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
