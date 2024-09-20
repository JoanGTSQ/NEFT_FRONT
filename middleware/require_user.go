package middleware

import (

	"net/http"

	"jgt.solutions/context"
	"jgt.solutions/logController"
	"jgt.solutions/models"
)

var mySigningKey = []byte("captainjacksparrowsayshi")

type User struct {
	models.UserService
}

func (mw *User) Apply(next http.Handler) http.Handler {
	return mw.ApplyFn(next.ServeHTTP)
}

func (mw *User) ApplyFn(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("remember_token")
		if err != nil {
			next(w, r)
			return
		}
		user, err := mw.UserService.ByRemember(cookie.Value)
		if err != nil {
			next(w, r)
			return
		}

		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		next(w, r)
	})
}

type RequireUser struct {
	User
}

func (mw *RequireUser) Apply(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next(w, r)
	})
}
func (mw *RequireUser) CheckUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			next(w, r)
			return
		} else {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

	})
}

func (mw *RequireUser) CheckPerm(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())

		if user == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		if user.IsAdmin {
			next(w, r)
			return
		} else {
			http.Redirect(w, r, "/404", http.StatusFound)
		}
	})
}
func LogMiddlware(next http.Handler) http.Handler{
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Registra la URL y el método HTTP de la solicitud
        logController.DebugLogger.Printf("Solicitud recibida: %s %s", r.Method, r.URL.Path)

        // Llama al siguiente manejador en la cadena
        next.ServeHTTP(w, r)
    })
}

