package middleware

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// claims токена
type Claims struct {
	Sub         uint   `json:"sub"` // user id
	Login       string `json:"login"`
	IsModerator bool   `json:"ismoderator"`
	jwt.RegisteredClaims
}

type JWTConfig struct {
	Secret string
	Issuer string
	TTL    time.Duration
}
