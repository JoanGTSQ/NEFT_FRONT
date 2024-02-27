package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"jgt.solutions/controllers"
	"jgt.solutions/errorController"
	"jgt.solutions/middleware"
	"jgt.solutions/models"
	"jgt.solutions/rand"
)

var (
	isProd        bool
	portWebServer string
	RedisServer   string
	debug         bool
	dbDirection   string
	dbUser        string
	dbPassword    string
	dbName        string
	Spinner       *spinner.Spinner
)

func init() {
	flag.BoolVar(&isProd, "isProd", false, "This will ensure csrf is in production enviroment")
	flag.StringVar(&portWebServer, "portWebServer", "443", "This will configure the default port in your web server")
	flag.StringVar(&RedisServer, "redisServer", "127.0.0.1", "This will configure global IP of redis")
	flag.BoolVar(&debug, "debug", false, "This will export all stats to file log.log")
	flag.StringVar(&dbDirection, "dbDirection", "127.0.0.1", "This will configure global IP of database")
	flag.StringVar(&dbUser, "dbUser", "root", "This will configure global user of database")
	flag.StringVar(&dbPassword, "dbPsswd", "", "This will configure global psswd of database")
	flag.StringVar(&dbName, "dbName", "project_dev", "This will configure global name of database")

}

func main() {
	flag.Parse()
	errorController.InitLog(true)

	Spinner = spinner.New(spinner.CharSets[11], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
	Spinner.Color("bgBlack", "bold", "fgGreen")
	Spinner.Suffix = "Starting server..."
	Spinner.Start()

	Spinner.Suffix = " Connecting to database."
	Spinner.Restart()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
		"flora.db.elephantsql.com", 5432, "mljgqygv", "ZfVD-ql9hLg4G6NZ6nCxYUKlgQTg3x_B", "mljgqygv")
	services, err := models.NewServices(psqlInfo)
	if err != nil {
		errorController.ErrorLogger.Println("ABORTING...\ndatabase error: ", err)
		errorController.WD.Content = err.Error()
		errorController.WD.Site = "Error generating services"
		errorController.WD.SendErrorWHWeb()
		os.Exit(0)
	}
	defer services.Close()
	// use DestructiveReset to restore DB
	// use AutoMigrate to create or mantain tables but not delete it
	err = services.AutoMigrate()
	if err != nil {
		fmt.Println(err)
	}
	Spinner.Suffix = " Configuring all static pages..."
	Spinner.Restart()

	staticC := controllers.NewStatic()
	Spinner.Suffix = " Configuring all contact pages..."
	Spinner.Restart()

	Spinner.Suffix = " Configuring all users pages..."
	Spinner.Restart()

	userC := controllers.NewUsers(services.User)

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
	Spinner.Suffix = " Configuring middleware..."
	Spinner.Restart()
	csrfMw := csrf.Protect(b, csrf.Secure(isProd))
	userMW := middleware.User{
		UserService: services.User,
	}
	requireUseMW := middleware.RequireUser{
		User: userMW,
	}

	Spinner.Suffix = " Applying routes..."
	Spinner.Restart()

	r.HandleFunc("/", requireUseMW.CheckPerm(staticC.NewHome)).Methods("GET")

	r.NotFoundHandler = staticC.NotFound
	r.Handle("/505", staticC.Error).Methods("GET")

	// Login And Register

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

	if debug {
		go middleware.PrintStats()
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000" // Default port if not specified
	}
	Spinner.Suffix = " Running web server on port " + port
	Spinner.Restart()
	http.ListenAndServe(":"+port, csrfMw(userMW.Apply(r)))

	// go runServer80()
	// runServerSSL(csrfMw(userMW.Apply(r)))
	// s.Suffix = " Running web server..."
	// s.Restart()
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
