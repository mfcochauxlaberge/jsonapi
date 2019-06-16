package jsonapi_test

import (
	"testing"

	. "github.com/mfcochauxlaberge/jsonapi"
	"github.com/mfcochauxlaberge/tchek"
)

// TestDocument ...
func TestDocument(t *testing.T) {
	pl1 := Document{}

	tchek.AreEqual(t, "empty", nil, pl1.Data)
}
