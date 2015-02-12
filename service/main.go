package main

import (
	"flag"
	"fmt"
	"github.com/fredcy/hrfacade"
	"net/http"
	"log"
)

var address = flag.String("address", ":8080", "Listen and serve at this address")

func contacthandler(w http.ResponseWriter, r *http.Request) {
	c, _ := hrfacade.PersonnelCount()
	fmt.Fprintf(w, "%d", c)
}

func main() {
	flag.Parse()
	http.HandleFunc("/contacts", contacthandler)
	log.Printf("Listening at %v", *address)
	log.Fatal(http.ListenAndServe(*address, nil))
}
