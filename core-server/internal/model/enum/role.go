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

// ParseRoleString 把接口传进来的字符串角色转成数字枚举
func ParseRoleString(s string) Role {
	switch s {
	case "admin":
		return RoleAdmin
	case "user":
		return RoleUser
	case "visitor":
		return RoleVisitor
	default:
		return RoleUnknown
	}
}
