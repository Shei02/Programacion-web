package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/lib/pq"
	db "tpe.com/tpcursada/db/sqlc"
)

// --- GET /miras?idu={idu} ---
// Lista todas las miras de un usuario específico
func (s *Server) ListMiras(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("idu")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	miras, err := s.queries.ListMira(ctx, int32(id))
	if err != nil {
		http.Error(w, "Error al listar miras", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(miras)
}

// --- GET /miras/{idp}/{idu} ---
// Obtiene una mira específica de un usuario
func (s *Server) GetMira(w http.ResponseWriter, r *http.Request) {
	parts := r.URL.Path[len("/miras/"):]
	ids := splitPath(parts)
	if len(ids) != 2 {
		http.Error(w, "URL inválida", http.StatusBadRequest)
		return
	}

	idp, err1 := strconv.Atoi(ids[0])
	idu, err2 := strconv.Atoi(ids[1])
	if err1 != nil || err2 != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	mira, err := s.queries.GetMira(ctx, db.GetMiraParams{
		Idp: int32(idp),
		Idu: int32(idu),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Mira no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al obtener mira", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mira)
}

// --- POST /miras ---
// Crea una nueva mira
func (s *Server) CreateMira(w http.ResponseWriter, r *http.Request) {
	var m db.CreateMiraParams
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}
	// validaciones opcionales, p.ej. calif entre 1 y 5
	if m.Calif < 1 || m.Calif > 5 {
		http.Error(w, "calif fuera de rango", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	if err := s.queries.CreateMira(ctx, m); err != nil {
		// detectar duplicado de PK y devolver 409
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505":
				http.Error(w, "La mira ya existe", http.StatusConflict)
				return
			case "23503":
				http.Error(w, "No se puede crear la mira: referencia a película o usuario inexistente.", http.StatusConflict)
				return
			}
		}
		http.Error(w, "Error al crear mira", http.StatusInternalServerError)
		return
	}

	// recuperar la mira creada (suponiendo que GetMira existe)
	mira, err := s.queries.GetMira(ctx, db.GetMiraParams{Idp: m.Idp, Idu: m.Idu})
	if err != nil {
		http.Error(w, "Error al obtener mira creada", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mira)
}

// --- DELETE /miras/{idp}/{idu} ---
// Borra una mira específica a través del path
func (s *Server) DeleteMiraByPath(w http.ResponseWriter, r *http.Request) {
	parts := r.URL.Path[len("/miras/"):]
	ids := splitPath(parts)
	if len(ids) != 2 {
		http.Error(w, "URL inválida", http.StatusBadRequest)
		return
	}
	idp, err1 := strconv.Atoi(ids[0])
	idu, err2 := strconv.Atoi(ids[1])
	if err1 != nil || err2 != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	if err := s.queries.DeleteMira(ctx, db.DeleteMiraParams{Idp: int32(idp), Idu: int32(idu)}); err != nil {
		http.Error(w, "Error al borrar mira", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- DELETE /miras ---
// Borra una mira específica
func (s *Server) DeleteMira(w http.ResponseWriter, r *http.Request) {
	var m db.DeleteMiraParams
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	if err := s.queries.DeleteMira(ctx, m); err != nil {
		http.Error(w, "Error al borrar mira", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- Helper para separar path ---
func splitPath(path string) []string {
	var parts []string
	for _, p := range []rune(path) {
		if p == '/' {
			parts = append(parts, "")
		} else {
			if len(parts) == 0 {
				parts = append(parts, string(p))
			} else {
				parts[len(parts)-1] += string(p)
			}
		}
	}
	return parts
}
