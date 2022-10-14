package analytics

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/http/controller/analytics"
	"goskeleton/app/http/validator/core/data_transfer"
	"goskeleton/app/utils/response"
)

type MatchLog struct {
	MatchId int32 `form:"match_id" json:"match_id"  binding:"required"`
	MatchType string `form:"match_type" json:"match_type"  binding:"required"`
	Source string `form:"source" json:"source"  binding:"required"`
}

func (n MatchLog) CheckParams(context *gin.Context) {
	if err := context.ShouldBind(&n); err != nil {
		response.ErrorParam(context, gin.H{
			"err":  err.Error(),
		})
		return
	}

	extraAddBindDataContext := data_transfer.DataAddContext(n, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		response.ErrorSystem(context, "验证器json化失败", "")
	} else {
		analytics.WriteLog(extraAddBindDataContext)
	}

}
