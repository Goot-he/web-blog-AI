package enum

type UserStatus string

const (
	StatusApproved UserStatus = "approved"
	StatusBlocked  UserStatus = "blocked"
	StatusUnknown  UserStatus = "unknown"
)

func (s UserStatus) String() string {
	return string(s)
}

func ParseStatus(s string) UserStatus {
	switch s {
	case "approved":
		return StatusApproved
	case "blocked":
		return StatusBlocked
	default:
		return StatusUnknown
	}
}
