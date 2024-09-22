package controllers

import (
	"jgt.solutions/logController"
	"jgt.solutions/views"
	"net/http"
)

func (c *Crm) Materials(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var es EssentialData
	var err error
	es.Materials, err = c.crm.GetAllMaterials()
	if err != nil {
		logController.ErrorLogger.Println("No se han podido obtener todos los materiales ", err)
		return
	}
	vd.Yield = es
	c.MaterialsView.Render(w, r, &vd)
}

func (c *Crm) FormNewMaterial(w http.ResponseWriter, r *http.Request) {

	c.NewMaterial.Render(w, r, nil)
}

type NewMaterialForm struct {
	Name     string  `schema:"name"`
	Color    string  `schema:"color"`
	Supplier string  `schema:"supplier"`
	Price    float64 `schema:"price"`
	Weight   float64 `schema:"weight"`
}

// Create Process the signup form
// POST /new-product
func (c *Crm) CreateMaterial(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form NewMaterialForm
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

	http.Redirect(w, r, "/materials", http.StatusFound)
}
