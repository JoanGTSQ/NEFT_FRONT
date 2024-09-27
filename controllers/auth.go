package controllers

import (
	"net/http"
	"time"

	"jgt.solutions/context"
	"jgt.solutions/logController"
	"jgt.solutions/models"
	"jgt.solutions/rand"
	"jgt.solutions/views"
)

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func NewAuth(us models.UserService) *Auth {
	return &Auth{
		NewView:      views.NewView("dashboard", "users/new"),
		LoginView:    views.NewView("dashboard", "users/login"),
		ForgotPwView: views.NewView("dashboard", "users/forgot_pw"),
		ResetPwView:  views.NewView("dashboard", "users/reset_pw"),
		us:           us,
	}
}

type Auth struct {
	NewView      *views.View
	LoginView    *views.View
	ForgotPwView *views.View
	ResetPwView  *views.View
	us           models.UserService
}

// Login GET /login
func (u *Auth) LoginNew(w http.ResponseWriter, r *http.Request) {
	u.LoginView.Render(w, r, nil)
}

// Login POST /login
func (u *Auth) Login(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form LoginForm
	vd.Yield = &form
	if err := ParseForm(r, &form); err != nil {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlError,
			Message: views.AlertMsgGeneric,
		}
		logController.ErrorLogger.Println(err)
		u.LoginView.Render(w, r, &vd)
		return
	}

	user, err := u.us.Authenticate(form.Email, form.Password)
	switch err {
	case models.ErrInvalidEmail, models.ErrEmailNotExist, models.ErrPasswordIncorrect:
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlError,
			Message: models.ErrNotMatchLogin.Error(),
		}
		u.LoginView.Render(w, r, &vd)
		return
	case nil:
	default:
		logController.ErrorLogger.Println(err)
		return
	}

	err = signIn(w, user)
	if err != nil {
		logController.ErrorLogger.Println(err)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

// Logout POST /logout
func (u *Auth) Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}

	user := context.User(r.Context())
	token, _ := rand.RememberToken()
	user.RememberToken = token
	err := user.Update()
	if err != nil {
		logController.ErrorLogger.Println("Error al actualizar el usuario:", err)
		http.Redirect(w, r, "/505", http.StatusFound)
		return
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

// Funciones para el manejo de restablecimiento de contraseña

type ResetPwForm struct {
	Email    string `schema:"email"`
	Token    string `schema:"token"`
	Password string `schema:"password"`
}

// POST /forgot
func (u *Auth) InitiateReset(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form ResetPwForm
	vd.Yield = &form
	if err := ParseForm(r, &form); err != nil {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlError,
			Message: err.Error(),
		}
		u.ForgotPwView.Render(w, r, &vd)
		return
	}

	user, err := u.us.ByEmail(form.Email)
	if err != nil {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlError,
			Message: models.ErrEmailNotExist.Error(),
		}
		u.ForgotPwView.Render(w, r, &vd)
		return
	}
	token, err := u.us.InitiateReset(user.ID)
	if err != nil {
		logController.ErrorLogger.Println(err)
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlError,
			Message: views.AlertMsgGeneric,
		}
		u.ForgotPwView.Render(w, r, &vd)
		return
	}
	_ = token
	vd.Alert = &views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Congratulations, you will receive instructions to reset your password via email",
	}
	u.ForgotPwView.Render(w, r, &vd)
}

// GET /reset
func (u *Auth) ResetPw(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form ResetPwForm
	vd.Yield = &form
	if err := parseURLParams(r, &form); err != nil {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlError,
			Message: views.AlertMsgGeneric,
		}
		logController.ErrorLogger.Println(err)
		u.ResetPwView.Render(w, r, &vd)
		return
	}
	u.ResetPwView.Render(w, r, &vd)
}

// POST /reset
func (u *Auth) CompleteReset(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form ResetPwForm
	vd.Yield = &form
	if err := ParseForm(r, &form); err != nil {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlError,
			Message: views.AlertMsgGeneric,
		}
		logController.ErrorLogger.Println(err)
		u.ResetPwView.Render(w, r, &vd)
		return
	}
	user, err := u.us.CompleteReset(form.Token, form.Password)
	if err != nil {
		if err == models.ErrSamePasswordReset {
			vd.Alert = &views.Alert{
				Level:   views.AlertLvlError,
				Message: models.ErrSamePasswordReset.Error(),
			}
		} else {
			logController.ErrorLogger.Println(err)
			vd.Alert = &views.Alert{
				Level:   views.AlertLvlError,
				Message: views.AlertMsgGeneric,
			}
		}
		u.ResetPwView.Render(w, r, &vd)
		return
	}
	signIn(w, user)
	http.Redirect(w, r, "/", http.StatusFound)
}

// Create Process the signup form
// POST /signup
func (u *Auth) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form SignupForm
	vd.Yield = &form
	if err := ParseForm(r, &form); err != nil {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlError,
			Message: views.AlertMsgGeneric,
		}
		u.NewView.Render(w, r, &vd)
		logController.ErrorLogger.Println(err)
		return
	}

	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	err := user.Create()
	switch err {
	case models.ErrInvalidEmail, models.ErrEmailIsRequired, models.ErrEmailIsNotValid, models.ErrEmailIsTaken,
		models.ErrPasswordTooShort, models.ErrPasswordRequired:
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlError,
			Message: err.Error(),
		}
		u.NewView.Render(w, r, &vd)
		return
	case nil:
	default:
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlError,
			Message: err.Error(),
		}
		logController.ErrorLogger.Println(err)
		u.NewView.Render(w, r, &vd)
		return
	}
	err = signIn(w, &user)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

// New GET /signup
func (u *Auth) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, r, nil)
}
