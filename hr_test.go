package hrfacade

import (
	"log"
	"testing"
	"os"
)

func TestSimple(t *testing.T) {
	dsn := os.Getenv("HR_DSN")
	c, err := PersonnelCount(dsn)
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
