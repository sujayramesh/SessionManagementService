package utils_test

import (
	"testing"
	"github.com/sujayramesh/SMS/utils"
)

const (
	uuid_length = 36
)

func TestUUIDGeneration(t *testing.T) {
	if leng := len(utils.GenerateUUID()); leng != uuid_length {
		t.Fatal("Incorrect length")
	}
}