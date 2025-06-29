package jwttokengen

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserDetails
	jwt.RegisteredClaims
}
type UserDetails struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Password string `json:"pwd"`
}
