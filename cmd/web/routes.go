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

	mux.Post("/contracts", standardMiddleware.ThenFunc(app.contractHandler.Create))
	mux.Get("/contracts/:id", standardMiddleware.ThenFunc(app.contractHandler.GetByID))
	mux.Get("/contracts/token/:token", standardMiddleware.ThenFunc(app.contractHandler.GetByToken))
	mux.Get("/contracts/company/:id", standardMiddleware.ThenFunc(app.contractHandler.GetByCompany))
	mux.Put("/contracts/:id", standardMiddleware.ThenFunc(app.contractHandler.Update))
	mux.Del("/contracts/:id", standardMiddleware.ThenFunc(app.contractHandler.Delete))

	mux.Post("/signatures", standardMiddleware.ThenFunc(app.signatureHandler.Create))
	mux.Get("/signatures/:id", standardMiddleware.ThenFunc(app.signatureHandler.GetByID))
	mux.Get("/signatures/contract/:id", standardMiddleware.ThenFunc(app.signatureHandler.GetByContractID))
	mux.Del("/signatures/:id", standardMiddleware.ThenFunc(app.signatureHandler.Delete))

	mux.Post("/contract-fields", standardMiddleware.ThenFunc(app.contractFieldHandler.Create))
	mux.Get("/contract-fields/:id", standardMiddleware.ThenFunc(app.contractFieldHandler.GetByContractID))

	mux.Post("/signature-fields/bulk", standardMiddleware.ThenFunc(app.signatureFieldValueHandler.CreateAll))
	mux.Get("/signature-fields/signature/:id", standardMiddleware.ThenFunc(app.signatureFieldValueHandler.GetBySignatureID))

	mux.Get("/stats/company/:id", standardMiddleware.ThenFunc(app.statisticsHandler.GetCompanyStats))

	return standardMiddleware.Then(mux)
}
