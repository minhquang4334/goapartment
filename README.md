# Goapartment
goapartment implements multi-tenancy database connection for golang

# Installation
```bash
go get github.com/minhquang4334/goapartment
```

# Usage
## Open Tenant Connection
```go
// mysql
db, err := sqlx.Open("mysql", dsn)
if err != nil {
  return "", err
}
apartment := goapartment.ProvideApartment(db)
conn, err := apartment.TenantExecConn(ctx, "tenantName")
if err != nil {
  return "", errors.New("can not open sqlx.Conn")
}
query := "SELECT DATABASE()"
row := conn.QueryRowContext(ctx, query)
var dbName string
if err := row.Scan(&dbName); err != nil {
  return "", err
}
return dbName, nil
```

## Open Tenant Transaction
```go
// mysql
db, err := sqlx.Open("mysql", dsn)
if err != nil {
  return "", err
}
apartment := goapartment.ProvideApartment(db)
tx, err := apartment.TenantExecTx(ctx, "tenantName", nil)
if err != nil {
  return "", errors.New("can not open sqlx.Conn")
}
query := "SELECT DATABASE()"
row := tx.QueryRowContext(ctx, query)
var dbName string
if err := row.Scan(&dbName); err != nil {
  _ = tx.Rollback()
  return "", err
}
if err = tx.Commit(); err != nil {
  return "", err
}
return dbName, nil
```
# License
Distributed under the MIT License. See LICENSE for more information.

# Acknowledgements
I acknowledge gratitude to the following resources in creating this repo
- @aereal with
  - [waitmysql](https://github.com/aereal/waitmysql)
  - [paramsenc](https://github.com/aereal/paramsenc)

