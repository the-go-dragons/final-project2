package tests

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	code := m.Run()
	CleanUp()
	os.Exit(code)
}
