package models

import (
	"github.com/jinzhu/gorm"
	"jgt.solutions/logController"
)

type CrmDB interface {
	CreateOrder(order *Order) error

	crmMaterials

	crmProducts

	ObtainIncome() (int, error)

	CreatePrinter(printer *Printer) error
	GetAllPrinters() ([]*Printer, error)
	SearchPrinterByID(id int64) (*Printer, error)
	UpdatePrinter(printer *Printer) error

	CountAllSales() (float64, error)
	CountAllSalesExpenses() (float64, error)
	GetAllOrders() ([]*Order, error)
	SearchOrderByID(id int) (*Order, error)

	GetAllCategories() ([]*Category, error)

	GetAllUsers() ([]*User, error)
	CreateCustomer(user *User) error
}

type CrmService interface {
	CrmDB
}
type crmGorm struct {
	db *gorm.DB
}

var DB *gorm.DB

func newCrmGorm(db *gorm.DB) (*crmGorm, error) {
	return &crmGorm{
		db: db,
	}, nil
}
func NewCrmService(gD *gorm.DB) CrmService {
	ug, err := newCrmGorm(gD)
	if err != nil {
		logController.ErrorLogger.Println(err)
		return nil
	}

	tv := newCrmValidator(ug)
	return &crmService{
		CrmDB: tv,
	}
}
func (tg *crmGorm) GetDB() *gorm.DB {
	return tg.db
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

func (tg *crmGorm) ObtainIncome() (int, error) {
	var totalSales int
	err := tg.db.Table("order_lines").Select("sum(price)").Row().Scan(&totalSales)

	if err != nil {
		logController.ErrorLogger.Println(err)
		return 0, err
	}
	return totalSales, nil
}

// Functions orders
func (tg *crmGorm) CreateOrder(order *Order) error {
	if err := tg.db.Create(order).Error; err != nil {
		return err
	}

	// Asocia los productos con el pedido
	//for _, product := range order.Products {
	//	if err := tg.db.Model(order).Association("Products").Append(product).Error; err != nil {
	//		return err
	//	}
	//}

	return nil
}

//func (tg *crmGorm) CreateOrderProductMaterial(orderProductMaterial []OrderProductMaterial) error {
//	return tg.db.Create(orderProductMaterial).Error
//}

type result struct {
	Total float64
}

func (tg *crmGorm) CountAllSales() (float64, error) {
	var result result
	err := tg.db.Table("orders").Select("sum(sale) as total ").Find(&result).Error
	if err != nil {
		logController.ErrorLogger.Println(result)
		return 0, err
	}
	return result.Total, nil
}
func (tg *crmGorm) CountAllSalesExpenses() (float64, error) {
	var result result
	err := tg.db.Table("orders").Select("sum(cost) as total ").Find(&result).Error
	if err != nil {
		logController.ErrorLogger.Println(result)
		return 0, err
	}

	return result.Total, nil
}

func (tg *crmGorm) SearchOrderByID(id int) (*Order, error) {
	var order Order
	err := tg.db.Where("id = ?", id).
		Preload("Customer").
		Preload("Products").
		Preload("Products.Product").
		Preload("Products.Material").
		Preload("Products.Printer").
		Find(&order).Error
	return &order, err
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

// Functions customer

func (tg *crmGorm) GetAllUsers() ([]*User, error) {
	var customers []*User
	err := tg.db.Find(&customers).Error
	if err != nil {
		return nil, err
	}
	return customers, nil
}
func (tg *crmGorm) CreateCustomer(user *User) error {
	return tg.db.Create(user).Error

}

// Functions material
func (tg *crmGorm) CreatePrinter(printer *Printer) error {
	return tg.db.Create(printer).Error
}
func (tg *crmGorm) GetAllPrinters() ([]*Printer, error) {
	var printers []*Printer
	err := tg.db.Find(&printers).Error
	if err != nil {
		return nil, err
	}
	return printers, nil
}
func (tg *crmGorm) UpdatePrinter(printer *Printer) error {
	err := tg.db.Save(printer).Error
	if err != nil {
		return err
	}
	return nil
}
func (tg *crmGorm) SearchPrinterByID(id int64) (*Printer, error) {
	var printer Printer
	err := tg.db.Where("id = ?", id).First(&printer).Error
	return &printer, err
}

type Category struct {
	ProtoModel
	Name        string `gorm:"not null"`
	Description string `gorm:"not null"`
}

type PrinterMaintenance struct {
	ProtoModel
	ExtrusorChange string `gorm:"not null"`
	OilChange      string `gorm:"not null"`
}
type Printer struct {
	ProtoModel
	Name         string               `gorm:"not null"`
	Maintenances []PrinterMaintenance `gorm:"many2many:printers_maintenance;"`
}
