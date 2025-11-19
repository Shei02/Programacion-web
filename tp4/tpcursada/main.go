package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/lib/pq" // driver PostgreSQL
	api "tpe.com/tpcursada/api"
	db "tpe.com/tpcursada/db/sqlc" // módulo SQLC
	"tpe.com/tpcursada/views"
)

func main() {
	// --- Conexión a la base de datos ---
	connStr := "host=db port=5432 user=postgres password=postgres dbname=platpelis sslmode=disable"
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

	// --- API JSON handlers (montados en /api/*) ---
	srv := api.NewServer(q)

	// Peliculas API
	http.HandleFunc("/api/peliculas", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			srv.ListPeliculas(w, r)
		case http.MethodPost:
			srv.CreatePelicula(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/api/peliculas/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			srv.GetPelicula(w, r)
		case http.MethodPut:
			srv.UpdatePelicula(w, r)
		case http.MethodDelete:
			srv.DeletePelicula(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	// Usuarios API
	http.HandleFunc("/api/usuarios", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			srv.ListUsuarios(w, r)
		case http.MethodPost:
			srv.CreateUsuario(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/api/usuarios/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			srv.GetUsuario(w, r)
		case http.MethodPut:
			srv.UpdateUsuario(w, r)
		case http.MethodDelete:
			srv.DeleteUsuario(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	// Miras API
	http.HandleFunc("/api/miras", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			srv.ListMiras(w, r)
		case http.MethodPost:
			srv.CreateMira(w, r)
		case http.MethodDelete:
			srv.DeleteMira(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/api/miras/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			srv.GetMira(w, r)
		case http.MethodPut:
			srv.UpdateMira(w, r)
		case http.MethodDelete:
			srv.DeleteMiraByPath(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

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
			"Content":   "index",
			"BodyClass": "page-inicio",
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
			log.Printf("error ejecutando template: %v", err)
			http.Error(w, "No se puede borrar pelicula/usuario porque tiene referencia en la tabla mira", http.StatusInternalServerError)
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
				"Content":   "peliculas",
				"BodyClass": "page-peliculas",
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			if err := views.ParseTemplates().ExecuteTemplate(w, "layout", data); err != nil {
				log.Printf("error ejecutando template peliculas: %v", err)
				http.Error(w, "No se puede borrar pelicula/usuario porque tiene referencia en la tabla mira", http.StatusInternalServerError)
			}
		case http.MethodPost:
			// Crear película (manejo POST en la misma ruta, patrón PRG)
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
				http.Error(w, "No se puede borrar pelicula/usuario porque tiene referencia en la tabla mira", http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
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
		// Also fetch the full list so the page still shows the list alongside the edit form
		pelis, err := q.ListPelis(r.Context())
		if err != nil {
			log.Printf("error listando peliculas para editar: %v", err)
			http.Error(w, "No se puede borrar pelicula/usuario porque tiene referencia en la tabla mira", http.StatusInternalServerError)
			return
		}
		data := map[string]interface{}{
			"Title":     "Editar Película",
			"Header":    "Editar Película",
			"Pelicula":  peli,
			"Peliculas": pelis,
			"Action":    "/peliculas/edit",
			"Content":   "peliculas",
			"BodyClass": "page-peliculas",
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := views.ParseTemplates().ExecuteTemplate(w, "layout", data); err != nil {
			log.Printf("error ejecutando template pelicula editar: %v", err)
			http.Error(w, "No se puede borrar pelicula/usuario porque tiene referencia en la tabla mira", http.StatusInternalServerError)
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
			http.Error(w, "No se puede borrar pelicula/usuario porque tiene referencia en la tabla mira", http.StatusInternalServerError)
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
			http.Error(w, "No se puede borrar pelicula/usuario porque tiene referencia en la tabla mira", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/peliculas", http.StatusSeeOther)
	})

	// Usuarios: GET page and POST (crear usuario)
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
				"Title":     "Usuarios",
				"Header":    "Usuarios",
				"Usuarios":  usuarios,
				"Content":   "usuarios",
				"BodyClass": "page-usuarios",
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			if err := views.ParseTemplates().ExecuteTemplate(w, "layout", data); err != nil {
				log.Printf("error ejecutando template usuarios: %v", err)
				http.Error(w, "No se puede borrar pelicula/usuario porque tiene referencia en la tabla mira", http.StatusInternalServerError)
			}
		case http.MethodPost:
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
				http.Error(w, "No se puede borrar pelicula/usuario porque tiene referencia en la tabla mira", http.StatusInternalServerError)
				return
			}
			// Mantenerse en la pestaña de usuarios después de crear (no volver al inicio)
			http.Redirect(w, r, "/usuarios", http.StatusSeeOther)
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
		// Also fetch the full usuarios list so the page shows the list alongside the edit form
		usuarios, err := q.ListUsuario(r.Context())
		if err != nil {
			log.Printf("error listando usuarios para editar: %v", err)
			http.Error(w, "No se puede borrar pelicula/usuario porque tiene referencia en la tabla mira", http.StatusInternalServerError)
			return
		}
		data := map[string]interface{}{
			"Title":     "Editar Usuario",
			"Header":    "Editar Usuario",
			"Usuario":   usu,
			"Usuarios":  usuarios,
			"Action":    "/usuarios/edit",
			"Content":   "usuarios",
			"BodyClass": "page-usuarios",
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := views.ParseTemplates().ExecuteTemplate(w, "layout", data); err != nil {
			log.Printf("error ejecutando template usuario editar: %v", err)
			http.Error(w, "No se puede borrar pelicula/usuario porque tiene referencia en la tabla mira", http.StatusInternalServerError)
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
			http.Error(w, "No se puede borrar pelicula/usuario porque tiene referencia en la tabla mira", http.StatusInternalServerError)
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
			http.Error(w, "No se puede borrar pelicula/usuario porque tiene referencia en la tabla mira", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/usuarios", http.StatusSeeOther)
	})

	// Miras page (GET and POST for crear)
	http.HandleFunc("/miras", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			pelis, err := q.ListPelis(r.Context())
			if err != nil {
				log.Printf("error listando peliculas para miras: %v", err)
				http.Error(w, "No se puede borrar pelicula/usuario porque tiene referencia en la tabla mira", http.StatusInternalServerError)
				return
			}
			usuarios, err := q.ListUsuario(r.Context())
			if err != nil {
				log.Printf("error listando usuarios para miras: %v", err)
				http.Error(w, "No se puede borrar pelicula/usuario porque tiene referencia en la tabla mira", http.StatusInternalServerError)
				return
			}
			miras, err := q.ListMira(r.Context())
			if err != nil {
				log.Printf("error listando miras: %v", err)
				http.Error(w, "No se puede borrar pelicula/usuario porque tiene referencia en la tabla mira", http.StatusInternalServerError)
				return
			}
			data := map[string]interface{}{
				"Title":     "Calificaciones",
				"Header":    "Calificaciones",
				"Peliculas": pelis,
				"Usuarios":  usuarios,
				"Miras":     miras,
				"Content":   "miras",
				"BodyClass": "page-miras",
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			if err := views.ParseTemplates().ExecuteTemplate(w, "layout", data); err != nil {
				log.Printf("error ejecutando template miras: %v", err)
				http.Error(w, "No se puede borrar pelicula/usuario porque tiene referencia en la tabla mira", http.StatusInternalServerError)
			}
		case http.MethodPost:
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
				http.Error(w, "No se puede borrar pelicula/usuario porque tiene referencia en la tabla mira", http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	// Editar mira (render form con datos)
	http.HandleFunc("/miras/editar", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		iduStr := r.URL.Query().Get("idu")
		idpStr := r.URL.Query().Get("idp")
		if iduStr == "" || idpStr == "" {
			http.Error(w, "idu y idp requeridos", http.StatusBadRequest)
			return
		}
		idu, err := strconv.Atoi(iduStr)
		if err != nil {
			http.Error(w, "idu inválido", http.StatusBadRequest)
			return
		}
		idp, err := strconv.Atoi(idpStr)
		if err != nil {
			http.Error(w, "idp inválido", http.StatusBadRequest)
			return
		}

		// Obtener recursos necesarios
		pelis, err := q.ListPelis(r.Context())
		if err != nil {
			log.Printf("error listando peliculas para editar mira: %v", err)
			http.Error(w, "No se puede borrar pelicula/usuario porque tiene referencia en la tabla mira", http.StatusInternalServerError)
			return
		}
		usuarios, err := q.ListUsuario(r.Context())
		if err != nil {
			log.Printf("error listando usuarios para editar mira: %v", err)
			http.Error(w, "No se puede borrar pelicula/usuario porque tiene referencia en la tabla mira", http.StatusInternalServerError)
			return
		}

		// Buscar la mira concreta
		miras, err := q.GetMira(r.Context(), int32(idu))
		if err != nil {
			log.Printf("error obteniendo miras para usuario: %v", err)
			http.Error(w, "No se puede borrar pelicula/usuario porque tiene referencia en la tabla mira", http.StatusInternalServerError)
			return
		}
		var found *db.Mira
		for _, m := range miras {
			if m.Idp == int32(idp) {
				mm := m
				found = &mm
				break
			}
		}
		if found == nil {
			http.Error(w, "No encontrado", http.StatusNotFound)
			return
		}

		data := map[string]interface{}{
			"Title":     "Editar Calificación",
			"Header":    "Editar Calificación",
			"Peliculas": pelis,
			"Usuarios":  usuarios,
			"Mira":      found,
			"Miras":     miras,
			"Action":    "/miras/edit",
			"Content":   "miras",
			"BodyClass": "page-miras",
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := views.ParseTemplates().ExecuteTemplate(w, "layout", data); err != nil {
			log.Printf("error ejecutando template mira editar: %v", err)
			http.Error(w, "No se puede borrar pelicula/usuario porque tiene referencia en la tabla mira", http.StatusInternalServerError)
		}
	})

	// Handler edit POST for miras
	http.HandleFunc("/miras/edit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		idp, _ := strconv.Atoi(r.FormValue("idp"))
		idu, _ := strconv.Atoi(r.FormValue("idu"))
		gustoono := r.FormValue("gustoono")
		calif, _ := strconv.Atoi(r.FormValue("calif"))

		// Validaciones: gustoono debe ser "sí" o "no", calif entre 1 y 5
		// Validación en edición también
		if gustoono != "sí" && gustoono != "no" {
			log.Printf("gustoono inválido en edit: %q", gustoono)
			http.Error(w, "gustoono inválido", http.StatusBadRequest)
			return
		}
		if calif < 1 || calif > 5 {
			log.Printf("calif fuera de rango en edit: %d", calif)
			http.Error(w, "calif debe estar entre 1 y 5", http.StatusBadRequest)
			return
		}

		params := db.UpdateMiraParams{
			Idp:      int32(idp),
			Idu:      int32(idu),
			Gustoono: gustoono,
			Calif:    int32(calif),
		}
		if err := q.UpdateMira(r.Context(), params); err != nil {
			log.Printf("error actualizando mira: %v", err)
			http.Error(w, "No se puede borrar pelicula/usuario porque tiene referencia en la tabla mira", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/miras", http.StatusSeeOther)
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
			http.Error(w, "No se puede borrar pelicula/usuario porque tiene referencia en la tabla mira", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/miras", http.StatusSeeOther)
	})

	// No servimos más archivos estáticos desde /static/ —
	// todos los estilos y scripts están embebidos en las plantillas
	// según el requerimiento del TP (renderizado en servidor).

	// (Los handlers POST para /peliculas, /usuarios y /miras
	// fueron unificados en sus respectivos manejadores arriba.)

	// --- Iniciar servidor ---
	port := 8080
	fmt.Printf("🚀 Servidor corriendo en http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
