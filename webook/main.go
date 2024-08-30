package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	cookie2 "github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"go_test/webook/internal/repository"
	"go_test/webook/internal/repository/dao"
	"go_test/webook/internal/service"
	"go_test/webook/internal/web"
	"go_test/webook/internal/web/middlewares"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {
	db := initDb()
	server := InitWebServer()
	initUserHandler(db, server)
	server.Run(":8080")
}

func initDb() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:970827@tcp(localhost:3306)/webook"))
	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

func InitWebServer() *gin.Engine {
	server := gin.Default()
	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type"},
		AllowOriginFunc: func(origin string) bool {
			return strings.HasPrefix(origin, "http://localhost") || strings.Contains(origin, "my_company.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	login := &middlewares.LoginMiddlewareBuilder{}
	cookie := cookie2.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("ssid", cookie), login.CheckLogin())
	return server
}

func initUserHandler(db *gorm.DB, server *gin.Engine) {
	ud := dao.NewUserDao(db)
	ur := repository.NewUserRepository(ud)
	us := service.NewUserService(ur)
	userHandler := web.NewUserHandler(us)
	userHandler.RegisterRouters(server)
}
