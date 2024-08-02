package admin

import (
	"Guang_Shun_Coin_Action/internal/auth"
	"Guang_Shun_Coin_Action/internal/response"
	"Guang_Shun_Coin_Action/pkg/logger"
	"net/http"
	"github.com/gin-gonic/gin"
	"time"
)

type UpdateProductPageRequest struct {
    Page     int    `json:"page" binding:"required"`
    Sort     string `json:"sort"`
    Category string `json:"category"`
	Price    string  `json:"price"`
}

type UpdateProductDetailRequest struct {
	ProductId   string    `json:"productId"`
	ProductName string    `json:"productName"`
	Category    string    `json:"category"`
	Price       float64   `json:"price"`
	MinBidPrice float64   `json:"minBidPrice"`
	StartAt    time.Time `json:"startAt"`
	EndedAt     time.Time `json:"endedAt"`
	Description    string    `json:"description"`
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

func UpdateProductDetail(c *gin.Context) {
	var err error

	// Create response
	r := response.New()

	// Validate token and set UUID in context (validation failed.)
	auth.ValidateToken(c)

	// Parse request body to JSON format
	var updateProductDetailRequest UpdateProductDetailRequest
	if err = c.ShouldBindJSON(&updateProductDetailRequest); err != nil {
		logger.Warn("[SHOP] " + err.Error())
		r.Message = err.Error()
		c.JSON(http.StatusBadRequest, r)
		return
	}

	err = updateProductDetail(updateProductDetailRequest)
	if err != nil {
		r.Message = err.Error()
		logger.Warn("[ADMIN] " + err.Error())
		c.JSON(http.StatusInternalServerError, r)
		return
	}

	r.Status = true
	r.Data = updateProductDetailRequest
	c.JSON(http.StatusOK, r)
}

func UpdateImage(c *gin.Context) {
	var err error

	// Create response
	r := response.New()

	// Validate token and set UUID in context (validation failed.)
	auth.ValidateToken(c)

    // Parse the form data to get the files and productID
    form, err := c.MultipartForm()
    if err != nil {
		r.Message = err.Error()
		logger.Error("[ADMIN] " + err.Error())
		c.JSON(http.StatusBadRequest, r)
        return
    }

    // Retrieve the list of files from the "files" field
    files := form.File["files"]

    // Retrieve the productID from the form data
    productID := c.PostForm("productID")
    if productID == "" {
		r.Message = "productID is empty"
		logger.Error("[ADMIN] productID is empty")
		c.JSON(http.StatusBadRequest, r)
        return
    }

	// Add image
	err = updateImage(files, productID)
	if err != nil {
		r.Message = err.Error()
		if r.Message == "File exceeds 10MB" {
			c.JSON(http.StatusOK, r)
			return
		}
		c.JSON(http.StatusInternalServerError, r)
		return
	}

    // Return a success message
	r.Status = true
    c.JSON(http.StatusOK, r)
}