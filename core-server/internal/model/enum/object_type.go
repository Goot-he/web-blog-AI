package enum

type ObjectType string

const (
	ObjectTypeArticle ObjectType = "article"
	ObjectTypeComment ObjectType = "comment"
	ObjectTypeLife    ObjectType = "life"
	ObjectTypeUnknown ObjectType = "unknown"
)

func (t ObjectType) String() string {
	return string(t)
}

func ParseObjectType(s string) ObjectType {
	switch s {
	case "article":
		return ObjectTypeArticle
	case "comment":
		return ObjectTypeComment
	case "life":
		return ObjectTypeLife
	default:
		return ObjectTypeUnknown
	}
}
