package examples_test

import (
	"fmt"

	"github.com/mfcochauxlaberge/jsonapi"
)

func ExampleUnmarshalDocument() {
	// See the schema example for more details and the definitions
	// of User and Article.
	schema := &jsonapi.Schema{}
	schema.AddType(jsonapi.MustBuildType(User{}))
	schema.AddType(jsonapi.MustBuildType(Article{}))
	_ = schema.Check()

	// This is the payload to be unmarshaled.
	payload := `
		{
			"data": {
				"attributes": {
					"registered-at": "2019-11-19T23:17:01-05:00",
					"username": "rob"
				},
				"id": "user1",
				"relationships": {
					"articles": {
						"data": [
							{
								"type": "articles",
								"id": "article1"
							}
						]
					}
				},
				"type": "users"
			},
			"jsonapi": {
				"version": "1.0"
			},
			"meta": {
				"meta": "meta_value"
			}
		}
	`

	// UnmarhsalDocument unmarshals a payload using a schema for some
	// validation and returns a document.
	doc, err := jsonapi.UnmarshalDocument([]byte(payload), schema)
	if err != nil {
		panic(err)
	}

	// If the data top-level field contains something valid, it will be
	// an object that implements the Resource or Collection interface.
	res, _ := doc.Data.(jsonapi.Resource)

	// Print the result.
	fmt.Printf("user.ID: %s\n", res.GetID())
	fmt.Printf("user.Username: %s\n", res.Get("username"))
	fmt.Printf("user.RegisteredAt: %s\n", res.Get("registered-at"))
	fmt.Printf("user.Articles: %s\n", res.GetToMany("articles"))
	// Output:
	// user.ID: user1
	// user.Username: rob
	// user.RegisteredAt: 2019-11-19 23:17:01 -0500 EST
	// user.Articles: [article1]
}
