package delivery

import (
	"github.com/gin-gonic/gin"
	customErrs "github.com/shakh9006/numbers-store/errors"
	"github.com/shakh9006/numbers-store/internal/apiserver/services"
	"net/http"
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
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "id is empty"})
		return
	}

	number, err := nc.numberService.GetById(id)
	if err != nil {
		if err == customErrs.ErrRowsNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "number not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "number": number.Number})
}
