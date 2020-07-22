package main

import (
	"fmt"
	"testing"
)

func TestIsValidURL(t *testing.T) {
	tables := []struct {
		in       string
		expected bool
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
		result := IsURLValid(table.in)
		if result != table.expected {
			t.Errorf("IsValidURL(%v): got: %v, want: %v", table.in, result, table.expected)
		}
	}
}

func TestGenerateID(t *testing.T) {
	id := GenerateKey(false)

	fmt.Println(id)

	if id == "" {
		t.Error("ID is empty")
	}
}

func BenchmarkGenerateID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateKey(false)
	}
}
