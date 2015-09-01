package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fredcy/hrfacade"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)

var version string = "unknown" // set with -ldflags "-X main.version 1.3"
var build string = "unknown"

func counthandler(w http.ResponseWriter, r *http.Request, dsn string) {
	c, _ := hrfacade.PersonnelCount(dsn)
	fmt.Fprintf(w, "%d", c)
}

// contacthandler returns data about all contacts
func contacthandler(w http.ResponseWriter, r *http.Request, dsn string) {
	all := r.FormValue("all") != ""
	cs, err := hrfacade.GetContacts(dsn, all)
	if err != nil {
		log.Printf("ERROR: GetContacts: %v", err)
		http.Error(w, "Unable to get contact data.", http.StatusInternalServerError)
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
		_, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			c.Empno, c.Active, c.Fname, c.Mi, c.Lname, c.Jobtitle, c.Email,
			c.Homephone, c.Busphone, c.Cellphone, c.Faxphone, c.Pagerphone,
			c.Level1, c.Level2, c.Level3, c.Level4, c.Superno, c.Sup2no)
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
		if !first {
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

var basicAuthRe = regexp.MustCompile(`^Basic (.*)$`)

func wrapLog(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		fn(w, r)
		endTime := time.Now()
		log.Printf("served %v to %v in %v",
			r.URL, r.RemoteAddr, endTime.Sub(startTime))
	}
}

// wrapAuthenticate adds authentication to the HTTP handler.
func wrapAuthenticate(fn http.HandlerFunc, authcode string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if authcode != "" {
			auth := r.Header.Get(http.CanonicalHeaderKey("Authorization"))
			if basicAuthRe.MatchString(auth) {
				auth64 := basicAuthRe.ReplaceAllString(auth, "$1")
				auth_b, err := base64.StdEncoding.DecodeString(auth64)
				if err != nil {
					log.Println(err)
				}
				auth = string(auth_b)
			}
			if auth != authcode {
				log.Printf("Invalid authorization value: \"%v\"", auth)
				http.Error(w, "Not authorized", http.StatusUnauthorized)
				return
			}
		}
		fn(w, r)
	}
}

// wrapDSN passes a dsn argument to the underlying handler function.
func wrapDSN(fn func(http.ResponseWriter, *http.Request, string), dsn string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, dsn)
	}
}

func main() {
	dsn := os.Getenv("HR_DSN")
	if dsn == "" {
		log.Fatal("HR_DSN is not set in environment")
	}
	var address = flag.String("address", ":8080", "Listen and serve at this address")
	var authcode = flag.String("authcode", "", "Authorization header value expected")
	flag.Parse()

	http.HandleFunc("/count", wrapAuthenticate(wrapDSN(counthandler, dsn), *authcode))
	http.HandleFunc("/contacts", wrapLog(wrapAuthenticate(wrapDSN(contacthandler, dsn), *authcode)))

	log.Printf("hrfacade/service, version %s (%s): Listening at %v", version, build, *address)
	log.Fatal(http.ListenAndServe(*address, nil))
}
