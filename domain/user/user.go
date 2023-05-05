package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Role string

const (
	Admin  Role = "admin"
	Member Role = "member"
	Guest  Role = "guest"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	ID            uuid.UUID `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Email         string    `bun:",unique" json:"email"`
	Role          Role      `bun:"type:varchar(6)" json:"role"`
	CreatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`
}
