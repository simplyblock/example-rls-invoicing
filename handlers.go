package main

import (
	"encoding/json"
	"fmt"
	"invoicing-example/transactional"
	"net/http"
	"time"

	"github.com/google/uuid"
	"goyave.dev/goyave/v5"
)

func listCustomers(response *goyave.Response, request *goyave.Request) {
	ctx := request.Context()
	tx := transactional.FromContext(ctx)

	rows, err := tx.Query(ctx, `
		SELECT customer_id, customer_name FROM customers`,
	)
	if err != nil {
		response.Error(fmt.Sprintf("db query error: %s", err.Error()))
		return
	}
	defer rows.Close()

	customers := make([]Customer, 0)
	for rows.Next() {
		var customer Customer
		if err := rows.Scan(&customer.Id, &customer.Name); err != nil {
			response.Error(fmt.Sprintf("db scan error: %s", err.Error()))
		}
		customers = append(customers, customer)
	}

	response.JSON(http.StatusOK, customers)
}

func showCustomer(response *goyave.Response, request *goyave.Request) {
	ctx := request.Context()
	tx := transactional.FromContext(ctx)

	customerId, err := uuid.Parse(request.RouteParams["customerId"])
	if err != nil {
		response.Status(http.StatusBadRequest)
		response.Error("empty or invalid customer id")
		return
	}

	rows, err := tx.Query(ctx, `
		SELECT customer_id, customer_name
		FROM customers
		WHERE customer_id = $1`,
		customerId.String(),
	)
	if err != nil {
		response.Error(fmt.Sprintf("db query error: %s", err.Error()))
		return
	}
	defer rows.Close()

	var customer Customer
	if rows.Next() {
		if err := rows.Scan(&customer.Id, &customer.Name); err != nil {
			response.Error(fmt.Sprintf("db scan error: %s", err.Error()))
		}
	}

	response.JSON(http.StatusOK, customer)
}

func createCustomer(response *goyave.Response, request *goyave.Request) {
	ctx := request.Context()
	tx := transactional.FromContext(ctx)

	var customer Customer
	if err := json.NewDecoder(request.Body()).Decode(&customer); err != nil {
		response.Status(http.StatusBadRequest)
		response.Error(fmt.Sprintf("invalid request body: %s", err.Error()))
		return
	}
	customer.Id = uuid.New()

	_, err := tx.Exec(ctx, `
		INSERT INTO customers (
			customer_id, customer_name
		) VALUES ($1, $2)`,
		customer.Id.String(), customer.Name,
	)
	if err != nil {
		response.Error(fmt.Sprintf("db query error: %s", err.Error()))
		return
	}

	response.JSON(http.StatusOK, customer)
}

func listInvoices(response *goyave.Response, request *goyave.Request) {
	ctx := request.Context()
	tx := transactional.FromContext(ctx)

	rows, err := tx.Query(ctx, `
		SELECT invoice_id, customer_id, invoice_date FROM invoices`,
	)
	if err != nil {
		response.Error(fmt.Sprintf("db query error: %s", err.Error()))
		return
	}
	defer rows.Close()

	invoices := make([]Invoice, 0)
	for rows.Next() {
		var inv Invoice
		if err := rows.Scan(&inv.Id, &inv.CustomerId, &inv.Date); err != nil {
			response.Error(fmt.Sprintf("db scan error: %s", err.Error()))
			return
		}

		invoices = append(invoices, inv)
	}

	for i := range invoices {
		posRows, err := tx.Query(ctx, `
			SELECT invoice_position_id, name, description, amount, price FROM invoice_positions
			WHERE invoice_id = $1
			ORDER BY position ASC`,
			invoices[i].Id,
		)
		if err != nil {
			response.Error(fmt.Sprintf("db query error: %s", err.Error()))
			return
		}

		var positions []InvoicePosition
		for posRows.Next() {
			var pos InvoicePosition
			if err := posRows.Scan(&pos.Id, &pos.Name, &pos.Desc, &pos.Amount, &pos.Price); err != nil {
				posRows.Close()
				response.Error(fmt.Sprintf("db scan error: %s", err.Error()))
				return
			}
			positions = append(positions, pos)
		}
		posRows.Close()
		invoices[i].Positions = positions
	}

	response.JSON(http.StatusOK, invoices)
}

func showInvoice(response *goyave.Response, request *goyave.Request) {
	ctx := request.Context()
	tx := transactional.FromContext(ctx)

	customerId, err := uuid.Parse(request.RouteParams["customerId"])
	if err != nil {
		response.Status(http.StatusBadRequest)
		response.Error("empty or invalid customer id")
		return
	}

	invoiceId, err := uuid.Parse(request.RouteParams["invoiceId"])
	if err != nil {
		response.Status(http.StatusBadRequest)
		response.Error("empty or invalid invoice id")
		return
	}

	rows, err := tx.Query(ctx, `
		SELECT invoice_id, customer_id, invoice_date
		FROM invoices
		WHERE customer_id = $1 AND invoice_id = $2`,
		customerId.String(), invoiceId.String(),
	)
	if err != nil {
		response.Error(fmt.Sprintf("db query error: %s", err.Error()))
		return
	}
	defer rows.Close()

	var invoice Invoice
	if rows.Next() {
		if err := rows.Scan(&invoice.Id, &invoice.CustomerId, &invoice.Date); err != nil {
			response.Error(fmt.Sprintf("db scan error: %s", err.Error()))
			return
		}
	}

	if invoice.Id == uuid.Nil {
		response.Status(http.StatusNotFound)
		response.Error("invoice not found")
	}

	posRows, err := tx.Query(ctx, `
			SELECT invoice_position_id, name, description, amount, price FROM invoice_positions
			WHERE invoice_id = $1
			ORDER BY position ASC`,
		invoice.Id,
	)
	if err != nil {
		response.Error(fmt.Sprintf("db query error: %s", err.Error()))
		return
	}

	var positions []InvoicePosition
	for posRows.Next() {
		var pos InvoicePosition
		if err := posRows.Scan(&pos.Id, &pos.Name, &pos.Desc, &pos.Amount, &pos.Price); err != nil {
			posRows.Close()
			response.Error(fmt.Sprintf("db scan error: %s", err.Error()))
			return
		}
		positions = append(positions, pos)
	}
	posRows.Close()
	invoice.Positions = positions

	response.JSON(http.StatusOK, invoice)
}

func createInvoice(response *goyave.Response, request *goyave.Request) {
	ctx := request.Context()
	tx := transactional.FromContext(ctx)

	var invoice Invoice
	if err := json.NewDecoder(request.Body()).Decode(&invoice); err != nil {
		response.Status(http.StatusBadRequest)
		response.Error(fmt.Sprintf("invalid request body: %s", err.Error()))
		return
	}
	invoice.Id = uuid.New()
	invoice.Date = time.Now()

	_, err := tx.Exec(ctx, `
		INSERT INTO invoices (
			invoice_id, customer_id, invoice_date
		) VALUES ($1, $2, $3)`,
		invoice.Id, invoice.CustomerId, invoice.Date,
	)
	if err != nil {
		response.Error(fmt.Sprintf("db insert error: %s", err.Error()))
		return
	}

	for i, item := range invoice.Positions {
		item.Id = uuid.New()
		_, err := tx.Exec(ctx, `
			INSERT INTO invoice_positions (
            	invoice_position_id, invoice_id, customer_id, name, description, amount, price, position
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			item.Id, invoice.Id, invoice.CustomerId, item.Name, item.Desc, item.Amount, item.Price, i+1,
		)
		if err != nil {
			response.Error(fmt.Sprintf("db insert error: %s", err.Error()))
			return
		}
	}

	response.JSON(http.StatusCreated, invoice)
}
