package models

type Attribute struct {
	ID         string   `gorm:"type:uuid;primaryKey"` // ID de la orden
	ProductID  string   `gorm:"type:varchar(255)"`
	FinishID   string   `gorm:"type:varchar(255)"`
	MaterialID string   `gorm:"type:varchar(255)"`
	PictureID  string   `gorm:"type:varchar(255)"`
	Price      int      `gorm:"type:int"`
	CostPrice  int      `gorm:"type:int"`
	OfferPrice int      `gorm:"type:int"`
	Minutes    int      `gorm:"type:int"`
	InOffer    bool     `gorm:"default:false"`
	IsActive   bool     `gorm:"default:true"`
	Product    Product  `gorm:"foreignKey:ProductID"` // Relación con Product
	Material   Material `gorm:"foreignKey:MaterialID"`
	Finish     Finish   `gorm:"foreignKey:FinishID"`
}
