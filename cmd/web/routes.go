package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders, makeResponseJSON)

	//dynamicMiddleware := alice.New()

	mux := pat.New()

	mux.Post("/companies/register", standardMiddleware.ThenFunc(app.companyHandler.Register))
	mux.Post("/companies/login", standardMiddleware.ThenFunc(app.companyHandler.Login))
	mux.Get("/companies", standardMiddleware.ThenFunc(app.companyHandler.GetAll))
	mux.Get("/companies/id/:id", standardMiddleware.ThenFunc(app.companyHandler.GetByID))
	mux.Get("/companies/phone/:phone", standardMiddleware.ThenFunc(app.companyHandler.GetByPhone))
	mux.Put("/companies/:id", standardMiddleware.ThenFunc(app.companyHandler.Update))
	mux.Del("/companies/:id", standardMiddleware.ThenFunc(app.companyHandler.Delete))

	mux.Post("/templates/upload", standardMiddleware.ThenFunc(app.templateHandler.Upload))
	mux.Get("/templates/:id", standardMiddleware.ThenFunc(app.templateHandler.GetByID))
	mux.Get("/templates/company/:id", standardMiddleware.ThenFunc(app.templateHandler.GetByCompany))
	mux.Put("/templates/:id", standardMiddleware.ThenFunc(app.templateHandler.Update))
	mux.Del("/templates/:id", standardMiddleware.ThenFunc(app.templateHandler.Delete))
	mux.Get("/templates/pdf/:id", standardMiddleware.ThenFunc(app.templateHandler.ServePDFByID))

	mux.Post("/contracts", standardMiddleware.ThenFunc(app.contractHandler.Create))
	mux.Post("/contracts/with-fields", standardMiddleware.ThenFunc(app.contractHandler.CreateWithFields))
	mux.Get("/contracts/:id", standardMiddleware.ThenFunc(app.contractHandler.GetByID))
	mux.Get("/contracts/token/:token", standardMiddleware.ThenFunc(app.contractHandler.GetByToken))
	mux.Get("/contracts/token/:token/details", standardMiddleware.ThenFunc(app.contractHandler.GetByTokenWithFields))
	mux.Get("/contracts/company/:id", standardMiddleware.ThenFunc(app.contractHandler.GetByCompany))
	mux.Get("/contracts/pdf/:id", standardMiddleware.ThenFunc(app.contractHandler.ServePDFByID))
	mux.Put("/contracts/:id", standardMiddleware.ThenFunc(app.contractHandler.Update))
	mux.Del("/contracts/:id", standardMiddleware.ThenFunc(app.contractHandler.Delete))

	mux.Post("/signatures", standardMiddleware.ThenFunc(app.signatureHandler.Create))
	mux.Get("/signatures/:id", standardMiddleware.ThenFunc(app.signatureHandler.GetByID))
	mux.Get("/signatures/contract/:id", standardMiddleware.ThenFunc(app.signatureHandler.GetByContractID))
	mux.Get("/signatures/company/:id", standardMiddleware.ThenFunc(app.signatureHandler.GetContractsByCompanyID))
	mux.Get("/signatures/admin/all", standardMiddleware.ThenFunc(app.signatureHandler.GetSignaturesAll))
	mux.Get("/signatures/admin/status-summary", standardMiddleware.ThenFunc(app.signatureHandler.GetStatusSummary))
	mux.Put("/signatures/qr/:id", standardMiddleware.ThenFunc(app.signatureHandler.UpdateQR))
	mux.Del("/signatures/:id", standardMiddleware.ThenFunc(app.signatureHandler.Delete))
	mux.Post("/signatures/sign", standardMiddleware.ThenFunc(app.signatureHandler.Sign))
	mux.Get("/signatures/pdf/:id", standardMiddleware.ThenFunc(app.signatureHandler.ServeSignedPDFByID))

	mux.Post("/contract-fields", standardMiddleware.ThenFunc(app.contractFieldHandler.Create))
	mux.Get("/contract-fields/:id", standardMiddleware.ThenFunc(app.contractFieldHandler.GetByContractID))

	mux.Post("/signature-fields/bulk", standardMiddleware.ThenFunc(app.signatureFieldValueHandler.CreateAll))
	mux.Get("/signature-fields/signature/:id", standardMiddleware.ThenFunc(app.signatureFieldValueHandler.GetBySignatureID))

	mux.Get("/stats/company/:id", standardMiddleware.ThenFunc(app.statisticsHandler.GetCompanyStats))

	mux.Post("/company-balances", standardMiddleware.ThenFunc(app.companyBalanceHandler.Create))
	mux.Get("/company-balances/:id", standardMiddleware.ThenFunc(app.companyBalanceHandler.GetByCompanyID))
	mux.Put("/company-balances/:id", standardMiddleware.ThenFunc(app.companyBalanceHandler.Update))
	mux.Del("/company-balances/:id", standardMiddleware.ThenFunc(app.companyBalanceHandler.Delete))

	mux.Post("/tariff-plans", standardMiddleware.ThenFunc(app.tariffPlanHandler.Create))
	mux.Get("/tariff-plans", standardMiddleware.ThenFunc(app.tariffPlanHandler.GetAll))
	mux.Get("/tariff-plans/:id", standardMiddleware.ThenFunc(app.tariffPlanHandler.GetByID))
	mux.Put("/tariff-plans/:id", standardMiddleware.ThenFunc(app.tariffPlanHandler.Update))
	mux.Del("/tariff-plans/:id", standardMiddleware.ThenFunc(app.tariffPlanHandler.Delete))

	mux.Post("/payment_requests", standardMiddleware.ThenFunc(app.paymentHandler.Create))
	mux.Get("/payment_requests/:id", standardMiddleware.ThenFunc(app.paymentHandler.GetByID))
	mux.Get("/payment_requests/company/:id", standardMiddleware.ThenFunc(app.paymentHandler.GetByCompany))
	mux.Get("/payment_requests/admin/all", standardMiddleware.ThenFunc(app.paymentHandler.GetAll))
	mux.Put("/payment_requests/:id", standardMiddleware.ThenFunc(app.paymentHandler.Update))
	mux.Del("/payment_requests/:id", standardMiddleware.ThenFunc(app.paymentHandler.Delete))

	mux.Post("/sms/send", standardMiddleware.ThenFunc(app.smsHandler.Send))

	return standardMiddleware.Then(mux)
}
