package models

type SetupAppBase struct {
	ID                 string `gorm:"type:uuid;primaryKey"`
	PriceKw            int    `gorm:"type:int"`
	PrinterConsumption int    `gorm:"type:int"`
	Comment            string `gorm:"type:text"`
}
