# Goapartment
goapartment implements multi-tenancy database connection for golang

# Installation
```bash
go get github.com/minhquang4334/goapartment
```

# Usage
```go
// mysql
db, err := sqlx.Open("mysql", dsn)
if err != nil {
  return nil, err
}
apartment := goapartment.ProvideApartment(db)
apartment.TenantExecConn(ctx, "tenantName", func(ctx context.Context, conn *sqlx.Conn) error {
  query := "SELECT DATABASE()"
  row := conn.QueryRowContext(ctx, query)
  var dbName string
  if err := row.Scan(&dbName); err != nil {
    return "", err
  }
  return dbName, nil
}
```
# License
Distributed under the MIT License. See LICENSE for more information.

# Acknowledgements
I acknowledge gratitude to the following resources in creating this repo
- @aereal with
  - [waitmysql](https://github.com/aereal/waitmysql)
  - [paramsenc](https://github.com/aereal/paramsenc)

