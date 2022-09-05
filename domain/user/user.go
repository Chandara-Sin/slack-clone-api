package user

import "time"

type Role string

const (
	ADMIN Role = "ADMIN"
	STAFF Role = "STAFF"
)

type User struct {
	Id             uint      `gorm:"primary_key" json:"id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Email          string    `gorm:"unique" json:"email"`
	Password       string    `gorm:"-" json:"password,omitempty"`
	HashedPassword string    `json:"-"`
	Role           Role      `json:"role"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
