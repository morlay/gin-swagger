package swagger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertGinPathToSwaggerPath(t *testing.T) {
	path := convertGinPathToSwaggerPath("/:id/edit/:postId")

	assert.Equal(t, "/{id}/edit/{postId}", path)
}
