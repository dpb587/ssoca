package ovpn

type ProfileElementType string

const (
	CommentProfileElementType   ProfileElementType = "comment"
	DirectiveProfileElementType ProfileElementType = "directive"
	EmbeddedProfileElementType  ProfileElementType = "embedded"
)

type ProfileElement interface {
	ProfileElementType() ProfileElementType
}
