package ovpn

type Profile struct {
	Name     string
	Elements []ProfileElement
}

func (p Profile) GetDirective(name string) []ProfileElement {
	found := []ProfileElement{}

	for _, e := range p.Elements {
		switch d := e.(type) {
		case DirectiveProfileElement:
			if d.Directive() != name {
				continue
			}

			found = append(found, e)
		}
	}

	return found
}

func (p Profile) GetEmbedded(name string) []ProfileElement {
	found := []ProfileElement{}

	for _, e := range p.Elements {
		switch d := e.(type) {
		case EmbeddedProfileElement:
			if d.Embed() != name {
				continue
			}

			found = append(found, e)
		}
	}

	return found
}
