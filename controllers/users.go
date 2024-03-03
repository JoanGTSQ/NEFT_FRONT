package controllers

import (
	"fmt"
	"net/http"
	"time"

	"jgt.solutions/context"
	"jgt.solutions/errorController"

	"jgt.solutions/models"
	"jgt.solutions/rand"
	"jgt.solutions/views"
)

func NewUsers(us models.UserService) *Users {
	return &Users{
		NewView:      views.NewView("dashboard", "users/new"),
		LoginView:    views.NewView("dashboard", "users/login"),
		ForgotPwView: views.NewView("dashboard", "users/forgot_pw"),
		ResetPwView:  views.NewView("dashboard", "users/reset_pw"),

		us: us,
	}
}

type Users struct {
	NewView      *views.View
	LoginView    *views.View
	ForgotPwView *views.View
	ResetPwView  *views.View
	us           models.UserService
}

// New GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, r, nil)
}

type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// New GET /login
func (u *Users) LoginNew(w http.ResponseWriter, r *http.Request) {
	u.LoginView.Render(w, r, nil)
}

// New POST /logout
func (u *Users) Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	user := context.User(r.Context())
	token, _ := rand.RememberToken()
	user.Remember = token
	err := u.us.Update(user)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

// Create Process the signup form
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form SignupForm
	vd.Yield = &form
	if err := ParseForm(r, &form); err != nil {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlError,
			Message: views.AlertMsgGeneric,
		}
		u.NewView.Render(w, r, &vd)
		errorController.ErrorLogger.Println(err)
		return
	}

	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	err := u.us.Create(&user)
	switch err {
	case models.ErrInvalidEmail, models.ErrEmailNotExist,
		models.ErrPasswordIncorrect, models.ErrEmailIsRequired,
		models.ErrEmailIsNotValid, models.ErrEmailIsTaken,
		models.ErrPasswordTooShort, models.ErrPasswordRequired,
		models.ErrRememberTooShort, models.ErrRememberRequired:
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
		errorController.ErrorLogger.Println(err)
		u.NewView.Render(w, r, &vd)
		return
	}
	err = u.signIn(w, &user)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// Login POST /login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form LoginForm
	vd.Yield = &form
	if err := ParseForm(r, &form); err != nil {
		fmt.Println("test")
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlError,
			Message: views.AlertMsgGeneric,
		}
		errorController.ErrorLogger.Println(err)
		u.LoginView.Render(w, r, &vd)
		return
	}

	user, err := u.us.Authenticate(form.Email, form.Password)
	switch err {
	case models.ErrInvalidEmail, models.ErrEmailNotExist,
		models.ErrPasswordIncorrect, models.ErrEmailIsRequired,
		models.ErrEmailIsNotValid, models.ErrEmailIsTaken,
		models.ErrPasswordTooShort, models.ErrPasswordRequired,
		models.ErrRememberTooShort, models.ErrRememberRequired:
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlError,
			Message: models.ErrNotMatchLogin.Error(),
		}
		u.LoginView.Render(w, r, &vd)
		return
	case nil:

	default:
		errorController.ErrorLogger.Println(err)
		return
	}

	err = u.signIn(w, user)
	if err != nil {
		errorController.ErrorLogger.Println(err)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

type ResetPwForm struct {
	Email    string `schema:"email"`
	Token    string `schema:"token"`
	Password string `schema:"password"`
}

// POST /forgot
func (u *Users) InitiateReset(w http.ResponseWriter, r *http.Request) {
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
		errorController.ErrorLogger.Println(err)
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
func (u *Users) ResetPw(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form ResetPwForm
	vd.Yield = &form
	if err := parseURLParams(r, &form); err != nil {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlError,
			Message: views.AlertMsgGeneric,
		}
		errorController.ErrorLogger.Println(err)

		u.ResetPwView.Render(w, r, &vd)
		return
	}
	u.ResetPwView.Render(w, r, &vd)
}

//POST /RESET

func (u *Users) CompleteReset(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form ResetPwForm
	vd.Yield = &form
	if err := ParseForm(r, &form); err != nil {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlError,
			Message: views.AlertMsgGeneric,
		}
		errorController.ErrorLogger.Println(err)

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
			errorController.ErrorLogger.Println(err)

			vd.Alert = &views.Alert{
				Level:   views.AlertLvlError,
				Message: views.AlertMsgGeneric,
			}
		}
		u.ResetPwView.Render(w, r, &vd)
		return
	}
	u.signIn(w, user)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (u *Users) signIn(w http.ResponseWriter, user *models.User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return err
		}
	}

	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	return nil
}
