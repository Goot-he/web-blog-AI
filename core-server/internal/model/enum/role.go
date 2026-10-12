package enum

type Role uint32

const (
	RoleUnknown Role = iota
	RoleAdmin
	RoleUser
	RoleVisitor // 游客
)

func (r Role) Int() uint32 {
	return uint32(r)
}

func ParseRole(num uint32) Role {
	switch num {
	case 1:
		return RoleAdmin
	case 2:
		return RoleUser
	case 3:
		return RoleVisitor
	default:
		return RoleUnknown
	}
}
