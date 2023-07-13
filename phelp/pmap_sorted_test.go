package phelp

import "testing"

func TestPSortedMap(t *testing.T) {
	m := map[string]string{
		"a": "aaa",
		"b": "bbb",
		"c": "ccc",
	}
	sm := NewPSortedMap[string, string]()
	sm.Append(m)

	l := sm.Length()
	for i := 0; i < l; i++ {
		k, v, err := sm.GetByIndex(i)
		if err != nil {
			t.Fatal(err)
		}
		if m[k] != v {
			t.Fatalf("key:%v, got value:%v, expect:%v", k, v, m[k])
		}
	}
}
