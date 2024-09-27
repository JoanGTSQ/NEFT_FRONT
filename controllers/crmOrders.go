package controllers

import (
	"jgt.solutions/logController"
	"jgt.solutions/models"

	"net/http"

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
	c.OrderIndex.Render(w, r, &vd)
}

func (c *Crm) ViewSingleOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	order := models.Order{
		ID: orderID,
	}

	err := order.ByID()
	if err != nil {
		logController.ErrorLogger.Println("Error al obtener el pedido:", err)
		http.Redirect(w, r, "/505", http.StatusFound)
		return
	}
	var vd views.Data
	vd.Yield = order
	c.OrderShow.Render(w, r, &vd)
}
