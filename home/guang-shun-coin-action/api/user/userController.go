package user

import (
	"Guang_Shun_Coin_Action/internal/auth"
	"Guang_Shun_Coin_Action/internal/response"
	"Guang_Shun_Coin_Action/pkg/logger"
	"net/http"
	// "regexp"
	// "strings"

	"github.com/gin-gonic/gin"
)

type registerRequest struct {
	UserID    string `json:"UserID"`
	Username  string `json:"Username" binding:"required"`
	Password  string `json:"Password" binding:"required"`
	Cellphone string `json:"Cellphone"`
	FbAccount string `json:"FbAccount"`
	Email     string `json:"Email"`
	Address   string `json:"Address"`
	Postcode  string `json:"Postcode"`
	RealName  string `json:"RealName"`
	NickName  string `json:"NickName"`
}

type loginRequest struct {
	Username string `json:"Username" binding:"required"`
	Password string `json:"Password" binding:"required"`
}

// type updateRequest struct {
// 	UUID     string `json:"UUID" binding:"required"`
// 	Username string `json:"Username" binding:"required"`
// 	Email    string `json:"Email" binding:"required"`
// 	Is_Pro   bool   `json:"Is_Pro"`
// }

func Register(c *gin.Context) {
	var err error

	// Create response
	r := response.New()

	// Parse request body to JSON format
	var registerRequest registerRequest
	if err = c.ShouldBindJSON(&registerRequest); err != nil {
		logger.Warn("[USER] " + err.Error())
		r.Message = err.Error()
		c.JSON(http.StatusBadRequest, r)
		return
	}

	// Register the user
	err = register(registerRequest)
	if err != nil {
		r.Message = err.Error()
		if r.Message == "username already exists" {
			c.JSON(http.StatusBadRequest, r)
			return
		}
		if r.Message == "invalid email address" {
			c.JSON(http.StatusBadRequest, r)
			return
		}
		c.JSON(http.StatusInternalServerError, r)
		return
	}

	// return UUID with formatted response
	r.Status = true
	// r.Data = response.RegisterResponse{username: registerRequest.username, password: registerRequest.password}
	c.JSON(http.StatusCreated, r)
}

func Login(c *gin.Context) {
	var err error

	// Create response
	r := response.New()

	// Parse request body to JSON format
	var loginRequest loginRequest
	if err = c.ShouldBindJSON(&loginRequest); err != nil {
		logger.Warn("[USER] " + err.Error())
		r.Message = err.Error()
		c.JSON(http.StatusBadRequest, r)
		return
	}

	// Login the user
	UUID, err := login(loginRequest)
	if err != nil {
		r.Message = err.Error()
		if r.Message == "user not found" {
			c.JSON(http.StatusNotFound, r)
			return
		} else if r.Message == "incorrect password" {
			c.JSON(http.StatusUnauthorized, r)
			return
		}
		c.JSON(http.StatusInternalServerError, r)
		return
	}

	// Generate token
	token, err := auth.GenerateToken(UUID, loginRequest.Username)
	if err != nil {
		r.Message = err.Error()
		c.JSON(http.StatusInternalServerError, r)
		return
	}

	// return UUID with formatted response
	r.Status = true
	r.Data = response.LoginResponse{UUID: UUID, Token: token}
	c.JSON(http.StatusOK, r)
}

// func Get(c *gin.Context) {
// 	var err error

// 	// Create response
// 	r := response.New()

// 	// Check UUID format
// 	if matched, _ := regexp.MatchString("^[a-z0-9-]{36}$", c.Param("uuid")); !matched {
// 		logger.Warn("[LEDGER] UUID is not valid format")
// 		r.Message = "UUID is not valid format"
// 		c.JSON(http.StatusBadRequest, r)
// 		return
// 	}

// 	// Get user info
// 	userInfo, err := get(c.Param("uuid"))
// 	if err != nil {
// 		if err.Error() == "sql: no rows in result set" {
// 			r.Message = "user not found"
// 			c.JSON(http.StatusNotFound, r)
// 			return
// 		}
// 		r.Message = err.Error()
// 		c.JSON(http.StatusInternalServerError, r)
// 		return
// 	}

// 	// return user info with formatted response
// 	r.Status = true
// 	r.Data = userInfo
// 	c.JSON(http.StatusOK, r)
// }

// func Update(c *gin.Context) {
// 	var err error

// 	// Create response
// 	r := response.New()

// 	// Parse request body to JSON format
// 	var updateRequest updateRequest
// 	if err = c.ShouldBindJSON(&updateRequest); err != nil {
// 		logger.Warn("[USER] " + err.Error())
// 		r.Message = err.Error()
// 		c.JSON(http.StatusBadRequest, r)
// 		return
// 	}

// 	// Can not update other user's info
// 	if updateRequest.UUID != c.MustGet("UUID") {
// 		logger.Warn("[USER] Can not update other user's info")
// 		r.Message = "Can not update other user's info"
// 		c.JSON(http.StatusUnauthorized, r)
// 		return
// 	}

// 	// Update user info
// 	err = update(updateRequest)
// 	if err != nil {
// 		r.Message = err.Error()
// 		if r.Message == "user not found" {
// 			c.JSON(http.StatusNotFound, r)
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, r)
// 		return
// 	}

// 	// return formatted response
// 	r.Status = true
// 	c.JSON(http.StatusOK, r)
// }


