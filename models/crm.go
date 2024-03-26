package models

import (
	"github.com/jinzhu/gorm"
	"jgt.solutions/logController"
    "time"
)

type CrmDB interface {
	CountAllSales() (float64, error)
	GetAllOrders() ([]*Order, error)
	SearchOrderByID(id int) (*Order, error)
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

	return &crmService{
		CrmDB: ug,
	}
}

type crmService struct {
	CrmDB
}

// Functions orders

type result struct {
	Total float64
}

func (tg *crmGorm) CountAllSales() (float64, error) {
	var result result
	err := tg.db.Table("orders").Select("sum(total_amount) as total ").Find(&result).Error
	if err != nil {
		logController.ErrorLogger.Println(err)
		return 0, err
	}
	return result.Total, nil
}

func (tg *crmGorm) GetAllOrders() ([]*Order, error) {
	var orders []*Order
	err := tg.db.
		//Preload("Customer").
		//Preload("Products").
		//Preload("Products.Product").
		//Preload("Products.Material").
		//Preload("Products.Printer").
		Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (tg *crmGorm) SearchOrderByID(id int) (*Order, error) {
	var order Order
	err := tg.db.Where("id = ?", id).Find(&order).Error
	return &order, err
}

// Functions Products

// Functions categories

// Functions material


type Brand struct {
    ID          string    `gorm:"column:id"`
    Name        string    `gorm:"column:name"`
    Description string    `gorm:"column:description"`
    IsActive    bool      `gorm:"column:is_active"`
    CreatedAt   time.Time `gorm:"column:created_at"`
    UpdatedAt   time.Time `gorm:"column:updated_at"`
}

type Category struct {
    ID          string    `gorm:"column:id"`
    Name        string    `gorm:"column:name"`
    Slug        string    `gorm:"column:slug"`
    Description string    `gorm:"column:description"`
    IsActive    bool      `gorm:"column:is_active"`
    CreatedAt   time.Time `gorm:"column:created_at"`
    UpdatedAt   time.Time `gorm:"column:updated_at"`
}

type Color struct {
    ID          string    `gorm:"column:id"`
    Name        string    `gorm:"column:name"`
    Description string    `gorm:"column:description"`
    Hex         string    `gorm:"column:hex"`
    IsActive    bool      `gorm:"column:is_active"`
    CreatedAt   time.Time `gorm:"column:created_at"`
    UpdatedAt   time.Time `gorm:"column:updated_at"`
}

type Country struct {
    ID         string `gorm:"column:id"`
    Numeric    string `gorm:"column:numeric"`
    Alfa2      string `gorm:"column:alfa2"`
    Alfa3      string `gorm:"column:alfa3"`
    Name       string `gorm:"column:name"`
    Phonecode  uint16 `gorm:"column:phonecode"`
    IsActive   bool   `gorm:"column:is_active"`
}

type FailedJob struct {
    ID          uint64    `gorm:"column:id"`
    UUID        string    `gorm:"column:uuid"`
    Connection  string    `gorm:"column:connection"`
    Queue       string    `gorm:"column:queue"`
    Payload     string    `gorm:"column:payload"`
    Exception   string    `gorm:"column:exception"`
    FailedAt    time.Time `gorm:"column:failed_at"`
}

type Finish struct {
    ID          string    `gorm:"column:id"`
    Name        string    `gorm:"column:name"`
    Description string    `gorm:"column:description"`
    IsActive    bool      `gorm:"column:is_active"`
    CreatedAt   time.Time `gorm:"column:created_at"`
    UpdatedAt   time.Time `gorm:"column:updated_at"`
}

type Job struct {
    ID           uint64 `gorm:"column:id"`
    Queue        string `gorm:"column:queue"`
    Payload      string `gorm:"column:payload"`
    Attempts     uint8  `gorm:"column:attempts"`
    ReservedAt   uint   `gorm:"column:reserved_at"`
    AvailableAt  uint   `gorm:"column:available_at"`
    CreatedAt    uint   `gorm:"column:created_at"`
}

type Material struct {
    ID       string `gorm:"column:id"`
    Name     string `gorm:"column:name"`
    Description string `gorm:"column:description"`
    CostKg   uint16 `gorm:"column:cost_kg"`
    VATType  string `gorm:"column:vat_type"`
    IsActive bool   `gorm:"column:is_active"`
}

type Migration struct {
    ID        uint   `gorm:"column:id"`
    Migration string `gorm:"column:migration"`
    Batch     int    `gorm:"column:batch"`
}

type PasswordResetToken struct {
    Email      string    `gorm:"column:email"`
    Token      string    `gorm:"column:token"`
    CreatedAt  time.Time `gorm:"column:created_at"`
}

type PersonalAccessToken struct {
    ID             uint64    `gorm:"column:id"`
    TokenableType string    `gorm:"column:tokenable_type"`
    TokenableID   uint64    `gorm:"column:tokenable_id"`
    Name           string    `gorm:"column:name"`
    Token          string    `gorm:"column:token"`
    Abilities     string    `gorm:"column:abilities"`
    LastUsedAt    time.Time `gorm:"column:last_used_at"`
    ExpiresAt     time.Time `gorm:"column:expires_at"`
    CreatedAt     time.Time `gorm:"column:created_at"`
    UpdatedAt     time.Time `gorm:"column:updated_at"`
}

type SetupApp struct {
    ID                 string    `gorm:"column:id"`
    PriceKw            uint8     `gorm:"column:price_kw"`
    PrinterConsumption uint16    `gorm:"column:printer_consumption"`
    Comment            string    `gorm:"column:comment"`
    CreatedAt          time.Time `gorm:"column:created_at"`
    UpdatedAt          time.Time `gorm:"column:updated_at"`
}

type ShippingType struct {
    ID          string    `gorm:"column:id"`
    Name        string    `gorm:"column:name"`
    Title       string    `gorm:"column:title"`
    Description string    `gorm:"column:description"`
    Price       int       `gorm:"column:price"`
    IsActive    bool      `gorm:"column:is_active"`
    IsDefault   bool      `gorm:"column:is_default"`
    CreatedAt   time.Time `gorm:"column:created_at"`
    UpdatedAt   time.Time `gorm:"column:updated_at"`
    DeletedAt   time.Time `gorm:"column:deleted_at"`
}

type Size struct {
    ID          string    `gorm:"column:id"`
    Name        string    `gorm:"column:name"`
    Description string    `gorm:"column:description"`
    IsActive    bool      `gorm:"column:is_active"`
    CreatedAt   time.Time `gorm:"column:created_at"`
    UpdatedAt   time.Time `gorm:"column:updated_at"`
}

type Tag struct {
    ID          string    `gorm:"type:char(36);primaryKey" json:"id"`
    Name        string    `gorm:"type:varchar(30);not null" json:"name"`
    Slug        string    `gorm:"type:varchar(150);not null" json:"slug"`
    Description string    `gorm:"type:varchar(255)" json:"description"`
    IsActive    bool      `gorm:"type:tinyint;not null;default:true" json:"is_active"`
    CreatedAt   time.Time `gorm:"type:timestamp" json:"created_at"`
    UpdatedAt   time.Time `gorm:"type:timestamp" json:"updated_at"`
}

// User model
type User struct {
    ID                string         `gorm:"type:char(36);primaryKey" json:"id"`
    Name              string         `gorm:"type:varchar(255);not null" json:"name"`
    LastName          string         `gorm:"type:varchar(255)" json:"last_name"`
    Email             string         `gorm:"type:varchar(255);not null;unique" json:"email"`
    Password          string         `gorm:"type:varchar(255);not null" json:"password"`
    NIF               string         `gorm:"type:varchar(15)" json:"nif"`
    PhoneNumber       string         `gorm:"type:varchar(15)" json:"phone_number"`
    IsAdmin           bool           `gorm:"type:tinyint;not null;default:false" json:"is_admin"`
    IsActive          bool           `gorm:"type:tinyint;not null;default:true" json:"is_active"`
    IsSupplier        bool           `gorm:"type:tinyint;not null;default:false" json:"is_supplier"`
    RememberToken     string         `gorm:"type:varchar(100)" json:"remember_token"`
    EmailVerifiedAt   *time.Time     `gorm:"type:timestamp" json:"email_verified_at"`
    CreatedAt         time.Time      `gorm:"type:timestamp" json:"created_at"`
    UpdatedAt         time.Time      `gorm:"type:timestamp" json:"updated_at"`
    StripeID          string         `gorm:"type:varchar(255)" json:"stripe_id"`
    PMType            string         `gorm:"type:varchar(255)" json:"pm_type"`
    PMLastFour        string         `gorm:"type:varchar(4)" json:"pm_last_four"`
    TrialEndsAt       *time.Time     `gorm:"type:timestamp" json:"trial_ends_at"`
}

// ColorMaterial model
type ColorMaterial struct {
    ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
    ColorID    string    `gorm:"type:char(36);not null" json:"color_id"`
    MaterialID string    `gorm:"type:char(36);not null" json:"material_id"`
    CreatedAt  time.Time `gorm:"type:timestamp" json:"created_at"`
    UpdatedAt  time.Time `gorm:"type:timestamp" json:"updated_at"`
}

// Community model
type Community struct {
    ID           string         `gorm:"type:char(36);primaryKey" json:"id"`
    Code         string         `gorm:"type:varchar(2);not null" json:"code"`
    CountryAlfa2 string         `gorm:"type:varchar(2);not null" json:"country_alfa2"`
    Name         string         `gorm:"type:varchar(255);not null" json:"name"`
    IsActive     bool           `gorm:"type:tinyint;not null;default:true" json:"is_active"`
    Country      Country        `json:"country"`
    Regions      []Region       `json:"regions"`
    CreatedAt    time.Time      `gorm:"type:timestamp" json:"created_at"`
    UpdatedAt    time.Time      `gorm:"type:timestamp" json:"updated_at"`
}

// Order model
type Order struct {
    ID               string            `gorm:"type:char(36);primaryKey" json:"id"`
    UserID           string            `gorm:"type:char(36);not null" json:"user_id"`
    PaymentType      string            `gorm:"type:enum('Tranferencia','Efectivo','Stripe','Paypal');not null" json:"payment_type"`
    Comment          string            `gorm:"type:varchar(255)" json:"comment"`
    ShipmentID       string            `gorm:"type:varchar(255)" json:"shipment_id"`
    ShipmentName     string            `gorm:"type:varchar(255)" json:"shipment_name"`
    ShipmentTrace    string            `gorm:"type:varchar(255)" json:"shipment_trace"`
    ShipmentInstructions string        `gorm:"type:varchar(255)" json:"shipment_instructions"`
    Shipment         uint16            `gorm:"type:smallint unsigned;not null;default:0" json:"shipment"`
    BaseAmount       uint32            `gorm:"type:mediumint unsigned;not null;default:0" json:"base_amount"`
    TotalAmount      uint32            `gorm:"type:mediumint unsigned;not null" json:"total_amount"`
    CreatedAt        time.Time         `gorm:"type:timestamp" json:"created_at"`
    UpdatedAt        time.Time         `gorm:"type:timestamp" json:"updated_at"`
    User             User              `json:"user"`
    OrderLines       []OrderLine       `json:"order_lines"`
    OrderPaymentStatuses []OrderPaymentStatus `json:"order_payment_statuses"`
    OrderStatuses    []OrderStatus     `json:"order_statuses"`
}

// Product model (continuación)
type Product struct {
    ID               string        `gorm:"type:char(36);primaryKey" json:"id"`
    CategoryID       string        `gorm:"type:char(36);not null" json:"category_id"`
    BrandID          string        `gorm:"type:char(36);not null" json:"brand_id"`
    Slug             string        `gorm:"type:varchar(255);not null" json:"slug"`
    Name             string        `gorm:"type:varchar(255);not null" json:"name"`
    Description      string        `gorm:"type:text" json:"description"`
    Price            uint32        `gorm:"type:mediumint unsigned;not null" json:"price"`
    Stock            uint32        `gorm:"type:mediumint unsigned;not null" json:"stock"`
    IsActive         bool          `gorm:"type:tinyint;not null;default:true" json:"is_active"`
    CreatedAt        time.Time     `gorm:"type:timestamp" json:"created_at"`
    UpdatedAt        time.Time     `gorm:"type:timestamp" json:"updated_at"`
    ProductImages    []ProductImage `json:"product_images"`
    ProductColors    []ProductColor `json:"product_colors"`
    ProductMaterials []ProductMaterial `json:"product_materials"`
    ProductTags      []ProductTag  `json:"product_tags"`
}

// ProductImage model
type ProductImage struct {
    ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
    ProductID  string    `gorm:"type:char(36);not null" json:"product_id"`
    URL        string    `gorm:"type:text;not null" json:"url"`
    CreatedAt  time.Time `gorm:"type:timestamp" json:"created_at"`
    UpdatedAt  time.Time `gorm:"type:timestamp" json:"updated_at"`
}

// ProductColor model
type ProductColor struct {
    ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
    ProductID  string    `gorm:"type:char(36);not null" json:"product_id"`
    ColorID    string    `gorm:"type:char(36);not null" json:"color_id"`
    CreatedAt  time.Time `gorm:"type:timestamp" json:"created_at"`
    UpdatedAt  time.Time `gorm:"type:timestamp" json:"updated_at"`
}

// ProductMaterial model
type ProductMaterial struct {
    ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
    ProductID   string    `gorm:"type:char(36);not null" json:"product_id"`
    MaterialID  string    `gorm:"type:char(36);not null" json:"material_id"`
    CreatedAt   time.Time `gorm:"type:timestamp" json:"created_at"`
    UpdatedAt   time.Time `gorm:"type:timestamp" json:"updated_at"`
}

// ProductTag model
type ProductTag struct {
    ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
    ProductID  string    `gorm:"type:char(36);not null" json:"product_id"`
    TagID      string    `gorm:"type:char(36);not null" json:"tag_id"`
    CreatedAt  time.Time `gorm:"type:timestamp" json:"created_at"`
    UpdatedAt  time.Time `gorm:"type:timestamp" json:"updated_at"`
}
// Region model
type Region struct {
    ID        string    `gorm:"type:char(36);primaryKey" json:"id"`
    Name      string    `gorm:"type:varchar(255);not null" json:"name"`
    CreatedAt time.Time `gorm:"type:timestamp" json:"created_at"`
    UpdatedAt time.Time `gorm:"type:timestamp" json:"updated_at"`
}

// OrderLine model
type OrderLine struct {
    ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
    OrderID      string    `gorm:"type:char(36);not null" json:"order_id"`
    ProductID    string    `gorm:"type:char(36);not null" json:"product_id"`
    Quantity     uint32    `gorm:"type:mediumint unsigned;not null" json:"quantity"`
    Price        uint32    `gorm:"type:mediumint unsigned;not null" json:"price"`
    TotalPrice   uint32    `gorm:"type:mediumint unsigned;not null" json:"total_price"`
    CreatedAt    time.Time `gorm:"type:timestamp" json:"created_at"`
    UpdatedAt    time.Time `gorm:"type:timestamp" json:"updated_at"`
}

// OrderPaymentStatus model
type OrderPaymentStatus struct {
    ID        string    `gorm:"type:char(36);primaryKey" json:"id"`
    Name      string    `gorm:"type:varchar(255);not null" json:"name"`
    CreatedAt time.Time `gorm:"type:timestamp" json:"created_at"`
    UpdatedAt time.Time `gorm:"type:timestamp" json:"updated_at"`
}

// OrderStatus model
type OrderStatus struct {
    ID        string    `gorm:"type:char(36);primaryKey" json:"id"`
    Name      string    `gorm:"type:varchar(255);not null" json:"name"`
    CreatedAt time.Time `gorm:"type:timestamp" json:"created_at"`
    UpdatedAt time.Time `gorm:"type:timestamp" json:"updated_at"`
}
