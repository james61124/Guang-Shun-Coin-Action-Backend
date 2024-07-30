package admin

import (
	"Guang_Shun_Coin_Action/internal/auth"
	"Guang_Shun_Coin_Action/internal/response"
	"Guang_Shun_Coin_Action/pkg/logger"
	"net/http"
	"github.com/gin-gonic/gin"
)

type UpdateProductPageRequest struct {
    Page     int    `json:"page" binding:"required"`
    Sort     string `json:"sort"`
    Category string `json:"category"`
	Price    string  `json:"price"`
}


func UpdateProductPage(c *gin.Context) {
	var err error

	// Create response
	r := response.New()

	// Validate token and set UUID in context (validation failed.)
	UUID := auth.ValidateToken(c)

	// Parse request body to JSON format
	var productRequest UpdateProductPageRequest
	if err = c.ShouldBindJSON(&productRequest); err != nil {
		logger.Warn("[SHOP] " + err.Error())
		r.Message = err.Error()
		c.JSON(http.StatusBadRequest, r)
		return
	}

	var totalPages int
	var ProductResponse []response.UpdateProductPageResponse
	var ProductList response.UpdateProductPageListResponse
	ProductResponse, totalPages, err = updateProductPage(productRequest, UUID)
	if err != nil {
		r.Message = err.Error()
		logger.Warn("[ADMIN] " + err.Error())
		c.JSON(http.StatusInternalServerError, r)
		return
	}
	ProductList.Products = ProductResponse
    ProductList.TotalPages = totalPages

	r.Status = true
	r.Data = ProductList
	c.JSON(http.StatusOK, r)
}