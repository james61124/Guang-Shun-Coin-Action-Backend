package member

import (
	"Guang_Shun_Coin_Action/pkg/logger"
	"Guang_Shun_Coin_Action/pkg/mariadb"
	"fmt"

	"github.com/google/uuid"
)

func addProduct(rr addProductRequest, ownerUUID string) error {
	var query string
	var err error

	shippingStatus := []string{"Unconfirmed Order", "Shipping in Progress", "Order Completed"}
	
	fmt.Printf("haha")
	fmt.Printf(ownerUUID)
	
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
		uuid.NewString(), 
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
		return err
	}


	logger.Info("[Product] Successfully add product with productname: " + rr.Name)

	return nil
}
