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

func validateCellphone(cellphone string) error {
	if cellphone == "" {
		logger.Warn("[USER] cellphone is empty")
		return errors.New("cellphone is empty")
	}

	pattern := `^09\d{8}$`
	regex, err := regexp.Compile(pattern)
	if err != nil {
		logger.Error("[USER] Error compiling regex: " + cellphone)
		return err
	}
	if !regex.MatchString(cellphone) {
		logger.Warn("[USER] Invalid phone number format")
		return errors.New("invalid phone number format")
	}

	return nil
}

func validatePassword(password string) error {
	if password == "" {
		logger.Warn("[USER] password is empty")
		return errors.New("password is empty")
	}

	lowercasePattern := `[a-z]`
	uppercasePattern := `[A-Z]`
	digitPattern := `\d`
	lowercaseRegex, err := regexp.Compile(lowercasePattern)
	if err != nil {
		logger.Error("[USER] Error compiling regex: " + password)
		return err
	}
	uppercaseRegex, err := regexp.Compile(uppercasePattern)
	if err != nil {
		logger.Error("[USER] Error compiling regex: " + password)
		return err
	}
	digitRegex, err := regexp.Compile(digitPattern)
	if err != nil {
		logger.Error("[USER] Error compiling regex: " + password)
		return err
	}
	if !lowercaseRegex.MatchString(password) || !uppercaseRegex.MatchString(password) || !digitRegex.MatchString(password) {
		return errors.New("invalid password format")
	}
	
	return nil
}

func register(rr registerRequest) error {
	var query, cellphone string
	var err error


	err = validateCellphone(rr.Cellphone)
	if err != nil {
		return err
	}
	err = validatePassword(rr.Password)
	if err != nil {
		return err
	}

	if rr.PasswordConfirm == "" {
		logger.Warn("[USER] confirmed password is empty")
		return errors.New("confirmed password is empty")
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

	// Check difference between password and passwordConfirm
	if rr.Password != rr.PasswordConfirm {
		logger.Warn("[USER] password is different from PasswordConfirm")
		return errors.New("password is different from PasswordConfirm")
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

	err = validateCellphone(lr.Cellphone)
	if err != nil {
		return "", err
	}
	err = validatePassword(lr.Password)
	if err != nil {
		return "", err
	}

	// Get user password
	query = "SELECT userId, userPasswd, realName FROM User WHERE cellphone = ?"
	err = mariadb.DB.QueryRow(query, lr.Cellphone).Scan(&UUID, &password, &realName)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			logger.Warn("[USER] User: " + lr.Cellphone + " not found")
			return "", errors.New("can't find the user with cellphone")
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

	return UUID, err
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
	var duplicatedCount int
	var originalCellphone, originalNickName, originalRealName string

	err = validateCellphone(rr.Cellphone)
	if err != nil {
		return err
	}

	query := "SELECT cellphone, realName, nickName FROM `User` WHERE userId = ?"
	err = mariadb.DB.QueryRow(query, UUID).Scan(&originalCellphone, &originalRealName, &originalNickName)
	if err != nil && err.Error() != "sql: no rows in result set" {
		logger.Error("[USER] " + err.Error())
		return err
	}

	// check if All fields remain unchanged
	if rr.RealName == originalRealName && rr.NickName == originalNickName && rr.Cellphone == originalCellphone {
		logger.Info("[User] user:" + rr.Cellphone + " didn't update any info")
		return err
	}

	// Check if new cellphone already exists in the database
	query = "SELECT COUNT(*) FROM `User` WHERE cellphone = ?"
	err = mariadb.DB.QueryRow(query, rr.Cellphone).Scan(&duplicatedCount)
	if err != nil {
		logger.Error("[USER] " + err.Error())
		return err
	}
	if duplicatedCount > 0 && rr.Cellphone != originalCellphone {
		logger.Error("[USER] the new cellphone already exists")
		return errors.New("the new cellphone already exists")
	} else if duplicatedCount == 0 {
		query = `UPDATE User SET realName = ?, nickName = ?, cellphone = ? WHERE userId = ?`
		_, err = mariadb.DB.Exec(query, rr.RealName, rr.NickName, rr.Cellphone, UUID)
		if err != nil {
			logger.Error("[User] " + err.Error())
			return err
		}
		logger.Info("[User] Successfully update user info: " + rr.Cellphone)
	}


	return err
}

func updatePassword(UUID string, rr updatePasswordRequest) error {
	var err error
	var cellphone, userPasswd string

	if rr.OriginPassword == "" {
		logger.Warn("[USER] original password is empty")
		return errors.New("original password is empty")
	}
	
	if rr.NewPassword == "" {
		logger.Warn("[USER] new password is empty")
		return errors.New("new password is empty")
	}

	if rr.ConfirmNewPassword == "" {
		logger.Warn("[USER] confirmed password is empty")
		return errors.New("confirmed password is empty")
	}
	
	// Check if new password has correct format
	err = validatePassword(rr.NewPassword)
	if err != nil {
		return err
	}

	// Check if password and passwordConfirm is the same
	if rr.NewPassword != rr.ConfirmNewPassword {
		logger.Warn("[USER] password is different from PasswordConfirm")
		return errors.New("password is different from PasswordConfirm")
	}

	// Get user password from userId
	query := "SELECT cellphone, userPasswd FROM User WHERE userId = ?"
	err = mariadb.DB.QueryRow(query, UUID).Scan(&cellphone, &userPasswd)
	if err != nil {
		logger.Error("[USER] " + err.Error())
		return err
	}

	// Check if original password is correct
	if !checkPasswordHash(rr.OriginPassword, userPasswd) {
		logger.Warn("[USER] Incorrect password")
		return errors.New("originalPassword is wrong")
	}

	//  check if the password remain unchanged
	if rr.OriginPassword == rr.NewPassword {
		logger.Info("[User] user:" + cellphone + " remain unchanged")
		return err
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

	logger.Info("[User] Successfully update user password: " + cellphone)

	return err
}

func resetPassword(rr resetPasswordRequest) error {
	var err error
	var query string

	// Check pass in phone field (phone number: start with 09 and 10 numbers in total)
	err = validateCellphone(rr.Cellphone)
	if err != nil {
		return err
	}

	// Check if new password has correct format
	err = validatePassword(rr.NewPassword)
	if err != nil {
		return err
	}

	if rr.ConfirmNewPassword == "" {
		logger.Warn("[USER] confirmed password is empty")
		return errors.New("confirmed password is empty")
	}

	// Check if password and passwordConfirm is the same
	if rr.NewPassword != rr.ConfirmNewPassword {
		logger.Warn("[USER] password is different from PasswordConfirm")
		return errors.New("password is different from PasswordConfirm")
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
		logger.Warn("[User] Unable to find the user based on the cellphone")
		return errors.New("unable to find the user based on the cellphone")
	}


	logger.Info("[User] Successfully reset user password: " + rr.Cellphone)

	return err
}
