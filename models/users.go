package models

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"jgt.solutions/hash"
	"jgt.solutions/rand"
)

var (
	ErrNotMatchLogin     = errors.New("The user or the password is incorrect, please try again or contact with our support")
	ErrSamePasswordReset = errors.New("The new password can not be the same")
	ErrNotFound          = errors.New("models: resource not found")
	ErrUserIDRequired    = errors.New("A user id must be valid and required")
	ErrInvalidEmail      = errors.New("The email provided is not valid")
	ErrEmailNotExist     = errors.New("The email provided don't exist")
	ErrPasswordIncorrect = errors.New("The password provided is incorrect")
	ErrIDInvalid         = errors.New("There was an error: Invalid ID" + " Please try again, and contact us if the problem persists.")

	ErrEmailIsRequired = errors.New("Email address is required")
	ErrEmailIsNotValid = errors.New("The email provided is not valid")
	ErrEmailIsTaken    = errors.New("The email provided is already registered")

	ErrPasswordTooShort = errors.New("Password must be at least 8 characters longer")
	ErrPasswordRequired = errors.New("Password is required")

	ErrRememberTooShort = errors.New("There was an error: Remember token must be at least 32 characters long" + " Please try again, and contact us if the problem persists.")
	ErrRememberRequired = errors.New("There was an error: Remember token is required" + " Please try again, and contact us if the problem persists.")
)

const (
	userPwPPepper = "joan"
	hmacScretKey  = "Sabadell"
)

type UserDB interface {
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error
}

type UserService interface {
	Authenticate(email, password string) (*User, error)

	InitiateReset(userID uint) (string, error)
	CompleteReset(token, newPw string) (*User, error)
	UserDB
}

func NewUserService(gD *gorm.DB) UserService {
	ug, err := newUserGorm(gD)
	if err != nil {
		return nil
	}
	hmac := hash.NewHMAC(hmacScretKey)
	uv := newUserValidator(ug, hmac)
	return &userService{
		UserDB:    uv,
		pwResetDB: newPwResetValidator(&pwResetGorm{db: gD}, hmac),
	}
}

type userService struct {
	UserDB
	pwResetDB pwResetDB
}

type userValidator struct {
	UserDB
	hmac       hash.HMAC
	emailRegex *regexp.Regexp
}

type userValFunc func(*User) error

func runUserValFuncs(user *User, fns ...userValFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

func newUserValidator(udb UserDB, hmac hash.HMAC) *userValidator {
	return &userValidator{
		UserDB:     udb,
		hmac:       hmac,
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
	}
}

func (uv *userValidator) ByEmail(email string) (*User, error) {
	user := User{
		Email: email,
	}
	if err := runUserValFuncs(&user, normalizeEmail, defaultify, hmacRemember); err != nil {
		return nil, err
	}

	return uv.UserDB.ByEmail(user.Email)
}

func bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}

	pwByte := []byte(user.Password + userPwPPepper)

	hashedBytes, err := bcrypt.GenerateFromPassword(pwByte, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}

func passwordMinLength(user *User) error {
	if user.Password == "" {
		return nil
	}
	if len(user.Password) < 8 {
		return ErrPasswordTooShort
	}
	return nil
}

func passwordHashRequired(user *User) error {
	if user.PasswordHash == "" {
		return ErrPasswordRequired
	}
	return nil
}

func passwordRequired(user *User) error {
	if user.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}

func hmacRemember(user *User) error {
	if user.Remember == "" {
		return nil
	}
	hmac := hash.NewHMAC(hmacScretKey)
	user.RememberHash = hmac.Hash(user.Remember)
	return nil
}

func defaultify(user *User) error {
	if user.Remember != "" {
		return nil
	}

	token, err := rand.RememberToken()
	if err != nil {
		return err
	}
	user.Remember = token
	return nil
}

func rememberMinBytes(user *User) error {
	if user.Remember == "" {
		return nil
	}
	n, err := rand.NBytes(user.Remember)
	if err != nil {
		return err

	}
	if n < 32 {
		return ErrRememberTooShort
	}
	return nil
}
func rememberHashRequired(user *User) error {
	if user.RememberHash == "" {
		return ErrPasswordRequired
	}
	return nil
}
func idGreaterThanZero(user *User) error {
	if user.ID <= 0 {
		return ErrIDInvalid
	}
	return nil
}

func normalizeEmail(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)

	return nil
}

func requireEmail(user *User) error {
	if user.Email == "" {
		return ErrEmailIsRequired
	}
	return nil
}

func emailFormat(user *User) error {
	if user.Email == "" {
		return nil
	}
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\]+@[a-z0-9.\-]+\.[a-z]{2,16}$`)
	if !emailRegex.MatchString(user.Email) {
		return ErrEmailIsNotValid
	}
	return nil
}

func (uv *userValidator) Create(user *User) error {

	if err := runUserValFuncs(user,
		passwordRequired,
		passwordMinLength,
		bcryptPassword,
		passwordHashRequired,
		defaultify,
		rememberMinBytes,
		hmacRemember,
		rememberHashRequired,
		normalizeEmail,
		requireEmail,
		emailFormat,
	); err != nil {
		return err
	}

	return uv.UserDB.Create(user)
}

func (uv *userValidator) Update(user *User) error {
	if err := runUserValFuncs(user,
		passwordMinLength,
		bcryptPassword,
		passwordHashRequired,
		rememberMinBytes,
		hmacRemember,
		rememberHashRequired,
		normalizeEmail,
		requireEmail,
		emailFormat,
	); err != nil {
		return err
	}

	return uv.UserDB.Update(user)
}

func (uv *userValidator) Delete(id uint) error {
	var user User
	user.ID = id
	err := runUserValFuncs(&user, idGreaterThanZero)
	if err != nil {
		return err
	}
	return uv.UserDB.Delete(id)
}

func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := User{
		Remember: token,
	}
	if err := runUserValFuncs(&user,
		hmacRemember,
		rememberHashRequired); err != nil {
		return nil, err
	}

	return uv.UserDB.ByRemember(user.RememberHash)
}

func newUserGorm(db *gorm.DB) (*userGorm, error) {
	return &userGorm{
		db: db,
	}, nil
}

var _ UserDB = &userGorm{}

type userGorm struct {
	db *gorm.DB
}

func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id).First(&user)
	err := first(db, &user)
	return &user, err
}

func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User

	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var user User
	err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *userService) InitiateReset(userID uint) (string, error) {
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
	if time.Now().Sub(pwr.CreatedAt) > (2 * time.Hour) {
		return nil, err
	}
	user, err := us.ByID(pwr.UserID)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(newPw+userPwPPepper))
	if err == nil {
		return nil, ErrSamePasswordReset
	}
	user.Password = newPw
	err = us.Update(user)
	if err != nil {
		return nil, err
	}
	us.pwResetDB.Delete(pwr.ID)

	return user, nil

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
		return nil, ErrEmailNotExist
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+userPwPPepper))
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

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	switch err {
	case nil:
		return nil
	case gorm.ErrRecordNotFound:
		return ErrNotFound
	default:
		return err
	}
}

func (ug *userGorm) Create(user *User) error {
	err := ug.db.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (ug *userGorm) Delete(id uint) error {
	user := User{ProtoModel: ProtoModel{ID: id}}
	return ug.db.Delete(&user).Error
}

func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

type User struct {
	ProtoModel
	Name         string    `gorm:"not null"`
	Email        string    `gorm:"not null;unique_index"`
	Password     string    `gorm:"-"`
	PasswordHash string    `gorm:"not null"`
	Remember     string    `gorm:"-"`
	RememberHash string    `gorm:"not null;unique_index"`
	PermLevel    string    `gorm:"not null;default:'User'"`
	Enabled      bool      `gorm:"not null;default:false"`
	Photo        string    `gorm:"default:null"`
	DOB          time.Time `gorm:""`
	Phone        string    `gorm:""`
	Instagram    string    `gorm:""`
	Direction    string    `gorm:""`
}
