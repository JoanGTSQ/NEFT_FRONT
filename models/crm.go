package models

import (
	"github.com/jinzhu/gorm"
	"jgt.solutions/errorController"
)

type CrmDB interface {
	CreateOrder(order *Order) error

	CreateMaterial(material *Material) error
	GetAllMaterials() ([]*Material, error)

	CountAllSales() (float64, error)
	CountAllSalesExpenses() (float64, error)
	GetAllOrders() ([]*Order, error)

	GetAllProducts() ([]*Product, error)
	CreateProduct(product *Product) error

	GetAllCategories() ([]*Category, error)
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
	err := tg.db.Preload("Customer").Preload("Material").Preload("Products").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// Functions Products
func (tg *crmGorm) GetAllProducts() ([]*Product, error) {
	var products []*Product
	err := tg.db.Preload("Category").Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}
func (tg *crmGorm) CreateProduct(product *Product) error {
	return tg.db.Create(product).Error
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

// Functions material
func (tg *crmGorm) CreateMaterial(material *Material) error {
	return tg.db.Create(material).Error
}
func (tg *crmGorm) GetAllMaterials() ([]*Material, error) {
	var materials []*Material
	err := tg.db.Find(&materials).Error
	if err != nil {
		return nil, err
	}
	return materials, nil
}

type Category struct {
	ProtoModel
	Name        string `gorm:"not null"`
	Description string `gorm:"not null"`
}

type Product struct {
	ProtoModel
	Name        string     `gorm:"not null"`
	Picture     string     `gorm:"not null"`
	Price       float64    `gorm:"not null"`
	Description string     `gorm:"not null"`
	Category    []Category `gorm:"many2many:products_category;"`
	Weight      int        `gorm:"not null"`
}
type Configurations struct {
	ProtoModel
	BedTemp      int    `gorm:"not null"`
	ExtrusorTemp int    `gorm:"not null"`
	Speed        int    `gorm:"not null"`
	CloackFan    bool   `gorm:"not null"`
	Adhesion     string `gorm:"not null"`
}

type Material struct {
	ProtoModel
	Name           string `gorm:"not null"`
	Color          string `gorm:"not null"`
	Supplier       string `gorm:"not null"`
	Configurations `gorm:"-"`
	Weight         int `gorm:"not null"`
	Price          float64 `gorm:"not null"`
}

type Customer struct {
	ProtoModel
	Name      string `gorm:"not null"`
	Email     string `gorm:"not null"`
	Direction string `gorm:"not null"`
	Phone     string `gorm:"not null"`
	Origin    string `gorm:"not null"`
}

type Order struct {
	ProtoModel
	MaterialID  int       `gorm:"-" json:"materialid"`
	Material    Material  `gorm:"foreignkey:materialID" json:"material"`
	CustomerID  int       `gorm:"" json:"customerid"`
	Customer    Customer  `gorm:"foreignkey:customerID" json:"customer"`
	Products    []Product `gorm:"many2many:products_whoknow" json:"products"`
	TimeMinutes int       `gorm:"not null"`
	Cost        float64   `gorm:"not null"`
	Sale        float64   `gorm:"not null"`
	Sent        bool      `gorm:"not null"`
	Quality     string    `gorm:"not null"`
}
