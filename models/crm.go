package models

import (
	"github.com/jinzhu/gorm"
	"jgt.solutions/errorController"
)

type CrmDB interface {
	CreateOrder(order *Order) error
	CreateMaterial(material *Material) error

	CountAllSales() (float64, error)
	CountAllSalesExpenses() (float64, error)
	GetAllOrders() ([]*Order, error)
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

// Functions orders
func (tg *crmGorm) CreateOrder(order *Order) error {
	return tg.db.Create(order).Error
}

type result struct {
	Total float64
}

func (tg *crmGorm) CountAllSales() (float64, error) {
	var result result
	err := tg.db.Table("orders").Select("sum(sale) as total ").Find(&result).Error
	if err != nil {
		errorController.ErrorLogger.Println(result)
		return 0, err
	}
	return result.Total, nil
}
func (tg *crmGorm) CountAllSalesExpenses() (float64, error) {
	var result result
	err := tg.db.Table("orders").Select("sum(cost) as total ").Find(&result).Error
	if err != nil {
		errorController.ErrorLogger.Println(result)
		return 0, err
	}

	return result.Total, nil
}
func (tg *crmGorm) GetAllOrders() ([]*Order, error) {
	var orders []*Order
	err := tg.db.Preload("Costumer").Preload("Material").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// Functions material
func (tg *crmGorm) CreateMaterial(material *Material) error {
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
	Price  int `gorm:"not null"`
}

type Costumer struct {
	gorm.Model
	Name      string `gorm:"not null"`
	Email     string `gorm:"not null"`
	Direction string `gorm:"not null"`
	Phone     string `gorm:"not null"`
	Origin    string `gorm:"not null"`
}

type Order struct {
	gorm.Model
	MaterialID  int      `gorm:"" json:"materialid"`
	Material    Material `gorm:"foreignkey:materialID" json:"material"`
	CostumerID  int      `gorm:"" json:"costumerid"`
	Costumer    Costumer `gorm:"foreignkey:costumerID" json:"costumer"`
	TimeMinutes int      `gorm:"not null"`
	Cost        float64  `gorm:"not null"`
	Sale        float64  `gorm:"not null"`
	Sent        bool     `gorm:"not null"`
	Category    string   `gorm:"not null"`
	Quality     string   `gorm:"not null"`
}
