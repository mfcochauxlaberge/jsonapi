package examples_test

import (
	"fmt"
	"time"

	"github.com/mfcochauxlaberge/jsonapi"
)

func ExampleSchema() {
	// A schema holds a list of types.
	schema := &jsonapi.Schema{}

	schema.AddType(jsonapi.MustBuildType(User{}))
	schema.AddType(jsonapi.MustBuildType(Article{}))

	// A schema can be checked. Some validation is performed
	// like checking the names and making sure relationships
	// point to types that exist.
	// It is possible to modify a schema at anytime, but a call
	// to Check should always be performed before using it.
	// If the data it contains in inconsistent, this library
	// can behave unexpectedly.
	_ = schema.Check()

	// Useful methods are offered, like HasType.
	has := schema.HasType("users")
	fmt.Println(has)
	// Output: true
}

// The following structs are defined and used in this file, but they are also
// used in other example.

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
	Author string `json:"author" api:"rel,author,articles"`
}
