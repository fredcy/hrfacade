package main

import (
	_ "fmt"
	"reflect"
	"testing"
)

func MustEqual(t *testing.T, got, want interface{}) {
	if want != got {
		twant := reflect.TypeOf(want)
		tgot := reflect.TypeOf(got)
		t.Errorf("want <<%v>> (%v), got <<%v>> (%v)", want, twant, got, tgot)
	}
}

func TestSimple(t *testing.T) {
	header := []string{"*/*"}
	avs := acceptValues(header)
	MustEqual(t, len(avs), 1)
	v, ok := avs["*/*"]
	if !ok {
		t.Errorf("Missing map key")
	} else {
		MustEqual(t, len(v), 0)
	}
}

func MapMustEqual(t *testing.T, m map[string]string, key, want string) {
	got, ok := m[key]
	if !ok {
		t.Errorf("Missing map key %v", key)
	} else {
		MustEqual(t, got, want)
	}
}

func TestAccept(t *testing.T) {
	header := []string{"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"}
	avs := acceptValues(header)
	MustEqual(t, len(avs), 5)

	v, ok := avs["text/html"]
	if !ok {
		t.Errorf("Missing key text/html")
	} else {
		MustEqual(t, len(v), 0)
	}

	v, ok = avs["application/xml"]
	if !ok {
		t.Errorf("Missing key application/xml")
	} else {
		MustEqual(t, len(v), 1)
		MapMustEqual(t, v, "q", "0.9")
	}
}
