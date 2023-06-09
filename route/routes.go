package route

import (
	"go_server/controllers"
	"go_server/logger"
	"go_server/middleware"
	"net/http"

	"github.com/gin-contrib/pprof"

	"github.com/gin-gonic/gin"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(mode)
	}
	r := gin.New()
	// 最重要的就是这个日志库
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	//v1 版本的路由
	v1 := r.Group("/api/v1")

	//swagger 接口文档
	// http://localhost:8080/swagger/index.html 可以看到接口文档
	// r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	// curl  http://127.0.0.1:19787/api/v1/register -X POST -d '{"UserName":"cc","Password":"123456","Email":"cc@cjic.ga","RePassword":"123456"}'
	v1.POST("/register", controllers.RegisterHandler)
	// curl  http://127.0.0.1:19787/api/v1/login -X POST -d '{"Password":"123456","Email":"cc@cjic.ga"}'
	v1.POST("/login", controllers.LoginHandler)

	// curl  http://127.0.0.1:19787/api/v1/profile  -H "auth-token:"
	v1.GET("/profile", middleware.JWTAuthMiddleWare(), controllers.ProfileHandler)

	v1.GET("/community", controllers.CommunityHandler)
	v1.GET("/community/:id", controllers.CommunityDetailHandler)

	v1.POST("/post", middleware.JWTAuthMiddleWare(), controllers.CreatePostHandler)
	// http://127.0.0.1:19787/api/v1/post/5820230320533868
	v1.GET("/post/:id", controllers.GetPostDetailHandler)
	// http://127.0.0.1:19787/api/v1/video/8820220121280003.mp4
	v1.GET("/video/:video", controllers.GetPostVideoHandler)

	v1.GET("/postlist", controllers.GetPostListHandler)
	// http://127.0.0.1:19787/api/v1/postlistcommunity?community=88&pageSize=15&pageNum=1
	v1.GET("/postlistcommunity", controllers.GetPostListByCommunityHandler)
	// 最新或者最热列表
	v1.GET("/postlist2", controllers.GetPostListHandler2)

	v1.POST("/like/", middleware.JWTAuthMiddleWare(), controllers.PostLikeHandler)

	v1.GET("/list", controllers.ListHandler)
	v1.GET("/listlast", controllers.ListLastHandler)
	v1.GET("/list/:id", controllers.ListDetailHandler)

	//验证jwt机制
	// v1.GET("/ping", middleware.JWTAuthMiddleWare(), func(context *gin.Context) {
	// 	// 这里post man 模拟的 将token auth-token
	// 	zap.L().Debug("ping", zap.String("ping-email", context.GetString("email")))
	// 	// controllers.ResponseSuccess(context, context)
	// 	controllers.LoginHandler
	// })

	r.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "ok")
	})
	// r.LoadHTMLFiles("./html/dist/index.html")
	// r.Static("/css", "./html/dist/css")
	// r.Static("/js", "./html/dist/js")
	// r.Static("/img", "./html/dist/img")

	// r.GET("/vue", func(context *gin.Context) {
	// 	context.HTML(http.StatusOK, "index.html", nil)
	// })

	// 注册pprof 相关路由
	pprof.Register(r)

	return r
}
