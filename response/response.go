package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Response(ctx *gin.Context, httpStatus int, data interface{}, msg string)  {
	ctx.JSON(httpStatus, gin.H{"code": httpStatus, "data": data, "msg": msg})
}

func Success(ctx *gin.Context, data interface{}, msg string)  {
	Response(ctx, http.StatusOK, data, msg)
}

func Failure(ctx *gin.Context, data interface{}, msg string)  {
	Response(ctx, http.StatusBadRequest, data, msg)
}