package swagger

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseEnum(t *testing.T) {
	doc, hasEnum := ParseEnum("swagger:enum \nasdasdasdad")
	assert.Equal(t, "asdasdasdad", doc)
	assert.Equal(t, true, hasEnum)
}

func TestParseStrfmt(t *testing.T) {
	doc, fmtName := ParseStrfmt("swagger:strfmt date-time\nasdasdasdad")
	assert.Equal(t, "asdasdasdad", doc)
	assert.Equal(t, "date-time", fmtName)
}

func TestGetCommonValidations(t *testing.T) {
	c := GetCommonValidations("@int[1,2]")
	assert.NotNil(t, c.Minimum)
	assert.NotNil(t, c.Maximum)
	assert.Equal(t, false, c.ExclusiveMinimum)
	assert.Equal(t, false, c.ExclusiveMaximum)
}

func TestGetCommonValidationsWithExclusive(t *testing.T) {
	c := GetCommonValidations("@int[^1,^2]")
	assert.NotNil(t, c.Minimum)
	assert.NotNil(t, c.Maximum)
	assert.Equal(t, true, c.ExclusiveMinimum)
	assert.Equal(t, true, c.ExclusiveMaximum)
}
