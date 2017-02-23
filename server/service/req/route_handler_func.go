package req

import "net/http"

type RouteHandlerFunc struct {
	Path string
	Func http.HandlerFunc
}

func (h RouteHandlerFunc) Route() string {
	return h.Path
}

func (h RouteHandlerFunc) Execute(w http.ResponseWriter, r *http.Request) {
	h.Func(w, r)
}
