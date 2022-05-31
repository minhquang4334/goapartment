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

	// ErrTenantIsRequired はtenantの名前を期待するのに、空文字列が渡された
	ErrTenantIsRequired = errors.New("Tenant is required")
)

// Apartment は各Tenantに接続する情報をもつ構造体です
type Apartment struct {
	db *sqlx.DB
}

// ProvideApartment はApartmentを生成する関数です
func ProvideApartment(db *sqlx.DB) (*Apartment, error) {
	if db == nil {
		return nil, ErrDBIsRequired
	}
	return &Apartment{
		db: db,
	}, nil
}

type queryHandler func(context.Context, *sqlx.Tx) error

// TenantExec はTenantのデータベスにアクセスしてQueryを実行するメソッドです
func (ap *Apartment) TenantExec(ctx context.Context, tenant string, handler queryHandler) error {
	if tenant == "" {
		return ErrTenantIsRequired
	}
	tx, err := ap.db.BeginTxx(ctx, nil)
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
