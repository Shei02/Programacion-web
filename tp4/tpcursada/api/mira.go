package api

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/lib/pq"
	db "tpe.com/tpcursada/db/sqlc"
)

// --- GET /miras ---
func (s *Server) ListMiras(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	idStr := r.URL.Query().Get("idu")

	if idStr == "" {
		// devolver toda la tabla
		miras, err := s.queries.ListMira(ctx)
		if err != nil {
			http.Error(w, "Error al listar miras", http.StatusInternalServerError)
			return
		}

		if miras == nil {
			miras = []db.Mira{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(miras)
		return
	}

	// si hay idu -> listar miras del usuario
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	miras, err := s.queries.GetMira(ctx, int32(id))
	if err != nil {
		http.Error(w, "Error al listar miras", http.StatusInternalServerError)
		return
	}

	if miras == nil {
		miras = []db.Mira{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(miras)
}

// --- GET /miras/{idu} ---
// Lista las miras de un usuario dado (ruta /miras/{idu})
func (s *Server) GetMira(w http.ResponseWriter, r *http.Request) {
	parts := r.URL.Path
	// Accept both "/miras/..." and "/api/miras/..." mounting
	parts = strings.TrimPrefix(parts, "/api")
	parts = strings.TrimPrefix(parts, "/miras/")
	ids := splitPath(parts)
	if len(ids) != 1 {
		http.Error(w, "URL inválida: use /miras/{idu}", http.StatusBadRequest)
		return
	}

	idu, err := strconv.Atoi(ids[0])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	miras, err := s.queries.GetMira(ctx, int32(idu))
	if err != nil {
		http.Error(w, "Error al obtener miras", http.StatusInternalServerError)
		return
	}

	if miras == nil {
		miras = []db.Mira{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(miras)
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

	// devolver la entidad creada (usamos el payload recibido)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(m)
}

// --- DELETE /miras/{idp}/{idu} ---
// Borra una mira específica a través del path
func (s *Server) DeleteMiraByPath(w http.ResponseWriter, r *http.Request) {
	parts := r.URL.Path
	parts = strings.TrimPrefix(parts, "/api")
	parts = strings.TrimPrefix(parts, "/miras/")
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
	// Delete the specific mira identified by idp + idu
	if err := s.queries.DeleteMiraByIDs(ctx, db.DeleteMiraByIDsParams{Idp: int32(idp), Idu: int32(idu)}); err != nil {
		http.Error(w, "Error al borrar mira", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- DELETE /miras ---
// Borra una mira específica
func (s *Server) DeleteMira(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Idu int32 `json:"idu"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	if err := s.queries.DeleteMira(ctx, body.Idu); err != nil {
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

// --- PUT /miras/{idp}/{idu} ---
// Actualiza una mira específica identificada por idp e idu en el path
func (s *Server) UpdateMira(w http.ResponseWriter, r *http.Request) {
	parts := r.URL.Path
	parts = strings.TrimPrefix(parts, "/api")
	parts = strings.TrimPrefix(parts, "/miras/")
	ids := splitPath(parts)
	if len(ids) != 2 {
		http.Error(w, "URL inválida: use /miras/{idp}/{idu}", http.StatusBadRequest)
		return
	}

	idp, err1 := strconv.Atoi(ids[0])
	idu, err2 := strconv.Atoi(ids[1])
	if err1 != nil || err2 != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var payload db.UpdateMiraParams
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("UpdateMira decode error: %v", err)
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}
	// Sobrescribir ids con los de la ruta para evitar inconsistencia
	payload.Idp = int32(idp)
	payload.Idu = int32(idu)

	ctx := context.Background()
	if err := s.queries.UpdateMira(ctx, payload); err != nil {
		log.Printf("UpdateMira DB error idp=%d idu=%d: %v", payload.Idp, payload.Idu, err)
		http.Error(w, "Error al actualizar mira", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Devolver el recurso actualizado (payload)
	json.NewEncoder(w).Encode(payload)
}
