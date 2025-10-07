package handlerJSON

import (
	"github.com/gin-gonic/gin"
)

func GetJson(ctx *gin.Context) {
	ctx.JSON(200, "Hello, server is runnig!!!")
}
