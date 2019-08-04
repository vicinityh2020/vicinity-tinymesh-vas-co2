package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"
	"vicinity-tinymesh-vas-co2/vas-co2-backend/config"
	"vicinity-tinymesh-vas-co2/vas-co2-backend/controller"
	"vicinity-tinymesh-vas-co2/vas-co2-backend/database"
	"vicinity-tinymesh-vas-co2/vas-co2-backend/model"
	"vicinity-tinymesh-vas-co2/vas-co2-backend/sms"
	"vicinity-tinymesh-vas-co2/vas-co2-backend/vicinity"
)

type Environment struct {
	Config  *config.Config
	DB      *gorm.DB
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
}

func (app *Environment) newLogWriter(logName string) *os.File {
	l, err := os.OpenFile(path.Join(app.LogPath, fmt.Sprintf("%s-%s.log", logName, time.Now().Format("2006-01-02"))), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatal("Could not create mainLogger logfile:", err.Error())
	}

	return l
}

func (app *Environment) run() {
	rand.Seed(time.Now().UnixNano())
	// Main logger
	mainLogger := app.newLogWriter("adapter")
	defer mainLogger.Close()

	// Gin logger
	ginLogger := app.newLogWriter("gin")
	defer ginLogger.Close()

	// DB logger
	dbLogger := app.newLogWriter("gorm")
	defer dbLogger.Close()

	kpiLogger := app.newLogWriter("tracker")
	defer kpiLogger.Close()

	log.SetOutput(mainLogger)

	// Database
	app.DB = database.New(app.Config.Database, dbLogger)
	defer app.DB.Close()

	//app.DB.DropTableIfExists(&model.Reading{}, &model.Sensor{}, &model.Notification{})
	app.DB.AutoMigrate(&model.Sensor{}, &model.Reading{}, &model.Notification{})
	app.DB.Model(&model.Reading{}).AddForeignKey("sensor_oid", "sensors(oid)", "CASCADE", "RESTRICT")

	// KPI Tracker
	kpiTracker := vicinity.NewKPITracker(app.Config.Vicinity, app.DB, kpiLogger)
	kpiTracker.Tick(10)
	defer kpiTracker.Stop()

	kpiTracker.GatherAndReport()

	// KeySMS
	keysmsClient := sms.New(app.Config.SMS, app.DB)
	if err := keysmsClient.Auth(); err != nil {
		log.Fatalln(err.Error())
	}

	// VICINITY
	vas := vicinity.New(app.Config.Vicinity, app.DB)

	// Controller
	server := controller.New(app.Config.Server, app.DB, vas, keysmsClient, ginLogger)
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
