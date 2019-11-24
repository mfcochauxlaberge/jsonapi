package examples_test

import (
	"fmt"
	"sort"
	"time"

	"github.com/mfcochauxlaberge/jsonapi"
)

func ExampleSchema() {
	// A schema holds a list of types.
	schema := &jsonapi.Schema{}

	// A type holds information about a type, like its name,
	// attributes, and relationships.
	comments := jsonapi.Type{
		Name: "comments",
	}

	// Attributes can be added.
	comments.AddAttr(jsonapi.Attr{
		Name:     "content",
		Type:     jsonapi.AttrTypeString,
		Nullable: false,
	})

	// Relationships can be added.
	comments.AddRel(jsonapi.Rel{
		FromType: "comments",
		FromName: "author",
		ToOne:    true,
		ToType:   "users",
		ToName:   "comments",
		FromOne:  false,
	})
	comments.AddRel(jsonapi.Rel{
		FromType: "comments",
		FromName: "article",
		ToOne:    true,
		ToType:   "articles",
		ToName:   "comments",
		FromOne:  false,
	})

	// Finally, the type is added to the schema. But it can
	// still be modified after.
	schema.AddType(comments)

	// Here, types are built from structs and added.
	schema.AddType(jsonapi.MustBuildType(User{}))
	schema.AddType(jsonapi.MustBuildType(Article{}))

	// Since a comments type was added dynamically, the two types
	// added above to not contain the necessary relationships, but
	// they can be added.
	schema.AddRel("users", jsonapi.Rel{
		FromType: "users",
		FromName: "comments",
		ToOne:    false,
		ToType:   "comments",
		ToName:   "author",
		FromOne:  true,
	})
	schema.AddRel("articles", jsonapi.Rel{
		FromType: "articles",
		FromName: "comments",
		ToOne:    false,
		ToType:   "comments",
		ToName:   "article",
		FromOne:  true,
	})

	// A schema can be checked. Some validation is performed
	// like checking the names and making sure relationships
	// point to types that exist.
	// It is possible to modify a schema at anytime, but a call
	// to Check should always be performed before using it.
	// If the data it contains in inconsistent, this library
	// can behave unexpectedly.
	_ = schema.Check()

	// This schema contains 0 errors and three types.
	out := []string{
		fmt.Sprint(len(schema.Check())), // 0
	}
	for _, typ := range schema.Types {
		out = append(out, typ.Name)
	}
	sort.Strings(out)
	for _, name := range out {
		fmt.Println(name)
	}
	// Output:
	// 0
	// articles
	// comments
	// users
}

// The following structs are defined and used in this file, but they are also
// used in other examples.

type User struct {
	// The ID field is mandatory and the api tag sets the type name.
	ID string `json:"id" api:"users"`

	// Attributes
	// They are defined by setting the api to tag "attr".
	Username     string    `json:"username" api:"attr"`
	RegisteredAt time.Time `json:"registered-at" api:"attr"`

	// Relationships
	// They are defined by setting the api to tag "rel," followed
	// by the name of the target type. Optionally, a third argument
	// can be given to specify a relationship on the target type
	// which indicates a two-way relationship.
	Articles []string `json:"articles" api:"rel,articles,author"`
}

type Article struct {
	// The ID field is mandatory and the api tag sets the type name.
	ID string `json:"id" api:"articles"`

	// Attributes
	Title      string    `json:"title" api:"attr"`
	Content    string    `json:"content" api:"attr"`
	CreratedAt time.Time `json:"created-at" api:"attr"`

	// Relationships
	Author string `json:"author" api:"rel,users,articles"`
}
