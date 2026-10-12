package enum

type Sex string

const (
	Male    Sex = "male"
	Female  Sex = "female"
	Unknown Sex = "unknown"
)

func (s Sex) String() string {
	return string(s)
}

func ParseSex(s string) Sex {
	switch s {
	case "male":
		return Male
	case "female":
		return Female
	default:
		return Unknown
	}
}
