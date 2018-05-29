// Handles connection to the database server.
package dbutil

import (
	"database/sql"
	_ "github.com/tsealex/pq"
	"github.com/jmoiron/sqlx"
)

var Instance *sqlx.DB

func init() {
	instance, err := sql.Open("postgres",
		"dbname=postgres host=localhost port=6000 sslmode=disable")
	if err != nil {
		panic(err)
	}
	Instance = sqlx.NewDb(instance, "postgres")
}
