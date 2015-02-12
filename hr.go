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

type Contact struct {
	Empno string
	Fname string
	Mi string
	Lname string
	Active string
	Jobtitle string
}

// GetContacts reads from the HR database and returns a channel of Contact structs
func GetContacts() (chan struct { Contact; Error error }, error) {
	cs := make(chan struct { Contact; Error error })
	dsn := os.Getenv("HR_DSN")
	db, err := sql.Open("mssql", dsn)
	if err != nil {
		log.Printf("ERROR: cannot open DSN = %v", dsn)
		return cs, err
	}
	defer db.Close()

	q := `select p_empno, p_active, p_fname, p_mi, p_lname
 , p_jobtitle
 from hrpersnl
 where p_active = 'A' order by p_lname, p_fname`
	rows, err := db.Query(q)
	if err != nil {
		log.Printf("ERROR: query failed")
		return cs, err
	}

	go func() {
		defer rows.Close()
		defer close(cs)
		
		for rows.Next() {
			c := Contact{}
			err := rows.Scan(&c.Empno, &c.Active, &c.Fname, &c.Mi, &c.Lname,
				&c.Jobtitle)
			cs <- struct { Contact; Error error }{ c, err }
		}
	}()
	return cs, nil
}
