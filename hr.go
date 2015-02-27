package hrfacade

import (
	"database/sql"
	"log"
	"os"
	_ "github.com/denisenkom/go-mssqldb"
	"regexp"
	"strings"
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
	Homephone string
	Busphone string
	Cellphone string
	Faxphone string
	Pagerphone string
}

// GetContacts reads from the HR database and returns a channel of Contact structs
func GetContacts() (chan Contact, error) {
	cs := make(chan Contact)
	dsn := os.Getenv("HR_DSN")
	db, err := sql.Open("mssql", dsn)
	if err != nil {
		log.Printf("ERROR: cannot open DSN = %v", dsn)
		return cs, err
	}
	defer db.Close()

	q := `select p_empno, p_active, p_fname, p_mi, p_lname, p_jobtitle
 , p_hphone, p_busphone, p_cellular, p_empfax, p_pager
 from hrpersnl
 where p_active = 'A' order by lower(p_lname), lower(p_fname)`
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
			var empno, active, fname, mi, lname, jobtitle, homephone, busphone, cellphone, faxphone, pagerphone sql.NullString
			err := rows.Scan(&empno, &active, &fname, &mi, &lname,
				&jobtitle, &homephone, &busphone, &cellphone, &faxphone, &pagerphone)
			if err != nil {
				log.Printf("ERROR: %v", err)
				continue
			}
			c.Empno = normalize(empno)
			c.Active = normalize(active)
			c.Fname = normalize(fname)
			c.Mi = normalize(mi)
			c.Lname = normalize(lname)
			c.Jobtitle = normalize(jobtitle)
			c.Homephone = phonecanon(normalize(homephone))
			c.Busphone = phonecanon(normalize(busphone))
			c.Cellphone = phonecanon(normalize(cellphone))
			c.Faxphone = phonecanon(normalize(faxphone))
			c.Pagerphone = phonecanon(normalize(pagerphone))
			cs <- c
		}
	}()
	return cs, nil
}

func normalize(s sql.NullString) string {
	if ! s.Valid {
		return ""
	}
	return strings.TrimSpace(s.String)
}

var phoneJunkRe = regexp.MustCompile(`[() -]`)

func phonecanon(s string) string {
	//return "{{" + phoneRe.ReplaceAllString(s, "$1$2$3") + "}}"
	return phoneJunkRe.ReplaceAllLiteralString(s, "")
}
