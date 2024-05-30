package member

import (
	"Guang_Shun_Coin_Action/pkg/logger"
	"Guang_Shun_Coin_Action/pkg/mariadb"

	"github.com/google/uuid"
)

func addProduct(rr addProductRequest) error {
	var query string
	var err error

	// Insert into user database
	query = "INSERT INTO Product (productId, productName, category, price, minBidPrice, createAt, endedAt) VALUES (?, ?, ?, ?, ?, ?, ?)"
	_, err = mariadb.DB.Exec(query, uuid.NewString(), rr.Name, rr.Category, rr.Price, rr.MinBidPrice, rr.StartDate, rr.EndDate)
	if err != nil {
		logger.Error("[Product] " + err.Error())
		return err
	}

	logger.Info("[Product] Successfully add product with productname: " + rr.Name)

	return nil
}
