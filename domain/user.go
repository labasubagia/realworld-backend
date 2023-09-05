package domain

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bu:"table:users,alias:u"`
	ID            int64     `bun:"id,pk,autoincrement"`
	Email         string    `bun:"email,notnull"`
	Username      string    `bun:"username,notnull"`
	Password      string    `bun:"password,notnull"`
	Image         string    `bun:"image"`
	Bio           string    `bun:"bio"`
	CreatedAt     time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
