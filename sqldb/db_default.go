//go:build !static
// +build !static

package sqldb

// Go database API http://go-database-sql.org/overview.html

// CGO based database drivers:
// Oracle drivers:
//   _ "gopkg.in/rana/ora.v4"
//   _ "github.com/mattn/go-oci8"
// MySQL driver:
//   _ "github.com/go-sql-driver/mysql"
// SQLite driver:
//  _ "github.com/mattn/go-sqlite3"

import (
	"database/sql"

	// mysql CGO driver
	_ "github.com/go-sql-driver/mysql"
	// sqlite CGO driver:
	_ "github.com/mattn/go-sqlite3"
)

func dbOpen(dbtype, dburi string) (*sql.DB, error) {
	return sql.Open(dbtype, dburi)
}
