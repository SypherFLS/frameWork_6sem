package db

import (
	"testing"
)

func TestITEMValidation(t *testing.T) {
	i := &Item{"1", "Name", 0}

	res := i.Validate()
	t.Error(res)
}