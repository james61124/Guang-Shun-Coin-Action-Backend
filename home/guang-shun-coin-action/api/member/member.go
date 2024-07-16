package member

import (
	"Guang_Shun_Coin_Action/pkg/logger"
	"Guang_Shun_Coin_Action/pkg/mariadb"
	"Guang_Shun_Coin_Action/internal/response"
	"github.com/google/uuid"
	"errors"
	"io"
    "mime/multipart"
    "os"
	"fmt"
	"path/filepath"
	"net/http"
	"time"
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

func addImage(files []*multipart.FileHeader, productID string) error {

	var imageID string

	// Create the directory to save the uploaded files
    uploadDir := "./assets"
    if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		logger.Error("[Member] Unable to create assets directory")
        return err
    }

    for seq, file := range files {

		imageID = uuid.NewString()
		var imageUrl = "/assets/" + imageID

        // Check if the file size exceeds the limit (10 MB)
        if file.Size > 10*1024*1024 { // 10 MB
			logger.Error("[Member] File exceeds 10MB")
            return errors.New("File exceeds 10MB")
        }

        // Open the uploaded file
        src, err := file.Open()
        if err != nil {
			logger.Error("[Member] Unable to open file")
        	return err
        }
        // defer src.Close()

        // Check the file type
        if err := checkFileType(src); err != nil {
			logger.Error("[Member] Invalid file type")
            return errors.New("Invalid file type")
        }

        // Prefix the filename with "<imageID>" and then the original filename
        newFilename := fmt.Sprintf("%s", imageID)
        dstPath := filepath.Join(uploadDir, newFilename)
        dst, err := os.Create(dstPath)
        if err != nil {
			logger.Error("[Member] Unable to create file")
        	return err
        }
        // defer dst.Close()

        // Copy the file contents to the destination file
        if _, err := io.Copy(dst, src); err != nil {
			logger.Error("[Member] Unable to save file")
        	return err
        }

		// Insert into database
		query := `
				INSERT INTO ProductImage (imageId, productId, imageUrl, seq)
				VALUES (?, ?, ?, ?)
				`
		_, err = mariadb.DB.Exec(
			query, 
			imageID, 
			productID,
			imageUrl,
			seq,
		)
	
		if err != nil {
			logger.Error("[Member] " + err.Error())
			return err
		}
    }

	logger.Info("[Member] Successfully add images")
	return nil
}

// checkFileType checks if the uploaded file is an image.
func checkFileType(file multipart.File) error {
    // Convert file to io.ReadSeeker to support Seek method
    readSeeker, ok := file.(io.ReadSeeker)
    if !ok {
		logger.Error("[Member] file does not support Seek")
        return fmt.Errorf("file does not support Seek")
    }

    // Read the first 512 bytes of the file to determine the MIME type
    buf := make([]byte, 512)
    if _, err := readSeeker.Read(buf); err != nil && err != io.EOF {
		logger.Error("[Member] " + err.Error())
        return err
    }
    
    // Reset the file read position
    _, err := readSeeker.Seek(0, io.SeekStart)
    if err != nil {
		logger.Error("[Member] " + err.Error())
        return err
    }

    // Check if the file type is an image
    if !isImage(buf) {
		logger.Error("[Member] Invalid file type")
        return fmt.Errorf("Invalid file type")
    }

    return nil
}

// isImage determines if the file type is an image based on its MIME type.
func isImage(buf []byte) bool {
    mimeType := http.DetectContentType(buf)
    return mimeType == "image/jpeg" || mimeType == "image/png" || mimeType == "image/gif"
}

func getHistoryBid(rr getHistoryBidRequest, ownerUUID string) ([]response.GetBidHistory, error) {

	// get history record
	query := "SELECT productId, bidPrice, bidTime, status FROM History WHERE userId = ? LIMIT ? OFFSET ?"
	productNums := 12
	offset := (rr.Page - 1) * productNums
	rows, err := mariadb.DB.Query(query, ownerUUID, productNums, offset)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			logger.Warn("[Member] No product is available currently.\n")
		}
		logger.Error("[Member] " + err.Error())
	}
	defer rows.Close()

	var historyBidList []response.GetBidHistory
	for rows.Next() {
		
		var historyBid response.GetBidHistory
		var bidTime string

		if err = rows.Scan(&historyBid.ProductID, &historyBid.BidPrice, &bidTime, &historyBid.Status); err != nil {
			logger.Error("[Member] " + err.Error())
			return historyBidList, err
		}
		historyBid.BidTime, err = time.Parse("2006-01-02 15:04:05", bidTime)
		if err != nil {
			logger.Error("[Member] " + err.Error())
			return historyBidList, err
		}

		// get productName
		query = `SELECT productName FROM Product WHERE productId = ?;`;
		err := mariadb.DB.QueryRow(query, historyBid.ProductID).Scan(&historyBid.ProductName)
		if err != nil {
			logger.Error("[Member] " + err.Error())
			return historyBidList, err
		}
		
		// get imageUrl
		query = `SELECT imageUrl FROM ProductImage WHERE productId = ? ORDER BY seq LIMIT 1;`;
		err = mariadb.DB.QueryRow(query, historyBid.ProductID).Scan(&historyBid.ImageUrl)
		if err != nil {
			logger.Error("[Member] " + err.Error())
			return historyBidList, err
		}

		historyBidList = append(historyBidList, historyBid)
	}

	logger.Info("[Member] Successfully get historyBidList")
	return historyBidList, nil
}

func totalPagesOfHistoryBid(UUID string) (int, error) {
	var err error

    query := `SELECT COUNT(*) AS product_count FROM History WHERE userId = ?`

    var productCount int
    err = mariadb.DB.QueryRow(query, UUID).Scan(&productCount)
    if err != nil {
        logger.Error("[Member] " + err.Error())
		return productCount, err
    }

	total := fmt.Sprintf("%d", productCount / 12 + 1)
	logger.Info("[Member] Successfully return total pages of history bid: " + total)
	
	return productCount / 12 + 1, err
}

func getUserInfo(UUID string) (response.GetUserInfoResponse, error) {
	var err error

    query := `SELECT realName, nickName, cellphone, fbAccount, email, postcode, shippingAddr, username FROM User WHERE userId = ?`

    var getUserInfoResponse response.GetUserInfoResponse
    err = mariadb.DB.QueryRow(query, UUID).Scan(&getUserInfoResponse.RealName, &getUserInfoResponse.NickName, &getUserInfoResponse.Cellphone, &getUserInfoResponse.FbAccount, &getUserInfoResponse.Email, &getUserInfoResponse.Postcode, &getUserInfoResponse.ShippingAddr, &getUserInfoResponse.Username)
    if err != nil {
        logger.Error("[Member] " + err.Error())
		return getUserInfoResponse, err
    }

	logger.Info("[Member] Successfully return user info: " + UUID)
	
	return getUserInfoResponse, err
}

