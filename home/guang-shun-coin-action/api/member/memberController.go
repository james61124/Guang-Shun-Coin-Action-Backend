package member

import (
	"Guang_Shun_Coin_Action/internal/response"
	"Guang_Shun_Coin_Action/internal/auth"
	"Guang_Shun_Coin_Action/pkg/logger"
	"net/http"
	"github.com/gin-gonic/gin"
	"time"
	"Guang_Shun_Coin_Action/pkg/mariadb"
)

type addProductRequest struct {
	Name    string `json:"name"`
	Category  string `json:"category"`
	Price  int `json:"price"`
	MinBidPrice int `json:"minBidPrice"`
	StartDate time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	ProductDescription string `json:"productDescription"`
}

func AddProduct(c *gin.Context) {
	var err error
	var query string
	var exists bool
	var productID string

	// Create response
	r := response.New()

	// Validate token and set UUID in context (validation failed.)
	auth.ValidateToken(c)

	// Get UUID from context
	uuid, exists := c.Get("UUID")
	if !exists {
		logger.Warn("[MEMBER] UUID not found from auth")
		return
	}

	// Type assert UUID to string
	UUID, ok := uuid.(string)
	if !ok {
		logger.Error("[MEMBER] UUID not a string")
		return
	}

	// Check if user exists
	query = "SELECT EXISTS(SELECT 1 FROM User WHERE userId = ?)"
    err = mariadb.DB.QueryRow(query, UUID).Scan(&exists)
	if err != nil && err.Error() != "sql: no rows in result set" {
		logger.Error("[MEMBER] " + err.Error())
		r.Message = err.Error()
		c.JSON(http.StatusBadRequest, r)
		return
	} else if exists == false {
		logger.Error("[MEMBER] user doesn't exists")
		r.Message = "user doesn't exists"
		c.JSON(http.StatusBadRequest, r)
		return
	}

	// Parse request body to JSON format
	var addProductRequest addProductRequest
	if err = c.ShouldBindJSON(&addProductRequest); err != nil {
		logger.Warn("[Product] " + err.Error())
		r.Message = err.Error()
		c.JSON(http.StatusBadRequest, r)
		return
	}

	// Add product
	productID, err = addProduct(addProductRequest, UUID)
	if err != nil {
		r.Message = err.Error()
		if r.Message == "productName is empty" || r.Message == "endDate earlier than startDate" {
			c.JSON(http.StatusOK, r)
			return
		}
		c.JSON(http.StatusInternalServerError, r)
		return
	}

	// return response
	r.Status = true
	r.Data = response.AddProductResponse{ProductId: productID}
	c.JSON(http.StatusOK, r)
}

// func AddImage(c *gin.Context) {

// 	file, err := c.FormFile("file")
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
// 		return
// 	}

	


// 	var err error

// 	// Create response
// 	r := response.New()

// 	// Validate token and set UUID in context (validation failed.)
// 	auth.ValidateToken(c)

// 	// Get UUID from context
// 	uuid, exists := c.Get("UUID")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "UUID not found from auth"})
// 		return
// 	}

// 	// Type assert UUID to string
// 	UUID, ok := uuid.(string)
// 	if !ok {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "UUID not a string"})
// 		return
// 	}

// 	// Add product
// 	err = addImage(c, UUID)
// 	if err != nil {
// 		r.Message = err.Error()
// 		c.JSON(http.StatusInternalServerError, r)
// 		return
// 	}

// 	// return response
// 	r.Status = true
// 	c.JSON(http.StatusCreated, r)
// }
