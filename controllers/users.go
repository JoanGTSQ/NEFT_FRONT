package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
	"jgt.solutions/logController"
	"jgt.solutions/models"
	"jgt.solutions/views"
)

type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func NewUsers(us models.UserService) *Users {
	return &Users{
		Index: views.NewView("dashboard", "crm/user/index"),
		us:    us,
	}
}

type Users struct {
	Index *views.View
	us    models.UserService
}

// New GET /signup
func (c *Crm) UserIndexFunc(w http.ResponseWriter, r *http.Request) {

	var vd views.Data
	var es EssentialData
	var err error
	es.Users, err = c.crm.GetAllUsers()
	if err != nil {
		logController.ErrorLogger.Println("No se han podido obtener todos los clientes ", err)
		return
	}
	vd.Yield = es
	c.UserIndex.Render(w, r, &vd)
}

// New GET /signup
func (c *Crm) UserShowFunc(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	vars := mux.Vars(r)

	user := models.User{
		ID: vars["id"],
	}
	err := user.ByID()
	//es.Users, err = c.crm.GetAllUsers()
	if err != nil {
		logController.ErrorLogger.Println("No se han podido obtener todos los clientes ", err)
		return
	}
	//GetDirections
	user.GetDirections()
	user.GetOrders()
	vd.Yield = user
	c.UserShow.Render(w, r, &vd)
}
