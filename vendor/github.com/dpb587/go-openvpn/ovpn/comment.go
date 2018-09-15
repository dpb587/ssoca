package ovpn

type CommentProfileElement struct {
	Comment string
}

func (CommentProfileElement) ProfileElementType() ProfileElementType {
	return CommentProfileElementType
}
