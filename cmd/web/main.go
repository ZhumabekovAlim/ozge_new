package main

import (
	"OzgeContract/internal/config"
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	cfg := config.LoadConfig()

	port := os.Getenv("PORT")
	if port != "" {
		port = ":" + port
	} else {
		port = ":4000"

	}

	addr := flag.String("addr", port, "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(cfg.Database.URL)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
		}
	}(db)

	app := initializeApp(cfg, db, errorLog, infoLog)

	//c := cors.New(cors.Options{
	//	AllowedOrigins:   []string{"*"},
	//	AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
	//	AllowCredentials: true,
	//	AllowedHeaders:   []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Accept", "Origin", "Cache-Control", "X-Requested-With"},
	//})

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      addSecurityHeaders(c.Handler(app.routes())),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err = srv.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}
	infoLog.Printf("Starting server on %s", *addr)
	select {}
}
