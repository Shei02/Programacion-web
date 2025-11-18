package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/lib/pq"          // driver PostgreSQL
	db "tpe.com/tpcursada/db/sqlc" // módulo SQLC
	"tpe.com/tpcursada/views"
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

	// Parsear plantillas (embebidas)
	// use views.ParseTemplates() which parses templates embedded via go:embed
	tmpl := views.ParseTemplates()

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

	// --- Páginas separadas para cada entidad (mejor UX) ---
	http.HandleFunc("/peliculas", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			pelis, err := q.ListPelis(r.Context())
			if err != nil {
				log.Printf("error listando peliculas: %v", err)
				http.Error(w, "Error al obtener películas", http.StatusInternalServerError)
				return
			}
			data := map[string]interface{}{
				"Title":     "Películas",
				"Header":    "Películas",
				"Peliculas": pelis,
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			if err := views.ParseTemplates().ExecuteTemplate(w, "layout", data); err != nil {
				log.Printf("error ejecutando template peliculas: %v", err)
				http.Error(w, "Error interno", http.StatusInternalServerError)
			}
		case http.MethodPost:
			// Creation handled earlier by /peliculas POST handler (PRG)
			http.Error(w, "Método no permitido aquí", http.StatusMethodNotAllowed)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	// Editar película (render form con datos)
	http.HandleFunc("/peliculas/editar", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "id requerido", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "id inválido", http.StatusBadRequest)
			return
		}
		peli, err := q.GetPeli(r.Context(), int32(id))
		if err != nil {
			log.Printf("error obteniendo pelicula: %v", err)
			http.Error(w, "No encontrado", http.StatusNotFound)
			return
		}
		data := map[string]interface{}{
			"Title":   "Editar Película",
			"Header":  "Editar Película",
			"Pelicula": peli,
			"Action":  "/peliculas/edit",
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := views.ParseTemplates().ExecuteTemplate(w, "layout", data); err != nil {
			log.Printf("error ejecutando template pelicula editar: %v", err)
			http.Error(w, "Error interno", http.StatusInternalServerError)
		}
	})

	// Handler edit POST
	http.HandleFunc("/peliculas/edit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		idp, _ := strconv.Atoi(r.FormValue("idp"))
		dur, _ := strconv.Atoi(r.FormValue("duracion"))
		ed, _ := strconv.Atoi(r.FormValue("edadmin"))
		anio, _ := strconv.Atoi(r.FormValue("anioestr"))
		params := db.UpdatePeliParams{
			Idp:      int32(idp),
			Titulo:   r.FormValue("titulo"),
			Duracion: int32(dur),
			Director: r.FormValue("director"),
			Actores:  r.FormValue("actores"),
			Edadmin:  int32(ed),
			Sinopsis: r.FormValue("sinopsis"),
			Anioestr: int32(anio),
		}
		if err := q.UpdatePeli(r.Context(), params); err != nil {
			log.Printf("error actualizando pelicula: %v", err)
			http.Error(w, "Error interno", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/peliculas", http.StatusSeeOther)
	})

	// Delete pelicula
	http.HandleFunc("/peliculas/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		idp, _ := strconv.Atoi(r.FormValue("idp"))
		if err := q.DeletePeli(r.Context(), int32(idp)); err != nil {
			log.Printf("error eliminando pelicula: %v", err)
			http.Error(w, "Error interno", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/peliculas", http.StatusSeeOther)
	})

	// Usuarios: GET page
	http.HandleFunc("/usuarios", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			usuarios, err := q.ListUsuario(r.Context())
			if err != nil {
				log.Printf("error listando usuarios: %v", err)
				http.Error(w, "Error al obtener usuarios", http.StatusInternalServerError)
				return
			}
			data := map[string]interface{}{
				"Title":    "Usuarios",
				"Header":   "Usuarios",
				"Usuarios": usuarios,
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			if err := views.ParseTemplates().ExecuteTemplate(w, "layout", data); err != nil {
				log.Printf("error ejecutando template usuarios: %v", err)
				http.Error(w, "Error interno", http.StatusInternalServerError)
			}
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	// Usuarios edit/delete handlers
	http.HandleFunc("/usuarios/editar", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "id requerido", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "id inválido", http.StatusBadRequest)
			return
		}
		usu, err := q.GetUsuario(r.Context(), int32(id))
		if err != nil {
			log.Printf("error obteniendo usuario: %v", err)
			http.Error(w, "No encontrado", http.StatusNotFound)
			return
		}
		data := map[string]interface{}{
			"Title":   "Editar Usuario",
			"Header":  "Editar Usuario",
			"Usuario": usu,
			"Action":  "/usuarios/edit",
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := views.ParseTemplates().ExecuteTemplate(w, "layout", data); err != nil {
			log.Printf("error ejecutando template usuario editar: %v", err)
			http.Error(w, "Error interno", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/usuarios/edit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		idu, _ := strconv.Atoi(r.FormValue("idu"))
		fechaStr := r.FormValue("fechanac")
		var fecha time.Time
		if fechaStr != "" {
			fecha, _ = time.Parse("2006-01-02", fechaStr)
		}
		params := db.UpdateUsuarioParams{
			Idu:         int32(idu),
			Nomusu:      r.FormValue("nomusu"),
			Contrasenia: r.FormValue("contrasenia"),
			Email:       r.FormValue("email"),
			Fechanac:    fecha,
		}
		if err := q.UpdateUsuario(r.Context(), params); err != nil {
			log.Printf("error actualizando usuario: %v", err)
			http.Error(w, "Error interno", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/usuarios", http.StatusSeeOther)
	})

	http.HandleFunc("/usuarios/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		idu, _ := strconv.Atoi(r.FormValue("idu"))
		if err := q.DeleteUsuario(r.Context(), int32(idu)); err != nil {
			log.Printf("error eliminando usuario: %v", err)
			http.Error(w, "Error interno", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/usuarios", http.StatusSeeOther)
	})

	// Miras page
	http.HandleFunc("/miras", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		pelis, err := q.ListPelis(r.Context())
		if err != nil {
			log.Printf("error listando peliculas para miras: %v", err)
			http.Error(w, "Error interno", http.StatusInternalServerError)
			return
		}
		usuarios, err := q.ListUsuario(r.Context())
		if err != nil {
			log.Printf("error listando usuarios para miras: %v", err)
			http.Error(w, "Error interno", http.StatusInternalServerError)
			return
		}
		miras, err := q.ListMira(r.Context())
		if err != nil {
			log.Printf("error listando miras: %v", err)
			http.Error(w, "Error interno", http.StatusInternalServerError)
			return
		}
		data := map[string]interface{}{
			"Title":     "Calificaciones",
			"Header":    "Calificaciones",
			"Peliculas": pelis,
			"Usuarios":  usuarios,
			"Miras":     miras,
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := views.ParseTemplates().ExecuteTemplate(w, "layout", data); err != nil {
			log.Printf("error ejecutando template miras: %v", err)
			http.Error(w, "Error interno", http.StatusInternalServerError)
		}
	})

	// Delete mira by user (sqlc DeleteMira deletes by Idu)
	http.HandleFunc("/miras/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		idu, _ := strconv.Atoi(r.FormValue("idu"))
		if err := q.DeleteMira(r.Context(), int32(idu)); err != nil {
			log.Printf("error eliminando miras: %v", err)
			http.Error(w, "Error interno", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/miras", http.StatusSeeOther)
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
