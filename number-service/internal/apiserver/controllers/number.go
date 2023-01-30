package controllers

import (
	"github.com/shakh9006/numbers-store/internal/apiserver/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NumberCtrl struct {
	numberService *services.NumberService
}

func NewNumberCtrl(numService *services.NumberService) *NumberCtrl {
	return &NumberCtrl{
		numberService: numService,
	}
}

func (nc *NumberCtrl) GetById(c *gin.Context) {
	//if binder := c.ShouldBindJSON(&user); binder != nil {
	//	c.JSON(http.StatusUnprocessableEntity, gin.H{
	//		"message": binder.Error(),
	//	})
	//	return
	//}

	c.JSON(http.StatusOK, gin.H{
		"message": "Hello",
	})

}
