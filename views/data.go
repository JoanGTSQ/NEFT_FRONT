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
	Yield  interface{}
}





func (d *Data) SetAlert(err string) {
	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: err,
	}
}

