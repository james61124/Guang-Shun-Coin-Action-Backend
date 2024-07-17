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
	Username  string `json:"Username"`
	Password  string `json:"Password"`
	PasswordConfirm  string `json:"PasswordConfirm"`
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

type updateUserInfoRequest struct {
	RealName        string      `json:"realName"`
	NickName        string      `json:"nickName"`
	Cellphone        string      `json:"cellphone"`
	FbAccount        string      `json:"fbAccount"`
	Email        string      `json:"email"`
	Postcode        string      `json:"postcode"`
	ShippingAddr        string      `json:"shippingAddr"`
	Username        string      `json:"username"`
}

type updatePasswordRequest struct {
	OriginPassword        string      `json:"originPassword"`
	NewPassword        string      `json:"newPassword"`
	ConfirmNewPassword        string      `json:"confirmNewPassword"`
}

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
		if r.Message == "username is empty" || r.Message == "password is empty" {
			c.JSON(http.StatusOK, r)
			return
		}
		if r.Message == "passwordConfirm is empty" || r.Message == "address is empty" {
			c.JSON(http.StatusOK, r)
			return
		}
		if r.Message == "cellphone is empty" || r.Message == "username already exists" {
			c.JSON(http.StatusOK, r)
			return
		}
		if r.Message == "invalid email address" || r.Message == "invalid phone number format" {
			c.JSON(http.StatusOK, r)
			return
		}
		if r.Message == "password and passwordConfirm is different" || r.Message == "invalid password format" {
			c.JSON(http.StatusOK, r)
			return
		}
		if r.Message == "cellphone already exists" {
			c.JSON(http.StatusOK, r)
			return
		}
		c.JSON(http.StatusInternalServerError, r)
		return
	}

	// return UUID with formatted response
	r.Status = true
	c.JSON(http.StatusOK, r)
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
		c.JSON(http.StatusOK, r)
		return
	}

	// Login the user
	UUID, err := login(loginRequest)
	if err != nil {
		r.Message = err.Error()
		if r.Message == "incorrent password" || r.Message == "username not found" {
			c.JSON(http.StatusOK, r)
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

func GetUserInfo(c *gin.Context) {
	var err error

	// Create response
	r := response.New()

	// Validate token and set UUID in context (validation failed.)
	UUID := auth.ValidateToken(c)

	userInfo, err := getUserInfo(UUID)
	if err != nil {
		r.Message = err.Error()
		logger.Warn("[User] " + err.Error())
		c.JSON(http.StatusInternalServerError, r)
		return
	}

	r.Status = true
	r.Data = userInfo
	c.JSON(http.StatusOK, r)
}

func UpdateUserInfo(c *gin.Context) {
	var err error

	// Create response
	r := response.New()

	// Validate token and set UUID in context (validation failed.)
	UUID := auth.ValidateToken(c)

	// Parse request body to JSON format
	var updateUserInfoRequest updateUserInfoRequest
	if err := c.ShouldBindJSON(&updateUserInfoRequest); err != nil {
		logger.Warn("[USER] " + err.Error())
		r.Message = err.Error()
		c.JSON(http.StatusOK, r)
		return
	}

	err = updateUserInfo(UUID, updateUserInfoRequest)
	if err != nil {
		r.Message = err.Error()
		if r.Message == "username is empty" || r.Message == "address is empty" {
			c.JSON(http.StatusOK, r)
			return
		}
		if r.Message == "cellphone is empty" || r.Message == "username already exists" {
			c.JSON(http.StatusOK, r)
			return
		}
		if r.Message == "invalid email address" || r.Message == "invalid phone number format" {
			c.JSON(http.StatusOK, r)
			return
		}
		if r.Message == "cellphone already exists" {
			c.JSON(http.StatusOK, r)
			return
		}
		logger.Warn("[USER] " + err.Error())
		c.JSON(http.StatusInternalServerError, r)
		return
	}

	r.Status = true
	c.JSON(http.StatusOK, r)
}

func UpdatePassword(c *gin.Context) {
	var err error

	// Create response
	r := response.New()

	// Validate token and set UUID in context (validation failed.)
	UUID := auth.ValidateToken(c)

	// Parse request body to JSON format
	var updatePasswordRequest updatePasswordRequest
	if err := c.ShouldBindJSON(&updatePasswordRequest); err != nil {
		logger.Warn("[USER] " + err.Error())
		r.Message = err.Error()
		c.JSON(http.StatusOK, r)
		return
	}

	err = updatePassword(UUID, updatePasswordRequest)
	if err != nil {
		r.Message = err.Error()
		expectedErrors := []string{
			"originalPassword is empty",
			"newPassword is empty",
			"confirmNewPassword is empty",
			"originalPassword is wrong",
			"newPassword is different from confirmNewPassword",
			"invalid newPassword format",
		}
		for _, e := range expectedErrors {
			if r.Message == e {
				c.JSON(http.StatusOK, r)
				return
			}
		}
		
		logger.Warn("[USER] " + err.Error())
		c.JSON(http.StatusInternalServerError, r)
		return
	}

	r.Status = true
	c.JSON(http.StatusOK, r)
}
