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
	//bindAddress := "localhost:8080"
	//r := gin.Default()
	//r.GET("/ws", ws)
	//r.Run(bindAddress)
	router := routers.InitApiRouter()
	router.GET("/", func(context *gin.Context) {
		response.Success(context, "health ok", nil)
	})
	_ = router.Run(variable.ConfigYml.GetString("HttpServer.Api.Port"))
}

//type WsMap struct {
//	Event  string      `json:"event" binding:"required"`
//	Params interface{} `json:"params" binding:"required"`
//}

//func ws(c *gin.Context) {
//	var wsMap WsMap
//	//升级get请求为webSocket协议
//	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
//	if err != nil {
//		return
//	}
//	defer ws.Close()
//	for {
//		//读取ws中的数据
//		mt, message, err := ws.ReadMessage()
//		if err != nil {
//			break
//		}
//		if string(message) == "ping" {
//			message = []byte("pong")
//		}
//		_ = json.Unmarshal(message, &wsMap)
//
//		//业务逻辑
//		var resp interface{}
//		switch wsMap.Event {
//		case "login":
//
//			resp = api.Login(c, wsMap.Params)
//			fmt.Println(resp)
//		}
//		returnFlag, _ := json.Marshal(resp)
//		//写入ws数据
//		err = ws.WriteMessage(mt, returnFlag)
//		if err != nil {
//			break
//		}
//	}
//}
