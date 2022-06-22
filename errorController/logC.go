package errorController

import(
	"io"
	"os"
	"log"
	"bytes"
)
 
 
var (
	Reset  = "\033[0m"
 Red    = "\033[31m"
 Green  = "\033[32m"
 Yellow = "\033[33m"
 Blue   = "\033[34m"
 Purple = "\033[35m"
 Cyan   = "\033[36m"
 Gray   = "\033[37m"
	White  = "\033[97m"
    WarningLogger *log.Logger
    InfoLogger    *log.Logger
		DebugLogger    *log.Logger
    ErrorLogger   *log.Logger
)
func InitLog(debugEnabled bool){
		f, err := os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)
	InfoLogger = log.New(wrt, White + "\nINFO: " + Reset, log.Ldate|log.Ltime|log.Lshortfile)

	
	WarningLogger = log.New(wrt, Yellow + "\nWARNING: " + Reset, log.Ldate|log.Ltime|log.Lshortfile)

	DebugLogger = log.New(wrt, Green + "\nDEBUG: " + Reset, log.Ldate|log.Ltime|log.Lshortfile)
	if debugEnabled {
		var buff bytes.Buffer
		DebugLogger.SetOutput(&buff)
	}
	
	ErrorLogger = log.New(wrt, Red + "\nERROR: " + Reset, log.Ldate|log.Ltime|log.Lshortfile)
}