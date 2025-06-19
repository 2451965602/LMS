package pack

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/2451965602/LMS/biz/model/model"
	"github.com/2451965602/LMS/pkg/errno"
)

func BuildBaseResp(err error) *model.BaseResp {
	if err == nil {
		return &model.BaseResp{
			Code: errno.SuccessCode,
			Msg:  errno.Success.ErrorMsg,
		}
	}
	Errno := errno.ConvertErr(err)
	return &model.BaseResp{
		Code: Errno.ErrorCode,
		Msg:  Errno.ErrorMsg,
	}
}

func SendFailResponse(c *app.RequestContext, err error) {
	resp := new(model.ErrorResp)
	resp.Base = BuildBaseResp(err)
	SendResponse(c, resp)
}

func SendResponse(c *app.RequestContext, data interface{}) {
	c.JSON(consts.StatusOK, data)
}
