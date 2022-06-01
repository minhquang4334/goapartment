package goapartment

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	// ErrDBIsRequired は sqlx.DBを期待するのに、nilが渡された
	ErrDBIsRequired = errors.New("sqlx.DB is required")

	// ErrTenantIsRequired is a error type, that raise when expect tenant name as parameter but recieved empty string.
	ErrTenantIsRequired = errors.New("Tenant is required")
)

// Apartment structure have a implemented of AparmentDB interface
type Apartment struct {
	DB *sqlx.DB
}

// ProvideApartment initialize a Apartment struct and return it
func ProvideApartment(db *sqlx.DB) (*Apartment, error) {
	if db == nil {
		return nil, ErrDBIsRequired
	}
	return &Apartment{
		DB: db,
	}, nil
}

type (
	TxHandler func(context.Context, *sqlx.Tx) error

	ConnHandler func(context.Context, *sqlx.Conn) error
)

// TenantExecTx open an transaction on tenant DB and run query handler on it. If error is raise when execute query, an error must be returned
func (ap *Apartment) TenantExecTx(ctx context.Context, tenant string, handler TxHandler) error {
	if tenant == "" {
		return ErrTenantIsRequired
	}
	tx, err := ap.DB.BeginTxx(ctx, nil)
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

// TenantExecConn connect to tenant database and run query handler on it. If error is raise when execute query, an error must be returned
func (ap *Apartment) TenantExecConn(ctx context.Context, tenant string, handler ConnHandler) error {
	if tenant == "" {
		return ErrTenantIsRequired
	}
	conn, err := ap.DB.Connx(ctx)
	if err != nil {
		return err
	}
	query := fmt.Sprintf("USE %s", tenant)
	if _, err = conn.ExecContext(ctx, query); err != nil {
		return err
	}
	if err = handler(ctx, conn); err != nil {
		return err
	}
	return nil
}
