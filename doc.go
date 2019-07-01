/*
Package jsonapi offers a library to marshal and unmarshal JSON:API payloads.

It also offers many utilies for developing a JSON:API backend.

The simplest way to start using jsonapi is to use the Marshal and Unmarshal functions.

	func Marshal(doc *Document, url *URL) ([]byte, error)
	func Unmarshal(payload []byte, url *URL, schema *Schema) (*Document, error)

A schema is collection of types where relationships can point to each other. A schema can also look at its types and return any errors.

A type is generally defined with a struct.

There needs to be an ID field of type string. The `api` tag represents the name of the type.

	type User struct {
		ID string `json:"id" api:"users"` // ID is mandatory and the api tag sets the type

		// Attributes
		Name string `json:"name" api:"attr"` // attr means it is an attribute
		BornAt time.Time `json:"born-at" api:"attr"`

		// Relationships
		Articles []string `json:"articles" api:"rel,articles"`
	}

A lot more is offered in this library. The best way to learn how to use it is to look at the source code and its comments.
*/
package jsonapi
