package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/lib/pq"
	db "tpe.com/tpcursada/db/sqlc"
)

type Server struct {
	queries *db.Queries
}

func NewServer(q *db.Queries) *Server {
	return &Server{queries: q}
}

// --- GET /peliculas ---
func (s *Server) ListPeliculas(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	pelis, err := s.queries.ListPelis(ctx)
	if err != nil {
		http.Error(w, "Error al listar películas", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pelis)
}

// --- GET /peliculas/{id} ---
func (s *Server) GetPelicula(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/peliculas/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	peli, err := s.queries.GetPeli(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Película no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al obtener película", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(peli)
}

// --- POST /peliculas ---
func (s *Server) CreatePelicula(w http.ResponseWriter, r *http.Request) {
	var p db.CreatePeliParams
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(p.Titulo) == "" {
		http.Error(w, "El título no puede estar vacío", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	peli, err := s.queries.CreatePeli(ctx, p)
	if err != nil {
		http.Error(w, "Error al crear película", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(peli)
}

// --- PUT /peliculas/{id} ---
func (s *Server) UpdatePelicula(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/peliculas/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var p db.UpdatePeliParams
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}
	p.Idp = int32(id)
	// Validación mínima
	if strings.TrimSpace(p.Titulo) == "" {
		http.Error(w, "El título no puede estar vacío", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	if err := s.queries.UpdatePeli(ctx, p); err != nil {
		http.Error(w, "Error al actualizar película", http.StatusInternalServerError)
		return
	}
	// Devolver la entidad actualizada
	peli, err := s.queries.GetPeli(ctx, p.Idp)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Película no encontrada tras actualización", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al obtener película actualizada", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(peli)
}

// --- DELETE /peliculas/{id} ---
func (s *Server) DeletePelicula(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/peliculas/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	if err := s.queries.DeletePeli(ctx, int32(id)); err != nil {
		// Detectar violación de clave foránea (hay miras relacionadas)
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23503" {
			http.Error(w, "No se puede borrar la película: existen recursos relacionados (ej.: miras).", http.StatusConflict)
			return
		}
		http.Error(w, "Error al borrar película", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
