package transactional

import (
    "context"

    "github.com/jackc/pgx/v5"
)

var transactionKey struct{}

func WithTransaction(ctx context.Context, tx pgx.Tx) context.Context {
    return context.WithValue(ctx, transactionKey, tx)
}

func FromContext(ctx context.Context) pgx.Tx {
    tx, _ := ctx.Value(transactionKey).(pgx.Tx)
    return tx
}
