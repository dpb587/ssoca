package ovpn

func ParseEmbedded(spec string, raw []byte) (ProfileElement, error) {
	return GenericEmbeddedProfileElement{
		embed: spec,
		data:  string(raw),
	}, nil
}
