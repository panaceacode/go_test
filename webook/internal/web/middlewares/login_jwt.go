package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go_test/webook/internal/web"
	"log"
	"net/http"
	"strings"
	"time"
)

type LoginJWTMiddlewareBuilder struct {
}

func (m *LoginJWTMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/user/login" || path == "/user/signup" {
			return
		}

		// token 在 Authorization 里取
		authCode := c.GetHeader("Authorization")
		if authCode == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		sgns := strings.Split(authCode, " ")
		if len(sgns) != 2 {
			// Authorization 不合法
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := sgns[1]
		var uc web.UserClaims
		token, err := jwt.ParseWithClaims(tokenString, &uc, func(token *jwt.Token) (interface{}, error) {
			return web.JWTKey, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if uc.UserAgent != c.GetHeader("User-Agent") {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		expireTime := uc.ExpiresAt
		if expireTime.Sub(time.Now()).Seconds() < 50 {
			uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute * 30))
			tokenString, err = token.SignedString(web.JWTKey)
			c.Header("X-Jwt-Token", tokenString)
			if err != nil {
				log.Println(err)
			}
		}

		c.Set("user", uc)
	}
}
