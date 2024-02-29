package controllers

import (
	"jgt.solutions/errorController"
	"jgt.solutions/models"
	"jgt.solutions/views"
	"net/http"
)

func (c *Crm) Customers(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var es EssentialData
	var err error
	es.Customers, err = c.crm.GetAllCustomers()
	if err != nil {
		errorController.ErrorLogger.Println("nope ", err)
	}
	vd.Yield = es
	c.CustomersView.Render(w, r, &vd)
}

func (c *Crm) FormNewCustomer(w http.ResponseWriter, r *http.Request) {

	c.NewCustomer.Render(w, r, nil)
}

type NewCustomerForm struct {
	Name      string  `schema:"name"`
	Email     string  `schema:"email"`
	Direction string  `schema:"direction"`
	Phone     string `schema:"phone"`
	Origin    string     `schema:"origin"`
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
		errorController.ErrorLogger.Println(err)
		return
	}

	customer := models.Customer{
		Name:      form.Name,
		Direction: form.Direction,
		Email:     form.Email,
		Phone:     form.Phone,
		Origin:    form.Origin,
	}
	err := c.crm.CreateCustomer(&customer)
	if err != nil {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlError,
			Message: views.AlertMsgGeneric,
		}
		c.NewProduct.Render(w, r, &vd)
		errorController.ErrorLogger.Println(err)
		return
	}

	http.Redirect(w, r, "/customers", http.StatusFound)
}