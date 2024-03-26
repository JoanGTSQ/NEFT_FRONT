package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	"jgt.solutions/context"
	"jgt.solutions/logController"
)

var mySigningKey = []byte("captainjacksparrowsayshi")

type User struct {
	// models.UserService
}

func (mw *User) Apply(next http.Handler) http.Handler {
	return mw.ApplyFn(next.ServeHTTP)
}

func (mw *User) ApplyFn(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// cookie, err := r.Cookie("remember_token")
		// if err != nil {
		// 	next(w, r)
		// 	return
		// }
		// user, err := mw.UserService.ByRemember(cookie.Value)
		// if err != nil {
		// 	next(w, r)
		// 	return
		// }

		// ctx := r.Context()
		// ctx = context.WithUser(ctx, user)
		// r = r.WithContext(ctx)
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
		// user := context.User(r.Context())

		// if user == nil {
		// 	http.Redirect(w, r, "/login", http.StatusFound)
		// 	return
		// }
		// if user.PermLevel == "Admin" || user.PermLevel == "Worker" {
		// 	next(w, r)
		// 	return
		// } else {
		// 	http.Redirect(w, r, "/404", http.StatusFound)
		// }
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

func PrintStats() {
	ticker := time.NewTicker(30 * time.Second)

	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			before, err := cpu.Get()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				return
			}
			time.Sleep(time.Duration(1) * time.Second)
			after, err := cpu.Get()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				return
			}
			memory, err := memory.Get()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				return
			}
			total := float64(after.Total - before.Total)

			textPre := ("Printing stats\n-----os------stats---")

			currentCPU := fmt.Sprintf("\ncpu system:  %s %%", fmt.Sprintf("%.2f", float64(after.System-before.System)/total*100))

			totalCPU := fmt.Sprintf("\ncpu idle:    %s %%", fmt.Sprintf("%.2f", float64(after.Idle-before.Idle)/total*100))

			currentRAM := fmt.Sprintf("\nmemory used: %d  mb", memory.Used/1024/1024)

			freeRAM := fmt.Sprintf("\nmemory free: %d  mb", memory.Free/1024/1024)

			textPost := ("\n--------------------")

			logController.DebugLogger.Println(textPre + currentCPU + totalCPU + currentRAM + freeRAM + textPost)
		case <-quit:
			ticker.Stop()
			return
		}
	}
}
