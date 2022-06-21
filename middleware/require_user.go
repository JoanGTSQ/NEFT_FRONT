package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
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
type Stat struct{
	s 	*spinner.Spinner
}


func PrintStats() {
	ticker := time.NewTicker(30 * time.Second)
	statCPU := Stat{
		s: spinner.New(spinner.CharSets[11], 100* time.Millisecond, 	spinner.WithWriter(os.Stderr)),
	}
	statCPU.s.Start()
	statCPUTotal := Stat{
		s: spinner.New(spinner.CharSets[11], 100* time.Millisecond, 	spinner.WithWriter(os.Stderr)),
	}
	statCPUTotal.s.Start()
	statRAM := Stat{
		s: spinner.New(spinner.CharSets[11], 100* time.Millisecond, 	spinner.WithWriter(os.Stderr)),
	}
	statRAM.s.Start()
	statRAMFree := Stat{
		s: spinner.New(spinner.CharSets[11], 100* time.Millisecond, 	spinner.WithWriter(os.Stderr)),
	}
	statRAMFree.s.Start()
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
			preStat := Stat{
				s: spinner.New(spinner.CharSets[11], 100* time.Millisecond, 	spinner.WithWriter(os.Stderr)),
			}
			preStat.s.Start()
			textPre := ("-----os------stats---")
			preStat.printStat(textPre)
			
			currentCPU := fmt.Sprintf("cpu system:  %s %%", fmt.Sprintf("%.2f", float64(after.System-before.System)/total*100))
			statCPU.s.Color("bgBlack", "bold", "fgGreen")
			statCPU.printStat(currentCPU)
			

			
			totalCPU := fmt.Sprintf("cpu idle:    %s %%", fmt.Sprintf("%.2f", float64(after.Idle-before.Idle)/total*100))
			statCPUTotal.printStat(totalCPU)

			
			currentRAM := fmt.Sprintf("memory used: %d  mb", memory.Used/1024/1024)
			statRAM.printStat(currentRAM)
			
			freeRAM := fmt.Sprintf("memory free: %d  mb", memory.Free/1024/1024)
			statRAMFree.printStat(freeRAM)
			textPost := ("--------------------")
			postStat := Stat{
				s: spinner.New(spinner.CharSets[11], 100* time.Millisecond, 	spinner.WithWriter(os.Stderr)),
			}
			postStat.s.Start()
			postStat.printStat(textPost)

		case <-quit:
			ticker.Stop()
			return
		}
	}  
}


func (stat *Stat)printStat(text string) {
	stat.s.Color("bgBlack", "bold", "fgGreen")
	stat.s.Suffix = text
	stat.s.Restart()
}