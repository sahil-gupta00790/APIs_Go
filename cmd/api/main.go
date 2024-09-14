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
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
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
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:WintheW0rld@localhost:5432/apis_go", "postgresql dsn")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgresQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgresQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgresQL max idle time")

	flag.Parse()
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDb(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	// Open the database connection

	// Ping the database to check if the connection is established

	logger.Println("Connected to PostgreSQL database!")

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

func openDb(cfg config) (*sql.DB, error) {
	//connStr := "postgres://postgres:WintheW0rld@localhost:5432/apis_go"
	db, err := sql.Open("pgx", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)
	//create a context with a 5 second tiemout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil

}
