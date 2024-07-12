package shop

import (
	"Guang_Shun_Coin_Action/pkg/logger"
	"Guang_Shun_Coin_Action/pkg/mariadb"
	"Guang_Shun_Coin_Action/internal/response"
	"time"
	"fmt"
	"github.com/google/uuid"
	"errors"
	"database/sql"
)

// return status, list of products
func product(pr ProductRequest) ([]response.ProductResponse, error) {
	var err error
	
	query := "SELECT productId, productName, category, price, minBidPrice, startAt, endedAt FROM Product WHERE category = ? LIMIT ? OFFSET ?"
	productNums := 12
	offset := (pr.Page - 1) * productNums
	rows, err := mariadb.DB.Query(query, pr.Category, productNums, offset)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			logger.Warn("[SHOP] No product is available currently.\n")
		}
		logger.Error("[SHOP] " + err.Error())
	}
	defer rows.Close()
	
	var startAt, endedAt string
	var products []response.ProductResponse
	for rows.Next() {
		var p response.ProductResponse
		err = rows.Scan(&p.ProductId, &p.ProductName, &p.Category, &p.Price, &p.MinBidPrice, &startAt, &endedAt)
		if err != nil {
			logger.Error("[SHOP] " + err.Error())
			return products, err
		}
		p.StartAt, err = time.Parse("2006-01-02 15:04:05", startAt)
		if err != nil {
			logger.Error("[SHOP] " + err.Error())
			return products, err
		}

		p.EndedAt, err = time.Parse("2006-01-02 15:04:05", endedAt)
		if err != nil {
			logger.Error("[SHOP] " + err.Error())
			return products, err
		}
		
		const query = `SELECT imageUrl FROM ProductImage WHERE productId = ? ORDER BY seq LIMIT 1;`;

		var imageUrl string
		err := mariadb.DB.QueryRow(query, p.ProductId).Scan(&imageUrl)
		if err != nil {
			logger.Error("[SHOP] " + err.Error())
			return products, err
		}
		p.ImgUrl = imageUrl

		products = append(products, p)
	}

	logger.Info("[SHOP] Successfully return all the products")
	return products, err
}

func totalProduct(pr TotalProductRequest) (int, error) {
	var err error

    query := `SELECT COUNT(*) AS product_count FROM Product WHERE category = ?`

    var productCount int
    err = mariadb.DB.QueryRow(query, pr.Category).Scan(&productCount)
    if err != nil {
        logger.Error("[SHOP] " + err.Error())
		return productCount, err
    }

	total := fmt.Sprintf("%d", productCount / 12 + 1)
	logger.Info("[SHOP] Successfully return total product number: " + total)
	
	return productCount, err
}

func detail(pr DetailRequest) (response.DetailResponse, error) {
	var err error
	var product response.DetailResponse
	var startAt, endedAt string
	
	// select the product information
	query := "SELECT productName, category, price, minBidPrice, startAt, endedAt FROM Product WHERE productId = ?"
	err = mariadb.DB.QueryRow(query, pr.ProductID).Scan(&product.Name, &product.Category, &product.Price, &product.MinBidPrice, &startAt, &endedAt)
	if err != nil {
		logger.Error("[SHOP] " + err.Error())
		return product, err
	}

	product.StartTime, err = time.Parse("2006-01-02 15:04:05", startAt)
	if err != nil {
		logger.Error("[SHOP] " + err.Error())
		return product, err
	}

	product.EndTime, err = time.Parse("2006-01-02 15:04:05", endedAt)
	if err != nil {
		logger.Error("[SHOP] " + err.Error())
		return product, err
	}
		
	// select image url
	query = "SELECT imageUrl FROM ProductImage WHERE productId = ? "
	rows, err := mariadb.DB.Query(query, pr.ProductID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			logger.Warn("[SHOP] No product is available currently.\n")
		}
		logger.Error("[SHOP] " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var imageUrl string
		if err := rows.Scan(&imageUrl); err != nil {
			logger.Error("[SHOP] " + err.Error())
			return product, err
		}
		product.ImageUrl = append(product.ImageUrl, imageUrl)
	}

	// get history record

	query = "SELECT userId, bidPrice, bidTime FROM History WHERE productId = ? "
	rows, err = mariadb.DB.Query(query, pr.ProductID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			logger.Warn("[SHOP] No product is available currently.\n")
		}
		logger.Error("[SHOP] " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		
		var history response.BidHistory
		var userId, bidTime string

		if err = rows.Scan(&userId, &history.BidPrice, &bidTime); err != nil {
			logger.Error("[SHOP] " + err.Error())
			return product, err
		}
		history.BidTime, err = time.Parse("2006-01-02 15:04:05", bidTime)
		if err != nil {
			logger.Error("[SHOP] " + err.Error())
			return product, err
		}
		
		// get username
		query = "SELECT username FROM User WHERE userId = ?"
		err = mariadb.DB.QueryRow(query, userId).Scan(&history.Username)
		if err != nil {
			logger.Error("[SHOP] " + err.Error())
			return product, err
		}

		product.History = append(product.History, history)
	}

	logger.Info("[SHOP] Successfully return product detail")
	return product, err
}

func bid(pr BidRequest, UUID string) error {


	historyId := uuid.NewString()

	var highestBidPrice sql.NullInt64

	query := `
		SELECT MAX(bidPrice)
		FROM History
		WHERE productId = ?;
	`
	err := mariadb.DB.QueryRow(query, pr.ProductID).Scan(&highestBidPrice)
	if err != nil {
		logger.Error("[SHOP] " + err.Error())
		return err
	}
	if !highestBidPrice.Valid {
		logger.Info("[SHOP] No historical bids for the given product.")
	} else if highestBidPrice.Int64 >= int64(pr.BidPrice) {
		logger.Warn("[SHOP] bidPrice is smaller or equal than highest bidPrice")
		return errors.New("bidPrice is smaller or equal than highest bidPrice")
	}

	// insert a new history
	query = `
        INSERT INTO History (historyId, productId, userId, bidPrice, bidTime)
        VALUES (?, ?, ?, ?, NOW());
    `
	_, err = mariadb.DB.Exec(query, historyId, pr.ProductID, UUID, pr.BidPrice)
    if err != nil {
        logger.Error("[SHOP] " + err.Error())
        return err
    }

	// set all the status to null
	query = `
		UPDATE History
		SET status = NULL
		WHERE status = '最高出價'
		AND productId = ?;
	`
	_, err = mariadb.DB.Exec(query, pr.ProductID)
	if err != nil {
        logger.Error("[SHOP] " + err.Error())
        return err
    }

	// find the highest bidPrice and set the status
	query = `
        UPDATE History
        SET status = '最高出價'
        WHERE historyId = (
            SELECT historyId
            FROM History
            WHERE productId = ?
            ORDER BY bidPrice DESC
            LIMIT 1
        );
    `
	_, err = mariadb.DB.Exec(query, pr.ProductID)
	if err != nil {
        logger.Error("[SHOP] " + err.Error())
        return err
    }

	return nil

	






	
}