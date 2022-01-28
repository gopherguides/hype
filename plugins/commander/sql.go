package commander

import _ "embed"

//go:embed sql/migration.sql
var migrationSQL string

//go:embed sql/fetch.sql
var fetchSQL string

//go:embed sql/insert.sql
var insertSQL string
