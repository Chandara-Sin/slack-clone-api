package login

type GrantType string

const (
	Password     GrantType = "password"
	RefreshToken GrantType = "refresh_token"
)

type Login struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	GrantType string `json:"grant_type"`
}
