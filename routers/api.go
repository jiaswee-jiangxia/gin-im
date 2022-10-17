package routers

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/http/controller/api"
	"goskeleton/app/http/middleware/cors"
	"goskeleton/app/http/middleware/jwt"
	"io"
	"os"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

// 该路由主要设置门户类网站等前台路由

func InitApiRouter() *gin.Engine {
	var router *gin.Engine
	// 非调试模式（生产模式） 日志写到日志文件
	if variable.ConfigYml.GetBool("AppDebug") == false {
		//1.将日志写入日志文件
		gin.DisableConsoleColor()
		f, _ := os.Create(variable.BasePath + variable.ConfigYml.GetString("Logs.GinLogName"))
		gin.DefaultWriter = io.MultiWriter(f)
		// 2.如果是有nginx前置做代理，接口访问根本不需要gin框架记录访问日志，开启下面 2 行代码，屏蔽上面的三行代码，性能提升 5%
		//gin.SetMode(gin.ReleaseMode)
		//gin.DefaultWriter = ioutil.Discard

		router = gin.Default()
	} else {
		// 调试模式，开启 pprof 包，便于开发阶段分析程序性能
		router = gin.Default()
		pprof.Register(router)
	}

	router.Use(cors.Next())

	vApi := router.Group("/app/api")
	{
		userApi := vApi.Group("/user")
		{
			userApi.POST("/login", api.Login)
			userApi.POST("/register", api.Register)
			jwtUserGroup := userApi.Use(jwt.JWT())
			{
				jwtUserGroup.GET("/profile", api.Profile)
			}
		}
		groupApi := vApi.Group("/group")
		{
			jwtGroupApi := groupApi.Use(jwt.JWT())
			{
				jwtGroupApi.POST("/create", api.CreateGroup)
				jwtGroupApi.POST("/admin/list", api.ListGroupAdmin)
				jwtGroupApi.POST("/member/list", api.ListGroupMember)
				jwtGroupApi.POST("/addmember", api.AddGroupMember)
				jwtGroupApi.POST("/setadmin", api.SetGroupAdmin)
				jwtGroupApi.POST("/setowner", api.SetGroupOwner)
				jwtGroupApi.POST("/member/remove", api.RemoveGroupMember)
				jwtGroupApi.POST("/disband", api.DisbandGroup)
			}
		}
	}
	return router
}
