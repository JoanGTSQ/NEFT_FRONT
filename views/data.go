package views

import (
	"html/template"

	"jgt.solutions/models"
)

const (
	AlertLvlError   = "danger"
	AlertLvlWarning = "warning"
	AlertLvlInfo    = "info"
	AlertLvlSuccess = "success"

	// AlertMsgGeneric is displayed when any random error is encountered
	AlertMsgGeneric     = "Something went wrong." + "\n" + "Please try again, and contact us if the problem persists."
	AlertVersionCreated = "The version request is created succesfully"
	AlertTestSuccesfull = "The request to run your test has been received corretly, we will advise you when it finish."
	AlertContactSent    = "The contact mail have been sent, you will receive a copy"
)

type Alert struct {
	Level   string
	Message string
}

type Data struct {
	User   *models.User
	Active string
	Alert  *Alert
	CSRF   template.HTML
	Status []*models.StatusCategory
	Yield  interface{}
}

type RenderValues struct {
	QueueSize      int
	ProcessedTotal int
	Worker1        string
	Worker2        string
	Worker3        string
	Worker4        string
	Status         []*models.StatusCategory
}
type DataLive struct {
	RenderValues *RenderValues
	Yield        interface{}
}
type InfoTest struct {
	ID           uint
	Name         string
	Status       string
	InfoTestCase *[]InfoTestCase
}
type InfoTestCase struct {
	Name        string
	Description string
	Status      string
}
type InfoVersion struct {
	ID       uint
	Name     string
	InfoTest *[]InfoTest
}

type InfoUser struct {
	ID          uint
	Name        string
	Photo       string
	Email       string
	Enabled     bool
	Perms       string
	DOB         string
	InfoVersion []*models.Version
}
type DataCerberus struct {
	CountBugs   int
	CountTests  int
	InfoVersion *InfoVersion
	ListBugs    []*InfoTestCase
	Yield       interface{}
}
type Messages struct {
	User *models.User
	Mssg *Message
}
type Message struct {
	CreatedAt string
	Text      template.HTML
}

func (d *Data) SetAlert(err string) {
	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: err,
	}
}

type DataChangelog struct {
	ID               uint
	Published        bool
	Date             string
	VersionChangelog *models.VersionChange
	Changes          []*models.ChangesVersion
	ChangesWeb       []*models.ChangesVersion
	ChangesInternal  []*models.ChangesVersion
	ChangesCerberus  []*models.ChangesVersion
}
