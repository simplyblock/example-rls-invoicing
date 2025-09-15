package main

import (
	"invoicing-example/tenantized"

	"github.com/jackc/pgx/v5/pgxpool"
	"goyave.dev/goyave/v5"
)

func createRoutes(pool *pgxpool.Pool) func(
	server *goyave.Server, router *goyave.Router,
) {
	return func(server *goyave.Server, router *goyave.Router) {
		router = router.Middleware(tenantized.NewTenantAwareMiddleware(pool))

		router.Get("/customers", listCustomers)
		router.Post("/customers", createCustomer)
		subrouter := router.Subrouter("/customers/{customerId}")
		subrouter.Get("/", showCustomer)
		subrouter.Get("/invoices", listInvoices)
		subrouter.Post("/invoices", createInvoice)
		subrouter.Get("/invoices/{invoiceId}", showInvoice)
	}
}
