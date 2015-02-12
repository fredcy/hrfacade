package hrfacade

import (
	"testing"
)

func TestSimple(t *testing.T) {
	c, err := persnl_count()
	if err != nil {
		t.Fatal(err)
	}
	if c < 1 || c > 1000 {
		t.Errorf("invalid c: %v", c)
	}
}
