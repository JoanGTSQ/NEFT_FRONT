package models

import "github.com/jinzhu/gorm"
import "time"

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
    // Eliminar las tablas existentes si existen
    if err := s.db.DropTableIfExists(&Material{}, &User{}, &pwReset{}, &Customer{}, &Category{}, &Product{}, &OrderProductMaterial{}, &Order{}).Error; err != nil {
        return err
    }

    // Volver a crear todas las tablas y aplicar migraciones
    return s.AutoMigrate()
}

func (s *Services) AutoMigrate() error {
    // Realizar migraciones automáticas para todos los modelos
    if err := s.db.AutoMigrate(&User{}, &pwReset{}, &Material{}, &Customer{}, &Category{}, &Product{}, &OrderProductMaterial{}, &Order{}).Error; err != nil {
        return err
    }
    return nil
}

type ProtoModel struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}
