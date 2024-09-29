package models

import "jgt.solutions/logController"

type crmCategories interface {
	GetAllCategories() ([]*Category, error)
}

// Functions categories
func (tg *crmGorm) GetAllCategories() ([]*Category, error) {
	var categories []*Category
	err := tg.db.Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (c *Category) GetSubcategories() {
	err := DB.Table("subcategories").Where("category_id = ?", c.ID).Scan(&c.Subcategories)
	if err != nil {
		logController.ErrorLogger.Println(err)
	}
}

type Category struct {
	ID            string        `gorm:"type:uuid;primaryKey"` // UUID como ID
	Name          string        `gorm:"type:varchar(255)"`    // Nombre de la categoría
	Slug          string        `gorm:"type:varchar(255)"`    // Slug
	Description   string        `gorm:"type:text"`            // Descripción
	IsActive      bool          `gorm:"type:boolean"`         // Estado de la categoría (activa o no)
	Subcategories []SubCategory `gorm:"-"`
}

type SubCategory struct {
	ID          string   `gorm:"type:uuid;primaryKey"` // UUID como ID
	CategoryID  string   `gorm:"type:varchar(255)"`    // ID de la categoría principal
	Name        string   `gorm:"type:varchar(255)"`    // Nombre de la subcategoría
	Slug        string   `gorm:"type:varchar(255)"`    // Slug
	Description string   `gorm:"type:text"`            // Descripción
	IsActive    bool     `gorm:"type:boolean"`         // Estado de la subcategoría (activa o no)
	Category    Category `gorm:"foreignKey:CategoryID"`
}
