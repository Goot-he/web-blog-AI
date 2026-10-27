package enum

type InteractionType string

const (
	InteractionTypeComment InteractionType = "comment"
	InteractionTypeLike    InteractionType = "like"
	InteractionTypeFavor   InteractionType = "favor"
	InteractionTypeFollow  InteractionType = "follow"
	InteractionTypeView    InteractionType = "view"
	InteractionTypeUnknown InteractionType = "unknown"
)

func (i InteractionType) String() string {
	return string(i)
}

func ParseInteractionType(s string) InteractionType {
	switch s {
	case "comment":
		return InteractionTypeComment
	case "like":
		return InteractionTypeLike
	case "favor":
		return InteractionTypeFavor
	case "follow":
		return InteractionTypeFollow
	case "view":
		return InteractionTypeView
	default:
		return InteractionTypeUnknown
	}
}
