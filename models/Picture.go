package models

// Picture representa el equivalente de PictureBase en PHP
type Picture struct {
	ID        string `gorm:"type:uuid;primaryKey"` // ID de la imagen (UUID)
	ProductID string `gorm:"type:uuid"`            // Relación con el producto
	Image     string `gorm:"type:varchar(255)"`    // Ruta o URL de la imagen
	IsMain    bool   `gorm:"default:false"`        // Indica si es la imagen principal
	IsActive  bool   `gorm:"default:true"`         // Indica si la imagen está activa
}
