package clause

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestSQLRecipe_Insert(t *testing.T) {
	q, err := SQL().Read("Hello", "World").Select("table")
	assert.Equal(t, "SELECT Hello,World FROM table", q)
	assert.NoError(t, err)


}
