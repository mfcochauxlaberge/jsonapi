package jsonapi_test

import (
	"sort"
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

func TestInclude(t *testing.T) {
	assert := assert.New(t)

	doc := &Document{}
	typ1 := &Type{Name: "t1"}
	typ2 := &Type{Name: "t2"}

	// Main data
	res := &SoftResource{}
	res.SetType(typ1)
	res.SetID("id1")
	doc.Data = res

	// Inclusions
	res = &SoftResource{}
	res.SetType(typ1)
	res.SetID("id1")
	doc.Include(res)

	res = &SoftResource{}
	res.SetType(typ1)
	res.SetID("id2")
	doc.Include(res)

	res = &SoftResource{}
	res.SetType(typ1)
	res.SetID("id3")
	doc.Include(res)

	res = &SoftResource{}
	res.SetType(typ1)
	res.SetID("id3")
	doc.Include(res)

	res = &SoftResource{}
	res.SetType(typ2)
	res.SetID("id1")
	doc.Include(res)

	// Check
	ids := []string{}
	for _, res := range doc.Included {
		ids = append(ids, res.GetType().Name+"-"+res.GetID())
	}
	sort.Strings(ids)

	expect := []string{
		"t1-id2",
		"t1-id3",
		"t2-id1",
	}
	assert.Equal(expect, ids)
}
