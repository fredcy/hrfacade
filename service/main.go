package main

import (
	"flag"
	"fmt"
	"github.com/fredcy/hrfacade"
	"encoding/json"
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
	for c := range cs {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s", c.Empno, c.Active, c.Fname, c.Mi, c.Lname,	c.Jobtitle)
		fmt.Fprintf(w, "\t%s\t%s\t%s", c.Homephone, c.Busphone, c.Cellphone)
		fmt.Fprintf(w, "\t%s\t%s", c.Faxphone, c.Pagerphone)
		fmt.Fprintf(w, "\n")
	}
}

func contacthandlerj(w http.ResponseWriter, r *http.Request) {
	cs, err := hrfacade.GetContacts()
	if err != nil {
		log.Printf("ERROR: GetContacts: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "[")
	enc := json.NewEncoder(w)
	first := true
	for c := range cs {
		if ! first {
			fmt.Fprint(w, ",")
		}
		if err := enc.Encode(&c); err != nil {
			log.Println(err)
		}
		first = false
	}
	fmt.Fprintln(w, "]")
}

func main() {
	flag.Parse()
	http.HandleFunc("/count", counthandler)
	http.HandleFunc("/contacts", contacthandler)
	http.HandleFunc("/contactsj", contacthandlerj)
	log.Printf("Listening at %v", *address)
	log.Fatal(http.ListenAndServe(*address, nil))
}
