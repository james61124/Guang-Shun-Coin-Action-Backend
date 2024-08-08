package shop

import (
	"Guang_Shun_Coin_Action/internal/auth"
	"Guang_Shun_Coin_Action/internal/response"
	"Guang_Shun_Coin_Action/pkg/logger"
	"net/http"
	"github.com/gin-gonic/gin"
)

type ProductRequest struct {
    Page     int    `json:"page" binding:"required"`
    Sort     string `json:"sort"`
    Category string `json:"category"`
	Price    string  `json:"price"`
}

type TotalPagesOfProductRequest struct {
    Category string `json:"category"`
}

type DetailRequest struct {
    ProductID string `json:"productID"`
	HistoryPage int `json:"historyPage"`
}

type BidRequest struct {
    ProductID string `json:"productID"`
	BidPrice int `json:"bidPrice"`
}

type StarRequest struct {
    ProductID string `json:"productID"`
	IsStar bool `json:"isStar"`
}


func Product(c *gin.Context) {
	var err error

	// Create response
	r := response.New()

	// Validate token and set UUID in context (validation failed.)
	UUID := auth.ValidateToken(c)

	// Parse request body to JSON format
	var productRequest ProductRequest
	if err = c.ShouldBindJSON(&productRequest); err != nil {
		logger.Warn("[SHOP] " + err.Error())
		r.Message = err.Error()
		c.JSON(http.StatusBadRequest, r)
		return
	}

	var totalPages int
	var ProductResponse []response.ProductResponse
	var ProductList response.ProductListResponse
	ProductResponse, totalPages, err = product(productRequest, UUID)
	if err != nil {
		r.Message = err.Error()
		logger.Warn("[SHOP] " + err.Error())
		c.JSON(http.StatusInternalServerError, r)
		return
	}
	ProductList.Products = ProductResponse
    ProductList.TotalPages = totalPages

	r.Status = true
	r.Data = ProductList
	c.JSON(http.StatusOK, r)
}

func TotalPagesOfProduct(c *gin.Context) {
	var err error
	var total int

	// Create response
	r := response.New()

	// Validate token and set UUID in context (validation failed.)
	auth.ValidateToken(c)

	// Parse request body to JSON format
	var totalPagesOfProductRequest TotalPagesOfProductRequest
	if err = c.ShouldBindJSON(&totalPagesOfProductRequest); err != nil {
		logger.Warn("[SHOP] " + err.Error())
		r.Message = err.Error()
		c.JSON(http.StatusBadRequest, r)
		return
	}

	total, err = totalPagesOfProduct(totalPagesOfProductRequest)
	if err != nil {
		r.Message = err.Error()
		logger.Warn("[SHOP] " + err.Error())
		c.JSON(http.StatusInternalServerError, r)
		return
	}

	r.Status = true
	r.Data = response.TotalPagesOfProductResponse{Total: total}
	c.JSON(http.StatusOK, r)
}

func Detail(c *gin.Context) {
	var err error

	// Create response
	r := response.New()

	// Validate token and set UUID in context (validation failed.)
	UUID := auth.ValidateToken(c)

	// Parse request body to JSON format
	var detailRequest DetailRequest
	if err = c.ShouldBindJSON(&detailRequest); err != nil {
		logger.Warn("[SHOP] " + err.Error())
		r.Message = err.Error()
		c.JSON(http.StatusBadRequest, r)
		return
	}

	detailInfo, err := detail(UUID, detailRequest)
	if err != nil {
		r.Message = err.Error()
		logger.Warn("[SHOP] " + err.Error())
		c.JSON(http.StatusInternalServerError, r)
		return
	}

	r.Status = true
	r.Data = detailInfo
	c.JSON(http.StatusOK, r)
}

func Bid(c *gin.Context) {
	var err error

	// Create response
	r := response.New()

	// Validate token and set UUID in context (validation failed.)
	UUID := auth.ValidateToken(c)

	// Parse request body to JSON format
	var bidRequest BidRequest
	if err = c.ShouldBindJSON(&bidRequest); err != nil {
		logger.Warn("[SHOP] " + err.Error())
		r.Message = err.Error()
		c.JSON(http.StatusBadRequest, r)
		return
	}

	err = bid(bidRequest, UUID)
	if err != nil {
		r.Message = err.Error()
		if r.Message == "bidPrice is smaller or equal than highest bidPrice" {
			c.JSON(http.StatusOK, r)
			return
		}
		logger.Warn("[SHOP] " + err.Error())
		c.JSON(http.StatusInternalServerError, r)
		return
	}

	r.Status = true
	c.JSON(http.StatusOK, r)
}

func Star(c *gin.Context) {
	var err error

	// Create response
	r := response.New()

	// Validate token and set UUID in context (validation failed.)
	UUID := auth.ValidateToken(c)

	// Parse request body to JSON format
	var starRequest StarRequest
	if err = c.ShouldBindJSON(&starRequest); err != nil {
		logger.Warn("[SHOP] " + err.Error())
		r.Message = err.Error()
		c.JSON(http.StatusBadRequest, r)
		return
	}

	err = star(starRequest, UUID)
	if err != nil {
		r.Message = err.Error()
		logger.Warn("[SHOP] " + err.Error())
		c.JSON(http.StatusInternalServerError, r)
		return
	}

	r.Status = true
	c.JSON(http.StatusOK, r)
}
