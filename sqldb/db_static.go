//go:build static
// +build static

package sqldb

// Go database API: http://go-database-sql.org/overview.html
// SQLite non-CGO drivers (to make static executable):
//  _ "github.com/glebarez/go-sqlite"
//  _ "gitlab.com/cznic/sqlite"
//  _ "modernc.org/sqlite"

import (
	"database/sql"

	// non-CGO sqlite driver
	_ "github.com/glebarez/go-sqlite"
)

func dbOpen(dbtype, dburi string) (*sql.DB, error) {
	return sql.Open(dbtype, dburi)
}
