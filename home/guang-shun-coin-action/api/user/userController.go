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

type loginRequest struct {
	Cellphone string `json:"Cellphone" binding:"required"`
	Password string `json:"Password" binding:"required"`
}

type registerRequest struct {
	UserID    string `json:"UserID"`
	Cellphone string `json:"Cellphone" binding:"required"`
	Password  string `json:"Password" binding:"required"`
	PasswordConfirm  string `json:"PasswordConfirm" binding:"required"`
	RealName  string `json:"RealName"`
	NickName  string `json:"NickName"`
}

type updateUserInfoRequest struct {
	RealName        string      `json:"realName" binding:"required"`
	NickName        string      `json:"nickName" binding:"required"`
	Cellphone        string     `json:"cellphone" binding:"required"`
}

type updatePasswordRequest struct {
	OriginPassword     string      `json:"originPassword" binding:"required"`
	NewPassword        string      `json:"newPassword" binding:"required"`
	ConfirmNewPassword string      `json:"confirmNewPassword" binding:"required"`
}
type resetPasswordRequest struct {
	Cellphone 		   string 	   `json:"cellphone" binding:"required"`
	NewPassword        string      `json:"newPassword" binding:"required"`
	ConfirmNewPassword string      `json:"confirmNewPassword" binding:"required"`
}

var errorMessages = map[string]bool {
	// cellphone
    "cellphone is empty":    true,
    "can't find the user with cellphone":  true,
	"invalid phone number format": true,
	"user already exists": true,
	"the new cellphone already exists": true,
	"unable to find the user based on the cellphone": true,
	// password
    "password is empty":     true,
	"invalid password format": true,
    "incorrect password":    true,
	"password is different from PasswordConfirm": true, 
	"originalPassword is wrong": true,
	
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
		errMessage, _ := r.Message.(string)
		if errorMessages[errMessage] {
			c.JSON(http.StatusOK, r)
			return
		}
		logger.Warn("[USER] " + err.Error())
		c.JSON(http.StatusInternalServerError, r)
		return
	}

	// Generate token
	token, err := auth.GenerateToken(UUID, loginRequest.Cellphone)
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
		errMessage, _ := r.Message.(string)
		if errorMessages[errMessage] {
			c.JSON(http.StatusOK, r)
			return
		}
		logger.Warn("[USER] " + err.Error())
		c.JSON(http.StatusInternalServerError, r)
		return
	}

	// return UUID with formatted response
	r.Status = true
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
		errMessage, _ := r.Message.(string)
		if errorMessages[errMessage] {
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
		errMessage, _ := r.Message.(string)
		if errorMessages[errMessage] {
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

func ResetPassword(c *gin.Context) {
	var err error

	// Create response
	r := response.New()

	// Parse request body to JSON format
	var resetPasswordRequest resetPasswordRequest
	if err := c.ShouldBindJSON(&resetPasswordRequest); err != nil {
		logger.Warn("[USER] " + err.Error())
		r.Message = err.Error()
		c.JSON(http.StatusOK, r)
		return
	}

	err = resetPassword(resetPasswordRequest)
	if err != nil {
		r.Message = err.Error()
		errMessage, _ := r.Message.(string)
		if errorMessages[errMessage] {
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
