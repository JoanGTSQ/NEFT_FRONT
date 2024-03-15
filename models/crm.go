package models

import (
	"github.com/jinzhu/gorm"
	"jgt.solutions/logController"
)

type CrmDB interface {
	CreateOrder(order *Order) error

	CreateMaterial(material *Material) error
	GetAllMaterials() ([]*Material, error)
	SearchMaterialByID(id int64) (*Material, error)
	UpdateMaterial(material *Material) error

	CountAllSales() (float64, error)
	CountAllSalesExpenses() (float64, error)
	GetAllOrders() ([]*Order, error)
	SearchOrderByID(id int) (*Order, error)

	GetAllProducts() ([]*Product, error)
	CreateProduct(product *Product) error
	SearchProductByID(ID int64) (*Product, error)

	GetAllCategories() ([]*Category, error)

	CreateCustomer(material *Customer) error
	GetAllCustomers() ([]*Customer, error)
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
		logController.ErrorLogger.Println(err)
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
	if err := tg.db.Create(order).Error; err != nil {
		return err
	}

	// Asocia los productos con el pedido
	for _, product := range order.Products {
		if err := tg.db.Model(order).Association("Products").Append(product).Error; err != nil {
			return err
		}
	}

	return nil
}

func (tg *crmGorm) CreateOrderProductMaterial(orderProductMaterial []OrderProductMaterial) error {
	return tg.db.Create(orderProductMaterial).Error
}

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
func (tg *crmGorm) GetAllOrders() ([]*Order, error) {
	var orders []*Order
	err := tg.db.Preload("Customer").Preload("Products").Preload("Products.Product").Preload("Products.Material").Preload("Printer").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (tg *crmGorm) SearchOrderByID(id int) (*Order, error) {
	var order Order
	err := tg.db.Where("id = ?", id).Preload("Customer").Preload("Products").Preload("Products.Product").Preload("Products.Material").Preload("Printer").Find(&order).Error
	return &order, err
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
func (tg *crmGorm) SearchProductByID(id int64) (*Product, error) {
	var product Product
	err := tg.db.Where("id = ?", id).First(&product).Error
	return &product, err
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

// Functions customer
func (tg *crmGorm) CreateCustomer(material *Customer) error {
	return tg.db.Create(material).Error
}
func (tg *crmGorm) GetAllCustomers() ([]*Customer, error) {
	var customers []*Customer
	err := tg.db.Find(&customers).Error
	if err != nil {
		return nil, err
	}
	return customers, nil
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
	Stl         string     `gorm:"not null"`
	Price       float64    `gorm:"not null"`
	Description string     `gorm:"not null"`
	Category    []Category `gorm:"many2many:products_category;"`
	Weight      float64    `gorm:"not null"`
	Quality     string     `gorm:"not null"`
	TimeMinutes int        `gorm:"not null"`
}
type PrinterMaintenance struct {
	ProtoModel
	ExtrusorChange string `gorm:"not null"`
	OilChange      string `gorm:"not null"`
}
type Printer struct {
	ProtoModel
	Name         int                  `gorm:"not null"`
	Maintenances []PrinterMaintenance `gorm:"many2many:printers_maintenance;"`
}

type Material struct {
	ProtoModel
	Name     string  `gorm:"not null"`
	Color    string  `gorm:"not null"`
	Supplier string  `gorm:"not null"`
	Weight   float64 `gorm:"not null"`
	Price    float64 `gorm:"not null"`
}

type Customer struct {
	ProtoModel
	Name      string `gorm:"not null"`
	Email     string `gorm:"not null"`
	Direction string `gorm:"not null"`
	Phone     string `gorm:"not null"`
	Origin    string `gorm:"not null"`
}

type OrderProductMaterial struct {
	ProtoModel
	OrderID    int      `gorm:"" json:"orderid"`
	Order      Order    `gorm:"foreignkey:orderID"`
	ProductID  int      `gorm:""`
	Product    Product  `gorm:"foreignkey:productID"`
	MaterialID int      `gorm:""`
	Material   Material `gorm:"foreignkey:materialID"`
	PrinterID  int      `gorm:""`
	Printer    Printer  `gorm:"foreignkey:printerID"`
	Quality    string   `gorm:"not null"`
}
type Order struct {
	ProtoModel
	CustomerID  int                     `gorm:"" json:"customerid"`
	Customer    Customer                `gorm:"foreignkey:customerID"`
	Products    []*OrderProductMaterial `json:"products"` // Elimina la opción `gorm:"many2many:order_product_materials"`
	TimeMinutes int                     `gorm:"not null"`
	Cost        float64                 `gorm:"not null"`
	Sale        float64                 `gorm:"not null"`
	Sent        bool                    `gorm:"not null"`
}
