package create_entity_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golangtestcases/clean-architecture/internal/domain/model"
	"github.com/google/uuid"
)

type EntityService interface {
	CreateEntity(ctx context.Context, entity model.Entity) (model.Entity, error)
}

type CreateEntityHandler struct {
	entityService EntityService
}

func NewCreateEntityHandler(entityService EntityService) *CreateEntityHandler {
	return &CreateEntityHandler{entityService: entityService}
}

func (h CreateEntityHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var createEntityRequest CreateEntityRequest

	if err := json.NewDecoder(r.Body).Decode(&createEntityRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userUUID, err := uuid.Parse(createEntityRequest.UserID)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	entity := model.Entity{
		Name:   createEntityRequest.Name,
		UserID: userUUID,
	}

	newEntity, err := h.entityService.CreateEntity(r.Context(), entity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	entityResponse := CreateEntityResponse{
		ID:     uint64(newEntity.ID),
		Name:   newEntity.Name,
		UserID: newEntity.UserID.String(),
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(&entityResponse); err != nil {
		fmt.Println("json.Encode failed", err)
		return
	}
}
