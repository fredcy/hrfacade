package main

import (
	"flag"
	"fmt"
	"github.com/fredcy/hrfacade"
	"net/http"
	"log"
)

var address = flag.String("address", ":8080", "Listen and serve at this address")

func counthandler(w http.ResponseWriter, r *http.Request) {
	c, _ := hrfacade.PersonnelCount()
	fmt.Fprintf(w, "%d", c)
}

// contacthandler displays a line for each contact
func contacthandler(w http.ResponseWriter, r *http.Request) {
	cs, err := hrfacade.GetContacts()
	if err != nil {
		log.Printf("ERROR: GetContacts: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for ce := range cs {
		if ce.Error != nil {
			log.Printf("ERROR: ce.Error = %v", ce.Error)
			http.Error(w, ce.Error.Error(), http.StatusInternalServerError)
			return
		} else {
			c := ce.Contact
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", c.Empno, c.Active, c.Fname, c.Mi, c.Lname,
				c.Jobtitle)
		}
	}
}

func main() {
	flag.Parse()
	http.HandleFunc("/count", counthandler)
	http.HandleFunc("/contacts", contacthandler)
	log.Printf("Listening at %v", *address)
	log.Fatal(http.ListenAndServe(*address, nil))
}
