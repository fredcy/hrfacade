package hrfacade

import (
	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
	"regexp"
	"strings"
)

func PersonnelCount(dsn string) (int, error) {
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
	Empno      string
	Fname      string
	Mi         string
	Lname      string
	Active     string
	Jobtitle   string
	Homephone  string
	Busphone   string
	Cellphone  string
	Faxphone   string
	Pagerphone string
	Email      string
	Level1     string
	Level2     string
	Level3     string
	Level4     string
	Superno    string
}

// GetContacts reads from the HR database and returns a channel of Contact structs
func GetContacts(dsn string, all bool) (chan Contact, error) {
	cs := make(chan Contact)
	db, err := sql.Open("mssql", dsn)
	if err != nil {
		log.Printf("ERROR: cannot open DSN = %v", dsn)
		return cs, err
	}
	defer db.Close()

	var whereClause string
	if !all {
		whereClause = " where p_active = 'A' "
	}
	q := `select p_empno, p_active, p_fname, p_mi, p_lname, p_jobtitle
 , p_hphone, p_busphone, p_cellular, p_empfax, p_pager, p_empemail
 , p_level1, p_level2, p_level3, p_level4, p_superno
 from hrpersnl ` + whereClause +
		` order by lower(p_lname), lower(p_fname)`
	rows, err := db.Query(q)
	if err != nil {
		log.Printf("ERROR: query failed: %v", q)
		return cs, err
	}

	go func() {
		defer rows.Close()
		defer close(cs)

		for rows.Next() {
			c := Contact{}
			var empno, active, fname, mi, lname, jobtitle, homephone, busphone, cellphone, faxphone, pagerphone, email sql.NullString
			var level1, level2, level3, level4, superno sql.NullString
			err := rows.Scan(&empno, &active, &fname, &mi, &lname,
				&jobtitle, &homephone, &busphone, &cellphone, &faxphone, &pagerphone, &email,
				&level1, &level2, &level3, &level4, &superno)
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
			c.Email = normalize(email)
			c.Level1 = normalize(level1)
			c.Level2 = normalize(level2)
			c.Level3 = normalize(level3)
			c.Level4 = normalize(level4)
			c.Superno = normalize(superno)
			cs <- c
		}
	}()
	return cs, nil
}

func normalize(s sql.NullString) string {
	if !s.Valid {
		return ""
	}
	return strings.TrimSpace(s.String)
}

var phoneJunkRe = regexp.MustCompile(`[() -]`)

// phonecanon removes the extraneous (non-numeric) characters from a phone number.
// E.g., "(800) 555-1212" becomes "8005551212".
func phonecanon(s string) string {
	return phoneJunkRe.ReplaceAllLiteralString(s, "")
}
