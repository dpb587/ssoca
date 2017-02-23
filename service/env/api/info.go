package api

type InfoResponse struct {
	Auth     InfoServiceResponse   `json:"auth"`
	Env      InfoEnvResponse       `json:"env"`
	Services []InfoServiceResponse `json:"services"`
	Version  string                `json:"version"`
}

type InfoEnvResponse struct {
	Banner   string            `json:"banner,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
	Name     string            `json:"name,omitempty"`
	Title    string            `json:"title,omitempty"`
	URL      string            `json:"url"`
}

type InfoEnvLinkResponse struct {
	Title string `json:"title,omitempty"`
	URL   string `json:"url,omitempty"`
}

type InfoServiceResponse struct {
	Metadata interface{} `json:"metadata,omitempty"`
	Name     string      `json:"name,omitempty"`
	Type     string      `json:"type"`
	Version  string      `json:"version"`
}
