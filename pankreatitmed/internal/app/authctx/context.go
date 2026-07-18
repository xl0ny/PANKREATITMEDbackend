package authctx

import (
	"github.com/gin-gonic/gin"
)

const key = "user"

func Set(c *gin.Context, u UserCtx) { c.Set(key, u) }
func Get(c *gin.Context) (UserCtx, bool) {
	v, ok := c.Get(key)
	if !ok {
		return UserCtx{}, false
	}
	u, ok := v.(UserCtx)
	return u, ok
}
