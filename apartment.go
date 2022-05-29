package goapartment

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	ErrDBIsRequired = errors.New("sqlx.DB is required")

	ErrTenantIsRequired = errors.New("Tenant is required")
)

type Apartment struct {
	db *sqlx.DB
}

func ProvideApartment(db *sqlx.DB) (*Apartment, error) {
	if db == nil {
		return nil, ErrDBIsRequired
	}
	return &Apartment{
		db: db,
	}, nil
}

func (ap *Apartment) SwitchTenant(ctx context.Context, tenant string) (*sqlx.Tx, error) {
	if tenant == "" {
		return nil, ErrTenantIsRequired
	}
	tx, err := ap.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf("USE %s", tenant)
	_, err = tx.ExecContext(ctx, query)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	return tx, nil
}
