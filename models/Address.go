package models

import (
	"time"
)

type Address struct {
	ID           string    `gorm:"type:uuid;primaryKey"`                    // ID principal (UUID)
	UserID       string    `gorm:"type:uuid"`                               // ID del usuario
	CountryID    string    `gorm:"type:uuid"`                               // ID del país
	RegionID     string    `gorm:"type:uuid"`                               // ID de la región
	Alias        string    `gorm:"type:varchar(255)"`                       // Alias de la dirección
	NIF          string    `gorm:"type:varchar(50)"`                        // Número de identificación fiscal
	VATNumber    string    `gorm:"type:varchar(50)"`                        // Número de IVA
	PhoneNumber  string    `gorm:"type:varchar(20)"`                        // Número de teléfono
	Name         string    `gorm:"type:varchar(255)"`                       // Nombre de la persona
	Company      string    `gorm:"type:varchar(255)"`                       // Nombre de la empresa
	AddressLine1 string    `gorm:"column:address_line_1;type:varchar(255)"` // Nombre de la columna en la base de datos
	AddressLine2 string    `gorm:"column:address_line_2;type:varchar(255)"`
	PostalCode   string    `gorm:"type:varchar(20)"`           // Código postal
	City         string    `gorm:"type:varchar(255)"`          // Ciudad
	IsActive     bool      `gorm:"type:boolean;default:true"`  // Indica si la dirección está activa
	IsDefault    bool      `gorm:"type:boolean;default:false"` // Indica si es la dirección predeterminada
	CreatedAt    time.Time // Fecha de creación
	UpdatedAt    time.Time // Fecha de actualización
}
