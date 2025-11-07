package util

import (
	"fmt"
	"testing"
)

func TestMin(t *testing.T) {
	fmt.Print(Min(1, 2))
	fmt.Print(Min(1.1, 2.1))
}
