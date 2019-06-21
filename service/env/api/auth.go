package api

type AuthResponse struct {
	ID         string            `json:"id"`
	Groups     []string          `json:"groups"`
	Attributes map[string]string `json:"attributes,omitempty"`
}
