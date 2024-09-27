package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var (
	ErrNotFound       = errors.New("models: resource not found")
	ErrUserIDRequired = errors.New("A user id must be valid and required")
	ErrIDInvalid      = errors.New("There was an error: Invalid ID. Please try again, and contact us if the problem persists.")
)

type User struct {
	ID              string     `gorm:"type:uuid;primaryKey"`     // El campo de UUID como clave primaria
	Name            string     `gorm:"type:varchar(255)"`        // Nombre
	LastName        string     `gorm:"type:varchar(255)"`        // Apellido
	Email           string     `gorm:"type:varchar(255);unique"` // Correo electrónico, único
	Password        string     `gorm:"type:varchar(255)"`        // Contraseña (hashed)
	NIF             string     `gorm:"type:varchar(50)"`         // Número de identificación fiscal (NIF)
	PhoneNumber     string     `gorm:"type:varchar(20)"`         // Número de teléfono
	IsAdmin         bool       `gorm:"default:false"`            // Es administrador
	IsActive        bool       `gorm:"default:true"`             // Está activo
	IsSupplier      bool       `gorm:"default:false"`            // Es proveedor
	EmailVerifiedAt *time.Time // Fecha de verificación del correo electrónico
	RememberToken   string     `gorm:"type:varchar(100)"` // Token para recordar el usuario (opcional)
	CreatedAt       time.Time  // Fecha de creación
	UpdatedAt       time.Time  // Fecha de actualización
	Addresses       []Address  `gorm:"foreignKey:UserID"`
	Orders          []Order    `gorm:"foreignKey:UserID"`
}

type UserDB interface {
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)
}

var _ UserDB = &userGorm{}

type userGorm struct {
	db *gorm.DB
}

func newUserGorm(db *gorm.DB) (*userGorm, error) {
	return &userGorm{
		db: db,
	}, nil
}

func (user *User) GetDirections() error {
	return DB.Where("user_id = ?", user.ID).Find(&user.Addresses).Error
}
func (user *User) GetOrders() error {
	return DB.Where("user_id = ?", user.ID).Find(&user.Orders).Error

}

func (user *User) ByID() error {
	return DB.Where("id = ?", user.ID).First(&user).Error
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

func (user *User) Create() error {
	err := DB.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (user *User) Update() error {
	return DB.Save(user).Error
}

func (user *User) Delete() error {
	return DB.Delete(&user).Error
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
