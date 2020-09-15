package auth

import (
	"fmt"
	"testing"
)

func TestGetCreds(t *testing.T) {
	result := GetCreds()
	fmt.Println(result)
	// t.Error()
}
