package authctx

import "github.com/gin-gonic/gin"

type UserCtx struct {
	ID          uint
	Login       string
	IsModerator bool
}

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := c.Get("user"); !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthenticated"})
			return
		}
		c.Next()
	}
}

func RequireModerator() gin.HandlerFunc {
	return func(c *gin.Context) {
		v, ok := c.Get("user")
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthenticated"})
			return
		}
		if v.(UserCtx).IsModerator != true {
			c.AbortWithStatusJSON(403, gin.H{"error": "not moderator"})
			return
		}
		c.Next()
	}
}
