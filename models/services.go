package models

import "github.com/jinzhu/gorm"

func NewServices(connectionInfo string) (*Services, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}

	db.LogMode(false)

	return &Services{
		User: NewUserService(db),
		Test: NewTestService(db),
		db:   db,
	}, nil
}

type Services struct {
	User UserService
	Test TestService
	db   *gorm.DB
}

func (s *Services) Close() error {
	return s.db.Close()
}

func (s *Services) DestructiveReset() error {
	if err := s.db.DropTableIfExists(&Incidences{}, &StatusCategory{}, &ChangesVersion{}, &Categories{}, &VersionChange{}, &Messages{}, &Tickets{}, &TestCase{}, &Test{}, &Version{}, &User{}, &pwReset{}).Error; err != nil {
		return err
	}
	return s.AutoMigrate()
}

func (s *Services) DestructiveStatic() error {
	if err := s.db.DropTableIfExists(&Incidences{}, &StatusCategory{}).Error; err != nil {
		return err
	}
	return s.AutoMigrate()
}

func (s *Services) AutoMigrate() error {
	if err := s.db.AutoMigrate(&User{}, &pwReset{}, &Version{}, &Test{}, &TestCase{}, &Tickets{}, &Messages{}, &VersionChange{}, &Categories{}, &ChangesVersion{}, &StatusCategory{}, &Incidences{}).Error; err != nil {
		return err
	}
	return nil
}
