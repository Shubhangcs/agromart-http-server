package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/shubhangcs/agromart-server/internal/app"
	"github.com/shubhangcs/agromart-server/internal/routes"
)

// @title           Agromart API
// @version         1.0
// @description     Agromart backend server
// @host            localhost:8080

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "specify server port")
	flag.Parse()

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	defer app.DB.Close()
	r := routes.SetupRoutes(app)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
	}

	app.Logger.Printf("server is running on port: %d\n", port)

	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatalln("failed to start server:", err)
	}
}
