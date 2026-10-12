package enum

type Status string

const (
	StatusApproved Status = "approved"
	StatusBlocked  Status = "blocked"
	StatusUnknown  Status = "unknown"
)

func (s Status) String() string {
	return string(s)
}

func ParseStatus(s string) Status {
	switch s {
	case "approved":
		return StatusApproved
	case "blocked":
		return StatusBlocked
	default:
		return StatusUnknown
	}
}
