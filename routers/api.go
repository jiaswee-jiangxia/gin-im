package routers

import (
	"goskeleton/app/global/consts"
	"goskeleton/app/global/variable"
	"goskeleton/app/http/controller/api"
	"goskeleton/app/http/middleware/cors"
	"goskeleton/app/http/middleware/jwt"
	validatorFactory "goskeleton/app/http/validator/core/factory"
	"io/ioutil"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

// 该路由主要设置门户类网站等前台路由

func InitApiRouter() *gin.Engine {
	var router *gin.Engine
	// 非调试模式（生产模式） 日志写到日志文件
	if variable.ConfigYml.GetBool("AppDebug") == false {
		//1.将日志写入日志文件
		//gin.DisableConsoleColor()
		//f, _ := os.Create(variable.BasePath + variable.ConfigYml.GetString("Logs.GinLogName"))
		//gin.DefaultWriter = io.MultiWriter(f)
		// 2.如果是有nginx前置做代理，接口访问根本不需要gin框架记录访问日志，开启下面 2 行代码，屏蔽上面的三行代码，性能提升 5%
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = ioutil.Discard

		router = gin.Default()
	} else {
		// 调试模式，开启 pprof 包，便于开发阶段分析程序性能
		router = gin.Default()
		pprof.Register(router)
	}

	router.Use(cors.Next())

	backend := router.Group("/admin/")
	{
		backend.GET("ws", validatorFactory.Create(consts.ValidatorPrefix+"WebsocketConnect"))
	}

	router.POST("/app/api/account/login", api.Login)
	router.POST("/app/api/account/register", api.Register)
	router.POST("/app/api/account/vcode/get_vcode", api.GetVcode)

	vApi := router.Group("/app/api")
	{
		//vApi.Use(sign_check.Next())
		imApi := vApi.Group("/im")
		{
			imApi.POST("/update-register", api.ImUpdateRegister)
		}
		accountApi := vApi.Group("/account")
		{
			jwtAccountApi := accountApi.Use(jwt.JWT())
			{
				jwtAccountApi.POST("/token_refresh", api.RefreshToken)
				jwtAccountApi.POST("/check_token", api.CheckToken)
				jwtAccountApi.GET("/profile", api.GetProfile)
				jwtAccountApi.POST("update/profile", api.UpdateProfile)
				jwtAccountApi.POST("update/token", api.UpdateToken)
				jwtAccountApi.POST("change_password", api.UpdatePassword)
				jwtAccountApi.GET("my_phone", api.GetPhoneNo)
			}
		}
		blackUserApi := vApi.Group("/black_user")
		{
			jwtBlackUserApi := blackUserApi.Use(jwt.JWT())
			{
				jwtBlackUserApi.POST("/create", api.CreateBlackUser)
				jwtBlackUserApi.POST("/remove", api.RemoveBlackUser)
				jwtBlackUserApi.GET("/query_black_list", api.GetQueryBlackList)
			}
		}
		groupsApi := vApi.Group("/group")
		{
			jwtGroupsApi := groupsApi.Use(jwt.JWT())
			{
				jwtGroupsApi.POST("/create", api.CreateGroup)
				jwtGroupsApi.GET("/admin/list", api.ListGroupAdmin)
				jwtGroupsApi.GET("/member/list", api.ListGroupMember)
				jwtGroupsApi.POST("/addmember", api.AddGroupMember)
				jwtGroupsApi.POST("/setadmin", api.SetGroupAdmin)
				jwtGroupsApi.POST("/setowner", api.SetGroupOwner)
				jwtGroupsApi.POST("/member/remove", api.RemoveGroupMember)
				jwtGroupsApi.POST("/disband", api.DisbandGroup)
			}
		}
		contactApi := vApi.Group("/contact")
		{
			jwtContactApi := contactApi.Use(jwt.JWT())
			{
				jwtContactApi.GET("/", api.ContactList)
				jwtContactApi.GET("/lookup", api.SearchContact)
				jwtContactApi.POST("/create", api.CreateContact)
				jwtContactApi.POST("/grouping", api.CreateGrouping)
				//jwtContactApi.POST("/accept", api.AcceptContact)
				jwtContactApi.POST("/remove", api.RemoveContact)
			}
		}
		followApi := vApi.Group("/follow")
		{
			jwtFollowApi := followApi.Use(jwt.JWT())
			{
				jwtFollowApi.POST("/focus", api.FollowUser)
				jwtFollowApi.POST("/remove", api.UnfollowUser)
				jwtFollowApi.GET("/query_follow_ids", api.GetMyFollowList)
				jwtFollowApi.POST("/follow_count", api.GetUserFollowCount)
				jwtFollowApi.POST("/follow_both/is_friend", api.CheckIsFriend)
			}
		}
	}
	return router
}
