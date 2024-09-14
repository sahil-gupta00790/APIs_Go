package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// appn version number
const version = "1.0.0"

//configuration setting

type config struct {
	port int
	env  string //dev , staging , produciton , etc.
	db   struct {
		dsn string
	}
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

	db, err := openDb()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	// Open the database connection

	// Ping the database to check if the connection is established

	fmt.Println("Connected to PostgreSQL database!")

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
	err = srv.ListenAndServe()
	logger.Fatal(err)
}

func openDb() (*sql.DB, error) {
	connStr := "postgres://postgres:WintheW0rld@localhost:5432/apis_go"
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}
	//create a context with a 5 second tiemout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil

}
