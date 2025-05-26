package main

import (
	handlers "OzgeContract/internal/handlers"
	repository "OzgeContract/internal/repositories"
	service "OzgeContract/internal/services"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"time"
)

type application struct {
	errorLog                   *log.Logger
	infoLog                    *log.Logger
	companyHandler             *handlers.CompanyHandler
	templateHandler            *handlers.TemplateHandler
	contractHandler            *handlers.ContractHandler
	signatureHandler           *handlers.SignatureHandler
	contractFieldHandler       *handlers.ContractFieldHandler
	signatureFieldValueHandler *handlers.SignatureFieldValueHandler
	statisticsHandler          *handlers.StatisticsHandler
}

func initializeApp(db *sql.DB, errorLog, infoLog *log.Logger) *application {

	// Company
	companyRepo := repository.NewCompanyRepository(db)
	companyService := service.NewCompanyService(companyRepo)
	companyHandler := handlers.NewCompanyHandler(companyService)

	// Template
	templateRepo := repository.NewTemplateRepository(db)
	templateService := service.NewTemplateService(templateRepo)
	templateHandler := handlers.NewTemplateHandler(templateService)

	// Contract
	contractRepo := repository.NewContractRepository(db)
	contractService := service.NewContractService(contractRepo)
	contractHandler := handlers.NewContractHandler(contractService)

	// Signatures
	signatureRepo := repository.NewSignatureRepository(db)
	signatureService := service.NewSignatureService(signatureRepo)
	signatureHandler := handlers.NewSignatureHandler(signatureService)

	contractFieldRepo := repository.NewContractFieldRepository(db)
	contractFieldService := service.NewContractFieldService(contractFieldRepo)
	contractFieldHandler := handlers.NewContractFieldHandler(contractFieldService)

	signatureFieldValueRepo := repository.NewSignatureFieldValueRepository(db)
	signatureFieldValueService := service.NewSignatureFieldValueService(signatureFieldValueRepo)
	signatureFieldValueHandler := handlers.NewSignatureFieldValueHandler(signatureFieldValueService)

	statsRepo := repository.NewStatisticsRepository(db)
	statsService := service.NewStatisticsService(statsRepo)
	statsHandler := handlers.NewStatisticsHandler(statsService)

	return &application{
		errorLog:                   errorLog,
		infoLog:                    infoLog,
		companyHandler:             companyHandler,
		templateHandler:            templateHandler,
		contractHandler:            contractHandler,
		signatureHandler:           signatureHandler,
		contractFieldHandler:       contractFieldHandler,
		signatureFieldValueHandler: signatureFieldValueHandler,
		statisticsHandler:          statsHandler,
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Printf("%v", err)
		panic("failed to connect to database")
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(35)
	if err = db.Ping(); err != nil {
		log.Printf("%v", err)
		panic("failed to ping the database")
		return nil, err
	}
	fmt.Println("successfully connected")

	return db, nil
}

func addSecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
		w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")
		next.ServeHTTP(w, r)
	})
}
