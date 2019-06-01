package api

type ListResponse struct {
	Files []ListFileResponse `json:"files"`
}

type ListFileResponse struct {
	Name   string                 `json:"name"`
	Size   int64                  `json:"size"`
	Digest ListFileDigestResponse `json:"digest"`
}

type ListFileDigestResponse struct {
	SHA1   string `json:"sha1,omitempty"`
	SHA256 string `json:"sha256,omitempty"`
	SHA512 string `json:"sha512,omitempty"`
}
