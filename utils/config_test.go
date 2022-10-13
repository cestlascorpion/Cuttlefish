package utils

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	conf := NewTestConfig()
	if !conf.Check() {
		t.FailNow()
	}
}
