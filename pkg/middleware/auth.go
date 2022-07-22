package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/gole-dev/gole/pkg/app"
	"github.com/gole-dev/gole/pkg/errcode"
)

// Auth authorize user
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the json web token.
		ctx, err := app.ParseRequest(c)
		if err != nil {
			app.NewResponse().Error(c, errcode.ErrInvalidToken)
			c.Abort()
			return
		}

		// set uid to context
		c.Set("uid", ctx.UserID)

		c.Next()
	}
}
