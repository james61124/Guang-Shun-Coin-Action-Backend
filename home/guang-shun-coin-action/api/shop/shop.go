package shop

import (
	"Guang_Shun_Coin_Action/pkg/logger"
	"Guang_Shun_Coin_Action/pkg/mariadb"
	"Guang_Shun_Coin_Action/internal/response"
	// "errors"
	// "strings"
	// "regexp"
	// "github.com/google/uuid"
	// "golang.org/x/crypto/bcrypt"
)

// return status, list of products
func product(pr ProductRequest) ([]response.ProductResponse, error) {
	var err error
	
	query := "SELECT productId, productName, category, price, minBidPrice, imgUrl, createAt, endedAt FROM Product WHERE category = ? LIMIT ? OFFSET ?"
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
	
	
	var products []response.ProductResponse
	for rows.Next() {
		var p response.ProductResponse
		err = rows.Scan(&p.ProductId, &p.ProductName, &p.Category, &p.Price, &p.MinBidPrice, &p.ImgUrl, &p.CreateAt, &p.EndedAt)
		// if err != nil {
		// 	return err
		// }
		products = append(products, p)
	}

	
	return products, err
}