package middlewares

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
}

func (m *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 注册一下这个类型
		gob.Register(time.Now())
		path := c.Request.URL.Path
		if path == "/user/login" || path == "/user/signup" {
			return
		}
		session := sessions.Default(c)
		userId := session.Get("userId")
		if userId == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		now := time.Now()

		// 每分钟刷新
		const updateTimeKey = "update_time"
		// 尝试取上次更新时间
		val := session.Get(updateTimeKey)
		lastUpdateTime, ok := val.(time.Time)
		if val == nil || !ok || now.Sub(lastUpdateTime) > time.Minute {
			// 第一次进入或大于1分钟了
			session.Set(updateTimeKey, now)
			session.Set("userId", userId)
			err := session.Save()
			if err != nil {
				fmt.Println(err)
			}

		}

	}
}
