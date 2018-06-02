package query

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestPrepareParameters(t *testing.T) {
	names := []string{"Hello", "Two"}
	arg0 := struct {
		Two int
	}{Two: 4}
	arg1 := map[interface{}]string{
		"Hello": "World",
		2: "Okay",
		"OhNo": "Yeah",
	}
	arg2 := map[int]string{1: "No"}
	assert.Equal(t, []interface{}{"World", 4},
		*PrepareParameters(&names, arg0, &arg0, arg2, &arg1))
}