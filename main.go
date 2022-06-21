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
	"github.com/briandowns/spinner"
	"time"
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
	Spinner				*spinner.Spinner
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
	Spinner = spinner.New(spinner.CharSets[11], 100* time.Millisecond, 	spinner.WithWriter(os.Stderr))
	Spinner.Color("bgBlack", "bold", "fgGreen")
	Spinner.Suffix = "Starting server..."
	Spinner.Start()
	
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
	Spinner.Suffix = " Connecting to database."
	Spinner.Restart()

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
	Spinner.Suffix = " Configuring all static pages..."
	Spinner.Restart()
	staticC := controllers.NewStatic()
	Spinner.Suffix = " Configuring all contact pages..."
	Spinner.Restart()
	contactC := controllers.NewContact()



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

	Spinner.Suffix = " Applying routes..."
	Spinner.Restart()
	
	r.HandleFunc("/", staticC.NewHome).Methods("GET")
	r.HandleFunc("/",contactC.ContactForm).Methods("POST")

	r.NotFoundHandler = staticC.NotFound
	r.Handle("/505", staticC.Error).Methods("GET")



	// Start server
	
	if enableStats {
		go middleware.PrintStats()
		r.Use(middleware.IPREMOTE)
	}
  port := os.Getenv("PORT")
  if port == "" {
      port = "9000" // Default port if not specified
  }
	Spinner.Suffix = " Running web server on " + port
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
