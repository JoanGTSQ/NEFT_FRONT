package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	"jgt.solutions/context"
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
			next(w,r)
			return
		}
		user, err := mw.UserService.ByRemember(cookie.Value)
		if err != nil {
			next(w,r)
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
			http.Redirect(w, r, "/cerberus", http.StatusFound)
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
		if user.PermLevel == "Admin" || user.PermLevel == "Worker" {
			next(w, r)
			return
		} else {
			http.Redirect(w, r, "/cerberus", http.StatusFound)
		}
	})
}

func IPREMOTE(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		if !strings.Contains(req.RequestURI, "css") &&
			!strings.Contains(req.RequestURI, "js") &&
			!strings.Contains(req.RequestURI, "fonts") &&
			!strings.Contains(req.RequestURI, "assets") {
			log.Println("NEW ACCESS: " + req.RemoteAddr + ":" + req.RequestURI)
		}
		next.ServeHTTP(res, req)
	})
}

func PrintStats() {
	ticker := time.NewTicker(2 * time.Minute)
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
			log.Println("-----os----stats----")
			log.Printf("cpu system:  %s %%\n", fmt.Sprintf("%.2f", float64(after.System-before.System)/total*100))
			log.Printf("cpu idle:    %s %%\n", fmt.Sprintf("%.2f", float64(after.Idle-before.Idle)/total*100))
			log.Printf("memory used: %d  mb\n", memory.Used/1024/1024)
			log.Printf("memory free: %d  mb\n", memory.Free/1024/1024)
			log.Println("--------------------")
		case <-quit:
			ticker.Stop()
			return
		}
	}  
}

