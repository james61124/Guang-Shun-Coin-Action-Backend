package api

import (
	"Guang_Shun_Coin_Action/api/admin"
	"Guang_Shun_Coin_Action/api/user"
	"Guang_Shun_Coin_Action/api/shop"
	"Guang_Shun_Coin_Action/api/member"
	"Guang_Shun_Coin_Action/config"
	"Guang_Shun_Coin_Action/internal/auth"
	// "Guang_Shun_Coin_Action/internal/validator"
	"Guang_Shun_Coin_Action/pkg/logger"
	"Guang_Shun_Coin_Action/pkg/mariadb"
	// "Guang_Shun_Coin_Action/pkg/mongodb"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"github.com/gin-gonic/gin"
	"net/smtp"
	"time"
	"fmt"
	// "strings"

	"bytes"
    "crypto/tls"
    "errors"
    // "fmt"
    "net"
    // "net/smtp"
    "text/template"

)

func Main() {

	// Init API
	apiInit()
	Quit := make(chan os.Signal, 1)

	// Create gin router
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // Allowed origin
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE") // Allowed HTTP methods
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization") // Allowed request headers
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true") // Allow credentials (e.g., cookies)

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	})

	// Auth middleware for all routes below
	// r.Use(auth.ValidateToken)

	r.StaticFS("/assets",http.Dir("./assets"))

	insertAdminInfo()
	go monitorProducts()

	// Users (no token validation)
	r.POST("/user/register", user.Register)
	r.POST("/user/login", user.Login)
	r.POST("/user/getUserInfo", user.GetUserInfo)
	r.POST("/user/updateUserInfo", user.UpdateUserInfo)
	r.POST("/user/updatePassword", user.UpdatePassword)
	r.POST("/user/resetPassword", user.ResetPassword)
	
	// Product
	r.POST("/shop/product", shop.Product)
	r.POST("/shop/totalPagesOfProduct", shop.TotalPagesOfProduct)
	r.POST("/shop/detail", shop.Detail)
	r.POST("/shop/bid", shop.Bid)
	r.POST("/shop/star", shop.Star)
	r.POST("/shop/myProduct", shop.MyProduct)

	// Member
	r.POST("/member/addProduct", member.AddProduct)
	r.POST("/member/addImage", member.AddImage)
	r.POST("/member/getHistoryBid", member.GetHistoryBid)
	r.POST("/member/totalPagesOfHistoryBid", member.TotalPagesOfHistoryBid)

	// Admin
	r.POST("/admin/updateProductPage", admin.UpdateProductPage)
	r.POST("/admin/updateProductDetail", admin.UpdateProductDetail)
	r.POST("/admin/updateImage", admin.UpdateImage)
	r.POST("/admin/deleteProduct", admin.DeleteProduct)
	

	

	// Start API service
	srv := &http.Server{
		Addr:    ":" + config.Viper.GetString("GO_SERVER_PORT"),
		Handler: r,
	}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("[API] " + err.Error())
		os.Exit(1)
	}

	// Graceful shutdown
	signal.Notify(Quit, syscall.SIGINT, syscall.SIGTERM)
	<-Quit
	logger.Info("[API] Shutting down server...")
	if err := srv.Shutdown(nil); err != nil {
		logger.Error("[API] Error shutting down API server: " + err.Error())
		os.Exit(1)
	}
	if err := mariadb.Disconnect(); err != nil {
		logger.Error("[MARIADB] Error closing DB: " + err.Error())
		os.Exit(1)
	}
	// if err := mongodb.Disconnect(); err != nil {
	// 	logger.Error("[MONGODB] Error closing DB: " + err.Error())
	// 	os.Exit(1)
	// }
	logger.Info("[API] Server exited properly")
}

func apiInit() {
	config.LoadConfig() // Load config
	auth.SetJWTKey()    // Set JWT key
	logger.InitLogger() // Init logger
	ginInit()           // Init gin

	// Connect to MariaDB
	var err error
	if err = mariadb.Connect(); err != nil {
		logger.Error("[MARIADB] " + err.Error())
		return
	}

	// // Connect to MongoDB
	// if err = mongodb.Connect(); err != nil {
	// 	logger.Error("[MONGODB] " + err.Error())
	// 	return
	// }
}

func ginInit() {
	// Set gin log path
	// gin.SetMode(gin.ReleaseMode)
	f, _ := os.Create(config.Viper.GetString("GIN_LOG_PATH"))
	gin.DefaultWriter = io.MultiWriter(f)
}

func insertAdminInfo() {

	query := `INSERT INTO User (
			userId, 
			userPasswd,
			realName, 
			cellphone, 
			nickName, 
			userRole, 
			loginStatus
	) VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := mariadb.DB.Exec(
		query, 
		"4dbe7ad6-776c-425a-8b78-73c9b031a9ae", 
		"$2a$14$4IsnUa723JQbHrWy0GtZoOxenTg.HhI9TrsQvluDiCdiWXZcwNuiW",
		"admin", 
		"0912345678", 
		"admin", 
		"admin", 
		"signIn",
	)

	if err != nil {
		logger.Error("[ADMIN] " + err.Error())
		return
	}
}

// monitorProducts continuously monitors the products in the database
// and checks if their end time has passed. It runs in an infinite loop
// and calls checkProducts every minute.
func monitorProducts() {
	for {
		checkProducts()
		time.Sleep(1 * time.Minute) // Sleep for 1 minute before checking again
	}
}

// checkProducts queries the database for products whose end time has passed
// and have not yet been notified. If such products are found, it sends an
// email notification and updates the product's notification status.
func checkProducts() {
	// query := `
	// 	SELECT productId, productName, endedAt
	// 	FROM Product
	// 	WHERE endedAt <= NOW() AND notified = false
	// `
	query := `
			SELECT 
				p.productId,
				p.productName,
				p.endedAt,
				h.bidPrice,
				u.cellphone
			FROM 
				Product p
			JOIN (
				SELECT 
					productId,
					MAX(bidPrice) AS maxBidPrice
				FROM 
					History
				GROUP BY 
					productId
			) AS max_bids ON p.productId = max_bids.productId
			JOIN 
				History h ON h.productId = max_bids.productId AND h.bidPrice = max_bids.maxBidPrice
			JOIN 
				User u ON h.userId = u.userId
			WHERE 
				p.endedAt <= NOW()
				AND p.notified = false
				AND h.bidPrice IS NOT NULL;
			`
	rows, err := mariadb.DB.Query(query)
	if err != nil {
		logger.Error("[MONITOR] Failed to query products: " + err.Error())
		return
	}
	defer rows.Close()

	// Iterate over the query results
	for rows.Next() {
		var id string
		var productName string
		var endedTime string
		var bidPrice int
		var cellphone string

		// Scan the row data into variables
		if err := rows.Scan(&id, &productName, &endedTime, &bidPrice, &cellphone); err != nil {
			logger.Error("[MONITOR] Failed to scan product row: " + err.Error())
			continue
		}

		// Send an email notification for the product
		if err := sendEmail(productName, bidPrice, cellphone); err != nil {
			logger.Error("[MONITOR] Failed to send email: " + err.Error())
			continue
		}

		// Update the product's notification status in the database
		if _, err := mariadb.DB.Exec(`
			UPDATE Product
			SET notified = true
			WHERE productId = ?
		`, id); err != nil {
			logger.Error("[MONITOR] Failed to update product notification: " + err.Error())
			continue
		}

		// Log the successful notification and update
		logger.Info(fmt.Sprintf("[MONITOR] Notified and updated product: %d (%s)", id, productName))
	}
}

// sendEmail sends an email notification for the given product.
// It uses the specified SMTP settings to send the email.
func sendEmail(productName string, bidPrice int, cellphone string) error {
	// Sender data.
	from := "james61124@gmail.com"
	password := "ABC123ABC456ABC78981461"

	// Receiver email address.
	to := []string{
		cellphone,
	}

	// smtp server configuration.
	smtpHost := "smtp.office365.com"
	smtpPort := "587"

	// Establish a connection to the SMTP server
	conn, err := net.Dial("tcp", smtpHost+":"+smtpPort)
	if err != nil {
		return fmt.Errorf("failed to connect to the server: %w", err)
	}
	defer conn.Close()

	// Create a new SMTP client
	tlsconfig := &tls.Config{
		ServerName: smtpHost,
	}
	c, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer c.Close()

	// Start TLS encryption
	if err = c.StartTLS(tlsconfig); err != nil {
		return fmt.Errorf("failed to start TLS: %w", err)
	}

	// Authenticate with the SMTP server
	auth := LoginAuth(from, password)
	if err = c.Auth(auth); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Parse and execute the email template
	t, err := template.ParseFiles("template.html")
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: 廣順錢幣得標通知 \n%s\n\n", mimeHeaders)))

	err = t.Execute(&body, struct {
		Name    string
		Price int
	}{
		Name:    productName,
		Price: bidPrice,
	})
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Sending email.
	err = c.Mail(from)
	if err != nil {
		return fmt.Errorf("failed to set sender address: %w", err)
	}

	for _, addr := range to {
		err = c.Rcpt(addr)
		if err != nil {
			return fmt.Errorf("failed to set recipient address: %w", err)
		}
	}

	wc, err := c.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}
	defer wc.Close()

	_, err = wc.Write(body.Bytes())
	if err != nil {
		return fmt.Errorf("failed to write email data: %w", err)
	}

	fmt.Println("Email Sent!")
	return nil
}

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("unknown from server")
		}
	}
	return nil, nil
}