package sign_check

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/helpers"
	"goskeleton/app/utils/response"
	"strings"
)

type HeaderParams struct {
	Authorization string `header:"Authorization" binding:"required,min=20"`
}

// Next 检查token完整性、有效性中间件
func Next() gin.HandlerFunc {
	return func(context *gin.Context) {
		var form interface{}
		headerParams := HeaderParams{}

		//  推荐使用 ShouldBindHeader 方式获取头参数
		if err := context.ShouldBindHeader(&headerParams); err != nil {
			response.ErrorParam(context, consts.JwtTokenMustValid+err.Error())
			return
		}
		//if err := context.ShouldBind(&form); err != nil {
		//	response.ErrorParam(context, consts.CaptchaCheckFailMsg)
		//	return
		//}
		fmt.Println(form)
		token := strings.Split(headerParams.Authorization, " ")
		var subStringToken string
		subStringToken = helpers.GetMD5Hash(token[0])
		if len(token) > 1 {
			subStringToken = helpers.GetMD5Hash(token[1])
		}

		fmt.Println(subStringToken, 123123)
	}
}
