package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/response"
	_ "goskeleton/bootstrap"
	"goskeleton/routers"
	"log"
)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	model.Setup()
	model.SetupRedis()
	//translation.Setup()
}

// 这里可以存放门户类网站入口
func main() {
	router := routers.InitApiRouter()
	router.GET("/", func(context *gin.Context) {
		response.Success(context, "health ok", nil)
	})
	//loc, err := time.LoadLocation("Asia/Singapore")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//time.Local = loc
	_ = router.Run(variable.ConfigYml.GetString("HttpServer.Api.Port"))
}
