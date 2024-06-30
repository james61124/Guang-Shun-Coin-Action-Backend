package member

import (
	"Guang_Shun_Coin_Action/internal/response"
	"Guang_Shun_Coin_Action/internal/auth"
	"Guang_Shun_Coin_Action/pkg/logger"
	"net/http"
	"github.com/gin-gonic/gin"
	"time"
	// "fmt"
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

	// fmt.Printf("add the auth validation\n")

	// Validate token and set UUID in context (validation failed.)
	auth.ValidateToken(c)

	// Get UUID from context
	uuid, exists := c.Get("UUID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "UUID not found from auth"})
		return
	}

	// Type assert UUID to string
	UUID, ok := uuid.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "UUID not a string"})
		return
	}

    // var UUID = "e85e507e-0b71-4ca2-8825-72429b41a222"

	// Add product
	err = addProduct(addProductRequest, UUID)
	if err != nil {
		r.Message = err.Error()
		c.JSON(http.StatusInternalServerError, r)
		return
	}

	// return response
	r.Status = true
	c.JSON(http.StatusCreated, r)
}
