package controllers

import (
	"jgt.solutions/errorController"
	"jgt.solutions/models"

	// "jgt.solutions/models"
	"net/http"

	"jgt.solutions/views"
)

func (c *Crm) Orders(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var es EssentialData
	var err error
	es.Orders, err = c.crm.GetAllOrders()
	if err != nil {
		errorController.ErrorLogger.Println("nope ", err)
	}
	vd.Yield = es
	c.OrdersView.Render(w, r, &vd)
}

func (c *Crm) FormNewOrder(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var es EssentialData
	var err error
	es.Customers, err = c.crm.GetAllCustomers()
	if err != nil {
		errorController.ErrorLogger.Println("nope ", err)
	}
	es.Products, err = c.crm.GetAllProducts()
	if err != nil {
		errorController.ErrorLogger.Println("nope ", err)
	}
	es.Materials, err = c.crm.GetAllMaterials()
	if err != nil {
		errorController.ErrorLogger.Println("nope ", err)
	}
	vd.Yield = es
	c.NewOrder.Render(w, r, &vd)
}

type NewOrderForm struct {
	Material   int64   `schema:"materialID"`
	Customer   int64   `schema:"customerID"`
	Cost       float64 `schema:"cost"`
	Sale       float64 `schema:"sale"`
	Origin     string  `schema:"origin"`
	ProductsID []int64 `schema:"products[]"`
}

// // Create Process the signup form
// // POST /new-product
func (c *Crm) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form NewOrderForm
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
	errorController.InfoLogger.Println(form)
	var products []*models.Product
	var totalWeightOrder float64
	for _, productID := range form.ProductsID {
		product, err := c.crm.SearchByID(productID)
		if err != nil {
			vd.Alert = &views.Alert{
				Level:   views.AlertLvlError,
				Message: views.AlertMsgGeneric,
			}
			c.NewProduct.Render(w, r, &vd)
			errorController.ErrorLogger.Println(err)
			return
		}
		totalWeightOrder += product.Weight
		products = append(products, product)
	}
	errorController.InfoLogger.Println(products)
	order := models.Order{
		MaterialID: int(form.Material),
		Cost:       form.Cost,
		Sale:       form.Sale,
		Sent:       true,
		Products:   products,
		CustomerID: int(form.Customer),
	}
	err := c.crm.CreateOrder(&order)
	if err != nil {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlError,
			Message: views.AlertMsgGeneric,
		}
		c.NewProduct.Render(w, r, &vd)
		errorController.ErrorLogger.Println(err)
		return
	}
	material, err := c.crm.SearchMaterialByID(form.Material)
	if err != nil {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlError,
			Message: views.AlertMsgGeneric,
		}
		c.NewProduct.Render(w, r, &vd)
		errorController.ErrorLogger.Println(err)
		return
	}
	material.Weight -= totalWeightOrder
	errorController.InfoLogger.Println(material)
	err = c.crm.UpdateMaterial(material)
	if err != nil {
		http.Redirect(w, r, "/orders", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/orders", http.StatusFound)
}
