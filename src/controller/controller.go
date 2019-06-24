package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"io"
	"log"
	"net/http"
	"time"
	"vicinity-tinymesh-vas-co2/src/config"
	"vicinity-tinymesh-vas-co2/src/vicinity"
)

type Server struct {
	config    *config.ServerConfig
	db        *gorm.DB
	vicinity  *vicinity.Client
	http      *http.Server
	ginLogger io.Writer
}

func (server *Server) setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(gin.LoggerWithWriter(server.ginLogger))

	r.GET("/objects", server.handleTD)
	r.PUT("/objects/:iid/publishers/:oid/events/:eid", server.vicinityEventHandler)

	return r
}

func New(serverConfig *config.ServerConfig, db *gorm.DB, vicinity *vicinity.Client, logWriter io.Writer) *Server {
	return &Server{
		vicinity:  vicinity,
		config:    serverConfig,
		ginLogger: logWriter,
		db: db,
	}
}

// Goroutine
func (server *Server) Listen() {
	router := server.setupRouter()

	server.http = &http.Server{
		Addr:              fmt.Sprintf(":%s", server.config.Port),
		Handler:           router,
		WriteTimeout:      10 * time.Second,
		ReadTimeout:       1 * time.Minute,
		ReadHeaderTimeout: 20 * time.Second,
	}

	err := server.http.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			panic(err.Error())
		}
	}
}

func (server *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := server.http.Shutdown(ctx); err != nil {
		log.Print("Server Shutdown error:", err.Error())
	}

	log.Println("Server shut down")
}
