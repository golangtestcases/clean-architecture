package create_entity_handler

type CreateEntityResponse struct {
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
	UserID string `json:"user_id"`
}