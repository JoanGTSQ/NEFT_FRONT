package models

import (
	"github.com/jinzhu/gorm"
	"jgt.solutions/errorController"
)

type CrmDB interface {
	Create(material *Material) error
}

type CrmService interface {
	CrmDB
}
type crmGorm struct {
	db *gorm.DB
}

func newCrmGorm(db *gorm.DB) (*crmGorm, error) {
	return &crmGorm{
		db: db,
	}, nil
}
func NewCrmService(gD *gorm.DB) CrmService {
	ug, err := newCrmGorm(gD)
	if err != nil {
		errorController.ErrorLogger.Println(err)
		return nil
	}

	tv := newCrmValidator(ug)
	return &crmService{
		CrmDB: tv,
	}
}

type crmService struct {
	CrmDB
}

func newCrmValidator(tb CrmDB) *crmValidator {
	return &crmValidator{
		CrmDB: tb,
	}
}

type crmValidator struct {
	CrmDB
}

func (tg *crmGorm) Create(material *Material) error {
	return tg.db.Create(material).Error
}

type Configurations struct {
	BedTemp      int    `gorm:"not null"`
	ExtrusorTemp int    `gorm:"not null"`
	Speed        int    `gorm:"not null"`
	CloackFan    bool   `gorm:"not null"`
	Adhesion     string `gorm:"not null"`
}

type Material struct {
	gorm.Model
	Name  string `gorm:"not null"`
	Color string `gorm:"not null"`
	Configurations
	Weight int `gorm:"not null"`
}
