package controllers

import (
	"jgt.solutions/logController"

	
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"jgt.solutions/views"
)

func (c *Crm) Orders(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var es EssentialData
	var err error
	es.Orders, err = c.crm.GetAllOrders()
	if err != nil {
		logController.ErrorLogger.Println(err)
		http.Redirect(w, r, "/505", http.StatusFound)
		return
	}
	vd.Yield = es
	c.OrdersView.Render(w, r, &vd)
}

func (c *Crm) FormNewOrder(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var es EssentialData
	var err error
	es.Customers, err = c.crm.GetAllUsers()
	if err != nil {
		logController.ErrorLogger.Println("No se han podido obtener todos los clientes ", err)
		http.Redirect(w, r, "/505", http.StatusFound)
		return
	}
	es.Products, err = c.crm.GetAllProducts()
	if err != nil {
		logController.ErrorLogger.Println("No se han podido obtener todos los productos ", err)
		http.Redirect(w, r, "/505", http.StatusFound)
		return
	}
	es.Materials, err = c.crm.GetAllMaterials()
	if err != nil {
		logController.ErrorLogger.Println("No se han podido obtener todos los materiales ", err)
		http.Redirect(w, r, "/505", http.StatusFound)
		return
	}
	es.Printers, err = c.crm.GetAllPrinters()
	if err != nil {
		logController.ErrorLogger.Println("No se han podido obtener todos los materiales ", err)
		http.Redirect(w, r, "/505", http.StatusFound)
		return
	}
	vd.Yield = es
	c.NewOrder.Render(w, r, &vd)
}

// CreateOrder procesa el formulario de creación de orden
// POST /new-order
func (c *Crm) CreateOrder(w http.ResponseWriter, r *http.Request) {


	http.Redirect(w, r, "/orders", http.StatusFound)
}

func (c *Crm) ViewSingleOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]
	orderIDint, err := strconv.ParseInt(orderID, 10, 64)
	if err != nil {
		// Manejar el error
		logController.ErrorLogger.Println("Error al obtener el ID del pedido:", err)
		http.Redirect(w, r, "/505", http.StatusFound)
		return
	}
	order, err := c.crm.SearchOrderByID(int(orderIDint))
	if err != nil {
		logController.ErrorLogger.Println("Error al obtener el pedido:", err)
		http.Redirect(w, r, "/505", http.StatusFound)
		return
	}
	var vd views.Data
	vd.Yield = order
	c.SingleOrder.Render(w, r, &vd)
}
