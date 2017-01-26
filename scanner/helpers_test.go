package scanner

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertGinPathToSwaggerPath(t *testing.T) {
	path := convertGinPathToSwaggerPath("/:id/edit/:postId")

	assert.Equal(t, "/{id}/edit/{postId}", path)
}
