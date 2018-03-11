package main

import (
	"testing"
)

func testEqual(t *testing.T, expected, actual uint64) {
	if expected == actual {
		return
	}
	t.Errorf("Greška: %d (očekivano) != %d (dobiveno)", expected, actual)
}

func TestFibonacci(t *testing.T) {
	testEqual(t, uint64(5), fibonacci(5))
	testEqual(t, uint64(8), fibonacci(6))
	testEqual(t, uint64(89), fibonacci(11))
}
