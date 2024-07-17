package user

import (
	"Guang_Shun_Coin_Action/pkg/logger"
	"Guang_Shun_Coin_Action/pkg/mariadb"
	"Guang_Shun_Coin_Action/internal/response"
	"errors"
	"strings"
	"regexp"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func register(rr registerRequest) error {
	var query, username, cellphone string
	var err error

	// Check whether username is empty
	if rr.Username == "" {
		logger.Warn("[USER] username is empty")
		return errors.New("username is empty")
	}

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

	// Check whether address is empty
	if rr.Address == "" {
		logger.Warn("[USER] address is empty")
		return errors.New("address is empty")
	}

	// Check whether cellphone is empty
	if rr.Cellphone == "" {
		logger.Warn("[USER] cellphone is empty")
		return errors.New("cellphone is empty")
	}

	// Check if cellphone already exists
	query = "SELECT cellphone FROM `User` WHERE cellphone = ?"
	err = mariadb.DB.QueryRow(query, rr.Username).Scan(&cellphone)
	if err != nil && err.Error() != "sql: no rows in result set" {
		logger.Error("[USER] " + err.Error())
		return err
	} else if cellphone != "" {
		logger.Warn("[USER] cellphone:" + rr.Username + " already exists")
		return errors.New("cellphone already exists")
	}

	// Check whether password is empty
	if rr.Password == "" {
		logger.Warn("[USER] password is empty")
		return errors.New("password is empty")
	}

	// Check pass in email field (Email has @ symbol)
	if rr.Email != "" && !strings.Contains(rr.Email, "@") {
		logger.Warn("[USER] Invalid email address")
		return errors.New("invalid email address")
	}

	// Check pass in phone field (phone number: start with 09 and 10 numbers in total)
	pattern := `^09\d{8}$`
	regex, err := regexp.Compile(pattern)
	if err != nil {
	    logger.Error("[USER] Error compiling regex:" + rr.Cellphone)
	    return err
	}
	if !regex.MatchString(rr.Cellphone) {
	    logger.Warn("[USER] Invalid phone number format")
	    return errors.New("invalid phone number format")
	}

	// Check if password contains uppercase letters, lowercase letters, and digits
	lowercasePattern := `[a-z]`
	uppercasePattern := `[A-Z]`
	digitPattern := `\d`
	lowercaseRegex, err := regexp.Compile(lowercasePattern)
	if err != nil {
	    logger.Error("[USER] Error compiling regex:" + rr.Password)
	    return err
	}
	uppercaseRegex, err := regexp.Compile(uppercasePattern)
	if err != nil {
	    logger.Error("[USER] Error compiling regex:" + rr.Password)
	    return err
	}
	digitRegex, err := regexp.Compile(digitPattern)
	if err != nil {
	    logger.Error("[USER] Error compiling regex:" + rr.Password)
	    return err
	}
	if !lowercaseRegex.MatchString(rr.Password) {
		return errors.New("invalid password format")
	}
	if !uppercaseRegex.MatchString(rr.Password) {
		return errors.New("invalid password format")
	}
	if !digitRegex.MatchString(rr.Password) {
		return errors.New("invalid password format")
	}

	// Check whether passwordConfirm is empty
	if rr.PasswordConfirm == "" {
		logger.Warn("[USER] passwordConfirm is empty")
		return errors.New("passwordConfirm is empty")
	}

	// Check difference between password and passwordConfirm
	if rr.Password != rr.PasswordConfirm {
		logger.Warn("[USER] password and passwordConfirm is different")
		return errors.New("password and passwordConfirm is different")
	}

	// Hash password
	if rr.Password, err = hashPassword(rr.Password); err != nil {
		logger.Error("[USER] " + err.Error())
		return err
	}

	userRole := []string{"buyer", "seller"}
	loginStatus := []string{"signIn", "signOut"}

	// Insert into user database
	query = `INSERT INTO User (
        userId,
        username,
        userPasswd,
        realName,
        cellphone,
        fbAccount,
        email,
        postcode,
        shippingAddr,
        userRole,
        loginStatus,
        nickName
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = mariadb.DB.Exec(
		query,
		uuid.NewString(),
		rr.Username,
		rr.Password,
		rr.RealName,
		rr.Cellphone,
		rr.FbAccount,
		rr.Email,
		rr.Postcode,
		rr.Address,
		userRole[0],
		loginStatus[0],
		rr.NickName,
	)
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
	query = "SELECT userId, userPasswd FROM User WHERE username = ?"
	err = mariadb.DB.QueryRow(query, lr.Username).Scan(&UUID, &password)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			logger.Warn("[USER] Username: " + lr.Username + " not found")
			return "", errors.New("username not found")
		}
		logger.Error("[USER] " + err.Error())
		return "", err
	}

	// Check if password is correct
	if !checkPasswordHash(lr.Password, password) {
		logger.Warn("[USER] Incorrect password for Username: " + lr.Username)
		return "", errors.New("incorrent password")
	}

	logger.Info("[USER] Successfully logged in user with Username: " + lr.Username)

	return UUID, nil
}

func getUserInfo(UUID string) (response.GetUserInfoResponse, error) {
	var err error

    query := `SELECT realName, nickName, cellphone, fbAccount, email, postcode, shippingAddr, username FROM User WHERE userId = ?`

    var getUserInfoResponse response.GetUserInfoResponse
    err = mariadb.DB.QueryRow(query, UUID).Scan(&getUserInfoResponse.RealName, &getUserInfoResponse.NickName, &getUserInfoResponse.Cellphone, &getUserInfoResponse.FbAccount, &getUserInfoResponse.Email, &getUserInfoResponse.Postcode, &getUserInfoResponse.ShippingAddr, &getUserInfoResponse.Username)
    if err != nil {
        logger.Error("[User] " + err.Error())
		return getUserInfoResponse, err
    }

	logger.Info("[User] Successfully return user info: " + UUID)
	
	return getUserInfoResponse, err
}

func updateUserInfo(UUID string, rr updateUserInfoRequest) error {
	var err error
	var username, cellphone string

	// Check whether username is empty
	if rr.Username == "" {
		logger.Warn("[USER] username is empty")
		return errors.New("username is empty")
	}

	// Check if user already exists if username has been modified
	query := "SELECT username FROM `User` WHERE userId = ?"
	err = mariadb.DB.QueryRow(query, UUID).Scan(&username)
	if err != nil && err.Error() != "sql: no rows in result set" {
		logger.Error("[USER] " + err.Error())
		return err
	}
	if username != rr.Username {
		query = "SELECT username FROM `User` WHERE username = ?"
		err = mariadb.DB.QueryRow(query, rr.Username).Scan(&username)
		if err != nil && err.Error() != "sql: no rows in result set" {
			logger.Error("[USER] " + err.Error())
			return err
		} else if username != "" {
			logger.Warn("[USER] username:" + rr.Username + " already exists")
			return errors.New("username already exists")
		}
	}

	
	// Check whether address is empty
	if rr.ShippingAddr == "" {
		logger.Warn("[USER] address is empty")
		return errors.New("address is empty")
	}

	// Check whether cellphone is empty
	if rr.Cellphone == "" {
		logger.Warn("[USER] cellphone is empty")
		return errors.New("cellphone is empty")
	}

	// Check if cellphone already exists if cellphone has been modified
	query = "SELECT cellphone FROM `User` WHERE userId = ?"
	err = mariadb.DB.QueryRow(query, UUID).Scan(&cellphone)
	if err != nil && err.Error() != "sql: no rows in result set" {
		logger.Error("[USER] " + err.Error())
		return err
	}
	if cellphone != rr.Cellphone {
		query = "SELECT cellphone FROM `User` WHERE cellphone = ?"
		err = mariadb.DB.QueryRow(query, rr.Cellphone).Scan(&cellphone)
		if err != nil && err.Error() != "sql: no rows in result set" {
			logger.Error("[USER] " + err.Error())
			return err
		} else if cellphone != "" {
			logger.Warn("[USER] cellphone:" + rr.Cellphone + " already exists")
			return errors.New("cellphone already exists")
		}
	}

	// Check pass in email field (Email has @ symbol)
	if rr.Email != "" && !strings.Contains(rr.Email, "@") {
		logger.Warn("[USER] Invalid email address")
		return errors.New("invalid email address")
	}

	// Check pass in phone field (phone number: start with 09 and 10 numbers in total)
	pattern := `^09\d{8}$`
	regex, err := regexp.Compile(pattern)
	if err != nil {
	    logger.Error("[USER] Error compiling regex:" + rr.Cellphone)
	    return err
	}
	if !regex.MatchString(rr.Cellphone) {
	    logger.Warn("[USER] Invalid phone number format")
	    return errors.New("invalid phone number format")
	}

    query = `UPDATE User SET realName = ?, nickName = ?, cellphone = ?, fbAccount = ?, email = ?, postcode = ?, shippingAddr = ?, username = ? WHERE userId = ?`
    _, err = mariadb.DB.Exec(query, rr.RealName, rr.NickName, rr.Cellphone, rr.FbAccount, rr.Email, rr.Postcode, rr.ShippingAddr, rr.Username, UUID)
    if err != nil {
        logger.Error("[User] " + err.Error())
        return err
    }

	logger.Info("[User] Successfully update user info: " + UUID)
	
	return err
}

func updatePassword(UUID string, rr updatePasswordRequest) error {
	var err error
	var password string

	// Check whether originalPassword is empty
	if rr.OriginPassword == "" {
		logger.Warn("[USER] originalPassword is empty")
		return errors.New("originalPassword is empty")
	}

	// Check whether newPassword is empty
	if rr.NewPassword == "" {
		logger.Warn("[USER] newPassword is empty")
		return errors.New("newPassword is empty")
	}

	// Check whether confirmNewPassword is empty
	if rr.ConfirmNewPassword == "" {
		logger.Warn("[USER] confirmNewPassword is empty")
		return errors.New("confirmNewPassword is empty")
	}

	// Get user password
	query := "SELECT userPasswd FROM User WHERE userId = ?"
	err = mariadb.DB.QueryRow(query, UUID).Scan(&password)
	if err != nil {
		logger.Error("[USER] " + err.Error())
		return err
	}

	// Check if password is correct
	if !checkPasswordHash(rr.OriginPassword, password) {
		logger.Warn("[USER] Incorrect password")
		return errors.New("originalPassword is wrong")
	}

	// Check if password and passwordConfirm is the same
	if rr.NewPassword != rr.ConfirmNewPassword {
		logger.Warn("[USER] newPassword is different from confirmNewPassword")
		return errors.New("newPassword is different from confirmNewPassword")
	}

	// Check if password contains uppercase letters, lowercase letters, and digits
	lowercasePattern := `[a-z]`
	uppercasePattern := `[A-Z]`
	digitPattern := `\d`
	lowercaseRegex, err := regexp.Compile(lowercasePattern)
	if err != nil {
	    logger.Error("[USER] Error compiling regex:" + rr.NewPassword)
	    return err
	}
	uppercaseRegex, err := regexp.Compile(uppercasePattern)
	if err != nil {
	    logger.Error("[USER] Error compiling regex:" + rr.NewPassword)
	    return err
	}
	digitRegex, err := regexp.Compile(digitPattern)
	if err != nil {
	    logger.Error("[USER] Error compiling regex:" + rr.NewPassword)
	    return err
	}
	if !lowercaseRegex.MatchString(rr.NewPassword) {
		return errors.New("invalid newPassword format")
	}
	if !uppercaseRegex.MatchString(rr.NewPassword) {
		return errors.New("invalid newPassword format")
	}
	if !digitRegex.MatchString(rr.NewPassword) {
		return errors.New("invalid newPassword format")
	}

	// Hash password
	if rr.NewPassword, err = hashPassword(rr.NewPassword); err != nil {
		logger.Error("[USER] " + err.Error())
		return err
	}

    query = `UPDATE User SET userPasswd = ? WHERE userId = ?`
    _, err = mariadb.DB.Exec(query, rr.NewPassword, UUID)
    if err != nil {
        logger.Error("[User] " + err.Error())
        return err
    }

	logger.Info("[User] Successfully update user info: " + UUID)
	
	return err
}