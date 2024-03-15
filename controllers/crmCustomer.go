package controllers

import (
	"jgt.solutions/logController"
	"jgt.solutions/models"
	"jgt.solutions/views"
	"net/http"
)

func (c *Crm) Customers(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var es EssentialData
	var err error
	es.Customers, err = c.crm.GetAllUsers()
	if err != nil {
		logController.ErrorLogger.Println("No se han podido obtener todos los clientes ", err)
		return
	}
	vd.Yield = es
	c.CustomersView.Render(w, r, &vd)
}

func (c *Crm) FormNewCustomer(w http.ResponseWriter, r *http.Request) {

	c.NewCustomer.Render(w, r, nil)
}

type NewCustomerForm struct {
	Name      string `schema:"name"`
	Email     string `schema:"email"`
	Direction string `schema:"direction"`
	Phone     string `schema:"phone"`
	Instagram string `schema:"instagram"`
}

// Create Process the signup form
// POST /new-product
func (c *Crm) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form NewCustomerForm
	vd.Yield = &form

	if err := ParseForm(r, &form); err != nil {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlError,
			Message: views.AlertMsgGeneric,
		}
		c.NewProduct.Render(w, r, &vd)
		logController.ErrorLogger.Println(err)
		return
	}

	customer := models.User{
		Name:      form.Name,
		Direction: form.Direction,
		Instagram: form.Instagram,
		Email:     form.Email,
		Phone:     form.Phone,
	}
	err := c.crm.CreateCustomer(&customer)
	if err != nil {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlError,
			Message: views.AlertMsgGeneric,
		}

		logController.ErrorLogger.Println(err)
		c.NewProduct.Render(w, r, &vd)
		return
	}

	http.Redirect(w, r, "/customers", http.StatusFound)
}
