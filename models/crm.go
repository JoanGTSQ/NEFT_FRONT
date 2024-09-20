package models

import (
	"github.com/jinzhu/gorm"
	"jgt.solutions/logController"
	"time"
)

type CrmDB interface {
	CreateOrder(order *Order) error

	CreateMaterial(material *Material) error
	GetAllMaterials() ([]*Material, error)
	SearchMaterialByID(id int64) (*Material, error)
	UpdateMaterial(material *Material) error

	CreatePrinter(printer *Printer) error
	GetAllPrinters() ([]*Printer, error)
	SearchPrinterByID(id int64) (*Printer, error)
	UpdatePrinter(printer *Printer) error

	CountAllSales() (float64, error)
	CountAllSalesExpenses() (float64, error)
	GetAllOrders() ([]*Order, error)
	SearchOrderByID(id int) (*Order, error)

	GetAllProducts() ([]*Product, error)
	CreateProduct(product *Product) error
	SearchProductByID(ID int64) (*Product, error)

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
func (tg *crmGorm) GetAllOrders() ([]*Order, error) {
	var orders []*Order
	err := tg.db.Preload("OrderLines").
		Preload("OrderLines.Attribute").
		Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
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
	Name         string               `gorm:"not null"`
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

type Order struct {
    ID             string         `gorm:"type:uuid;primaryKey"`              // ID de la orden
    UserID      string    `gorm:"type:varchar(255)"`       // ID del usuario (UUID)
    Comment     string    `gorm:"type:text"`               // Comentario
    BaseAmount  int64     `gorm:"type:bigint"`             // Monto base
    TotalAmount int64     `gorm:"type:bigint"`             // Monto total
    IsCompleted bool      `gorm:"default:false"`           // Indica si está completada
    CreatedAt   time.Time `gorm:"autoCreateTime"`          // Fecha de creación
    UpdatedAt   time.Time `gorm:"autoUpdateTime"`          // Fecha de actualización
    DeletedAt   *time.Time `gorm:"index"`                   // Fecha de eliminación (soft delete)
	OrderLines []OrderLine `gorm:"foreignKey:OrderID"`
}
type OrderLine struct {
    ID             string         `gorm:"type:uuid;primaryKey"`              // ID de la orden
    OrderID      string `gorm:"type:varchar(255)"`       // ID de la orden (UUID)
    AttributeID  string `gorm:"type:varchar(255)"`       // ID del atributo (UUID)
    ShipmentID   string `gorm:"type:varchar(255)"`       // ID del envío (UUID)
    Description  string `gorm:"type:text"`               // Descripción
    Quantity     int    `gorm:"type:int"`                // Cantidad
    BasePrice    int64  `gorm:"type:bigint"`             // Precio base
    Discount      int64  `gorm:"type:bigint"`             // Descuento
    Vat          float64 `gorm:"type:decimal(10,2)"`      // VAT (Impuesto sobre el valor añadido)
    Price        int64  `gorm:"type:bigint"`             // Precio final
	Attribute    Attribute `gorm:"foreignKey:AttributeID"` // Relación con Attribute
}
type OrderAddress struct {
    ID             string         `gorm:"type:uuid;primaryKey"`              // ID de la orden
    OrderID        string `gorm:"type:varchar(255)"`       // ID de la orden (UUID)
    Type           string `gorm:"type:varchar(50)"`        // Tipo de dirección
    Name           string `gorm:"type:varchar(255)"`       // Nombre
    AddressLine1   string `gorm:"type:varchar(255)"`       // Línea de dirección 1
    AddressLine2   string `gorm:"type:varchar(255)"`       // Línea de dirección 2
    PostalCode     string `gorm:"type:varchar(20)"`        // Código postal
    City           string `gorm:"type:varchar(100)"`       // Ciudad
    Region         string `gorm:"type:varchar(100)"`       // Región
    Country        string `gorm:"type:varchar(100)"`       // País
    Nif            string `gorm:"type:varchar(50)"`        // NIF
    PhoneNumber    string `gorm:"type:varchar(20)"`        // Número de teléfono
    Instructions   string `gorm:"type:text"`               // Instrucciones
}
type OrderPaymentStatus struct {
    ID             string         `gorm:"type:uuid;primaryKey"`              // ID de la orden
    OrderID           string `gorm:"type:varchar(255)"`       // ID de la orden (UUID)
    PaymentStatus     string `gorm:"type:varchar(50)"`        // Estado del pago
    PaymentType       string `gorm:"type:varchar(50)"`        // Tipo de pago
    Comment           string `gorm:"type:text"`               // Comentario
    Amount            int64  `gorm:"type:bigint"`             // Monto total
    AmountReceived    int64  `gorm:"type:bigint"`             // Monto recibido
    ChargeID          string `gorm:"type:varchar(255)"`       // ID del cargo
    Pi                string `gorm:"type:varchar(255)"`       // ID de la transacción
    Pm                string `gorm:"type:varchar(255)"`       // Método de pago
    PmCardType        string `gorm:"type:varchar(50)"`        // Tipo de tarjeta
    ReceiptURL        string `gorm:"type:varchar(255)"`       // URL del recibo
    ErrorID           string `gorm:"type:varchar(255)"`       // ID del error
    ErrorCode         string `gorm:"type:varchar(50)"`        // Código del error
    ErrorMessage      string `gorm:"type:text"`               // Mensaje del error
}
type OrderStatus struct {
    ID             string         `gorm:"type:uuid;primaryKey"`              // ID de la orden
    OrderID string `gorm:"type:varchar(255)"`       // ID de la orden (UUID)
    Status  string `gorm:"type:varchar(50)"`        // Estado de la orden
}
type Attribute struct {
    ID             string         `gorm:"type:uuid;primaryKey"`              // ID de la orden
    ProductID   string `gorm:"type:varchar(255)"`
    FinishID    string `gorm:"type:varchar(255)"`
    MaterialID  string `gorm:"type:varchar(255)"`
    PictureID   string `gorm:"type:varchar(255)"`
    Price       int    `gorm:"type:int"`
    CostPrice   int    `gorm:"type:int"`
    OfferPrice  int    `gorm:"type:int"`
    Minutes     int    `gorm:"type:int"`
    InOffer     bool   `gorm:"default:false"`
    IsActive    bool   `gorm:"default:true"`
}