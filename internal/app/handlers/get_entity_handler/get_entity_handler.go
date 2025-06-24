package get_entity_handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/golangtestcases/clean-architecture/internal/domain/model"
)

type EntityService interface {
	GetEntitiesByID(ctx context.Context, id model.EntityID) ([]model.Entity, error)
}

type GetEntityHandler struct {
	entityService EntityService
}

func NewGetEntityHandler(entityService EntityService) *GetEntityHandler {
	return &GetEntityHandler{entityService: entityService}
}

func (h GetEntityHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	idRaw := r.PathValue("id")
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		http.Error(w, "id must be a number", http.StatusBadRequest)
		return
	}

	if id < 1 {
		http.Error(w, "id must be more than zero", http.StatusBadRequest)
		return
	}

	entities, err := h.entityService.GetEntitiesByID(r.Context(), model.EntityID(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var responses []GetEntityResponse
	for _, entity := range entities {
		responses = append(responses, GetEntityResponse{
			ID:     uint64(entity.ID),
			Name:   entity.Name,
			UserID: entity.UserID.String(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}
