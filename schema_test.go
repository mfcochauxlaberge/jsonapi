package jsonapi_test

import (
	"testing"

	. "github.com/mfcochauxlaberge/jsonapi"

	"github.com/stretchr/testify/assert"
)

func TestSchemaTypes(t *testing.T) {
	assert := assert.New(t)

	// Add a type
	schema := &Schema{}
	err := schema.AddType(Type{Name: "type1"})
	assert.NoError(err)
	assert.True(schema.HasType("type1"))
	assert.False(schema.HasType("type2"))

	// Add an invalid type (no name)
	schema = &Schema{}
	err = schema.AddType(Type{})
	assert.Error(err)

	// Add two types with the same name
	schema = &Schema{}
	_ = schema.AddType(Type{Name: "type1"})
	err = schema.AddType(Type{Name: "type1"})
	assert.Error(err)

	// Remove a type
	schema = &Schema{}
	_ = schema.AddType(Type{Name: "type1"})
	schema.RemoveType("type1")
	typ := schema.GetType("type1")
	assert.Equal("", typ.Name)

	// Add and remove an attribute
	schema = &Schema{}
	_ = schema.AddType(Type{Name: "type1"})
	attr := Attr{
		Name:     "attr1",
		Type:     AttrTypeString,
		Nullable: false,
	}
	err = schema.AddAttr("type1", attr)
	assert.NoError(err)

	typ = schema.GetType("type1")
	assert.Contains(typ.Attrs, "attr1")
	assert.Equal(attr, typ.Attrs["attr1"])
	schema.RemoveAttr("type1", "attr1")
	assert.NotContains(schema.GetType("type1").Attrs, "attr1")

	// Add an invalid attribute (no name)
	schema = &Schema{}
	_ = schema.AddType(Type{Name: "type1"})
	err = schema.AddAttr("type1", Attr{Name: ""})
	assert.Error(err)

	// Add an invalid attribute (type does not exist)
	schema = &Schema{}
	_ = schema.AddType(Type{Name: "type1"})
	err = schema.AddAttr("type2", Attr{Name: "attr1"})
	assert.Error(err)

	// Add and remove an relationship
	schema = &Schema{}
	_ = schema.AddType(Type{Name: "type1"})
	rel := Rel{
		FromName: "rel1",
		ToOne:    true,
		ToType:   "type1",
	}
	err = schema.AddRel("type1", rel)
	assert.NoError(err)

	typ = schema.GetType("type1")
	assert.Contains(typ.Rels, "rel1")
	assert.Equal(rel, typ.Rels["rel1"])
	schema.RemoveRel("type1", "rel1")
	assert.NotContains(schema.GetType("type1").Rels, "rel1")

	// Add an invalid relationship (no name)
	schema = &Schema{}
	_ = schema.AddType(Type{Name: "type1"})
	err = schema.AddRel("type1", Rel{FromName: ""})
	assert.Error(err)

	// Add an invalid relationship (type does not exist)
	schema = &Schema{}
	_ = schema.AddType(Type{Name: "type1"})
	err = schema.AddRel("type2", Rel{FromName: "rel1"})
	assert.Error(err)
}

func TestSchemaCheck(t *testing.T) {
	assert := assert.New(t)

	schema := &Schema{}

	type1 := Type{
		Name:  "type1",
		Attrs: map[string]Attr{},
		Rels: map[string]Rel{
			"rel1": {
				FromName: "rel1",
				ToType:   "type2",
			},
			"rel2": {
				FromName: "rel2-invalid",
				ToType:   "nonexistent",
			},
			"rel3": {
				FromName: "rel3",
				ToType:   "type1",
			},
		},
	}
	err := schema.AddType(type1)
	assert.NoError(err)

	type2 := Type{
		Name:  "type2",
		Attrs: map[string]Attr{},
		Rels: map[string]Rel{
			"rel1": {
				FromName: "rel1",
				FromType: "type1",
				ToName:   "rel1",
				ToType:   "type1",
			},
			"rel2": {
				FromName: "rel2",
				FromType: "type2",
				ToName:   "rel3",
				ToType:   "type1",
			},
		},
	}
	err = schema.AddType(type2)
	assert.NoError(err)

	// assert.NotEmpty(schema.Types)
	// assert.NotEmpty(schema.GetType("type1").Rels)

	// Check schema
	errs := schema.Check()
	errsStr := []string{}

	for _, err := range errs {
		errsStr = append(errsStr, err.Error())
	}

	assert.Len(errs, 3)
	assert.Contains(
		errsStr,
		"jsonapi: field ToType of relationship \"rel2-invalid\" of type \"type1\" does not exist",
	)
	assert.Contains(
		errsStr,
		"jsonapi: field FromType of relationship \"rel1\" "+
			"must be its type's name (\"type2\", not \"type1\")",
	)
	assert.Contains(
		errsStr,
		"jsonapi: relationship \"rel2\" of type \"type2\" and its inverse do not point each other",
	)
}

func TestSchemaRels(t *testing.T) {
	assert := assert.New(t)

	schema := &Schema{}

	users := Type{
		Name: "users",
		Rels: map[string]Rel{
			"posts": {
				FromName: "posts",
				FromType: "users",
				ToOne:    false,
				ToName:   "author",
				ToType:   "messages",
				FromOne:  true,
			},
			"favorites": {
				FromName: "favorites",
				FromType: "users",
				ToOne:    false,
				ToName:   "",
				ToType:   "messages",
				FromOne:  false,
			},
		},
	}
	_ = schema.AddType(users)

	messages := Type{
		Name: "messages",
		Rels: map[string]Rel{
			"author": {
				FromName: "author",
				FromType: "messages",
				ToOne:    true,
				ToName:   "posts",
				ToType:   "users",
				FromOne:  false,
			},
		},
	}
	_ = schema.AddType(messages)

	rels := schema.Rels()
	assert.Len(rels, 2)
	assert.Equal(messages.Rels["author"], rels[0])
	assert.Equal(users.Rels["favorites"], rels[1])
}
