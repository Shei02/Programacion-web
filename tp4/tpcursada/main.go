package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/lib/pq"          // driver PostgreSQL
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

	// Parsear plantillas
	tmpl := template.Must(template.ParseGlob("views/*.templ"))

	// --- Ruta raíz: renderiza la página compuesta ---
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		// Obtener datos desde la BD
		pelis, err := q.ListPelis(r.Context())
		if err != nil {
			log.Printf("error listando peliculas: %v", err)
			http.Error(w, "Error al obtener películas", http.StatusInternalServerError)
			return
		}

		usuarios, err := q.ListUsuario(r.Context())
		if err != nil {
			log.Printf("error listando usuarios: %v", err)
			http.Error(w, "Error al obtener usuarios", http.StatusInternalServerError)
			return
		}

		miras, err := q.ListMira(r.Context())
		if err != nil {
			log.Printf("error listando miras: %v", err)
			http.Error(w, "Error al obtener calificaciones", http.StatusInternalServerError)
			return
		}

		// Estructura de datos para la plantilla
		data := map[string]interface{}{
			"Title":     "Plataforma de películas",
			"Header":    "Plataforma de películas",
			"Intro":     "Bienvenido a la plataforma - gestioná películas, usuarios y calificaciones.",
			"Peliculas": pelis,
			"Usuarios":  usuarios,
			"Miras":     miras,
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
			log.Printf("error ejecutando template: %v", err)
			http.Error(w, "Error interno", http.StatusInternalServerError)
			return
		}
	})

	// Servir archivos estáticos (CSS/JS/otras páginas si aún existen)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// --- Iniciar servidor ---
	port := 8080
	fmt.Printf("🚀 Servidor corriendo en http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
