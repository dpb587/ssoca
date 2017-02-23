package api

type InfoResponse struct {
	Username   string                 `json:"username"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}
