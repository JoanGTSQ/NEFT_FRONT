package controllers

import (
	"net/http"

	"jgt.solutions/logController"
	"jgt.solutions/models"
	"jgt.solutions/views"
)

func NewCrm(crm models.CrmService) *Crm {
	return &Crm{
		HomeDashboard: views.NewView("dashboard", "crm/home"),
		ProductsView:  views.NewView("dashboard", "crm/products"),
		NewProduct:    views.NewView("dashboard", "crm/addProduct"),
		MaterialsView: views.NewView("dashboard", "crm/materials"),
		NewMaterial:   views.NewView("dashboard", "crm/addMaterial"),
		CustomersView: views.NewView("dashboard", "crm/customers"),
		NewCustomer:   views.NewView("dashboard", "crm/addCustomer"),
		OrdersView:    views.NewView("dashboard", "crm/orders"),
		NewOrder:      views.NewView("dashboard", "crm/addOrder"),
		SingleOrder:   views.NewView("dashboard", "crm/singleOrder"),
		PrintersView:  views.NewView("dashboard", "crm/printers"),
		NewPrinter:    views.NewView("dashboard", "crm/addPrinter"),
		crm:           crm,
	}
}

type Crm struct {
	HomeDashboard *views.View
	ProductsView  *views.View
	NewProduct    *views.View
	MaterialsView *views.View
	NewMaterial   *views.View
	PrintersView  *views.View
	NewPrinter    *views.View
	CustomersView *views.View
	NewCustomer   *views.View
	OrdersView    *views.View
	NewOrder      *views.View
	SingleOrder   *views.View
	crm           models.CrmService
}
type FormFile struct {
	Id string
}
type EssentialData struct {
	TotalSales         float64
	TotalOrderExpenses float64
	Orders             []*models.Order
	Products           []*models.Product
	Materials          []*models.Material
	Customers          []*models.User
	FormFiles          []*FormFile
	Printers           []*models.Printer
}

func (c *Crm) Home(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var es EssentialData
	var err error
	es.TotalSales, err = c.crm.CountAllSales()
	if err != nil {
		logController.ErrorLogger.Println("no se han podido obtener todas las ventas ", err)
		return
	}
	es.TotalOrderExpenses, err = c.crm.CountAllSalesExpenses()
	if err != nil {
		logController.ErrorLogger.Println("no se han podido obtener todas los costes ", err)
	}

	es.Orders, err = c.crm.GetAllOrders()
	if err != nil {
		logController.ErrorLogger.Println("no se han podido obtener todos los pedidos ", err)
	}
	es.Products, err = c.crm.GetAllProducts()
	if err != nil {
		logController.ErrorLogger.Println("no se han podido obtener todos los productos ", err)
	}

	vd.Yield = es
	c.HomeDashboard.Render(w, r, &vd)
}
