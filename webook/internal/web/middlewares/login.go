package middlewares

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginMiddlewareBuilder struct {
}

func (m LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/user/login" || path == "/user/signup" {
			return
		}
		session := sessions.Default(c)
		if session.Get("userId") == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
