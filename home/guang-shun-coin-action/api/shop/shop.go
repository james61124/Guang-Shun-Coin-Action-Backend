package shop

import (
	"Guang_Shun_Coin_Action/pkg/logger"
	"Guang_Shun_Coin_Action/pkg/mariadb"
	"Guang_Shun_Coin_Action/internal/response"
	"time"
	// "errors"
	// "strings"
	// "regexp"
	// "github.com/google/uuid"
	// "golang.org/x/crypto/bcrypt"
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
			logger.Warn("[shop] No product is available currently.\n")
		}
		logger.Error("[USER] " + err.Error())
	}
	defer rows.Close()
	
	var startAt, endedAt string
	var products []response.ProductResponse
	for rows.Next() {
		var p response.ProductResponse
		err = rows.Scan(&p.ProductId, &p.ProductName, &p.Category, &p.Price, &p.MinBidPrice, &startAt, &endedAt)
		if err != nil {
			return products, err
		}
		p.StartAt, err = time.Parse("2006-01-02 15:04:05", startAt)
		if err != nil {
			return products, err
		}

		p.EndedAt, err = time.Parse("2006-01-02 15:04:05", endedAt)
		if err != nil {
			return products, err
		}
		p.ImgUrl = "/static/"
		products = append(products, p)
	}

	
	return products, err
}