package sign_check

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	consts "goskeleton/app/global/response"
	"goskeleton/app/helpers"
	"goskeleton/app/utils/response"
	"sort"
	"strconv"
	"strings"
	"time"
)

type HeaderParams struct {
	Authorization string `header:"Authorization" binding:"required,min=20"`
}

// Next 检查token完整性、有效性中间件
func Next() gin.HandlerFunc {
	return func(context *gin.Context) {
		//var form interface{}
		headerParams := HeaderParams{}

		//  推荐使用 ShouldBindHeader 方式获取头参数
		if err := context.ShouldBindHeader(&headerParams); err != nil {
			response.SuccessButFail(context, err.Error(), consts.Unauthorized, nil)
			context.Abort()
		}

		var query string
		if context.Request.Method == "POST" {
			var rawData map[string]string
			jsonData, err := context.GetRawData()
			if err != nil {
				response.SuccessButFail(context, err.Error(), consts.Unauthorized, nil)
				context.Abort()
			}
			err = json.Unmarshal(jsonData, &rawData)
			if err != nil {
				response.SuccessButFail(context, err.Error(), consts.Unauthorized, nil)
				context.Abort()
			}
			keys := sortAlpha(rawData)
			for _, val := range keys {
				if query == "" {
					query += val + "=" + rawData[val]
				} else {
					query += "&" + val + "=" + rawData[val]
				}
			}
		} else {
			var rawData map[string]string
			var rawDataSet map[string][]string
			jsonData, err := json.Marshal(context.Request.URL.Query())
			if err != nil {
				response.SuccessButFail(context, err.Error(), consts.Unauthorized, nil)
				context.Abort()
			}
			err = json.Unmarshal(jsonData, &rawDataSet)
			if err != nil {
				response.SuccessButFail(context, err.Error(), consts.Unauthorized, nil)
				context.Abort()
			}
			_ = json.Unmarshal(jsonData, &rawData)
			for k, v := range rawDataSet {
				if len(v) > 0 {
					rawData[k] = v[0]
				}
			}
			keys := sortAlpha(rawData)
			for _, val := range keys {
				if query == "" {
					query += val + "=" + rawData[val]
				} else {
					query += "&" + val + "=" + rawData[val]
				}
			}
		}
		t := time.Now().Unix()
		if query != "" {
			query += "&t=" + strconv.Itoa(int(t))
		} else {
			query += "t=" + strconv.Itoa(int(t))
		}

		token := strings.Split(headerParams.Authorization, " ")
		var subStringToken string
		subStringToken = helpers.GetMD5Hash(token[0])
		if len(token) > 1 {
			subStringToken = helpers.GetMD5Hash(token[1])
		}

		signKey := helpers.EncryptAES([]byte(subStringToken), query)
		// 添加签名判断
		fmt.Println(signKey)
		context.Next()
	}
}

func sortAlpha(convertText map[string]string) []string {
	keys := make([]string, 0, len(convertText))
	for k := range convertText {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
