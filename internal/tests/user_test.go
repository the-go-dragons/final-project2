package tests

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Setup()

	code := m.Run()

	// db.Close()

	os.Exit(code)
}

func TestExample1(t *testing.T) {
	// Test function 1
}
