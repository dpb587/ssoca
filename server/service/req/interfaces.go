package req

//go:generate counterfeiter . RouteHandler
type RouteHandler interface {
	Route() string
	// Execute(...) error
}
