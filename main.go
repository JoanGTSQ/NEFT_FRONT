package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"


	"github.com/gorilla/csrf"
	"jgt.solutions/errorController"
	"jgt.solutions/rand"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"jgt.solutions/controllers"
	"jgt.solutions/middleware"
	"jgt.solutions/models"
)

var (
	isProd        bool
	portWebServer string
	RedisServer   string
	enableStats   bool
  dbDirection   string
  dbUser        string
  dbPassword    string
  dbName        string
)

func init() {
	flag.BoolVar(&isProd, "isProd", false, "This will ensure csrf is in production enviroment")
	flag.StringVar(&portWebServer, "portWebServer", "443", "This will configure the default port in your web server")
	flag.StringVar(&RedisServer, "redisServer", "127.0.0.1", "This will configure global IP of redis")
	flag.BoolVar(&enableStats, "enableStats", false, "This will export all stats to file log.log")
  flag.StringVar(&dbDirection, "dbDirection", "127.0.0.1", "This will configure global IP of database")
  flag.StringVar(&dbUser, "dbUser", "root", "This will configure global user of database")
  flag.StringVar(&dbPassword, "dbPsswd", "", "This will configure global psswd of database")
  flag.StringVar(&dbName, "dbName", "project_dev", "This will configure global name of database")
}

func main() {
	flag.Parse()
	if enableStats {
		f, err := os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()
		wrt := io.MultiWriter(os.Stdout, f)
		log.SetOutput(wrt)
	}

	log.Println("Starting....")
	log.Println("Redis server configured.")

	controllers.InitializeViper()
	controllers.InitializeOAuthGoogle()
	log.Println("Connecting to database.")
  
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
	dbDirection, 5432, dbUser, dbPassword, dbName)
	services, err := models.NewServices(psqlInfo)
	if err != nil {
		log.Println("ABORTING...\ndatabase error: ", err)
		errorController.WD.Content = err.Error()
					errorController.WD.Site = "Error generating services"
		errorController.WD.SendErrorWHWeb()
		os.Exit(0)
	}
	defer services.Close()
	err = services.AutoMigrate()
	if err != nil {
		fmt.Println(err)
	}
	log.Println("Configuring all static pages...")
	staticC := controllers.NewStatic()
  log.Println("Configuring all contact pages...")
	contactC := controllers.NewContact()
  log.Println("Configuring all users pages...")
	userC := controllers.NewUsers(services.User)
 //  log.Println("Configuring all cerberus pages...")
	// CerberusC := cerberusController.NewCerberus(services.Test)
 //  log.Println("Configuring all gia pages...")
	// giaC := giaController.NewGia(services.Test)

	r := mux.NewRouter()
	cssHandler := http.FileServer(http.Dir("./css/"))
	cssHandler = http.StripPrefix("/css/", cssHandler)
	r.PathPrefix("/css/").Handler(cssHandler)

	jsHandler := http.FileServer(http.Dir("./js/"))
	jsHandler = http.StripPrefix("/js/", jsHandler)
	r.PathPrefix("/js/").Handler(jsHandler)

	imgHandler := http.FileServer(http.Dir("./images/"))
	imgHandler = http.StripPrefix("/images/", imgHandler)
	r.PathPrefix("/images/").Handler(imgHandler)

	fontsHandler := http.FileServer(http.Dir("./fonts/"))
	fontsHandler = http.StripPrefix("/fonts/", fontsHandler)
	r.PathPrefix("/fonts/").Handler(fontsHandler)

	vendorHandler := http.FileServer(http.Dir("./vendor_web/"))
	vendorHandler = http.StripPrefix("/vendor/", vendorHandler)
	r.PathPrefix("/vendor/").Handler(vendorHandler)

	assetsHandler := http.FileServer(http.Dir("./assets/"))
	assetsHandler = http.StripPrefix("/assets/", assetsHandler)
	r.PathPrefix("/assets/").Handler(assetsHandler)

  assetsAppHandler := http.FileServer(http.Dir("./assets/"))
  assetsAppHandler = http.StripPrefix("/app-assets/", assetsAppHandler)
	r.PathPrefix("/app-assets/").Handler(assetsAppHandler)

  
	b, err := rand.Bytes(32)
	if err != nil {
		return
	}
	log.Println("Configuring middleware")
	csrfMw := csrf.Protect(b, csrf.Secure(isProd))
	userMW := middleware.User{
		UserService: services.User,
	}
	requireUseMW := middleware.RequireUser{
		User: userMW,
	}

	log.Println("Applying routes")
	r.HandleFunc("/", staticC.NewHome).Methods("GET")
	r.HandleFunc("/", contactC.ContactForm).Methods("POST")

	// // /cerberus
	// r.HandleFunc("/cerberus", requireUseMW.Apply(CerberusC.New)).Methods("GET")

	// // /cerberus/users
	// r.HandleFunc("/cerberus/users", requireUseMW.CheckPerm(CerberusC.GetUsers)).Methods("GET")
	// r.HandleFunc("/cerberus/users", requireUseMW.CheckPerm(CerberusC.PostUsers)).Methods("POST")
 //  r.HandleFunc("/cerberus/users/{id:[0-9]+}", requireUseMW.Apply(CerberusC.GetUser)).Methods("GET")
	// // /cerberus/versions
	// r.HandleFunc("/cerberus/versions", requireUseMW.Apply(CerberusC.GetVersions)).Methods("GET")
	// r.HandleFunc("/cerberus/versions", requireUseMW.Apply(CerberusC.PostVersions)).Methods("POST")
	// r.HandleFunc("/cerberus/versions/{id:[0-9]+}", requireUseMW.Apply(CerberusC.SeeVersions)).Methods("GET")

	// // /cerberus/ticketss
	// r.Handle("/cerberus/tickets", requireUseMW.Apply(CerberusC.GetTickets)).Methods("GET")
	// r.Handle("/cerberus/tickets", requireUseMW.Apply(CerberusC.PostTicket)).Methods("POST")
	// r.Handle("/cerberus/tickets/{ticketId:[0-9]+}", requireUseMW.Apply(CerberusC.SeeTicket)).Methods("GET")
	// r.Handle("/cerberus/tickets/{ticketId:[0-9]+}", requireUseMW.Apply(CerberusC.PostTickets)).Methods("POST")
	// // /cerberus/test
	// r.HandleFunc("/cerberus/test", requireUseMW.Apply(CerberusC.GetTest)).Methods("GET")
	// r.HandleFunc("/cerberus/test/{testId:[0-9]+}", requireUseMW.Apply(CerberusC.SeeTest)).Methods("GET")

	// //GIA
	// r.HandleFunc("/cerberus/gia", requireUseMW.Apply(giaC.New)).Methods("GET")
	// r.HandleFunc("/cerberus/gia", requireUseMW.Apply(giaC.POST)).Methods("POST")
	// // /cerberus/changelog
	// r.HandleFunc("/cerberus/changelog", requireUseMW.CheckPerm(CerberusC.GetAdminChangelog)).Methods("GET")
	// r.HandleFunc("/cerberus/changelog", requireUseMW.CheckPerm(CerberusC.PostAdminChangelog)).Methods("POST")
	// r.HandleFunc("/cerberus/changelog/{changelog:[0-9]+}", requireUseMW.CheckPerm(CerberusC.GetIDChangeLog)).Methods("GET")
	// r.HandleFunc("/cerberus/changelog/{changelog:[0-9]+}", requireUseMW.CheckPerm(CerberusC.PostIDChangeLog)).Methods("POST")
	// // STATIC WEBPAGE
	// r.HandleFunc("/changelog", CerberusC.GetChangelog).Methods("GET")
	// r.HandleFunc("/faq", CerberusC.NewFaq).Methods("GET")

	r.NotFoundHandler = staticC.NotFound
	r.Handle("/505", staticC.Error).Methods("GET")
	// Login And Register
	r.HandleFunc("/newgoogle", controllers.HandleGoogleRegister)
	r.HandleFunc("/registergoogle", userC.CallBackFromGoogle)
	r.HandleFunc("/logingoogle", controllers.HandleGoogleLogin)
	r.HandleFunc("/logingl", userC.CallBackLoginFromGoogle)
	r.HandleFunc("/signup", requireUseMW.CheckUser(userC.New)).Methods("GET")
	r.HandleFunc("/signup", requireUseMW.CheckUser(userC.Create)).Methods("POST")
	r.HandleFunc("/login", requireUseMW.CheckUser(userC.LoginNew)).Methods("GET")
	r.HandleFunc("/login", requireUseMW.CheckUser(userC.Login)).Methods("POST")
	r.HandleFunc("/logout", userC.Logout).Methods("POST")
	r.Handle("/forgot", userC.ForgotPwView).Methods("GET")
	r.HandleFunc("/forgot", userC.InitiateReset).Methods("POST")
	r.HandleFunc("/reset", userC.ResetPw).Methods("GET")
	r.HandleFunc("/reset", userC.CompleteReset).Methods("POST")


	// Start server
	log.Println("Starting web servers....")
	if enableStats {
		go middleware.PrintStats()
		r.Use(middleware.IPREMOTE)
	}
  port := os.Getenv("PORT")
  if port == "" {
      port = "9000" // Default port if not specified
  }
  http.ListenAndServe(":" + port,csrfMw(userMW.Apply(r)))
	// go runServer80()
	// runServerSSL(csrfMw(userMW.Apply(r)))
}
func runServerSSL(handler http.Handler) {
	log.Println("Starting web server at ::" + portWebServer)
	err := http.ListenAndServeTLS(":"+portWebServer, "cert.cer", "key.key", handler)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Server 443 running")
	}
}
func runServer80() {
	log.Println("Starting web server at ::80")
	err := http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://"+r.Host+r.URL.String(), http.StatusMovedPermanently)
	}))
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("Server 80 running")
	}
}
