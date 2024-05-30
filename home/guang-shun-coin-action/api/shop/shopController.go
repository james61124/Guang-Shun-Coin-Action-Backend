package shop

import (
	// "Guang_Shun_Coin_Action/internal/auth"
	"Guang_Shun_Coin_Action/internal/response"
	"Guang_Shun_Coin_Action/pkg/logger"
	"net/http"
	// "regexp"
	// "strings"
	"github.com/gin-gonic/gin"
)

type ProductRequest struct {
    Page     int    `json:"page" binding:"required"`
    Sort     string `json:"sort"`
    Category string `json:"category"`
}


func Product(c *gin.Context) {
	var err error

	// Create response
	r := response.New()

	// Parse request body to JSON format
	var productRequest ProductRequest
	if err = c.ShouldBindJSON(&productRequest); err != nil {
		logger.Warn("[SHOP] " + err.Error())
		r.Message = err.Error()
		c.JSON(http.StatusBadRequest, r)
		return
	}

	var ProductResponse []response.ProductResponse
	ProductResponse, err = product(productRequest)
	if err != nil {
		r.Message = err.Error()
		c.JSON(http.StatusInternalServerError, r)
		return
	}

	
	r.Status = true
	r.Data = ProductResponse
	c.JSON(http.StatusCreated, r)
}
