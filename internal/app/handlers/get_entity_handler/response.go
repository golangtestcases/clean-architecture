package get_entity_handler

type GetEntityResponse struct {
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
	UserID string `json:"user_id"`
}