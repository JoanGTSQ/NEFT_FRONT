package models

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"jgt.solutions/hash"
	"jgt.solutions/logController"
)

var (
	ErrNotMatchLogin     = errors.New("The user or the password is incorrect, please try again or contact with our support")
	ErrSamePasswordReset = errors.New("The new password can not be the same")
	ErrInvalidEmail      = errors.New("The email provided is not valid")
	ErrEmailNotExist     = errors.New("The email provided don't exist")
	ErrPasswordIncorrect = errors.New("The password provided is incorrect")

	ErrEmailIsRequired  = errors.New("Email address is required")
	ErrEmailIsNotValid  = errors.New("The email provided is not valid")
	ErrEmailIsTaken     = errors.New("The email provided is already registered")
	ErrPasswordTooShort = errors.New("Password must be at least 8 characters long")
	ErrPasswordRequired = errors.New("Password is required")
	ErrRememberTooShort = errors.New("Remember token must be at least 32 characters long")
	ErrRememberRequired = errors.New("Remember token is required")
)

const (
	userPwPPepper = "joan"
	hmacScretKey  = "Sabadell"
)

type UserService interface {
	Authenticate(email, password string) (*User, error)
	InitiateReset(userID string) (string, error)
	CompleteReset(token, newPw string) (*User, error)
	UserDB
}

type userService struct {
	UserDB
	pwResetDB pwResetDB
}

func NewUserService(gD *gorm.DB) UserService {
	ug, err := newUserGorm(gD)
	if err != nil {
		return nil
	}
	hmac := hash.NewHMAC(hmacScretKey)
	//uv := newUserValidator(ug, hmac)
	return &userService{
		UserDB:    ug,
		pwResetDB: newPwResetValidator(&pwResetGorm{db: gD}, hmac),
	}
}

func (us *userService) InitiateReset(userID string) (string, error) {
	pwr := pwReset{
		UserID: userID,
	}
	if err := us.pwResetDB.Create(&pwr); err != nil {
		return "", err
	}
	return pwr.TokenHash, nil
}

func (us *userService) CompleteReset(token, newPw string) (*User, error) {
	pwr, err := us.pwResetDB.ByToken(token)
	if err != nil {
		return nil, err
	}
	if time.Since(pwr.CreatedAt) > (2 * time.Hour) {
		return nil, err
	}
	user := User{
		ID: pwr.UserID,
	}
	err = user.ByID()
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(newPw+userPwPPepper))
	if err == nil {
		return nil, ErrSamePasswordReset
	}
	user.Password = newPw
	err = user.Update()
	if err != nil {
		return nil, err
	}
	us.pwResetDB.Delete(pwr.ID)

	return &user, nil
}

func (us *userService) Authenticate(email, password string) (*User, error) {
	if email == "" {
		return nil, ErrEmailIsRequired
	}
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\]+@[a-z0-9.\-]+\.[a-z]{2,16}$`)
	if !emailRegex.MatchString(email) {
		return nil, ErrEmailIsNotValid
	}
	email = strings.ToLower(email)
	email = strings.TrimSpace(email)

	foundUser, err := us.ByEmail(email)
	if err != nil {
		logController.DebugLogger.Println("user not found")
		return nil, ErrEmailNotExist
	}
	logController.DebugLogger.Println("user found", foundUser.Password)
	logController.DebugLogger.Println("Our Password hash", []byte(password))
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(password))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrPasswordIncorrect
		default:
			return nil, err
		}
	}

	return foundUser, nil
}
