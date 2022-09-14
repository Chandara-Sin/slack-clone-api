package user

import "time"

type Role string

const (
	Admin  Role = "admin"
	Member Role = "member"
	Guest  Role = "guest"
)

type User struct {
	Id             uint      `gorm:"primary_key" json:"id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Email          string    `gorm:"unique" json:"email"`
	Password       string    `gorm:"-" json:"password,omitempty"`
	HashedPassword string    `json:"-"`
	Role           Role      `gorm:"type:varchar(6)" json:"role"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
