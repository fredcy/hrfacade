package main

import (
	"flag"
	"fmt"
	"github.com/fredcy/hrfacade"
	"encoding/json"
	"net/http"
	"log"
	"strings"
)

var address = flag.String("address", ":8080", "Listen and serve at this address")

func counthandler(w http.ResponseWriter, r *http.Request) {
	c, _ := hrfacade.PersonnelCount()
	fmt.Fprintf(w, "%d", c)
}

// acceptValues parses an Accept header value, returning a map from mimetype to
// a key-value map
func acceptValues(accept []string) map[string](map[string]string) {
	vs := make(map[string](map[string]string))
	for _, line := range accept {
		parts := strings.Split(line, ",")
		for _, p := range parts {
			values := strings.Split(p, ";")
			mimetype := strings.TrimSpace(values[0])
			var m map[string]string
			var ok bool
			if m, ok = vs[mimetype]; !ok {
				m = make(map[string]string)
				vs[mimetype] = m
			}
			for _, kv := range values[1:] {
				kvs := strings.Split(kv, "=")
				k := kvs[0]
				v := kvs[1]
				m[k] = v
			}
		}
	}
	return vs
}

// contacthandler returns data abouve all contacts
func contacthandler(w http.ResponseWriter, r *http.Request) {
	cs, err := hrfacade.GetContacts()
	if err != nil {
		log.Printf("ERROR: GetContacts: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	accept := r.Header[http.CanonicalHeaderKey("Accept")]
	avs := acceptValues(accept)
	if _, ok := avs["application/json"]; ok {
		contacthandlerjson(w, r, cs)
	} else {
		contacthandlertext(w, r, cs)
	}
}

func contacthandlertext(w http.ResponseWriter, r *http.Request, cs chan hrfacade.Contact) {
	for c := range cs {
		_, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			c.Empno, c.Active, c.Fname, c.Mi, c.Lname,	c.Jobtitle,
			c.Homephone, c.Busphone, c.Cellphone, c.Faxphone, c.Pagerphone)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func contacthandlerjson(w http.ResponseWriter, r *http.Request, cs chan hrfacade.Contact) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintln(w, "[")
	enc := json.NewEncoder(w)
	first := true
	for c := range cs {
		if ! first {
			fmt.Fprint(w, ",")
		}
		if err := enc.Encode(&c); err != nil {
			log.Println(err)
			return
		}
		first = false
	}
	fmt.Fprintln(w, "]")
}

func main() {
	flag.Parse()
	http.HandleFunc("/count", counthandler)
	http.HandleFunc("/contacts", contacthandler)
	log.Printf("Listening at %v", *address)
	log.Fatal(http.ListenAndServe(*address, nil))
}
