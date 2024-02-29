package controllers

import (
	"net/http"
	"os"

	"io"
	"jgt.solutions/errorController"
	"jgt.solutions/models"
	"jgt.solutions/views"
	"math/rand"
	"strconv"
	"time"
)

func NewCrm(crm models.CrmService) *Crm {
	return &Crm{
		HomeDashboard: views.NewView("dashboard", "crm/home"),
		ProductsView:  views.NewView("dashboard", "crm/products"),
		NewProduct:    views.NewView("dashboard", "crm/addProduct"),
		crm:           crm,
	}
}

type Crm struct {
	HomeDashboard *views.View
	ProductsView  *views.View
	NewProduct    *views.View
	crm           models.CrmService
}
type EssentialData struct {
	TotalSales         float64
	TotalOrderExpenses float64
	Orders             []*models.Order
	Products           []*models.Product
}

func (c *Crm) Home(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var es EssentialData
	var err error
	es.TotalSales, err = c.crm.CountAllSales()
	if err != nil {
		errorController.ErrorLogger.Println("nope ", err)
	}
	es.TotalOrderExpenses, err = c.crm.CountAllSalesExpenses()
	if err != nil {
		errorController.ErrorLogger.Println("nope ", err)
	}

	es.Orders, err = c.crm.GetAllOrders()
	if err != nil {
		errorController.ErrorLogger.Println("nope ", err)
	}
	es.Products, err = c.crm.GetAllProducts()
	if err != nil {
		errorController.ErrorLogger.Println("nope ", err)
	}

	vd.Yield = es
	c.HomeDashboard.Render(w, r, &vd)
}

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

	c.NewProduct.Render(w, r, nil)
}

type NewProductForm struct {
	Name        string  `schema:"name"`
	Picture     string  `schema:"picture"`
	Price       float64 `schema:"price"`
	Description string  `schema:"description"`
	Weight      int     `schema:"weight"`
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

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		errorController.ErrorLogger.Println("Error Retrieving the File")
		errorController.ErrorLogger.Println(err)
		return
	}
	defer file.Close()
	errorController.DebugLogger.Printf("Uploaded File: %+v\n", handler.Filename)
	errorController.DebugLogger.Printf("File Size: %+v\n", handler.Size)
	errorController.DebugLogger.Printf("MIME Header: %+v\n", handler.Header)
	rand.Seed(time.Now().UnixNano())

	// Genera un número entero aleatorio entre 0 y 100.
    // Verificar si el directorio existe
    //TODO cambiar por el directorio de la carpeta public
    if _, err := os.Stat("/home/runner/NEFTFRONT-2/assets/images/products/"); os.IsNotExist(err) {
        // Si no existe, crear el directorio
        if err := os.MkdirAll("/home/runner/NEFTFRONT-2/assets/images/products/", os.ModePerm); err != nil {
            // Manejar el error si no se puede crear el directorio
            errorController.ErrorLogger.Println("Error al crear directorio:", err)
            return
        }
    }
	numPicture := rand.Intn(1000000)
	namePicture := "upload-" + strconv.Itoa(numPicture) + ".png"
	newPicture, err := os.Create("/home/runner/NEFTFRONT-2/assets/images/products/" + namePicture)
	if err != nil {
		errorController.ErrorLogger.Println(err)
	}
	defer newPicture.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		errorController.ErrorLogger.Println(err)
	}
	// write this byte array to our temporary file
	newPicture.Write(fileBytes)
	// return that we have successfully uploaded our file!
	product := models.Product{
		Name:        form.Name,
		Picture:     namePicture,
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
