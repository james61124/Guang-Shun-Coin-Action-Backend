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

func AddImage(c *gin.Context) {
	var err error
	var query string
	var exists bool

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

    // Parse the form data to get the files and productID
    form, err := c.MultipartForm()
    if err != nil {
		r.Message = err.Error()
		logger.Error("[MEMBER] " + err.Error())
		c.JSON(http.StatusBadRequest, r)
        return
    }

    // Retrieve the list of files from the "files" field
    files := form.File["files"]
	if len(files) == 0 {
		logger.Error("[MEMBER] No files uploaded")
        c.JSON(http.StatusBadRequest, gin.H{"error": "No files uploaded"})
        return
    }

    // Retrieve the productID from the form data
    productID := c.PostForm("productID")
    if productID == "" {
		r.Message = "productID is empty"
		logger.Error("[MEMBER] productID is empty")
		c.JSON(http.StatusBadRequest, r)
        return
    }

	// Add image
	err = addImage(files, productID)
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
    c.JSON(http.StatusOK, gin.H{"message": "Files uploaded successfully"})
}