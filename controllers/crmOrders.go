package controllers

import (
	"jgt.solutions/logController"

	"fmt"
	"github.com/gorilla/mux"
	"jgt.solutions/models"
	"jgt.solutions/views"
	"net/http"
	"strconv"
)

func (c *Crm) Orders(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var es EssentialData
	var err error
	es.Orders, err = c.crm.GetAllOrders()
	if err != nil {
		logController.ErrorLogger.Println("nope ", err)
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
		logController.ErrorLogger.Println("nope ", err)
	}
	es.Products, err = c.crm.GetAllProducts()
	if err != nil {
		logController.ErrorLogger.Println("nope ", err)
	}
	es.Materials, err = c.crm.GetAllMaterials()
	if err != nil {
		logController.ErrorLogger.Println("nope ", err)
	}
	vd.Yield = es
	c.NewOrder.Render(w, r, &vd)
}

// CreateOrder procesa el formulario de creación de orden
// POST /new-order
func (c *Crm) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form NewOrderForm
	var totalTime int
	var totalCost float64
	var totalSale float64

	vd.Yield = &form

	// Parsear el formulario
	if err := ParseForm(r, &form); err != nil {
		logController.ErrorLogger.Println("Error al parsear el formulario:", err)
		return
	}

	var products []*models.OrderProductMaterial
	for i := 0; ; i++ {
		productID := r.FormValue(fmt.Sprintf("products[%d][productID]", i))
		if productID == "" {
			break // Salir del bucle si no hay más productos
		}
		materialID := r.FormValue(fmt.Sprintf("products[%d][materialID]", i))
		cost := r.FormValue(fmt.Sprintf("products[%d][cost]", i))
		sale := r.FormValue(fmt.Sprintf("products[%d][sale]", i))
		quality := r.FormValue(fmt.Sprintf("products[%d][quality]", i))
		// Convertir los ID de productos y materiales a int64
		productIDInt, err := strconv.ParseInt(productID, 10, 64)
		if err != nil {
			// Manejar el error
			http.Error(w, fmt.Sprintf("Error al obtener el ID del producto %d", i), http.StatusInternalServerError)
			return
		}
		materialIDInt, err := strconv.ParseInt(materialID, 10, 64)
		if err != nil {
			// Manejar el error
			http.Error(w, fmt.Sprintf("Error al obtener el ID del material del producto %d", i), http.StatusInternalServerError)
			return
		}
		productCost, err := strconv.ParseFloat(cost, 64)
		if err != nil {
			// Manejar el error
			http.Error(w, fmt.Sprintf("Error al obtener el ID del material del producto %d", i), http.StatusInternalServerError)
			return
		}
		productSale, err := strconv.ParseFloat(sale, 64)
		if err != nil {
			// Manejar el error
			http.Error(w, fmt.Sprintf("Error al obtener el ID del material del producto %d", i), http.StatusInternalServerError)
			return
		}
		product, err := c.crm.SearchProductByID(productIDInt)
		if err != nil {
			logController.ErrorLogger.Println("Error al buscar el producto:", err)
			return // Continuar con el siguiente producto si hay un error
		}

		material, err := c.crm.SearchMaterialByID(materialIDInt)
		if err != nil {
			logController.ErrorLogger.Println("Error al buscar el material:", err)
			return // Continuar con el siguiente producto si hay un error
		}
		material.Weight -= product.Weight
		products = append(products, &models.OrderProductMaterial{

			Product:  *product,
			Material: *material,
			Quality:  quality,
		})
		totalTime += product.TimeMinutes
		// TODO: Calcular el costo y la venta adecuadamente
		totalCost += productCost
		totalSale += productSale
		// Actualizar el peso del material
		err = c.crm.UpdateMaterial(material)
		if err != nil {
			logController.ErrorLogger.Println("Error al actualizar el material:", err)
			continue // Continuar con el siguiente producto si hay un error
		}
	}

	// Crear la orden en la base de datos
	order := models.Order{
		CustomerID:  int(form.Customer),
		Products:    products,
		TimeMinutes: totalTime,
		Cost:        totalCost,
		Sale:        totalSale,
		Sent:        true,
	}

	err := c.crm.CreateOrder(&order)
	if err != nil {
		logController.ErrorLogger.Println("Error al crear la orden:", err)
		return
	}

	http.Redirect(w, r, "/orders", http.StatusFound)
}

func (c *Crm) ViewSingleOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]
    orderIDint, err := strconv.ParseInt(orderID, 10, 64)
    if err != nil {
        // Manejar el error
        http.Error(w, fmt.Sprintf("Error al obtener el ID del pedido"), http.StatusInternalServerError)
        return
    }
    order, err := c.crm.SearchOrderByID(int(orderIDint))
    if err != nil {
        // Manejar el error
        http.Error(w, fmt.Sprintf("Error al obtener el ID del producto"), http.StatusInternalServerError)
        return
    }
    var vd views.Data
    vd.Yield = order
    c.SingleOrder.Render(w, r, &vd)
}
