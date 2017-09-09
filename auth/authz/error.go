package authz

type Error struct {
	error
}

func NewError(err error) Error {
	return Error{error: err}
}
