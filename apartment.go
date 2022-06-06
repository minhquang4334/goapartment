package goapartment

import (
	"context"
	"database/sql"
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

// TenantExecTx open an transaction on tenant DB and return a instance of *sqlx.Tx
// when error is raised, an error must be returned
func (ap *Apartment) TenantExecTx(ctx context.Context, tenant string, txOptions *sql.TxOptions) (*sqlx.Tx, error) {
	if tenant == "" {
		return nil, ErrTenantIsRequired
	}
	tx, err := ap.DB.BeginTxx(ctx, txOptions)
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf("USE %s", tenant)
	if _, err = tx.ExecContext(ctx, query); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	return tx, nil
}

// TenantExecConn open an connection on tenant DB and return a instance of *sqlx.Conn
// when error is raised, an error must be returned
func (ap *Apartment) TenantExecConn(ctx context.Context, tenant string) (*sqlx.Conn, error) {
	if tenant == "" {
		return nil, ErrTenantIsRequired
	}
	conn, err := ap.DB.Connx(ctx)
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf("USE %s", tenant)
	if _, err = conn.ExecContext(ctx, query); err != nil {
		return nil, err
	}
	return conn, nil
}
