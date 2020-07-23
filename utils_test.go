package main

import (
	"testing"
)

func TestIsValidURL(t *testing.T) {
	tables := []struct {
		url  string
		want bool
	}{
		{"http://www.myunsecuredwebsite.com", true},
		{"https://www.mysecuredwebsite.com", true},
		{"https://www.mywebsite.com/with/path/", true},
		{"https://www.mywebsite.com/with/path?and=also&query=params", true},
		{"https://www.mywebsite.com/with/trailing/slash/", true},
		{"https://www.mywebsitewithport.com:42/with/path", true},
		{"www.mywebsitewithnoprotocol.com", false},
		{"mywebsitewithonlythedomain.com", false},
		{"notawebsite", false},
		{"", false},
	}

	for _, table := range tables {
		valid := IsURLValid(table.url)
		if valid != table.want {
			t.Errorf("IsValidURL(%v): got: %v, want: %v", table.url, valid, table.want)
		}
	}
}

func TestGenerateID(t *testing.T) {
	id := GenerateKey(false)

	if len(id) != 18 {
		t.Error("ID is empty")
	}
}

func BenchmarkGenerateID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateKey(false)
	}
}
