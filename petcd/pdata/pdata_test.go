package pdata

import "testing"

func TestPData(t *testing.T) {
	mapData := map[string]string{
		"a":       "1",
		"/a":      "2",
		"/a/b":    "3",
		"/a/b/c":  "4",
		"a/d/e/f": "5",
	}

	root := NewPEtcdRoot()
	root.SetAll(mapData)

	r := getRegPath()
	for k, v := range mapData {
		if n := root.Get(k); n != nil {
			if n.GetValue() != v {
				t.Fatalf("key:%q's value not match,expect:%q,got:%q", k, v, n.GetValue())
			}
		} else {
			if r.MatchString(k) {
				t.Fatalf("key:%q not found", k)
			}
		}
	}
}
