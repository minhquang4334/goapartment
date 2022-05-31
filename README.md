# Goapartment
goapartment implements multi-tenancy database connection for golang

# Installation
```bash
go get github.com/minhquang4334/goapartment
```

# Usage
```go
// mysql
db, err := sql.Open("mysql", dsn)
if err != nil {
  return nil, err
}
apartment := goapartment.ProvideApartment(db)
apartment.TenantExec(ctx, "tenantName", func(ctx context.Context, tx *sql.Tx) error {
  query := "SELECT DATABASE()"
  row := tx.QueryRow(query)
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

