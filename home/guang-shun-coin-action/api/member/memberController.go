package member

import (
	"Guang_Shun_Coin_Action/internal/response"
	"Guang_Shun_Coin_Action/pkg/logger"
	"net/http"
	"github.com/gin-gonic/gin"
	"time"
)

type addProductRequest struct {
	Name    string `json:"name"`
	Category  string `json:"category"`
	Price  int `json:"price"`
	MinBidPrice int `json:"minBidPrice"`
	StartDate time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
}

func AddProduct(c *gin.Context) {
	var err error

	// Create response
	r := response.New()

	// Parse request body to JSON format
	var addProductRequest addProductRequest
	if err = c.ShouldBindJSON(&addProductRequest); err != nil {
		logger.Warn("[Product] " + err.Error())
		r.Message = err.Error()
		c.JSON(http.StatusBadRequest, r)
		return
	}

	// Add product
	err = addProduct(addProductRequest)
	if err != nil {
		r.Message = err.Error()
		c.JSON(http.StatusInternalServerError, r)
		return
	}

	// return response
	r.Status = true
	c.JSON(http.StatusCreated, r)
}
