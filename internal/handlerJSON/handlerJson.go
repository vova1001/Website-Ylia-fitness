package handlerJSON

import (
	"github.com/gin-gonic/gin"

	m "github.com/vova1001/Website-Ylia-fitness/internal/model"
)

func GetJson(ctx *gin.Context) {
	ctx.JSON(200, "Hello, server is runnig!!!")
}

func PostNewUser(ctx *gin.Context) {
	var NewUser m.User
	err := ctx.ShouldBindJSON(&NewUser)
	if err != nil {
		ctx.JSON(400, gin.H{"err": "err json"})
	}

}
