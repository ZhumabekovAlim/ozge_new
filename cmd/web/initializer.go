package main

import (
	"OzgeContract/internal/config"
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
	companyBalanceHandler      *handlers.CompanyBalanceHandler
	tariffPlanHandler          *handlers.TariffPlanHandler
	paymentHandler             *handlers.PaymentRequestHandler
	smsHandler                 *handlers.SMSHandler
	adminHandler               *handlers.AdminHandler
}

func initializeApp(cfg config.Config, db *sql.DB, errorLog, infoLog *log.Logger) *application {

	// Company
	companyBalanceRepo := repository.NewCompanyBalanceRepository(db)
	companyRepo := repository.NewCompanyRepository(db)
	companyService := service.NewCompanyService(companyRepo, companyBalanceRepo)
	companyHandler := handlers.NewCompanyHandler(companyService)

	// Template
	templateRepo := repository.NewTemplateRepository(db)
	templateService := service.NewTemplateService(templateRepo)
	templateHandler := handlers.NewTemplateHandler(templateService)

	// Contract
	contractRepo := repository.NewContractRepository(db)
	contractFieldRepo := repository.NewContractFieldRepository(db)
	contractService := service.NewContractService(contractRepo, contractFieldRepo)
	contractHandler := handlers.NewContractHandler(contractService)

	// Signatures
	signatureRepo := repository.NewSignatureRepository(db)
	signatureService := service.NewSignatureService(signatureRepo, contractRepo, companyBalanceRepo)
	signatureFieldValueRepo := repository.NewSignatureFieldValueRepository(db)
	signatureFieldValueService := service.NewSignatureFieldValueService(signatureFieldValueRepo)
	signatureHandler := handlers.NewSignatureHandler(signatureService, signatureFieldValueService)

	contractFieldService := service.NewContractFieldService(contractFieldRepo)
	contractFieldHandler := handlers.NewContractFieldHandler(contractFieldService)

	signatureFieldValueHandler := handlers.NewSignatureFieldValueHandler(signatureFieldValueService)

	statsRepo := repository.NewStatisticsRepository(db)
	statsService := service.NewStatisticsService(statsRepo)
	statsHandler := handlers.NewStatisticsHandler(statsService)

	companyBalanceService := service.NewCompanyBalanceService(companyBalanceRepo)
	companyBalanceHandler := handlers.NewCompanyBalanceHandler(companyBalanceService)

	tariffPlanRepo := repository.NewTariffPlanRepository(db)
	tariffPlanService := service.NewTariffPlanService(tariffPlanRepo)
	tariffPlanHandler := handlers.NewTariffPlanHandler(tariffPlanService)

	paymentRepo := repository.NewPaymentRequestRepository(db)
	paymentService := service.NewPaymentRequestService(paymentRepo, tariffPlanRepo, companyBalanceRepo)
	paymentHandler := handlers.NewPaymentRequestHandler(paymentService)

	adminRepo := repository.NewAdminRepository(db)
	adminService := service.NewAdminService(adminRepo)
	adminHandler := handlers.NewAdminHandler(adminService)

	smsService := service.NewSMSService(cfg.Mobizon.APIKey)
	waService := service.NewWhatsAppSMSC(
		cfg.SMSC.Login,
		cfg.SMSC.Password,
		cfg.SMSC.SenderWA,
	)
	smsHandler := handlers.NewSMSHandler(smsService, waService)

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
		companyBalanceHandler:      companyBalanceHandler,
		tariffPlanHandler:          tariffPlanHandler,
		paymentHandler:             paymentHandler,
		smsHandler:                 smsHandler,
		adminHandler:               adminHandler,
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
		//w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		//w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
		//w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")

		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		next.ServeHTTP(w, r)
	})
}
