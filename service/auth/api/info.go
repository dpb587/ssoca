package api

type InfoResponse struct {
	ID         string            `json:"id"`
	Groups     []string          `json:"groups"`
	Attributes map[string]string `json:"attributes,omitempty"`
}
