package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"jgt.solutions/controllers"
	"jgt.solutions/logController"
	"jgt.solutions/middleware"
	"jgt.solutions/models"
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
)

func init() {
	flag.BoolVar(&isProd, "isProd", false, "This will ensure csrf is in production enviroment")
	flag.StringVar(&portWebServer, "portWebServer", "443", "This will configure the default port in your web server")
	flag.StringVar(&RedisServer, "redisServer", "127.0.0.1", "This will configure global IP of redis")
	flag.BoolVar(&debug, "debug", false, "This will export all stats to file log.log")
	flag.StringVar(&dbDirection, "dbDirection", "", "This will configure global IP of database")
	flag.StringVar(&dbUser, "dbUser", "", "This will configure global user of database")
	flag.StringVar(&dbPassword, "dbPsswd", "", "This will configure global psswd of database")
	flag.StringVar(&dbName, "dbName", "", "This will configure global name of database")
}

func main() {
	flag.Parse()
	logController.InitLog(debug)

	logController.InfoLogger.Println("Starting server")

	//psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
	//dbDirection, 5432, dbUser, dbPassword, dbName)
	dsn := "qaiq735:PuroVici!1@tcp(qaiq735.protogt.com:3306)/qaiq735?charset=utf8mb4&parseTime=True&loc=Local"

	services, err := models.NewServices(dsn)
	if err != nil {
		logController.ErrorLogger.Println(err)
		os.Exit(0)
	}
	defer services.Close()

	// use DestructiveReset to restore DB
	// use AutoMigrate to create or mantain tables but not delete it
	logController.DebugLogger.Println("Configuring Database connection")
	// err = services.AutoMigrate()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	logController.DebugLogger.Println("Configuring all controllers")

	staticC := controllers.NewStatic()
	crmC := controllers.NewCrm(services.Crm)

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

	// b, err := rand.Bytes(32)
	// if err != nil {
	// 	return
	// }

	// Middleware configuration
	logController.DebugLogger.Println("Configuring middleware...")

	// csrfMw := csrf.Protect(b, csrf.Secure(isProd))
	// userMW := middleware.User{
	// 	// UserService: services.User,
	// }
	// requireUseMW := middleware.RequireUser{
	// 	User: userMW,
	// }

	// Routes configuration
	logController.DebugLogger.Println("Applying routes...")

	r.Use(middleware.LogMiddlware)

	r.HandleFunc("/", crmC.Home).Methods("GET")
	// r.HandleFunc("/products", requireUseMW.CheckPerm(crmC.Products)).Methods("GET")
	// r.HandleFunc("/new-product", requireUseMW.CheckPerm(crmC.FormNewProduct)).Methods("GET")
	// r.HandleFunc("/new-product", requireUseMW.CheckPerm(crmC.CreateProduct)).Methods("POST")
	// r.HandleFunc("/materials", requireUseMW.CheckPerm(crmC.Materials)).Methods("GET")
	// r.HandleFunc("/new-material", requireUseMW.CheckPerm(crmC.FormNewMaterial)).Methods("GET")
	// r.HandleFunc("/new-material", requireUseMW.CheckPerm(crmC.CreateMaterial)).Methods("POST")
	// r.HandleFunc("/customers", requireUseMW.CheckPerm(crmC.Customers)).Methods("GET")
	// r.HandleFunc("/new-customer", requireUseMW.CheckPerm(crmC.FormNewCustomer)).Methods("GET")
	// r.HandleFunc("/new-customer", requireUseMW.CheckPerm(crmC.CreateCustomer)).Methods("POST")
	// r.HandleFunc("/orders", requireUseMW.CheckPerm(crmC.Orders)).Methods("GET")
	// r.HandleFunc("/new-order", requireUseMW.CheckPerm(crmC.FormNewOrder)).Methods("GET")
	// r.HandleFunc("/new-order", requireUseMW.CheckPerm(crmC.CreateOrder)).Methods("POST")
	// r.HandleFunc("/orders/{id}", requireUseMW.CheckPerm(crmC.ViewSingleOrder)).Methods("GET")
	// r.HandleFunc("/printers", requireUseMW.CheckPerm(crmC.Printers)).Methods("GET")
	// r.HandleFunc("/new-printer", requireUseMW.CheckPerm(crmC.FormNewPrinter)).Methods("GET")
	// r.HandleFunc("/new-printer", requireUseMW.CheckPerm(crmC.CreatePrinter)).Methods("POST")
	r.NotFoundHandler = staticC.NotFound
	r.Handle("/505", staticC.Error).Methods("GET")

	// Login And Register

	// r.HandleFunc("/signup", requireUseMW.CheckUser(userC.New)).Methods("GET")
	// r.HandleFunc("/signup", requireUseMW.CheckUser(userC.Create)).Methods("POST")
	// r.HandleFunc("/login", requireUseMW.CheckUser(userC.LoginNew)).Methods("GET")
	// r.HandleFunc("/login", requireUseMW.CheckUser(userC.Login)).Methods("POST")
	// r.HandleFunc("/logout", userC.Logout).Methods("POST")
	// r.Handle("/forgot", userC.ForgotPwView).Methods("GET")
	// r.HandleFunc("/forgot", userC.InitiateReset).Methods("POST")
	// r.HandleFunc("/reset", userC.ResetPw).Methods("GET")
	// r.HandleFunc("/reset", userC.CompleteReset).Methods("POST")

	// Start server

	if debug {
		go middleware.PrintStats()
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000" // Default port if not specified
	}
	logController.InfoLogger.Println("Running web server on port " + port)

	http.ListenAndServe(":"+port, r)

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
