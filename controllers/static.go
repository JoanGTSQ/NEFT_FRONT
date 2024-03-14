package controllers

import (
	"net/http"
	"regexp"
	"strings"

	"jgt.solutions/errorController"
	"jgt.solutions/models"

	"jgt.solutions/views"
)

func NewStatic() *Static {
	return &Static{
		NotFound: views.NewView("dashboard", "static/404"),
		Error:    views.NewView("dashboard", "static/505"),
	}
}

type Static struct {
	NotFound *views.View
	Error    *views.View
}

func NewContact() *Contact {
	return &Contact{
		HomeView:  views.NewView("dashboard", "static/home"),
		LoginView: views.NewView("dashboard", "users/login"),
	}
}

type Contact struct {
	HomeView  *views.View
	LoginView *views.View
}

type ContactForm struct {
	Name    string `schema:"name"`
	Email   string `schema:"email"`
	Subject string `schema:"subject"`
	Message string `schema:"message"`
}

// Create Process the contact form
// POST /
func (c *Contact) ContactForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\]+@[a-z0-9.\-]+\.[a-z]{2,16}$`)
	var vd views.Data
	var form ContactForm
	vd.Yield = &form
	if err := ParseForm(r, &form); err != nil {
		errorController.ErrorLogger.Println(err)
		return
	}
	form.Email = strings.ToLower(form.Email)
	form.Email = strings.TrimSpace(form.Email)

	if !emailRegex.MatchString(form.Email) {
		vd.SetAlert(models.ErrEmailIsNotValid.Error())
		c.HomeView.Render(w, r, &vd)
		return
	}

	vd.Alert = &views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: views.AlertContactSent,
	}
	form = ContactForm{}
	c.HomeView.Flush(w, r, &vd)
}
