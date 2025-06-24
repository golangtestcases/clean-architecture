package create_entity_handler

type CreateEntityRequest struct {
	Name   string `json:"name"`
	UserID string `json:"user_id"`
}