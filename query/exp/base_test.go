package exp

import (
	"testing"
	"bytes"
	"github.com/stretchr/testify/assert"
)

func TestLiteral(t *testing.T) {
	l := LiteralExp{Value: 12}
	b := bytes.Buffer{}
	l.ToSQL(nil, &b)
	assert.Equal(t, "12", b.String())

	l = LiteralExp{Value: false}
	b = bytes.Buffer{}
	l.ToSQL(nil, &b)
	assert.Equal(t, "false", b.String())

	l = LiteralExp{Value: 13.5}
	b = bytes.Buffer{}
	l.ToSQL(nil, &b)
	assert.Equal(t, "13.5", b.String())

	l = LiteralExp{Value: "string"}
	b = bytes.Buffer{}
	l.ToSQL(nil, &b)
	assert.Equal(t, `'string'`, b.String())

}

