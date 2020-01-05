package ovpn

type EmbeddedProfileElement interface {
	ProfileElement

	Embed() string
	Data() string
}

type GenericEmbeddedProfileElement struct {
	embed string
	data  string
}

func (GenericEmbeddedProfileElement) ProfileElementType() ProfileElementType {
	return EmbeddedProfileElementType
}

func (e GenericEmbeddedProfileElement) Embed() string {
	return e.embed
}

func (e GenericEmbeddedProfileElement) Data() string {
	return e.data
}
