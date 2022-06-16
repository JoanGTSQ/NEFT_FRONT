package models

import (
	"errors"
	"log"

	"github.com/jinzhu/gorm"
)

func newTestGorm(db *gorm.DB) (*testGorm, error) {
	return &testGorm{
		db: db,
	}, nil
}

var _ UserDB = &userGorm{}

type testGorm struct {
	db *gorm.DB
}

type TestDB interface {
	ByID(id string) (*Test, error)
	ByUserID(id uint) (*User, error)
	ByVersionID(id string, idUser uint) (*Version, error)
	ByEmail(email string) (*Test, error)
	Create(test *Test) error
	Update(test *Test) error
	Delete(id uint) error
	DeleteTicket(ticket *Tickets) error
	UpdateTicket(ticket *Tickets) error
	CreateTicket(ticket *Tickets) error
	CreateMessage(mssg *Messages) error
	ByTicketID(id string) Tickets
	RetrieveTicket(idT uint) ([]*Messages, error)
	GetTestCase(idT, idC string, idU uint) (TestCase, error)
	RetrieveTestCase(idT string, idU uint) ([]*TestCase, error)
	UpdateUser(user *User) error
	UserByEmail(email string) *User
	UserByID(id uint) *User
	SearchUsers() []*User
	SearchWorkers() []*User
	ListBugs(id uint) int
	ListTestCase(id uint) []*TestCase
	LastTest(id uint) Test
	CountTest(id uint) []*Test
	LastVersion(id uint) Version
	ByIDVersion(id uint) []*Test
	ByIDVersionLast10(id uint) []*Test
  GetStatsPass(id uint) int
  GetStatsFailed(id uint) int
	ByIDTest(id uint) []*TestCase
	ByUserVersion(id uint) []*Version
	CreateVersion(version *Version) error
	UpdateVersion(version *Version) error
	DeleteVersion(id uint) error
	ByCreatorID(id uint) []*Tickets
	GetTickets() []*Tickets
	GetVersionChange() []*VersionChange
	GetChanges(id uint, filter int) []*ChangesVersion
	GetAllChanges(id uint) []*ChangesVersion
	CreateVersionChange(version *VersionChange) error
	DeleteVersionChange(version *VersionChange) error
	SearchVersionChange(id string) VersionChange
	CreateChange(change *ChangesVersion) error
	SearchChange(id int) ChangesVersion
	UpdateChange(change *ChangesVersion) error
	UpdateVersionChange(change *VersionChange) error
	DeleteChange(id uint) error
	GetStatus() []*StatusCategory
	GetIncidences() []*Incidences
}

type TestService interface {
	TestDB
}

func NewTestService(gD *gorm.DB) TestService {
	ug, err := newTestGorm(gD)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	tv := newTestValidator(ug)
	return &testService{
		TestDB: tv,
	}
}

type testService struct {
	TestDB
}

func newTestValidator(tb TestDB) *testValidator {
	return &testValidator{
		TestDB: tb,
	}
}

type testValidator struct {
	TestDB
}

func (tg *testGorm) ByUserID(id uint) (*User, error) {
	var user User
	db := tg.db.Where("id=?", id)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (tg *testGorm) RetrieveTestCase(idT string, idU uint) ([]*TestCase, error) {
	var version Version
	var test Test
	tg.db.Where("id = ?", idT).Find(&test)
	if test == (Test{}) {
		return nil, errors.New("Empty")
	}

	tg.db.Where("id = ? AND user_id = ?", test.VersionID, idU).Find(&version)
	if version == (Version{}) {
		return nil, errors.New("Empty")
	}

	var testCase []*TestCase
	tg.db.Where("test_id = ?", idT).Find(&testCase)

	if testCase == nil {
		return nil, errors.New("Empty")
	}

	return testCase, nil
}


func (tg *testGorm) GetStatsPass(id uint) int {
  var count int
  tg.db.Table("tests").Where("version_id = ? and status = ?", id, "Completed").Count(&count)
  return count
}

func (tg *testGorm) GetStatsFailed(id uint) int {
  var count int
  tg.db.Table("tests").Where("version_id = ? and status = ?", id, "Failed").Count(&count)
  return count
}


func (tg *testGorm) GetTestCase(idT, idC string, idU uint) (TestCase, error) {
	var testCase TestCase
	var version Version
	var test Test
	tg.db.Where("id = ?", idT).Find(&test)
	if test == (Test{}) {
		return testCase, errors.New("Empty")
	}

	tg.db.Where("id = ? AND user_id = ?", test.VersionID, idU).Find(&version)
	if version == (Version{}) {
		return testCase, errors.New("Empty")
	}
	tg.db.Where("id = ? AND test_id = ?", idC, idT).First(&testCase)

	if testCase == (TestCase{}) {
		return testCase, errors.New("Empty")
	}

	return testCase, nil
}
func (tg *testGorm) SearchWorkers() []*User {
	var user []*User
	tg.db.Select("name,email,perm_level,enabled").Where("perm_level=?", "Worker").Find(&user)
	return user
}
func (tg *testGorm) UserByEmail(email string) *User {
	var user User
	db := tg.db.Where("email = ?", email)
	err := first(db, &user)
	if err != nil {
		return nil
	}
	return &user
}
func (tg *testGorm) UserByID(id uint) *User {
	var user User
	db := tg.db.Where("id = ?", id)
	err := first(db, &user)
	if err != nil {
		return nil
	}
	return &user
}
func (tg *testGorm) SearchUsers() []*User {
	var user []*User
	tg.db.Select("id,name,email,perm_level,enabled,photo,dob").Find(&user)
	return user
}
func (tg *testGorm) ByID(id string) (*Test, error) {
	var test Test
	db := tg.db.Where("id = ?", id).First(&test)
	err := first(db, &test)
	return &test, err
}
func (tg *testGorm) ByVersionID(id string, idUser uint) (*Version, error) {
	var version Version
	db := tg.db.Where("id = ? AND user_id = ?", id, idUser).First(&version)
	err := first(db, &version)
	return &version, err
}
func (tg *testGorm) ByIDVersion(id uint) []*Test {
	var test []*Test
	tg.db.Where("version_id = ?", id).Order("id desc").Find(&test)
	return test
}
func (tg *testGorm) ByIDVersionLast10(id uint) []*Test {
	var test []*Test
	tg.db.Where("version_id = ?", id).Order("id desc").Limit("10").Find(&test)
	return test
}
func (tg *testGorm) ByIDTest(id uint) []*TestCase {
	var testCase []*TestCase
	db := tg.db.Where("test_id = ?", id)
	db.Find(&testCase)

	return testCase
}
func (tg *testGorm) LastVersion(id uint) Version {
	var version Version
	db := tg.db.Where("user_id=?", id)
	db.Last(&version)

	return version
}
func (tg *testGorm) LastTest(id uint) Test {
	var test Test
	db := tg.db.Where("version_id=?", id)
	db.Last(&test)

	return test
}
func (tg *testGorm) ListTestCase(id uint) []*TestCase {
	var test []*TestCase
	db := tg.db.Order("id desc").Where("test_id=?", id)
	db.Find(&test)

	return test
}
func (tg *testGorm) ListBugs(id uint) int {
	var test []TestCase
	tg.db.Where("test_id=? and status=?", id, "Failed").Count(&test)
	return len(test)
}
func (tg *testGorm) ByUserVersion(id uint) []*Version {
	var version []*Version
	db := tg.db.Where("user_id=?", id).Order("id desc")
	db.Find(&version)

	return version
}
func (tg *testGorm) ByEmail(email string) (*Test, error) {
	var test Test

	db := tg.db.Where("email = ?", email)
	err := first(db, &test)
	return &test, err
}
func (tg *testGorm) Create(test *Test) error {
	return tg.db.Create(test).Error
}
func (tg *testGorm) Delete(id uint) error {
	test := Test{Model: gorm.Model{ID: id}}
	return tg.db.Delete(&test).Error
}
func (tg *testGorm) Update(test *Test) error {
	return tg.db.Save(test).Error
}
func (tg *testGorm) UpdateUser(user *User) error {
	return tg.db.Save(user).Error
}

type Version struct {
	gorm.Model
	Name     string `gorm:"not null"`
	CodeTest string `gorm:""`
	UserID   uint   `sql:"type:int REFERENCES users(id)"`
	User     User   `gorm:"ForeignKey:VUFK;AssociationForeignKey:UserVersionFK"`
}

func (tg *testGorm) CreateVersion(version *Version) error {
	return tg.db.Create(version).Error
}
func (tg *testGorm) DeleteVersion(id uint) error {
	version := Version{Model: gorm.Model{ID: id}}
	return tg.db.Delete(&version).Error
}
func (tg *testGorm) UpdateVersion(version *Version) error {
	return tg.db.Save(version).Error
}
func (tg *testGorm) CountTest(id uint) []*Test {
	var versions []Version
	var tests []*Test
	db := tg.db.Where("user_id = ?", id).Order("id desc")
	db.Find(&versions)
	for _, v := range versions {
		db2 := tg.db.Where("version_id = ?", v.ID)
		db2.Find(&tests)
	}
	return tests
}
func (tg *testGorm) ByTestID(id uint) []*TestCase {
	var testCase []*TestCase
	db := tg.db.Where("test_id=?", id).Order("id desc")
	db.Find(&testCase)

	return testCase
}
func (tg *testGorm) ByCreatorID(id uint) []*Tickets {
	var tickets []*Tickets
	tg.db.Where("creator=? AND (status=? OR status=?)", id, "Open", "In Progress").Find(&tickets)
	return tickets
}
func (tg *testGorm) ByTicketID(id string) Tickets {
	var ticket Tickets
	tg.db.Where("id = ?", id).First(&ticket)
	return ticket
}
func (tg *testGorm) RetrieveTicket(idT uint) ([]*Messages, error) {
	var ticket Tickets

	tg.db.Where("id = ?", idT).First(&ticket)
	if ticket == (Tickets{}) {
		return nil, ErrNotFound
	}
	var messages []*Messages
	tg.db.Order("id desc").Where("ticket_id=?", idT).Find(&messages)
	return messages, nil
}
func (tg *testGorm) CreateTicket(ticket *Tickets) error {
	return tg.db.Create(ticket).Error
}
func (tg *testGorm) UpdateTicket(ticket *Tickets) error {
	err := tg.db.Save(ticket).Error
	if err != nil {
		return err
	}
	return nil
}
func (tg *testGorm) DeleteTicket(ticket *Tickets) error {
	return tg.db.Delete(ticket).Error
}
func (tg *testGorm) CreateMessage(mssg *Messages) error {
	return tg.db.Create(mssg).Error
}
func (tg *testGorm) GetTickets() []*Tickets {
	var tickets []*Tickets
	tg.db.Find(&tickets)
	return tickets
}
func (tg *testGorm) GetVersionChange() []*VersionChange {
	var versions []*VersionChange
	tg.db.Order("id desc").Find(&versions)
	return versions
}
func (tg *testGorm) GetChanges(id uint, filter int) []*ChangesVersion {
	var changes []*ChangesVersion
	tg.db.Where("version_change = ? AND categorie = ? ", id, filter).Find(&changes)
	return changes
}
func (tg *testGorm) CreateVersionChange(version *VersionChange) error {
	return tg.db.Create(version).Error
}
func (tg *testGorm) DeleteVersionChange(version *VersionChange) error {
	return tg.db.Delete(version).Error
}
func (tg *testGorm) SearchVersionChange(id string) VersionChange {
	var version VersionChange
	tg.db.Where("id = ?", id).First(&version)
	return version
}
func (tg *testGorm) GetAllChanges(id uint) []*ChangesVersion {
	var changes []*ChangesVersion
	tg.db.Where("version_change = ?", id).Find(&changes)
	return changes
}
func (tg *testGorm) CreateChange(change *ChangesVersion) error {
	return tg.db.Create(change).Error
}
func (tg *testGorm) UpdateVersionChange(change *VersionChange) error {
	return tg.db.Save(change).Error
}
func (tg *testGorm) UpdateChange(change *ChangesVersion) error {
	return tg.db.Save(change).Error
}
func (tg *testGorm) SearchChange(id int) ChangesVersion {
	var change ChangesVersion
	tg.db.Where("id=?", id).First(&change)
	return change
}
func (tg *testGorm) DeleteChange(id uint) error {
	change := ChangesVersion{Model: gorm.Model{ID: id}}
	return tg.db.Delete(&change).Error
}

type Test struct {
	gorm.Model
	Name      string  `gorm:"not null"`
	Status    string  `gorm:"not null"`
	VersionID int     `sql:"type:int REFERENCES versions(id)"`
	Versions  Version `gorm:"ForeignKey:VTFK;AssociationForeignKey:VersionTestFK"`
}

type TestCase struct {
	gorm.Model
	Name        string  `gorm:"not null"`
	Description string  `gorm:"not null"`
	Status      string  `gorm:"not null"`
	DescError   string  `gorm:"default:null"`
	Evidence    string  `gorm:"default:null"`
	TestID      int     `sql:"type:int REFERENCES tests(id)"`
	TestCase    Version `gorm:"ForeignKey:VTFK;AssociationForeignKey:CaseTestFK"`
}

type Tickets struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string `gorm:"not null"`
	Status      string `gorm:"not null"`
	Creator     uint   `sql:"type:int REFERENCES users(id)"`
	Creators    User   `gorm:"ForeignKey:UTFK;AssociationForeignKey:TicketsCreatorFK"`
}

type Messages struct {
	gorm.Model
	Text     string  `gorm:"not null"`
	Sender   uint    `sql:"type:int REFERENCES users(id)"`
	Senders  User    `gorm:"ForeignKey:WTFK;AssociationForeignKey:TicketsSenderFK"`
	TicketID int     `sql:"type:int REFERENCES tickets(id)"`
	Ticket   Tickets `gorm:"ForeignKey:MTFK;AssociationForeignKey:TicketMessageFK"`
}
type VersionChange struct {
	gorm.Model
	Title     string `gorm:"not null"`
	Published bool   `gorm:"not null;default:false"`
}
type Categories struct {
	gorm.Model
	Title string `gorm:"not null"`
}
type ChangesVersion struct {
	gorm.Model
	Title          string        `gorm:"not null"`
	Categorie      int           `sql:"type:int REFERENCES categories(id)"`
	Categories     Categories    `gorm:"ForeignKey:CFK;AssociationForeignKey:CategoriesVersionChangeFK"`
	VersionChange  int           `sql:"type:int REFERENCES version_changes(id)"`
	VersionChanges VersionChange `gorm:"ForeignKey:CVCFK;AssociationForeignKey:ChangeVersionChangeFK"`
}

type StatusCategory struct {
	gorm.Model
	Name   string `gorm:"not null"`
	Status string `gorm:"not null"`
}
type Incidences struct {
	gorm.Model
	Level      string         `gorm:"not null"`
	Message    string         `gorm:"not null"`
	Category   int            `sql:"type:int REFERENCES status_categories(id)"`
	Categories StatusCategory `gorm:"ForeignKey:SIK;AssociationForeignKey:CategoriesStatusIncidenceFK"`
}

func (tg *testGorm) GetStatus() []*StatusCategory {
	var status []*StatusCategory
	tg.db.Find(&status)
	return status
}
func (tg *testGorm) GetIncidences() []*Incidences {
	var incidences []*Incidences
	tg.db.Find(&incidences).Where("category desc")
	return incidences
}
