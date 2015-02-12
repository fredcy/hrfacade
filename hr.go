package hrfacade

import (
	"database/sql"
	"log"
	"os"
	_ "github.com/denisenkom/go-mssqldb"
)


func PersonnelCount() (int, error) {
	dsn := os.Getenv("HR_DSN")
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

	var count int
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	}
	return count, nil
}
