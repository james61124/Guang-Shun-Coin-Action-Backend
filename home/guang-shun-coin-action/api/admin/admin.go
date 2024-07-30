package admin

import (
	"Guang_Shun_Coin_Action/pkg/logger"
	"Guang_Shun_Coin_Action/pkg/mariadb"
	"Guang_Shun_Coin_Action/internal/response"
	"time"
	"strings"
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