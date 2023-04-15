package leak

import (
	"testing"

	"go.uber.org/goleak"
)

func TestLeak(t *testing.T) {

	if err := leak(); err != nil {
		t.Fatal("error not expexted")
	}
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
