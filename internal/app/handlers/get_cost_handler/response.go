package get_cost_handler

type GetCostResponse struct {
	TotalCost int `json:"total_cost"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}