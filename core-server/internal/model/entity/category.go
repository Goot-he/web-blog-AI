package entity

import (
	"database/sql"
	"time"
)

type Category struct {
	ID        uint64       `db:"id"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
	ParentID  uint64       `db:"parent_id"`
	Category  string       `db:"category"`
}
