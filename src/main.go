package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"
	"vicinity-tinymesh-vas-co2/src/config"
	"vicinity-tinymesh-vas-co2/src/controller"
	"vicinity-tinymesh-vas-co2/src/vicinity"
)

type Environment struct {
	Config  *config.Config
	DB      string
	LogPath string
}

var app Environment

func (app *Environment) init() {
	// loads values from .env into the system

	app.LogPath = path.Join(".", "logs")
	if err := os.MkdirAll(app.LogPath, os.ModePerm); err != nil {
		log.Fatal("could not create path:", app.LogPath)
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	if err := godotenv.Load(); err != nil {
		log.Fatalln("No .env file found")
	}

	app.Config = config.New()

	// open bolt db
	//db, err := storm.Open(".db", storm.BoltOptions(0600, &bolt.Options{Timeout: 1 * time.Second}))
	//if err != nil {
	//	log.Fatalln(err.Error())
	//}
	//
	//app.DB = db
}

func (app *Environment) run() {
	// Logger
	mainLogger, err := os.OpenFile(path.Join(app.LogPath, fmt.Sprintf("adapter-%s.log", time.Now().Format("2006-01-02"))), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatal("Could not create mainLogger logfile:", err.Error())
	}
	defer mainLogger.Close()

	ginLogger, err := os.OpenFile(path.Join(app.LogPath, fmt.Sprintf("gin-%s.log", time.Now().Format("2006-01-02"))), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Could not create GIN trace logfile:", err.Error())
	}
	defer ginLogger.Close()

	log.SetOutput(mainLogger)

	//defer app.DB.Close()

	// VICINITY
	vas := vicinity.New(app.Config.Vicinity, app.DB)

	// Controller
	server := controller.New(app.Config.Server, vas, ginLogger)
	go server.Listen()
	defer server.Shutdown()

	// INT handler
	quit := make(chan os.Signal, 1)
	defer close(quit)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("VAS shutting down...")
}

// init is invoked before main automatically
func init() {
	app.init()
}

func main() {
	app.run()
}
