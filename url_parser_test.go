package jsonapi_test

import (
	"testing"

	. "github.com/mfcochauxlaberge/jsonapi"

	"github.com/stretchr/testify/assert"
)

func TestURLParserParse(t *testing.T) {
	assert := assert.New(t)

	parser := URLParser{
		Schema: newMockSchema(),
	}

	// Invalid URL
	_, err := parser.Parse(":")
	assert.Error(err)

	// Empty URL
	_, err = parser.Parse("")
	assert.Error(err)

	_, err = parser.Parse("/")
	assert.Error(err)

	// Invalid SimpleURL
	_, err = parser.Parse("/some/path?invalidParam=hello")
	assert.Error(err)

	// Valid URL
	u, err := parser.Parse("https://example.com/mocktypes1/mc1-1")
	u.Params = nil
	// Params is set to nil because this is
	// not being tested here, so no need for
	// bloating this file.
	assert.NoError(err)
	assert.Equal(&URL{
		Fragments:       []string{"mocktypes1", "mc1-1"},
		Route:           "/mocktypes1/:id",
		BelongsToFilter: BelongsToFilter{},
		ResType:         "mocktypes1",
		ResID:           "mc1-1",
	}, u)
}

func TestURLParserPathPrefix(t *testing.T) {
	assert := assert.New(t)

	// Empty
	parser := URLParser{}

	assert.Equal("/", parser.Path("/"))
	assert.Equal("/type", parser.Path("/type"))
	assert.Equal("/type/id", parser.Path("/type/id"))
	assert.Equal("/type/id/relationship", parser.Path("/type/id/relationship"))
	assert.Equal("/type/id/relationship/id", parser.Path("/type/id/relationship/id"))

	// Slash
	parser = URLParser{PathPrefix: "/"}

	assert.Equal("/", parser.Path("/"))
	assert.Equal("/type", parser.Path("/type"))
	assert.Equal("/type/id", parser.Path("/type/id"))
	assert.Equal("/type/id/relationship", parser.Path("/type/id/relationship"))
	assert.Equal("/type/id/relationship/id", parser.Path("/type/id/relationship/id"))

	// Path no slash
	parser = URLParser{PathPrefix: "api"}

	assert.Equal("/", parser.Path("/"))
	assert.Equal("/type", parser.Path("/type"))
	assert.Equal("/type/id", parser.Path("/type/id"))
	assert.Equal("/type/id/relationship", parser.Path("/type/id/relationship"))
	assert.Equal("/type/id/relationship/id", parser.Path("/type/id/relationship/id"))
	assert.Equal("/", parser.Path("/api/"))
	assert.Equal("/type", parser.Path("/api/type"))
	assert.Equal("/type/id", parser.Path("/api/type/id"))
	assert.Equal("/type/id/relationship", parser.Path("/api/type/id/relationship"))
	assert.Equal("/type/id/relationship/id", parser.Path("/api/type/id/relationship/id"))
}
