package controllers

import (
	"math"
	"net/http"

	"jgt.solutions/errorController"
	"jgt.solutions/models"
	"jgt.solutions/views"
)

func NewCrm(crm models.CrmService) *Crm {
	return &Crm{
		HomeDashboard: views.NewView("dashboard", "crm/home"),

		crm: crm,
	}
}

type Crm struct {
	HomeDashboard *views.View
	crm           models.CrmService
}
type EssentialData struct {
	TotalSales         float64
	TotalOrderExpenses float64
	Profit             float64
	ProfitPercentage   float64
	Orders             []*models.Order
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
	es.Profit = es.TotalSales - es.TotalOrderExpenses
	es.ProfitPercentage = math.Round((es.TotalSales/es.TotalOrderExpenses)*100) / 10

    es.Orders, err = c.crm.GetAllOrders()
	if err != nil {
		errorController.ErrorLogger.Println("nope ", err)
	}

	vd.Yield = es
	c.HomeDashboard.Render(w, r, &vd)
}
