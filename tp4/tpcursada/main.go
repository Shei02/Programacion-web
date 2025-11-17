package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

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

	// --- Handlers para procesar formularios (PRG pattern) ---
	http.HandleFunc("/peliculas", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		if err := r.ParseForm(); err != nil {
			log.Printf("error parseando form pelicula: %v", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		titulo := r.FormValue("titulo")
		duracion, _ := strconv.Atoi(r.FormValue("duracion"))
		director := r.FormValue("director")
		actores := r.FormValue("actores")
		edadmin, _ := strconv.Atoi(r.FormValue("edadmin"))
		anioestr, _ := strconv.Atoi(r.FormValue("anioestr"))
		sinopsis := r.FormValue("sinopsis")

		_, err := q.CreatePeli(r.Context(), db.CreatePeliParams{
			Titulo:   titulo,
			Duracion: int32(duracion),
			Director: director,
			Actores:  actores,
			Edadmin:  int32(edadmin),
			Sinopsis: sinopsis,
			Anioestr: int32(anioestr),
		})
		if err != nil {
			log.Printf("error creando pelicula: %v", err)
			http.Error(w, "Error interno", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/usuarios", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		if err := r.ParseForm(); err != nil {
			log.Printf("error parseando form usuario: %v", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		nom := r.FormValue("nomusu")
		pass := r.FormValue("contrasenia")
		email := r.FormValue("email")
		fechaStr := r.FormValue("fechanac")
		var fecha time.Time
		if fechaStr != "" {
			fecha, _ = time.Parse("2006-01-02", fechaStr)
		}

		_, err := q.CreateUsuario(r.Context(), db.CreateUsuarioParams{
			Nomusu:      nom,
			Contrasenia: pass,
			Email:       email,
			Fechanac:    fecha,
		})
		if err != nil {
			log.Printf("error creando usuario: %v", err)
			http.Error(w, "Error interno", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/miras", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		if err := r.ParseForm(); err != nil {
			log.Printf("error parseando form mira: %v", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		idp, _ := strconv.Atoi(r.FormValue("idp"))
		idu, _ := strconv.Atoi(r.FormValue("idu"))
		gustoono := r.FormValue("gustoono")
		calif, _ := strconv.Atoi(r.FormValue("calif"))

		err := q.CreateMira(r.Context(), db.CreateMiraParams{
			Idp:      int32(idp),
			Idu:      int32(idu),
			Gustoono: gustoono,
			Calif:    int32(calif),
		})
		if err != nil {
			log.Printf("error creando mira: %v", err)
			http.Error(w, "Error interno", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	// --- Iniciar servidor ---
	port := 8080
	fmt.Printf("🚀 Servidor corriendo en http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
