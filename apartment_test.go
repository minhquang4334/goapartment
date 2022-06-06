package goapartment

import (
	"context"
	"fmt"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
)

func setupDB() error {
	db, err := openDB()
	if err != nil {
		return fmt.Errorf("failed to open DB: %w", err)
	}
	q1 := `CREATE DATABASE IF NOT EXISTS tenant_1_test`
	q2 := `CREATE DATABASE IF NOT EXISTS tenant_2_test`
	q3 := `CREATE DATABASE IF NOT EXISTS tenant_3_test`
	if _, err = db.Exec(q1); err != nil {
		return fmt.Errorf("failed to create database tenant1: %w", err)
	}
	if _, err = db.Exec(q2); err != nil {
		return fmt.Errorf("failed to create database tenant2: %w", err)
	}
	if _, err = db.Exec(q3); err != nil {
		return fmt.Errorf("failed to create database tenant3: %w", err)
	}
	return nil
}

func TestMain(m *testing.M) {
	if err := setupDB(); err != nil {
		fmt.Fprintf(os.Stderr, "! %+v\n", err)
		os.Exit(2)
	}
	os.Exit(m.Run())
}

func TestProvideApartment(t *testing.T) {
	dummyDB, err := openDB()
	if err != nil {
		t.Fatalf("can not open db: %v", err)
	}

	testCases := []struct {
		name    string
		db      *sqlx.DB
		wantErr bool
	}{
		{
			"ok with existing sqlx db",
			dummyDB,
			false,
		},
		{
			"ok with not existed sqlx db",
			nil,
			true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ProvideApartment(tc.db)
			gotErr := err != nil
			if gotErr != tc.wantErr {
				t.Fatalf("wantErr=%v but got=%v", tc.wantErr, gotErr)
			}
		})
	}
}

var (
	testCases = []struct {
		name    string
		tenant  string
		wantErr bool
	}{
		{
			"false with empty tenant name",
			"",
			true,
		},
		{
			"false with not existed tenant name",
			"not_existed_tenant_database",
			true,
		},
		{
			"true with created tenant 1",
			"tenant_1_test",
			false,
		},
		{
			"true with created tenant 2",
			"tenant_2_test",
			false,
		},
		{
			"true with created tenant 3",
			"tenant_3_test",
			false,
		},
	}
)

func TestTenantExecTx(t *testing.T) {
	db, err := openDB()
	if err != nil {
		t.Fatalf("can not open db: %v", err)
	}
	apartment := Apartment{
		DB: db,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			ctx := context.Background()
			tx, err := apartment.TenantExecTx(ctx, tc.tenant)
			gotErr := err != nil
			if gotErr != tc.wantErr {
				t.Fatalf("wantErr=%v but gotErr=%v, err=%v", tc.wantErr, gotErr, err)
			}
			if !tc.wantErr {
				gotTenant, err := currentTenantTx(tx)
				if err != nil {
					t.Fatalf("can not execute query on current tenant: %v", err)
				}
				if diff := cmp.Diff(tc.tenant, gotTenant); diff != "" {
					t.Errorf("-want, +got:\n%s", diff)
				}
			}
		})
	}
}

func TestTenantExecConn(t *testing.T) {
	db, err := openDB()
	if err != nil {
		t.Fatalf("can not open db: %v", err)
	}
	apartment := Apartment{
		DB: db,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			ctx := context.Background()
			conn, err := apartment.TenantExecConn(ctx, tc.tenant)
			gotErr := err != nil
			if gotErr != tc.wantErr {
				t.Fatalf("wantErr=%v but gotErr=%v, err=%v", tc.wantErr, gotErr, err)
			}
			if !tc.wantErr {
				gotTenant, err := currentTenantConn(conn)
				if err != nil {
					t.Fatalf("can not execute query on current tenant: %v", err)
				}
				if diff := cmp.Diff(tc.tenant, gotTenant); diff != "" {
					t.Errorf("-want, +got:\n%s", diff)
				}
			}
		})
	}
}

func currentTenantTx(tx *sqlx.Tx) (string, error) {
	query := "SELECT DATABASE()"
	row := tx.QueryRow(query)
	var dbName string
	if err := row.Scan(&dbName); err != nil {
		return "", err
	}
	return dbName, nil
}

func currentTenantConn(conn *sqlx.Conn) (string, error) {
	query := "SELECT DATABASE()"
	row := conn.QueryRowContext(context.Background(), query)
	var dbName string
	if err := row.Scan(&dbName); err != nil {
		return "", err
	}
	return dbName, nil
}

func openDB() (*sqlx.DB, error) {
	dsn := "root@tcp(127.0.0.1:3306)/"
	db, err := sqlx.Open("mysql", dsn)
	return db, err
}
