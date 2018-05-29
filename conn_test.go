package dbutil

import (
	"testing"
	"fmt"
)

func Test(t *testing.T) {
	_, err := Instance.Query("SELECT 102, 104")
	fmt.Println(err)
}