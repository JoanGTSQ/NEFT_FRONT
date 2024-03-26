package controllers

// import (
// 	"net/http"

// 	"jgt.solutions/logController"
// 	"jgt.solutions/models"
// 	"jgt.solutions/views"
// )

// func (c *Crm) Printers(w http.ResponseWriter, r *http.Request) {
// 	var vd views.Data
// 	var es EssentialData
// 	var err error
// 	es.Printers, err = c.crm.GetAllPrinters()
// 	if err != nil {
// 		logController.ErrorLogger.Println("No se han podido obtener todos los printeres ", err)
// 		return
// 	}
// 	vd.Yield = es
// 	c.PrintersView.Render(w, r, &vd)
// }

// func (c *Crm) FormNewPrinter(w http.ResponseWriter, r *http.Request) {

// 	c.NewPrinter.Render(w, r, nil)
// }

// type NewPrinterForm struct {
// 	Name     string  `schema:"name"`
// 	Color    string  `schema:"color"`
// 	Supplier string  `schema:"supplier"`
// 	Price    float64 `schema:"price"`
// 	Weight   float64 `schema:"weight"`
// }

// // Create Process the signup form
// // POST /new-product
// func (c *Crm) CreatePrinter(w http.ResponseWriter, r *http.Request) {
// 	var vd views.Data
// 	var form NewPrinterForm
// 	vd.Yield = &form

// 	if err := ParseForm(r, &form); err != nil {
// 		vd.Alert = &views.Alert{
// 			Level:   views.AlertLvlError,
// 			Message: views.AlertMsgGeneric,
// 		}
// 		c.NewProduct.Render(w, r, &vd)
// 		logController.ErrorLogger.Println(err)
// 		return
// 	}

// 	printer := models.Printer{
// 		Name: form.Name,
// 	}
// 	err := c.crm.CreatePrinter(&printer)
// 	if err != nil {
// 		vd.Alert = &views.Alert{
// 			Level:   views.AlertLvlError,
// 			Message: views.AlertMsgGeneric,
// 		}
// 		c.NewProduct.Render(w, r, &vd)
// 		logController.ErrorLogger.Println(err)
// 		return
// 	}

// 	http.Redirect(w, r, "/printers", http.StatusFound)
// }
