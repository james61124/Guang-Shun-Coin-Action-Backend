package auth

import (
	"Guang_Shun_Coin_Action/config"
	"Guang_Shun_Coin_Action/internal/response"
	"Guang_Shun_Coin_Action/pkg/logger"
	"Guang_Shun_Coin_Action/pkg/mariadb"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var jwtSecretKey []byte

type authClaims struct {
	UUID string `json:"UUID"`
	jwt.StandardClaims
}

func SetJWTKey() {
	jwtSecretKey = []byte(config.Viper.GetString("JWT_SECRET_KEY"))
}

func GenerateToken(UUID, email string) (string, error) {
	// Set JWT claims fields
	expiresAt := time.Now().Add(30 * 24 * time.Hour).Unix()	// 30 * 24 hours
    token := jwt.NewWithClaims(jwt.SigningMethodHS512, authClaims{
        UUID: UUID,
		StandardClaims: jwt.StandardClaims{
            Subject:   email,
            ExpiresAt: expiresAt,
        },
    })

	// Sign the token with our secret key
    tokenString, err := token.SignedString(jwtSecretKey)
    if err != nil {
		logger.Error("[AUTH] Failed to generate token: " + err.Error())
        return "", err
    }

	logger.Info("[AUTH] Generated token for user: " + email)

    return tokenString, nil
}

func ValidateToken(c *gin.Context) string {
	var query string

	// Get token from header
	auth := c.GetHeader("Authorization")
	if auth == "" {
		r := response.New()
		r.Message = "Authorization header is missing"
		logger.Warn("[AUTH] Received request without Bearer authorization header")
		c.JSON(http.StatusOK, r)
		c.Abort()
		return ""
	}
    token := strings.Split(auth, "Bearer ")[1]

	// Parse token
	tokenClaims, err := jwt.ParseWithClaims(token, &authClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return jwtSecretKey, nil
	})
	// Check for token validation errors
	if err != nil {
		var r = response.New()
		if ve, ok := err.(*jwt.ValidationError); ok {

			// testing error message: 1
			fmt.Printf("%d", ve.Errors)

			if ve.Errors & jwt.ValidationErrorMalformed != 0 {
				r.Message = "token is not correctly formatted as a JWT (missing or invalid segments)"
			} else if ve.Errors & jwt.ValidationErrorUnverifiable != 0{
				r.Message = "token cannot be verified due to problems with the token's signature"
			} else if ve.Errors & jwt.ValidationErrorSignatureInvalid != 0 {
				r.Message = "signature validation failed (token's content has been tampered with)"
			} else if ve.Errors & jwt.ValidationErrorExpired != 0 {
				r.Message = "token is expired"
			} else {
				r.Message = "can not handle this token"
			}
		}
		c.JSON(http.StatusUnauthorized, r)
		c.Abort()
		return ""
	}

	// Check if token is valid -> continue
	if claims, ok := tokenClaims.Claims.(*authClaims); ok && tokenClaims.Valid {
		c.Set("UUID", claims.UUID)
		c.Next()
	} else {
		c.Abort()
		return ""
	}

	// Get UUID from context
	uuid, exists := c.Get("UUID")
	if !exists {
		logger.Warn("[AUTH] UUID not found from auth")
		return ""
	}

	// Type assert UUID to string
	UUID, ok := uuid.(string)
	if !ok {
		logger.Error("[AUTH] UUID not a string")
		return ""
	}

	// Check if user exists
	query = "SELECT EXISTS(SELECT 1 FROM User WHERE userId = ?)"
    err = mariadb.DB.QueryRow(query, UUID).Scan(&exists)
	if err != nil && err.Error() != "sql: no rows in result set" {
		var r = response.New()
		logger.Error("[AUTH] " + err.Error())
		r.Message = err.Error()
		c.JSON(http.StatusBadRequest, r)
		return ""
	} else if exists == false {
		var r = response.New()
		logger.Error("[AUTH] user doesn't exists")
		r.Message = "user doesn't exists"
		c.JSON(http.StatusBadRequest, r)
		return ""
	}

	return UUID
}
