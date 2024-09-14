package main

import (
	"github.com/gin-contrib/cors"
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
	//server := gin.Default()
	//server.GET("/hello", func(ctx *gin.Context) {
	//	ctx.String(http.StatusOK, "hello, world")
	//})
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
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"X-Jwt-Token"},
		AllowOriginFunc: func(origin string) bool {
			return strings.HasPrefix(origin, "http://localhost") || strings.Contains(origin, "my_company.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	//redisClient := redis.NewClient(&redis.Options{
	//	Addr: "localhost:6379",
	//})
	//
	//server.Use(ratelimit.NewBuilder(limiter.NewRedisSlidingWindowLimiter(redisClient, time.Second, 1)).Build())

	UseJWT(server)
	// UserSession(server)
	return server
}

func UseJWT(server *gin.Engine) {
	login := &middlewares.LoginJWTMiddlewareBuilder{}
	server.Use(login.CheckLogin())
}

func UseSession(server *gin.Engine) {
	//login := &middlewares.LoginMiddlewareBuilder{}

	// session基于cookie的实现
	// cookie := cookie2.NewStore([]byte("secret"))

	// 分布式下基于redis的实现
	//store, err := redis.NewStore(16, "tcp", "localhost:6379", "",
	//	[]byte("bem36sguyjucw77teum4064f3lgjw5a4"), // Authentication key
	//	[]byte("lnij8x9s609gdaqiqweik4v656rry696")) // Encryption key
	//if err != nil {
	//	panic(err)
	//}
	//server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
}

func initUserHandler(db *gorm.DB, server *gin.Engine) {
	ud := dao.NewUserDao(db)
	ur := repository.NewUserRepository(ud)
	us := service.NewUserService(ur)
	userHandler := web.NewUserHandler(us)
	userHandler.RegisterRouters(server)
}
