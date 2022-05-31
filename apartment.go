package goapartment

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var (
	// ErrTenantIsRequired is a error type, that raise when expect tenant name as parameter but recieved empty string.
	ErrTenantIsRequired = errors.New("Tenant is required")
)

// ApartmentDB is a interface that provide database query transaction method
type ApartmentDB interface {
	// BeginTx begins a transaction and returns an *sql.Tx
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

// Apartment structure have a implemented of AparmentDB interface
type Apartment struct {
	DB ApartmentDB
}

// ProvideApartment initialize a Apartment struct and return it
func ProvideApartment(db ApartmentDB) *Apartment {
	return &Apartment{ DB: db }
}

// QueryHandler is a query handler that will run on tenant database. If error is raise when execute query, an error must be returned and transaction will be rollback
type QueryHandler func(context.Context, *sql.Tx) error

// TenantExec connect to tenant database and run query handler on it. If error is raise when execute query, an error must be returned
func (ap *Apartment) TenantExec(ctx context.Context, tenant string, handler QueryHandler) error {
	if tenant == "" {
		return ErrTenantIsRequired
	}
	tx, err := ap.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	query := fmt.Sprintf("USE %s", tenant)
	if _, err = tx.ExecContext(ctx, query); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err = handler(ctx, tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	_ = tx.Commit()

	return nil
}
