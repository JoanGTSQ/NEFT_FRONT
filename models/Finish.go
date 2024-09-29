package models

import "time"

type Finish struct {
	ID          string    `gorm:"type:uuid;primaryKey"` // ID del acabado (UUID)
	Name        string    `gorm:"type:varchar(255)"`    // Nombre del acabado
	Description string    `gorm:"type:text"`            // Descripción del acabado
	IsActive    bool      `gorm:"default:true"`         // Indica si está activo
	CreatedAt   time.Time // Fecha de creación
	UpdatedAt   time.Time // Fecha de actualización
}
