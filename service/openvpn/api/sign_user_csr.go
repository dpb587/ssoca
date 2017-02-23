package api

type SignUserCSRRequest struct {
	CSR string `json:"csr"`
}

type SignUserCSRResponse struct {
	Certificate string `json:"certificate"`
	Profile     string `json:"profile"`
}
