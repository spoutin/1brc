package main

import (
	"os"
	"testing"
)

func TestReadFile(t *testing.T) {
	filename := "../../../../measurements_small.txt"
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	fileInfo, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}
	cities := readFile("../../../../measurements_small.txt", 0, fileInfo.Size()-1)
	if len(cities) < 1 {
		t.Fatal("No cities found")
	}
}

func BenchmarkReadFile(t *testing.B) {
	filename := "../../../../measurements_small.txt"
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	fileInfo, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}
	cities := readFile("../../../../measurements_small.txt", 0, fileInfo.Size()-1)
	if len(cities) < 1 {
		t.Fatal("No cities found")
	}
}
