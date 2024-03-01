package controllers

import (
	"net/http"

	"jgt.solutions/errorController"
	"jgt.solutions/models"
	"jgt.solutions/views"
)

func (c *Crm) Products(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var es EssentialData
	var err error
	es.Products, err = c.crm.GetAllProducts()
	if err != nil {
		errorController.ErrorLogger.Println("nope ", err)
	}
	vd.Yield = es
	c.ProductsView.Render(w, r, &vd)
}
func (c *Crm) FormNewProduct(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var es EssentialData
	formPicture := FormFile{Id: "productPicture"}
	formSTL := FormFile{Id: "productSTL"}
	es.FormFiles = append(es.FormFiles, &formSTL)
	es.FormFiles = append(es.FormFiles, &formPicture)
	vd.Yield = es
	c.NewProduct.Render(w, r, &vd)
}

type NewProductForm struct {
	Name        string  `schema:"name"`
	Picture     string  `schema:"picture"`
	Price       float64 `schema:"price"`
	Description string  `schema:"description"`
	Weight      float64 `schema:"weight"`
}

// Create Process the signup form
// POST /new-product
func (c *Crm) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form NewProductForm
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

	namePicture, err := uploadPicture(r, "productPicture", "productPicture", form.Name)
	if err != nil {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlError,
			Message: views.AlertMsgGeneric,
		}
		c.NewProduct.Render(w, r, &vd)
		errorController.ErrorLogger.Println(err)
		return
	}
	nameStl, err := uploadPicture(r, "productSTL", "productSTL", form.Name)
	if err != nil {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlError,
			Message: views.AlertMsgGeneric,
		}
		c.NewProduct.Render(w, r, &vd)
		errorController.ErrorLogger.Println(err)
		return
	}
	// return that we have successfully uploaded our file!
	product := models.Product{
		Name:        form.Name,
		Picture:     namePicture,
		Stl:         nameStl,
		Price:       form.Price,
		Description: form.Description,
		Weight:      form.Weight,
	}
	err = c.crm.CreateProduct(&product)
	if err != nil {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlError,
			Message: views.AlertMsgGeneric,
		}
		c.NewProduct.Render(w, r, &vd)
		errorController.ErrorLogger.Println(err)
		return
	}

	http.Redirect(w, r, "/products", http.StatusFound)
}
