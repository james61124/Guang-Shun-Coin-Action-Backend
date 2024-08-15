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
	"strings"
)

func getProductQuery(pr ProductRequest) (string, string, []interface{}) {
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
            AND (p.endedAt IS NULL OR p.endedAt > NOW())
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
                AND (p.endedAt IS NULL OR p.endedAt > NOW())
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
func product(pr ProductRequest, UUID string) ([]response.ProductResponse, int, error) {
	var err error

	query, countQuery, args := getProductQuery(pr)

	// Execute the count query to get the total number of products matching the criteria
    var totalProductCount int
	var totalPages int
    err = mariadb.DB.QueryRow(countQuery, args[0], args[1], args[2], args[3]).Scan(&totalProductCount)
    if err != nil {
        logger.Error("[SHOP] " + err.Error())
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
			logger.Warn("[SHOP] No product is available currently.\n")
			return nil, totalPages, err
		}
		logger.Error("[SHOP] " + err.Error())
		return nil, totalPages, err
	}
	defer rows.Close()
	
	var startAt, endedAt string
	var products []response.ProductResponse
	for rows.Next() {
		var p response.ProductResponse
		err = rows.Scan(&p.ProductId, &p.ProductName, &p.Category, &p.Price, &p.MinBidPrice, &startAt, &endedAt, &p.BidCount)
		if err != nil {
			logger.Error("[SHOP] " + err.Error())
			return products, totalPages, err
		}
		p.StartAt, err = time.Parse("2006-01-02 15:04:05", startAt)
		if err != nil {
			logger.Error("[SHOP] " + err.Error())
			return products, totalPages, err
		}

		p.EndedAt, err = time.Parse("2006-01-02 15:04:05", endedAt)
		if err != nil {
			logger.Error("[SHOP] " + err.Error())
			return products, totalPages, err
		}
		
		query := `SELECT imageUrl FROM ProductImage WHERE productId = ? ORDER BY seq LIMIT 1;`;
		var imageUrl string
		err := mariadb.DB.QueryRow(query, p.ProductId).Scan(&imageUrl)
		if err != nil {
			logger.Error("[SHOP] " + err.Error())
			return products, totalPages, err
		}
		p.ImgUrl = imageUrl

		query = `
			SELECT EXISTS(
				SELECT 1
				FROM TrackingList
				WHERE userId = ? AND productId = ?
			) AS isTracking;
		`
		err = mariadb.DB.QueryRow(query, UUID, p.ProductId).Scan(&p.IsStar)
		if err != nil {
			logger.Error("[SHOP] " + err.Error())
			return products, totalPages, err
		}
		products = append(products, p)
	}

	logger.Info("[SHOP] Successfully return all the products")
	return products, totalPages, err
}

func totalPagesOfProduct(pr TotalPagesOfProductRequest) (int, error) {
	var err error
	var totalPages int

	// Updated query to include the condition for endedAt
	query := `SELECT COUNT(*) AS product_count FROM Product WHERE category = ? AND (endedAt IS NULL OR endedAt > NOW())`

	var productCount int
	err = mariadb.DB.QueryRow(query, pr.Category).Scan(&productCount)
	if err != nil {
		logger.Error("[SHOP] " + err.Error())
		return productCount, err
	}

	if ( productCount % 12 == 0 ) {
		totalPages = productCount / 12
	} else {
		totalPages = productCount / 12 + 1
	}

	total := fmt.Sprintf("%d", productCount / 12 + 1)
	logger.Info("[SHOP] Successfully return total product number: " + total)

	return totalPages, err
}

func detail(UUID string, pr DetailRequest) (response.DetailResponse, error) {
	var err error
	var product response.DetailResponse
	var startAt, endedAt string

	// select the product information
	query := "SELECT productName, category, price, minBidPrice, startAt, endedAt, productDescription FROM Product WHERE productId = ?"
	err = mariadb.DB.QueryRow(query, pr.ProductID).Scan(&product.Name, &product.Category, &product.Price, &product.MinBidPrice, &startAt, &endedAt, &product.Description)
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

	// count total number of history records
	query = "SELECT COUNT(*) FROM History WHERE productId = ?"
	var totalRecords int
	err = mariadb.DB.QueryRow(query, pr.ProductID).Scan(&totalRecords)
	if err != nil {
		logger.Error("[SHOP] " + err.Error())
		return product, err
	}

	// calculate total pages of history
	product.TotalPageOfHistory = (totalRecords + 9) / 10 // each page contains 10 records

	// get history record with pagination
	offset := (pr.HistoryPage - 1) * 10
	query = "SELECT userId, bidPrice, bidTime, status FROM History WHERE productId = ? ORDER BY bidTime DESC LIMIT 10 OFFSET ?"
	rows, err = mariadb.DB.Query(query, pr.ProductID, offset)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			logger.Warn("[SHOP] No product is available currently.\n")
		}
		logger.Error("[SHOP] " + err.Error())
	}
	defer rows.Close()

	var highestBidPrice float64 = product.Price

	for rows.Next() {
		var history response.BidHistory
		var userId, bidTime string

		if err = rows.Scan(&userId, &history.BidPrice, &bidTime, &history.Status); err != nil {
			logger.Error("[SHOP] " + err.Error())
			return product, err
		}
		history.BidTime, err = time.Parse("2006-01-02 15:04:05", bidTime)
		if err != nil {
			logger.Error("[SHOP] " + err.Error())
			return product, err
		}

		if history.BidPrice > highestBidPrice {
			highestBidPrice = history.BidPrice
		}

		// get username
		query = "SELECT nickName FROM User WHERE userId = ?"
		err = mariadb.DB.QueryRow(query, userId).Scan(&history.Username)
		if err != nil {
			logger.Error("[SHOP] " + err.Error())
			return product, err
		}

		product.History = append(product.History, history)
	}

	// set current price
	product.CurrentPrice = highestBidPrice

	query = `
		SELECT EXISTS(
			SELECT 1
			FROM TrackingList
			WHERE userId = ? AND productId = ?
		) AS isTracking;
	`
	err = mariadb.DB.QueryRow(query, UUID, pr.ProductID).Scan(&product.IsStar)
	if err != nil {
		logger.Error("[SHOP] " + err.Error())
		return product, err
	}

	logger.Info("[SHOP] Successfully return product detail")
	return product, err
}

func bid(pr BidRequest, UUID string) error {

	historyId := uuid.NewString()

	var highestBidPrice sql.NullInt64
	var midBidPrice sql.NullInt64

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

	query = `
		SELECT minBidPrice
		FROM Product
		WHERE productId = ?;
	`
	err = mariadb.DB.QueryRow(query, pr.ProductID).Scan(&midBidPrice)
	if err != nil {
		logger.Error("[SHOP] " + err.Error())
		return err
	}

	// history has no record
	if !highestBidPrice.Valid {
		query = `
			SELECT price
			FROM Product
			WHERE productId = ?;
		`
		err = mariadb.DB.QueryRow(query, pr.ProductID).Scan(&highestBidPrice)
		if err != nil {
			logger.Error("[SHOP] " + err.Error())
			return err
		}
	}
	
	if highestBidPrice.Int64 >= int64(pr.BidPrice) {
		logger.Warn("[SHOP] bidPrice is smaller or equal than highest bidPrice")
		return errors.New("bidPrice is smaller or equal than highest bidPrice")
	} else if (int64(pr.BidPrice) - highestBidPrice.Int64) < midBidPrice.Int64 {
		logger.Warn("[SHOP] bidPrice is smaller than midBidPrice")
		return errors.New("bidPrice is smaller than midBidPrice")
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

	logger.Info("[SHOP] Successfully bid")
	return nil
	
}

func star(pr StarRequest, UUID string) error {

	var query string
	trackingID := uuid.NewString()

	if pr.IsStar == true {
		// insert a new TrackingList
		query = `
			INSERT INTO TrackingList (trackingId, userId, productId)
			VALUES (?, ?, ?);
			`
		_, err := mariadb.DB.Exec(query, trackingID, UUID, pr.ProductID)
		if err != nil {
			logger.Error("[SHOP] " + err.Error())
			return err
		}
	} else {
		query := `
			DELETE FROM TrackingList 
			WHERE userId = ? AND productId = ?
			`
		_, err := mariadb.DB.Exec(query, UUID, pr.ProductID)
		if err != nil {
			logger.Error("[SHOP] " + err.Error())
			return err
		}

	}
	
	logger.Info("[SHOP] Successfully insert star info")
	return nil
	
}