package jsonapi_test

import (
	"testing"

	. "github.com/mfcochauxlaberge/jsonapi"
	"github.com/stretchr/testify/assert"
)

// TestDocument ...
func TestDocument(t *testing.T) {
	assert := assert.New(t)

	pl1 := Document{}
	assert.Equal(nil, pl1.Data, "empty")
}
