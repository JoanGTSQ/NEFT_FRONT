package models

import "time"

func (tg *crmGorm) GetAllOrders() ([]*Order, error) {
	var orders []*Order
	err := tg.db.Preload("User").
		Preload("OrderLines.Attribute.Product").
		Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

type Order struct {
	ID          string      `gorm:"type:uuid;primaryKey"` // ID de la orden
	UserID      string      `gorm:"type:varchar(255)"`    // ID del usuario (UUID)
	Comment     string      `gorm:"type:text"`            // Comentario
	BaseAmount  int         `gorm:"type:bigint"`          // Monto base
	TotalAmount int         `gorm:"type:bigint"`          // Monto total
	IsCompleted bool        `gorm:"default:false"`        // Indica si está completada
	CreatedAt   time.Time   `gorm:"autoCreateTime"`       // Fecha de creación
	UpdatedAt   time.Time   `gorm:"autoUpdateTime"`       // Fecha de actualización
	DeletedAt   *time.Time  `gorm:"index"`                // Fecha de eliminación (soft delete)
	OrderLines  []OrderLine `gorm:"foreignKey:OrderID"`
	User        User        `gorm:"foreignKey:UserID"`
	OrderStatus OrderStatus `gorm:"foreignKey:OrderID"`
}
type OrderLine struct {
	ID          string    `gorm:"type:uuid;primaryKey"`   // ID de la orden
	OrderID     string    `gorm:"type:varchar(255)"`      // ID de la orden (UUID)
	AttributeID string    `gorm:"type:varchar(255)"`      // ID del atributo (UUID)
	ShipmentID  string    `gorm:"type:varchar(255)"`      // ID del envío (UUID)
	Description string    `gorm:"type:text"`              // Descripción
	Quantity    int       `gorm:"type:int"`               // Cantidad
	BasePrice   int       `gorm:"type:bigint"`            // Precio base
	Discount    int       `gorm:"type:bigint"`            // Descuento
	Vat         float64   `gorm:"type:decimal(10,2)"`     // VAT (Impuesto sobre el valor añadido)
	Price       int       `gorm:"type:bigint"`            // Precio final
	Attribute   Attribute `gorm:"foreignKey:AttributeID"` // Relación con Attribute
}
type OrderAddress struct {
	ID           string `gorm:"type:uuid;primaryKey"` // ID de la orden
	OrderID      string `gorm:"type:varchar(255)"`    // ID de la orden (UUID)
	Type         string `gorm:"type:varchar(50)"`     // Tipo de dirección
	Name         string `gorm:"type:varchar(255)"`    // Nombre
	AddressLine1 string `gorm:"type:varchar(255)"`    // Línea de dirección 1
	AddressLine2 string `gorm:"type:varchar(255)"`    // Línea de dirección 2
	PostalCode   string `gorm:"type:varchar(20)"`     // Código postal
	City         string `gorm:"type:varchar(100)"`    // Ciudad
	Region       string `gorm:"type:varchar(100)"`    // Región
	Country      string `gorm:"type:varchar(100)"`    // País
	Nif          string `gorm:"type:varchar(50)"`     // NIF
	PhoneNumber  string `gorm:"type:varchar(20)"`     // Número de teléfono
	Instructions string `gorm:"type:text"`            // Instrucciones
}
type OrderPaymentStatus struct {
	ID             string `gorm:"type:uuid;primaryKey"` // ID de la orden
	OrderID        string `gorm:"type:varchar(255)"`    // ID de la orden (UUID)
	PaymentStatus  string `gorm:"type:varchar(50)"`     // Estado del pago
	PaymentType    string `gorm:"type:varchar(50)"`     // Tipo de pago
	Comment        string `gorm:"type:text"`            // Comentario
	Amount         int64  `gorm:"type:bigint"`          // Monto total
	AmountReceived int64  `gorm:"type:bigint"`          // Monto recibido
	ChargeID       string `gorm:"type:varchar(255)"`    // ID del cargo
	Pi             string `gorm:"type:varchar(255)"`    // ID de la transacción
	Pm             string `gorm:"type:varchar(255)"`    // Método de pago
	PmCardType     string `gorm:"type:varchar(50)"`     // Tipo de tarjeta
	ReceiptURL     string `gorm:"type:varchar(255)"`    // URL del recibo
	ErrorID        string `gorm:"type:varchar(255)"`    // ID del error
	ErrorCode      string `gorm:"type:varchar(50)"`     // Código del error
	ErrorMessage   string `gorm:"type:text"`            // Mensaje del error
}
type OrderStatus struct {
	ID      string `gorm:"type:uuid;primaryKey"` // ID de la orden
	OrderID string `gorm:"type:varchar(255)"`    // ID de la orden (UUID)
	Status  string `gorm:"type:varchar(50)"`     // Estado de la orden
}
