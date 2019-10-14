package jsonapi_test

import (
	"testing"

	. "github.com/mfcochauxlaberge/jsonapi"

	"github.com/stretchr/testify/assert"
)

func TestCheck(t *testing.T) {
	assert := assert.New(t)

	err := Check("not a struct")
	assert.EqualError(err, "jsonapi: not a struct")

	err = Check(emptyIDAPItag{})
	assert.EqualError(err, "jsonapi: ID field's api tag is empty")

	err = Check(invalidAttributeType{})
	assert.EqualError(
		err,
		"jsonapi: attribute \"Attr\" of type \"typename\" is of unsupported type",
	)

	err = Check(invalidRelAPITag{})
	assert.EqualError(
		err,
		"jsonapi: api tag of relationship \"Rel\" of struct \"invalidRelAPITag\" is invalid",
	)

	err = Check(invalidReType{})
	assert.EqualError(
		err,
		"jsonapi: relationship \"Rel\" of type \"typename\" is not string or []string",
	)
}

func TestBuildType(t *testing.T) {
	assert := assert.New(t)

	assert.Panics(func() {
		MustBuildType("invalid")
	})

	mock := mockType1{
		ID:    "abc13",
		Str:   "string",
		Int:   -42,
		Uint8: 12,
	}
	typ, err := BuildType(mock)
	assert.NoError(err)
	assert.Equal(true, Equal(Wrap(&mockType1{}), typ.New()))

	// Build type from pointer to struct
	typ, err = BuildType(&mock)
	assert.NoError(err)
	assert.Equal(true, Equal(Wrap(&mockType1{}), typ.New()))

	// Build from invalid struct
	_, err = BuildType(invalidRelAPITag{})
	assert.Error(err)
}

func TestIDAndType(t *testing.T) {
	assert := assert.New(t)

	mt := mocktype{
		ID: "abc123",
	}
	id, typ := IDAndType(mt)
	assert.Equal("abc123", id)
	assert.Equal("mocktype", typ)

	// Resource
	id, typ = IDAndType(Wrap(&mt))
	assert.Equal("abc123", id)
	assert.Equal("mocktype", typ)

	// Missing ID field
	id, typ = IDAndType(missingID{})
	assert.Equal("", id)
	assert.Equal("", typ)

	// Not a struct
	id, typ = IDAndType("not a struct")
	assert.Equal("", id)
	assert.Equal("", typ)
}

type emptyIDAPItag struct {
	ID string `json:"id"`
}

type invalidAttributeType struct {
	ID   string `json:"id" api:"typename"`
	Attr error  `json:"attr" api:"attr"`
}

type invalidRelAPITag struct {
	ID  string `json:"id" api:"typename"`
	Rel string `json:"rel" api:"rel,but,it,is,invalid"`
}

type invalidReType struct {
	ID  string `json:"id" api:"typename"`
	Rel int    `json:"rel" api:"rel,target,reverse"`
}

type missingID struct{}
