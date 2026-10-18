package entity

import (
	"core-server/internal/model/enum"
	"database/sql"
	"time"
)

type User struct {
	ID        uint64       `db:"id"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
	BirthDate sql.NullTime `db:"birth_date"`
	Username  string       `db:"username"`
	Password  string       `db:"password"`
	Email     string       `db:"email"`
	Phone     string       `db:"phone"`
	Avatar    string       `db:"avatar"`
	Sex       enum.Sex     `db:"sex"`
	Role      enum.Role    `db:"role"`
	Signature string       `db:"signature"`
	Status    enum.Status  `db:"status"`
}
