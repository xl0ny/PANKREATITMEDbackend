package middleware

import (
	"net/http"
	"strings"
	"time"

	"pankreatitmed/internal/app/authctx"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(jwtCfg JWTConfig, blacklist *RedisBlacklist) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(strings.ToLower(auth), "bearer ") {
			c.Next()
			return
		}

		raw := strings.TrimSpace(auth[len("bearer "):])

		if ok, _ := blacklist.IsBlacklisted(raw); ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "token is blacklisted"})
			return
		}

		parser := jwt.NewParser(jwt.WithIssuedAt(), jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
		claims := &Claims{}
		token, err := parser.ParseWithClaims(raw, claims, func(token *jwt.Token) (any, error) {
			return []byte(jwtCfg.Secret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		if jwtCfg.Issuer != "" && claims.Issuer != jwtCfg.Issuer {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid issuer"})
			return
		}
		if claims.ExpiresAt == nil || time.Now().After(claims.ExpiresAt.Time) {
			c.AbortWithStatusJSON(401, gin.H{"error": "token expired"})
			return
		}
		println(claims.Sub, claims.Login, claims.IsModerator)
		c.Set("user", authctx.UserCtx{ID: claims.Sub, Login: claims.Login, IsModerator: claims.IsModerator})
		c.Next()
	}
}
