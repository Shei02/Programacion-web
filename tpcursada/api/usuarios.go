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

// --- GET /usuarios ---
func (s *Server) ListUsuarios(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	users, err := s.queries.ListUsuario(ctx)
	if err != nil {
		http.Error(w, "Error al listar usuarios", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// --- GET /usuarios/{id} ---
func (s *Server) GetUsuario(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/usuarios/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	u, err := s.queries.GetUsuario(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Usuario no encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al obtener usuario", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

// --- POST /usuarios ---
func (s *Server) CreateUsuario(w http.ResponseWriter, r *http.Request) {
	var u db.CreateUsuarioParams
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(u.Nomusu) == "" || strings.TrimSpace(u.Email) == "" {
		http.Error(w, "nomusu y email son obligatorios", http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	usuario, err := s.queries.CreateUsuario(ctx, u)
	if err != nil {
		http.Error(w, "Error al crear usuario", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(usuario)
}

// --- PUT /usuarios/{id} ---
func (s *Server) UpdateUsuario(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/usuarios/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var u db.UpdateUsuarioParams
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}
	u.Idu = int32(id)

	// Validaciones mínimas
	if strings.TrimSpace(u.Nomusu) == "" || strings.TrimSpace(u.Email) == "" {
		http.Error(w, "nomusu y email son obligatorios", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	if err := s.queries.UpdateUsuario(ctx, u); err != nil {
		http.Error(w, "Error al actualizar usuario", http.StatusInternalServerError)
		return
	}
	// Devolver el usuario actualizado
	usuario, err := s.queries.GetUsuario(ctx, u.Idu)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Usuario no encontrado tras actualización", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al obtener usuario actualizado", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usuario)
}

// --- DELETE /usuarios/{id} ---
func (s *Server) DeleteUsuario(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/usuarios/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	if err := s.queries.DeleteUsuario(ctx, int32(id)); err != nil {
		// Detectar violación de clave foránea y devolver un 409 en vez de 500
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23503" {
			http.Error(w, "No se puede borrar el usuario: existen recursos relacionados (ej.: miras).", http.StatusConflict)
			return
		}
		http.Error(w, "Error al borrar usuario", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
