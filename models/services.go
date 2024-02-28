package models

import "github.com/jinzhu/gorm"

func NewServices(connectionInfo string) (*Services, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}

	db.LogMode(true)

	return &Services{
		User: NewUserService(db),
		Crm:  NewCrmService(db),
		db:   db,
	}, nil
}

type Services struct {
	User UserService
	Crm  CrmService
	db   *gorm.DB
}

func (s *Services) Close() error {
	return s.db.Close()
}

func (s *Services) DestructiveReset() error {
	if err := s.db.DropTableIfExists(&Material{}, &User{}, &pwReset{}).Error; err != nil {
		return err
	}
	return s.AutoMigrate()
}

func (s *Services) AutoMigrate() error {
	if err := s.db.AutoMigrate(&User{}, &pwReset{}, &Material{}, &Costumer{}, &Category{}, &Product{}, &Order{}).Error; err != nil {
		return err
	}
	return nil
}
