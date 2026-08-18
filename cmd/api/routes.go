package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// routes returns the HTTP router configured with all handlers, route-specific middleware,
// and global middleware.
func (app *application) routes() http.Handler {
	router := httprouter.New()

	// BACKEND

	// Defined handlers for 404 and 205 status code
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// Healthcheck route
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	// DATABASE SCHEMA ROUTES

	// consumer routes
	router.HandlerFunc(http.MethodPost, "/v1/consumers", app.createConsumerHandler)
	router.HandlerFunc(http.MethodGet, "/v1/consumers/:id", app.showConsumerHandler)
	router.HandlerFunc(http.MethodPut, "/v1/consumers/:id", app.updateConsumerHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/consumers/:id", app.deleteConsumerHandler)

	// API key routes
	router.HandlerFunc(http.MethodPost, "/v1/api-keys", app.createAPIKeyHandler)
	router.HandlerFunc(http.MethodGet, "/v1/api-keys/:id", app.showAPIKeyHandler)
	router.HandlerFunc(http.MethodPut, "/v1/api-keys/:id", app.updateAPIKeyHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/api-keys/:id", app.deleteAPIKeyHandler)

	// job routes
	router.HandlerFunc(http.MethodPost, "/v1/jobs", app.createJobHandler)
	router.HandlerFunc(http.MethodGet, "/v1/jobs/:public_id", app.showJobHandler)
	router.HandlerFunc(http.MethodPut, "/v1/jobs/:public_id", app.updateJobHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/jobs/:id", app.deleteJobHandler)

	// consumer activity report routes
	router.HandlerFunc(http.MethodPost, "/v1/reports", app.createReportHandler)

	// GLOBAL MIDDLEWARE

	return app.requestLogger(
		router, // last middleware
	)
}
