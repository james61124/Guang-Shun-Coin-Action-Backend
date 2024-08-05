package admin

import (
	"Guang_Shun_Coin_Action/pkg/logger"
	"Guang_Shun_Coin_Action/pkg/mariadb"
	"Guang_Shun_Coin_Action/internal/response"
	"time"
	"strings"
	"mime/multipart"
    "os"
	"fmt"
	"path/filepath"
	"io"
	"net/http"
	"errors"
	"github.com/google/uuid"
)

func getProductQuery(pr UpdateProductPageRequest) (string, string, []interface{}) {
    // Base query without the ORDER BY clause
    baseQuery := `
        SELECT 
            p.productId, 
            p.productName, 
            p.category, 
            COALESCE(MAX(h.bidPrice), p.price) AS actualPrice, 
            p.minBidPrice, 
            p.startAt, 
            p.endedAt,
            COUNT(h.historyId) AS bidCount
        FROM 
            Product p
        LEFT JOIN 
            History h ON p.productId = h.productId
        WHERE 
            p.category = ?
        GROUP BY 
            p.productId, p.productName, p.category, p.price, p.minBidPrice, p.startAt, p.endedAt
        HAVING 
            (COALESCE(MAX(h.bidPrice), p.price) BETWEEN ? AND ? OR ? = '50000+' AND COALESCE(MAX(h.bidPrice), p.price) > 50000)
    `

	// Query to count the total products matching the criteria
    countQuery := `
        SELECT COUNT(*)
        FROM (
            SELECT 
                p.productId
            FROM 
                Product p
            LEFT JOIN 
                History h ON p.productId = h.productId
            WHERE 
                p.category = ?
            GROUP BY 
                p.productId, p.price, p.endedAt
            HAVING 
                (COALESCE(MAX(h.bidPrice), p.price) BETWEEN ? AND ? OR ? = '50000+' AND COALESCE(MAX(h.bidPrice), p.price) > 50000)
        ) AS subquery
    `

    // Determine the ORDER BY clause based on pr.Sort
    var orderByClause string
    switch pr.Sort {
    case "目前價格由大至小":
        orderByClause = "ORDER BY actualPrice DESC, bidCount DESC"
    case "目前價格由小至大":
        orderByClause = "ORDER BY actualPrice ASC, bidCount DESC"
    case "截標日期由近至遠":
        orderByClause = "ORDER BY p.endedAt ASC, bidCount DESC"
    case "截標日期由遠至近":
        orderByClause = "ORDER BY p.endedAt DESC, bidCount DESC"
    default:
        orderByClause = "ORDER BY actualPrice DESC, bidCount DESC" // Default ordering
    }

    // Add LIMIT and OFFSET
    limitClause := "LIMIT ? OFFSET ?"

    // Construct the final query
    query := baseQuery + orderByClause + " " + limitClause

    // Calculate pagination
    productNums := 12
    offset := (pr.Page - 1) * productNums

    // Prepare arguments
    var args []interface{}
    args = append(args, pr.Category)

    // Assuming we only handle one price range for simplicity
    var minPrice, maxPrice string
	prices := strings.Split(pr.Price, "-")
	minPrice = prices[0]
	maxPrice = prices[1]

    args = append(args, minPrice, maxPrice, maxPrice) // Add price range arguments
    args = append(args, productNums, offset) // Add limit and offset arguments

    return query, countQuery, args
}

// return status, list of products
func updateProductPage(pr UpdateProductPageRequest, UUID string) ([]response.UpdateProductPageResponse, int, error) {
	var err error

	query, countQuery, args := getProductQuery(pr)

	// Execute the count query to get the total number of products matching the criteria
    var totalProductCount int
	var totalPages int
    err = mariadb.DB.QueryRow(countQuery, args[0], args[1], args[2], args[3]).Scan(&totalProductCount)
    if err != nil {
        logger.Error("[ADMIN] " + err.Error())
        return nil, totalProductCount, err
    }
	if ( totalProductCount % 12 == 0 ) {
		totalPages = totalProductCount / 12
	} else {
		totalPages = totalProductCount / 12 + 1
	}

	rows, err := mariadb.DB.Query(query, args...)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			logger.Warn("[ADMIN] No product is available currently.\n")
			return nil, totalPages, err
		}
		logger.Error("[ADMIN] " + err.Error())
		return nil, totalPages, err
	}
	defer rows.Close()
	
	var startAt, endedAt string
	var products []response.UpdateProductPageResponse
	for rows.Next() {
		var p response.UpdateProductPageResponse
		err = rows.Scan(&p.ProductId, &p.ProductName, &p.Category, &p.Price, &p.MinBidPrice, &startAt, &endedAt, &p.BidCount)
		if err != nil {
			logger.Error("[ADMIN] " + err.Error())
			return products, totalPages, err
		}
		p.StartAt, err = time.Parse("2006-01-02 15:04:05", startAt)
		if err != nil {
			logger.Error("[ADMIN] " + err.Error())
			return products, totalPages, err
		}

		p.EndedAt, err = time.Parse("2006-01-02 15:04:05", endedAt)
		if err != nil {
			logger.Error("[ADMIN] " + err.Error())
			return products, totalPages, err
		}
		
		query := `SELECT imageUrl FROM ProductImage WHERE productId = ? ORDER BY seq LIMIT 1;`;
		var imageUrl string
		err := mariadb.DB.QueryRow(query, p.ProductId).Scan(&imageUrl)
		if err != nil {
			logger.Error("[ADMIN] " + err.Error())
			return products, totalPages, err
		}
		p.ImgUrl = imageUrl
		products = append(products, p)
	}

	logger.Info("[ADMIN] Successfully return all the products")
	return products, totalPages, err
}

func updateProductDetail(pr UpdateProductDetailRequest) (error) {
	var err error

	query := `UPDATE Product SET productName = ?, category = ?, price = ?, minBidPrice = ?, startAt = ?, endedAt = ?, productDescription = ? WHERE productId = ?`
    _, err = mariadb.DB.Exec(query, pr.ProductName, pr.Category, pr.Price, pr.MinBidPrice, pr.StartAt, pr.EndedAt, pr.Description, pr.ProductId)
    if err != nil {
        logger.Error("[ADMIN] " + err.Error())
        return err
    }

	logger.Info("[ADMIN] Successfully update all the products")
	return err
}

func updateImage(files []*multipart.FileHeader, productID string) error {

	var imageID string

	// Create the directory to save the uploaded files
	uploadDir := "./assets"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		logger.Error("[Member] Unable to create assets directory")
		return err
	}

	// Step 1: Find and delete existing images for the productID from the database
	query := `SELECT imageId, imageUrl FROM ProductImage WHERE productId = ?`
	rows, err := mariadb.DB.Query(query, productID)
	if err != nil {
		logger.Error("[Member] Unable to fetch existing images: " + err.Error())
		return err
	}
	defer rows.Close()

	var existingImages []struct {
		ImageID  string
		ImageUrl string
	}

	for rows.Next() {
		var img struct {
			ImageID  string
			ImageUrl string
		}
		if err := rows.Scan(&img.ImageID, &img.ImageUrl); err != nil {
			logger.Error("[Member] Unable to scan row: " + err.Error())
			return err
		}
		existingImages = append(existingImages, img)
	}

	// Delete the images from the upload directory
	for _, img := range existingImages {
		filePath := filepath.Join(uploadDir, filepath.Base(img.ImageID))
		if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
			logger.Error("[Member] Unable to delete file: " + err.Error())
			return err
		}
	}

	// Delete the images from the database
	query = `DELETE FROM ProductImage WHERE productId = ?`
	if _, err := mariadb.DB.Exec(query, productID); err != nil {
		logger.Error("[Member] Unable to delete images from database: " + err.Error())
		return err
	}

	// Step 2: Upload new images and update the database
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
		defer src.Close()

		// Check the file type
		if err := checkFileType(src); err != nil {
			logger.Error("[Member] Invalid file type")
			return errors.New("Invalid file type")
		}

		// Prefix the filename with "<imageID>"
		newFilename := fmt.Sprintf("%s", imageID)
		dstPath := filepath.Join(uploadDir, newFilename)
		dst, err := os.Create(dstPath)
		if err != nil {
			logger.Error("[Member] Unable to create file")
			return err
		}
		defer dst.Close()

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

	logger.Info("[Member] Successfully updated images")
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

func deleteProduct(pr DeleteProductRequest) (error) {
	var err error
	uploadDir := "./assets"

	// Find and delete existing images for the productID from the database
	query := `SELECT imageId, imageUrl FROM ProductImage WHERE productId = ?`
	rows, err := mariadb.DB.Query(query, pr.ProductId)
	if err != nil {
		logger.Error("[ADMIN] Unable to fetch existing images: " + err.Error())
		return err
	}
	defer rows.Close()

	var existingImages []struct {
		ImageID  string
		ImageUrl string
	}

	for rows.Next() {
		var img struct {
			ImageID  string
			ImageUrl string
		}
		if err := rows.Scan(&img.ImageID, &img.ImageUrl); err != nil {
			logger.Error("[ADMIN] Unable to scan row: " + err.Error())
			return err
		}
		existingImages = append(existingImages, img)
	}

	// Delete the images from the upload directory
	for _, img := range existingImages {
		filePath := filepath.Join(uploadDir, filepath.Base(img.ImageID))
		if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
			logger.Error("[ADMIN] Unable to delete file: " + err.Error())
			return err
		}
	}

	// Define queries to delete related records
    deleteProductImagesQuery := `DELETE FROM ProductImage WHERE productId = ?`
    deleteHistoryQuery := `DELETE FROM History WHERE productId = ?`
    deleteTrackingListQuery := `DELETE FROM TrackingList WHERE productId = ?`
    deleteProductQuery := `DELETE FROM Product WHERE productId = ?`

    _, err = mariadb.DB.Exec(deleteProductImagesQuery, pr.ProductId)
    if err != nil {
        logger.Error("[ADMIN] " + err.Error())
        return err
    }

	_, err = mariadb.DB.Exec(deleteHistoryQuery, pr.ProductId)
    if err != nil {
        logger.Error("[ADMIN] " + err.Error())
        return err
    }

	_, err = mariadb.DB.Exec(deleteTrackingListQuery, pr.ProductId)
    if err != nil {
        logger.Error("[ADMIN] " + err.Error())
        return err
    }

	_, err = mariadb.DB.Exec(deleteProductQuery, pr.ProductId)
    if err != nil {
        logger.Error("[ADMIN] " + err.Error())
        return err
    }

	logger.Info("[ADMIN] Successfully delete the product")
	return err
}