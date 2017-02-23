package api

type ListResponse struct {
	Files []ListFileResponse `json:"files"`
}

type ListFileResponse struct {
	Name   string `json:"name"`
	Size   int64  `json:"size"`
	Digest string `json:"digest"`
}
