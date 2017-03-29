package errors

type Exit struct {
	error
	Code int
}

func (e Exit) Error() string {
	if e.error != nil {
		return e.error.Error()
	}

	return ""
}
