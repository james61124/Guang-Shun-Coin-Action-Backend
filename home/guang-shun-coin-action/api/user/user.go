package user

import (
	"Guang_Shun_Coin_Action/internal/response"
	"Guang_Shun_Coin_Action/pkg/logger"
	"Guang_Shun_Coin_Action/pkg/mariadb"

	// "embed"
	"errors"
	"regexp"

	// "strings"

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
	var query, cellphone string
	var err error

	// Check if user already exists
	query = "SELECT cellphone FROM `User` WHERE cellphone = ?"
	err = mariadb.DB.QueryRow(query, rr.Cellphone).Scan(&cellphone)
	if err != nil && err.Error() != "sql: no rows in result set" {
		logger.Error("[USER] " + err.Error())
		return err
	} else if cellphone != "" {
		logger.Warn("[USER] username:" + rr.Cellphone + " already exists")
		return errors.New("username already exists")
	}

	// Check whether cellphone is empty
	if rr.Cellphone == "" {
		logger.Warn("[USER] cellphone is empty")
		return errors.New("cellphone is empty")
	}

	// Check if cellphone already exists
	query = "SELECT cellphone FROM `User` WHERE cellphone = ?"
	err = mariadb.DB.QueryRow(query, rr.Cellphone).Scan(&cellphone)
	if err != nil && err.Error() != "sql: no rows in result set" {
		logger.Error("[USER] " + err.Error())
		return err
	} else if cellphone != "" {
		logger.Warn("[USER] cellphone:" + rr.Cellphone + " already exists")
		return errors.New("cellphone already exists")
	}

	// Check whether password is empty
	if rr.Password == "" {
		logger.Warn("[USER] password is empty")
		return errors.New("password is empty")
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
        userPasswd,
        realName,
        cellphone,
        userRole,
        loginStatus,
        nickName
    ) VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err = mariadb.DB.Exec(
		query,
		uuid.NewString(),
		rr.Password,
		rr.RealName,
		rr.Cellphone,
		userRole[0],
		loginStatus[0],
		rr.NickName,
	)
	if err != nil {
		logger.Error("[USER] " + err.Error())
		return err
	}

	logger.Info("[USER] Successfully registered user : " + rr.RealName)

	return nil
}

func login(lr loginRequest) (string, error) {
	var query, UUID, password, realName string
	var err error

	// Get user password
	query = "SELECT userId, userPasswd, realName FROM User WHERE cellphone = ?"
	err = mariadb.DB.QueryRow(query, lr.Cellphone).Scan(&UUID, &password, &realName)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			logger.Warn("[USER] Username: " + lr.Cellphone + " not found")
			return "", errors.New("username not found")
		}
		logger.Error("[USER] " + err.Error())
		return "", err
	}

	// Check if password is correct
	if !checkPasswordHash(lr.Password, password) {
		logger.Warn("[USER] Incorrect password for User: " + realName + " " + lr.Cellphone)
		return "", errors.New("incorrent password")
	}

	logger.Info("[USER] Successfully logged in user with Username: " + realName)

	return UUID, nil
}

func getUserInfo(UUID string) (response.GetUserInfoResponse, error) {
	var err error

	query := `SELECT realName, nickName, cellphone FROM User WHERE userId = ?`

	var getUserInfoResponse response.GetUserInfoResponse
	err = mariadb.DB.QueryRow(query, UUID).Scan(&getUserInfoResponse.RealName, &getUserInfoResponse.NickName, &getUserInfoResponse.Cellphone)
	if err != nil {
		logger.Error("[User] " + err.Error())
		return getUserInfoResponse, err
	}

	logger.Info("[User] Successfully return user info: " + UUID)

	return getUserInfoResponse, err
}

func updateUserInfo(UUID string, rr updateUserInfoRequest) error {
	var err error
	var cellphone string

	// Check whether cellphone is empty
	if rr.Cellphone == "" {
		logger.Warn("[USER] cellphone is empty")
		return errors.New("cellphone is empty")
	}

	// Check pass in phone field (phone number: start with 09 and 10 numbers in total) -> correct format
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

	// Check if new cellphone already exists if cellphone has been modified
	query := "SELECT cellphone FROM `User` WHERE userId = ?"
	err = mariadb.DB.QueryRow(query, UUID).Scan(&cellphone)
	if err != nil && err.Error() != "sql: no rows in result set" {
		logger.Error("[USER] " + err.Error())
		return err
	}
	if cellphone != rr.Cellphone { // need to update
		query = "SELECT cellphone FROM `User` WHERE cellphone = ?" // check if the new cellphone exists
		err = mariadb.DB.QueryRow(query, rr.Cellphone).Scan(&cellphone)
		if err != nil && err.Error() != "sql: no rows in result set" {
			logger.Error("[USER] " + err.Error())
			return err
		} else if cellphone != "" { // the new cellphone not exists in database, update information
			query = `UPDATE User SET realName = ?, nickName = ?, cellphone = ? WHERE userId = ?`
			_, err = mariadb.DB.Exec(query, rr.RealName, rr.NickName, rr.Cellphone, UUID)
			if err != nil {
				logger.Error("[User] " + err.Error())
				return err
			}
			logger.Info("[User] Successfully update user info: " + UUID)
		}
	}

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

	logger.Info("[User] Successfully update user password: " + UUID)

	return err
}

func resetPassword(rr resetPasswordRequest) error {
	var err error
	var query string

	// Check whether cellphone is empty
	if rr.Cellphone == "" {
		logger.Warn("[USER] cellphone is empty")
		return errors.New("cellphone is empty")
	}

	// check cellphone format

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

	// Check if password and passwordConfirm is the same
	if rr.NewPassword != rr.ConfirmNewPassword {
		logger.Warn("[USER] newPassword is different from confirmNewPassword")
		return errors.New("newPassword is different from confirmNewPassword")
	}

	// Hash password
	if rr.NewPassword, err = hashPassword(rr.NewPassword); err != nil {
		logger.Error("[USER] " + err.Error())
		return err
	}

	query = `UPDATE User SET userPasswd = ? WHERE cellphone = ?`
	result, e := mariadb.DB.Exec(query, rr.NewPassword, rr.Cellphone)
	if e != nil {
		logger.Error("[User] " + e.Error())
		return e
	}

	rowsAffected, rowCheckError := result.RowsAffected()
	if rowCheckError != nil {
		logger.Error("[User] " + rowCheckError.Error())
		return rowCheckError
	}

	if rowsAffected == 0 {
		logger.Warn("[User] Unable to find the row based on the cellphone")
		return errors.New("unable to find the row based on the cellphone")
	}

	logger.Info("[User] Successfully reset user password: " + rr.Cellphone)

	return err
}
