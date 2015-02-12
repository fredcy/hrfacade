package hrfacade

import (
	"database/sql"
	"log"
	"os"
	_ "github.com/denisenkom/go-mssqldb"
)


func persnl_count() (int, error) {
	dsn := os.Getenv("HR_DSN")
	log.Printf("dsn = %v", dsn)
	db, err := sql.Open("mssql", dsn)
	if err != nil {
		log.Printf("ERROR: cannot open DSN = %v", dsn)
		return 0, err
	}
	defer db.Close()
	q := "select count(*) c from hrpersnl"

	rows, err := db.Query(q)
	if err != nil {
		log.Printf("ERROR: query failed")
		return 0, err
	}
	log.Printf("rows = %v", rows)
	return 99, nil
}
