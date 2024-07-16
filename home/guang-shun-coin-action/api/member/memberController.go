package member

import (
	"Guang_Shun_Coin_Action/internal/response"
	"Guang_Shun_Coin_Action/internal/auth"
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
	ProductDescription string `json:"productDescription"`
}

type getHistoryBidRequest struct {
	Sorting string `json:"sorting"`
	Page int `json: page`
}

func AddProduct(c *gin.Context) {
	var err error
	var productID string

	// Create response
	r := response.New()

	// Validate token and set UUID in context (validation failed.)
	UUID := auth.ValidateToken(c)

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

	// Create response
	r := response.New()

	// Validate token and set UUID in context (validation failed.)
	auth.ValidateToken(c)

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
	// if len(files) == 0 {
	// 	logger.Error("[MEMBER] No files uploaded")
    //     c.JSON(http.StatusBadRequest, gin.H{"error": "No files uploaded"})
    //     return
    // }

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
	r.Status = true
    c.JSON(http.StatusOK, r)
}

func GetHistoryBid(c *gin.Context) {
	var err error

	// Create response
	r := response.New()

	// Validate token and set UUID in context (validation failed.)
	UUID := auth.ValidateToken(c)

	// Parse request body to JSON format
	var getHistoryBidRequest getHistoryBidRequest
	if err = c.ShouldBindJSON(&getHistoryBidRequest); err != nil {
		logger.Warn("[MEMBER] " + err.Error())
		r.Message = err.Error()
		c.JSON(http.StatusBadRequest, r)
		return
	}

	// Add product
	historyBidList, err := getHistoryBid(getHistoryBidRequest, UUID)
	if err != nil {
		r.Message = err.Error()
		c.JSON(http.StatusInternalServerError, r)
		return
	}

	// return response
	r.Status = true
	r.Data = historyBidList
	c.JSON(http.StatusOK, r)
}

func TotalPagesOfHistoryBid(c *gin.Context) {
	var err error

	// Create response
	r := response.New()

	// Validate token and set UUID in context (validation failed.)
	UUID := auth.ValidateToken(c)

	total, err := totalPagesOfHistoryBid(UUID)
	if err != nil {
		r.Message = err.Error()
		logger.Warn("[MEMBER] " + err.Error())
		c.JSON(http.StatusInternalServerError, r)
		return
	}

	r.Status = true
	r.Data = response.TotalPagesOfProductResponse{Total: total}
	c.JSON(http.StatusOK, r)
}

func GetUserInfo(c *gin.Context) {
	var err error

	// Create response
	r := response.New()

	// Validate token and set UUID in context (validation failed.)
	UUID := auth.ValidateToken(c)

	userInfo, err := getUserInfo(UUID)
	if err != nil {
		r.Message = err.Error()
		logger.Warn("[MEMBER] " + err.Error())
		c.JSON(http.StatusInternalServerError, r)
		return
	}

	r.Status = true
	r.Data = userInfo
	c.JSON(http.StatusOK, r)
}