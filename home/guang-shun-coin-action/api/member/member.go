package member

import (
	"Guang_Shun_Coin_Action/pkg/logger"
	"Guang_Shun_Coin_Action/pkg/mariadb"
	"github.com/google/uuid"
	"errors"
)

func addProduct(rr addProductRequest, ownerUUID string) (string, error) {
	var query string
	var err error
	var productID string
	shippingStatus := []string{"Unconfirmed Order", "Shipping in Progress", "Order Completed"}
	productID = uuid.NewString()

	// Check whether productName is empty
	if rr.Name == "" {
		logger.Warn("[Member] productName is empty")
		return "", errors.New("productName is empty")
	}

	// Check whether EndDate earlier than StartDate
	if rr.EndDate.Before(rr.StartDate) {
		logger.Warn("[Member] endDate earlier than startDate")
		return "", errors.New("endDate earlier than startDate")
	}
	
	// Insert into user database
	query = `INSERT INTO Product (
			productId, 
			userId,
			productName, 
			category, 
			price, 
			minBidPrice, 
			startAt, 
			endedAt,
			shippingStatus, 
			productDescription
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = mariadb.DB.Exec(
		query, 
		productID, 
		ownerUUID,
		rr.Name, 
		rr.Category, 
		rr.Price, 
		rr.MinBidPrice, 
		rr.StartDate, 
		rr.EndDate,
		shippingStatus[0],
		rr.ProductDescription,
	)

	if err != nil {
		logger.Error("[Product] " + err.Error())
		return "", err
	}

	logger.Info("[Product] Successfully add product with productname: " + rr.Name)
	return productID, nil
}

// func addImage(c *gin.Context, ownerUUID string) error {
// 	var query string
// 	var err error

// 	// check the extended file name
// 	allowedExtensions := map[string]bool{".png": true, ".jpg": true, ".jpeg": true, ".gif": true}
// 	fileExt := filepath.Ext(file.Filename)
// 	if !allowedExtensions[fileExt] {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "File type not allowed"})
// 		return
// 	}

// 	// check file path exists
// 	uploadPath := "./uploads"
// 	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
// 		os.Mkdir(uploadPath, os.ModePerm)
// 	}

// 	// 保存文件
// 	filePath := filepath.Join(uploadPath, file.Filename)
// 	if err := c.SaveUploadedFile(file, filePath); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save the file"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "file_path": filePath})
// 	logger.Info("[Product] Successfully add product with productname: " + rr.Name)

// 	return nil
// }