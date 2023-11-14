package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (r *Response) Responses(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, r)
	return
}
