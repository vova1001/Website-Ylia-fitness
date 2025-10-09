package handlerJSON

import (
	"github.com/gin-gonic/gin"

	h "github.com/vova1001/Website-Ylia-fitness/internal/handler"
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
		return
	}
	err = h.RegisterNewUser(NewUser)
	if err != nil {
		ctx.JSON(401, gin.H{"err": "err handler register"})
		return
	}
	ctx.JSON(200, gin.H{"Sucseseful": "User add table"})

}
