package models

type crmProducts interface {
	GetAllProducts() ([]*Product, error)
	CreateProduct(product *Product) error
	SearchProductByID(ID int64) (*Product, error)
	GetTopProducts() ([]*Product, error)
}

// Functions Products
func (tg *crmGorm) GetAllProducts() ([]*Product, error) {
	var products []*Product
	err := tg.db.Preload("Pictures").Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}
func (tg *crmGorm) GetTopProducts() ([]*Product, error) {
	var products []*Product

	// Realizamos un JOIN entre order_lines y products para sumar la cantidad total vendida por producto
	result := tg.db.Table("order_lines").
		Select("products.id, products.name, SUM(order_lines.quantity) as total_sold").
		Joins("JOIN attributes ON attributes.id = order_lines.attribute_id").
		Joins("JOIN products ON products.id = attributes.product_id").
		Group("products.id").
		Order("total_sold DESC").
		Limit(10).
		Scan(&products)

	if result.Error != nil {
		return nil, result.Error
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

func (product *Product) GetMainPicture() string {
	var picture Picture

	if err := DB.Where("product_id = ? AND is_main = ?", product.ID, true).First(&picture).Error; err != nil {
		return "default.jpg" // Devuelve una imagen por defecto si hay un error
	}
	return picture.Image
}

type Product struct {
	ID               string `gorm:"type:uuid;primaryKey"`
	CategoryID       string `gorm:"type:varchar(255)"`
	SubcategoryID    string `gorm:"type:varchar(255)"`
	BrandID          string `gorm:"type:varchar(255)"`
	Slug             string `gorm:"type:varchar(255)"`
	Name             string `gorm:"type:varchar(255)"`
	ShortDescription string `gorm:"type:varchar(255)"`
	LongDescription  string `gorm:"type:varchar(255)"`
	UsedMaterial     int
	VatType          string `gorm:"type:varchar(255)"` // Usa tu enumeración si la tienes
	IsCustomizable   bool
	IsActive         bool
	TestedInLab      bool
	InternalNotes    string    `gorm:"type:varchar(255)"`
	Pictures         []Picture `gorm:"foreignKey:ProductID"`
	TotalSold        int       `gorm:"-"`
}
