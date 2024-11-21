package custom_error

type DTO struct {
	Message string   `json:"error"`
	Details []string `json:"details"`
}
