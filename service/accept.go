package main

import (
	"encoding/json"
	"log"
	"strings"
)

func pprint(v interface{}) {
	log.Printf("v=%v", v)
	b, _ := json.MarshalIndent(v, "", "  ")
	log.Print(string(b))
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
