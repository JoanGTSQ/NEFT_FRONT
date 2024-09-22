package models

import "time"

type crmMaterials interface {
	GetAllMaterials() ([]*Material, error)
	SearchMaterialByID(id int64) (*Material, error)
	UpdateMaterial(material *Material) error
}

// Functions material
func (tg *crmGorm) GetAllMaterials() ([]*Material, error) {
	var materials []*Material
	err := tg.db.Find(&materials).Error
	if err != nil {
		return nil, err
	}
	return materials, nil
}
func (tg *crmGorm) UpdateMaterial(material *Material) error {
	err := tg.db.Save(material).Error
	if err != nil {
		return err
	}
	return nil
}
func (tg *crmGorm) SearchMaterialByID(id int64) (*Material, error) {
	var material Material
	err := tg.db.Where("id = ?", id).First(&material).Error
	return &material, err
}

// Material representa el equivalente de MaterialBase en PHP
type Material struct {
	ID          string    `gorm:"type:uuid;primaryKey"` // HasUuids en PHP es UUID en GORM
	Name        string    `gorm:"type:varchar(255)"`    // Nombre del material
	Description string    `gorm:"type:text"`            // Descripción del material
	CostKg      int       `gorm:"type:int"`             // Costo por kilogramo
	VatID       string    `gorm:"type:uuid"`            // Relación con el tipo de IVA (si hay)
	IsActive    bool      `gorm:"default:true"`         // Indica si está activo o no
	CreatedAt   time.Time `gorm:"autoCreateTime"`       // Timestamp para la creación
}
