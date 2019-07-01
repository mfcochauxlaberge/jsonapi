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

	/*
	 * Main data is a resource
	 */
	doc.Data = newResource(typ1, "id1")

	// Inclusions
	doc.Include(newResource(typ1, "id1"))
	doc.Include(newResource(typ1, "id2"))
	doc.Include(newResource(typ1, "id3"))
	doc.Include(newResource(typ1, "id3"))
	doc.Include(newResource(typ2, "id1"))

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

	/*
	 * Main data is a collection
	 */
	doc = &Document{}

	// Collection
	col := &SoftCollection{}
	col.SetType(typ1)
	col.Add(newResource(typ1, "id1"))
	col.Add(newResource(typ1, "id2"))
	col.Add(newResource(typ1, "id3"))
	doc.Data = Collection(col)

	// Inclusions
	doc.Include(newResource(typ1, "id1"))
	doc.Include(newResource(typ1, "id2"))
	doc.Include(newResource(typ1, "id3"))
	doc.Include(newResource(typ1, "id4"))
	doc.Include(newResource(typ2, "id1"))
	doc.Include(newResource(typ2, "id1"))
	doc.Include(newResource(typ2, "id2"))

	// Check
	ids = []string{}
	for _, res := range doc.Included {
		ids = append(ids, res.GetType().Name+"-"+res.GetID())
	}
	sort.Strings(ids)

	expect = []string{
		"t1-id4",
		"t2-id1",
		"t2-id2",
	}
	assert.Equal(expect, ids)
}

func newResource(typ *Type, id string) Resource {
	res := &SoftResource{}
	res.SetType(typ)
	res.SetID(id)
	return res
}
