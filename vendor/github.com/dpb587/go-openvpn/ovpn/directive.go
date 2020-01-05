package ovpn

type DirectiveProfileElement interface {
	ProfileElement

	Directive() string
	Args() []string
}

type GenericDirectiveProfileElement struct {
	directive string
	args      []string
}

func (e GenericDirectiveProfileElement) Directive() string {
	return e.directive
}

func (e GenericDirectiveProfileElement) Args() []string {
	return e.args
}

func (GenericDirectiveProfileElement) ProfileElementType() ProfileElementType {
	return DirectiveProfileElementType
}
