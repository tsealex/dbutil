package clause

import (
	"testing"
	"fmt"
)

func TestSQLRecipe_Insert(t *testing.T) {
	q, err := SQL().Read("Hello", "World").Select("table")
	fmt.Println(q)
	fmt.Println(err)


}
