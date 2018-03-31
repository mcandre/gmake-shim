package gmake_test

import (
	"testing"
	"github.com/mcandre/gmake-shim"
)

func TestVersion(t *testing.T) {
	if gmake.Version == "" {
		t.Errorf("Expected %v to be non-empty", gmake.Version)
	}
}
