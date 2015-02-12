package hrfacade

import (
	"log"
	"testing"
)

func TestSimple(t *testing.T) {
	c, err := PersonnelCount()
	if err != nil {
		t.Fatal(err)
	}
	if testing.Verbose() {
		log.Printf("persnl_count returned c = %v", c)
	}
	if c < 1 || c > 1000 {
		t.Errorf("invalid c: %v", c)
	}
}
