package models

import (
	"github.com/jinzhu/gorm"
	"time"
    _ "github.com/go-sql-driver/mysql"
)

func NewServices(connectionInfo string) (*Services, error) {

	db, err := gorm.Open("mysql", connectionInfo)
	if err != nil {
		return nil, err
	}

	db.LogMode(false)

	return &Services{
		Crm:  NewCrmService(db),
		db:   db,
	}, nil
}

type Services struct {
	Crm  CrmService
	db   *gorm.DB
}

func (s *Services) Close() error {
	return s.db.Close()
}


type ProtoModel struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}
