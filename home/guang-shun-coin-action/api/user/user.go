package user

import (
	"Guang_Shun_Coin_Action/pkg/logger"
	"Guang_Shun_Coin_Action/pkg/mariadb"
	"errors"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserID string `json:"UserID"`
	Username string `json:"Username"`
	Password string `json:"Password"`
	Cellphone string `json:"Cellphone"`
	FbAccount string `json:"FbAccount"`
	Email string `json:"Email"`
	Address string `json:"Address"`
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func register(rr registerRequest) (error) {
	var query, username string
	var err error

	// Check if user already exists
	query = "SELECT username FROM `User` WHERE username = ?"
	err = mariadb.DB.QueryRow(query, rr.Username).Scan(&username)
	if err != nil && err.Error() != "sql: no rows in result set" {
		logger.Error("[USER] " + err.Error())
		return err
	} else if username != "" {
		logger.Warn("[USER] username:" + rr.Username + " already exists")
		return errors.New("username already exists")
	}

	// Check pass in fields (Email has @ symbol)
	if !strings.Contains(rr.Email, "@") {
		logger.Warn("[USER] Invalid email address")
		return errors.New("invalid email address")
	}

	// Hash password
	if rr.Password, err = hashPassword(rr.Password); err != nil {
		logger.Error("[USER] " + err.Error())
		return err
	}

	// Insert into user database
	query = "INSERT INTO User (user_id, username, password, cellphone, fb_account, email, address) VALUES (?, ?, ?, ?, ?, ?, ?)"
	_, err = mariadb.DB.Exec(query, uuid.NewString(), rr.Username, rr.Password, rr.Cellphone, rr.FbAccount, rr.Email, rr.Address)
	if err != nil {
		logger.Error("[USER] " + err.Error())
		return err
	}

	logger.Info("[USER] Successfully registered user with username: " + rr.Username)

	return nil
}

func login(lr loginRequest) (string, error) {
	var query, UUID, password string
	var err error

	// Get user password
	query = "SELECT user_id, password FROM User WHERE username = ?"
	err = mariadb.DB.QueryRow(query, lr.Username).Scan(&UUID, &password)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			logger.Warn("[USER] Username: " + lr.Username + " not found")
			return "", errors.New("user not found")
		}
		logger.Error("[USER] " + err.Error())
		return "", err
	}

	// Check if password is correct
	if !checkPasswordHash(lr.Password, password) {
		logger.Warn("[USER] Incorrect password for Username: " + lr.Username)
		return "", errors.New("incorrect password")
	}

	logger.Info("[USER] Successfully logged in user with Username: " + lr.Username)

	return UUID, nil
}

// func get(UUID string) (UserInfo, error) {
// 	var query string
// 	var err error
// 	var userInfo UserInfo

// 	// Get user info (UUID, Username, Is_Pro)
// 	query = "SELECT UUID, Username, Is_Pro FROM User_Info WHERE UUID = ?"
// 	err = mariadb.DB.QueryRow(query, UUID).Scan(&userInfo.UUID, &userInfo.Username, &userInfo.Is_Pro)
// 	if err != nil {
// 		if err.Error() == "sql: no rows in result set" {
// 			logger.Warn("[USER] UUID: " + UUID + " not found")
// 			return userInfo, err
// 		}
// 		logger.Error("[USER] " + err.Error())
// 		return userInfo, err
// 	}

// 	// Get user email
// 	query = "SELECT Email FROM User WHERE UUID = ?"
// 	err = mariadb.DB.QueryRow(query, UUID).Scan(&userInfo.Email)
// 	if err != nil {
// 		logger.Error("[USER] " + err.Error())
// 		return userInfo, err
// 	}

// 	logger.Info("[USER] Successfully retrieved user info for UUID: " + UUID)

// 	return userInfo, nil
// }

// func update(ur updateRequest) error {
// 	var query string

// 	// Update user info
// 	query = "UPDATE User_Info SET Username = ?, Is_Pro = ? WHERE UUID = ?"
// 	result, err := mariadb.DB.Exec(query, ur.Username, ur.Is_Pro, ur.UUID)
// 	if err != nil {
// 		logger.Error("[USER] " + err.Error())
// 		return err
// 	}

// 	// Check if UUID exists
// 	rowsaffected, err := result.RowsAffected()
// 	if err != nil {
// 		logger.Error("[USER] " + err.Error())
// 		return err
// 	} else if rowsaffected == 0 {
// 		logger.Warn("[USER] UUID: " + ur.UUID + " not found")
// 		return errors.New("user not found")
// 	}

// 	// Update user email
// 	query = "UPDATE User SET Email = ? WHERE UUID = ?"
// 	_, err = mariadb.DB.Exec(query, ur.Email, ur.UUID)
// 	if err != nil {
// 		logger.Error("[USER] " + err.Error())
// 		return err
// 	}

// 	logger.Info("[USER] Successfully updated user info for UUID: " + ur.UUID)

// 	return nil
// }


