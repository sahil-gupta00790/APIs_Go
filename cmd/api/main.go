package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// appn version number
const version = "1.0.0"

//configuration setting

type config struct {
	port int
	env  string //dev , staging , produciton , etc.
}

// Dependecny Injection
type application struct {
	config config
	logger *log.Logger
}

func main() {
	var cfg config
	//read in flags that are needed to populate our config
	flag.IntVar(&cfg.port, "port", 4000, "API SERVER PORT")
	flag.StringVar(&cfg.env, "env", "development", "Environment(development | staging | produciton)")
	flag.Parse()
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	//createing an instance of appn struct
	app := &application{
		config: cfg,
		logger: logger,
	}
	//creating our server now
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthChech", app.healthCheckHandler)
	//create http server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	//start server
	logger.Printf("Starting %s server on %s", cfg.env, srv.Addr)
	err := srv.ListenAndServe()
	logger.Fatal(err)
}
