package tenantized

import (
	"context"
	"fmt"
	"invoicing-example/transactional"
	"log"
	"net/http"
	"regexp"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"goyave.dev/goyave/v5"
)

var accountManagerRole = regexp.MustCompile("^account_manager_[0-9]+$")

type TenantAwareMiddleware struct {
	goyave.Component
	pool *pgxpool.Pool
}

func NewTenantAwareMiddleware(pool *pgxpool.Pool) *TenantAwareMiddleware {
	return &TenantAwareMiddleware{pool: pool}
}

func (m *TenantAwareMiddleware) Handle(next goyave.Handler) goyave.Handler {
	return func(response *goyave.Response, request *goyave.Request) {
		ctx := request.Context()

		tx, err := m.pool.Begin(ctx)
		if err != nil {
			response.Error(fmt.Sprintf("could not start transaction: %s", err.Error()))
			return
		}

		// Inject transaction into request context
		ctx = transactional.WithTransaction(ctx, tx)
		request = request.WithContext(ctx)

		// Apply user role and tenant to transaction
		if !applyRowLevelSecurity(tx, ctx, response, request) {
			return
		}

		next(response, request)

		status := response.GetStatus()
		if status >= 200 && status < 300 {
			if err := tx.Commit(context.Background()); err != nil {
				log.Println("commit failed:", err)
			}
		} else {
			if err := tx.Rollback(context.Background()); err != nil {
				log.Println("rollback failed:", err)
			}
		}
	}
}

func applyRowLevelSecurity(
	tx pgx.Tx, ctx context.Context,
	response *goyave.Response, request *goyave.Request,
) bool {

	userRole := request.Header().Get("X-User-Role")
	if userRole != "administrator" &&
		!accountManagerRole.MatchString(userRole) {
		userRole = "customer"
	}

	// Since we made sure the role is valid, we can safely use string format
	_, err := tx.Exec(ctx, fmt.Sprintf("SET ROLE %s", userRole))
	if err != nil {
		response.Error("failed to set tenant")
		return false
	}

	if userRole == "customer" {
		// Apply tenant since we assume a customer transaction
		if !applyTenantPermission(tx, ctx, response, request) {
			return false
		}
	}

	return true
}

func applyTenantPermission(
	tx pgx.Tx, ctx context.Context,
	response *goyave.Response, request *goyave.Request,
) bool {
	// Extract tenant ID from request
	// In the real world, this would want to be a JWT token
	tenant := request.Header().Get("X-Customer-ID")
	if tenant == "" {
		response.Status(http.StatusBadRequest)
		response.Error("missing X-Customer-ID header")
		return false
	}
	tenantId, err := uuid.Parse(tenant)
	if err != nil {
		response.Status(http.StatusBadRequest)
		response.Error("invalid X-Customer-ID")
		return false
	}

	// Since we made sure the UUID is valid, we can safely use string format
	_, err = tx.Exec(ctx,
		fmt.Sprintf("SET LOCAL app.current_customer_id = '%s'", tenantId.String()),
	)
	if err != nil {
		response.Error("failed to set tenant")
		return false
	}
	return true
}
