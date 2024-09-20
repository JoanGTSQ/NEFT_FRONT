package models

import (
	"errors"
	"regexp"
	"strings"
	"time"
	"jgt.solutions/logController"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"jgt.solutions/hash"
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
	ByID(id string) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	Create(user *User) error
	Update(user *User) error
	Delete(id string) error
}

type UserService interface {
	Authenticate(email, password string) (*User, error)

	InitiateReset(userID string) (string, error)
	CompleteReset(token, newPw string) (*User, error)
	UserDB
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

func newUserValidator(udb UserDB, hmac hash.HMAC) *userValidator {
	return &userValidator{
		UserDB:     udb,
		hmac:       hmac,
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
	}
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

func (ug *userGorm) ByID(id string) (*User, error) {
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

func (ug *userGorm) ByRemember(rememberToken string) (*User, error) {
	var user User
	err := first(ug.db.Where("remember_token = ?", rememberToken), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
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
	if time.Now().Sub(pwr.CreatedAt) > (2 * time.Hour) {
		return nil, err
	}
	user, err := us.ByID(pwr.UserID)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(newPw+userPwPPepper))
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

func (ug *userGorm) Delete(id string) error {
	user := User{ID: id}
	return ug.db.Delete(&user).Error
}

func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

type User struct {
	ID             string         `gorm:"type:uuid;primaryKey"`  // El campo de UUID como clave primaria
	Name           string         `gorm:"type:varchar(255)"`      // Nombre
	LastName       string         `gorm:"type:varchar(255)"`      // Apellido
	Email          string         `gorm:"type:varchar(255);unique"` // Correo electrónico, único
	Password       string         `gorm:"type:varchar(255)"`      // Contraseña (hashed)
	NIF            string         `gorm:"type:varchar(50)"`       // Número de identificación fiscal (NIF)
	PhoneNumber    string         `gorm:"type:varchar(20)"`       // Número de teléfono
	IsAdmin        bool           `gorm:"default:false"`          // Es administrador
	IsActive       bool           `gorm:"default:true"`           // Está activo
	IsSupplier     bool           `gorm:"default:false"`          // Es proveedor
	EmailVerifiedAt *time.Time     // Fecha de verificación del correo electrónico
	RememberToken  string         `gorm:"type:varchar(100)"`      // Token para recordar el usuario (opcional)
	CreatedAt      time.Time      // Fecha de creación
	UpdatedAt      time.Time      // Fecha de actualización
	DeletedAt      time.Time  `gorm:"index"`                  // Soft deletes con GORM
}