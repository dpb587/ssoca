package api

type SignPublicKeyRequest struct {
	PublicKey string `json:"public_key"`
}

type SignPublicKeyResponse struct {
	Target      *SignPublicKeyTargetResponse `json:"target,omitempty"`
	Certificate string                       `json:"certificate"`
}

type SignPublicKeyTargetResponse struct {
	Host string `json:"host,omitempty"`
	User string `json:"user,omitempty"`
	Port int    `json:"port,omitempty"`
}
