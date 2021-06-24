package examples_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mfcochauxlaberge/jsonapi"
)

func ExampleMarshalDocument() {
	// See the schema example for more details and the definitions
	// of User and Article.
	schema := &jsonapi.Schema{}
	schema.AddType(jsonapi.MustBuildType(&User{}))
	schema.AddType(jsonapi.MustBuildType(&Article{}))
	_ = schema.Check()

	now, _ := time.Parse(time.RFC3339, "2019-11-19T23:17:01-05:00")

	// Two objects are created to use in our payload.
	user := &User{
		ID:           "user1",
		Username:     "rob",
		RegisteredAt: now,
		Articles:     []string{"article1"},
	}

	article := &Article{
		ID:         "article1",
		Title:      "How to make pizza",
		Content:    "Buy one.",
		CreratedAt: now,
		Author:     "user1",
	}

	doc := &jsonapi.Document{}
	// user is wrapped because its type does not implement the
	// Resource interface. Wrapping is useful to quickly get
	// started, but implementing the interface is necessary if
	// performance is an issue.
	doc.Data = jsonapi.Wrap(user)
	// article is also wrapped.
	doc.Include(jsonapi.Wrap(article))

	// As an example, some meta data is added.
	doc.Meta = jsonapi.Meta{
		"meta": "meta_value",
	}

	// This tells the document that relationship data for articles
	// must be included in the payload.
	// The relationship links are always included.
	doc.RelData = map[string][]string{
		"users": {"articles"},
	}

	// A URL represents a JSON:API compliant URL. Query
	// parameters are also properly handled.
	// A schema is given for validation. For example, it makes
	// sure the type mentioned in the path exists.
	url, _ := jsonapi.NewURLFromRaw(schema, `/users/user1?include=articles`)

	// MarhsalDocument marshals the document into a JSON:API
	// compliant payload and uses the given URL to add links
	// and know which fields to include in the result (through
	// the fields[...] query parameters).
	payload, _ := jsonapi.MarshalDocument(doc, url)

	// Beautify the output for clarity.
	out := &bytes.Buffer{}
	json.Indent(out, payload, "", "\t")

	// Print the result.
	fmt.Println(string(out.Bytes()))
	// Output:
	// {
	// 	"data": {
	// 		"attributes": {
	// 			"registered-at": "2019-11-19T23:17:01-05:00",
	// 			"username": "rob"
	// 		},
	// 		"id": "user1",
	// 		"links": {
	// 			"self": "/users/user1"
	// 		},
	// 		"relationships": {
	// 			"articles": {
	// 				"data": [
	// 					{
	// 						"id": "article1",
	// 						"type": "articles"
	// 					}
	// 				],
	// 				"links": {
	// 					"related": "/users/user1/articles",
	// 					"self": "/users/user1/relationships/articles"
	// 				}
	// 			}
	// 		},
	// 		"type": "users"
	// 	},
	// 	"included": [
	// 		{
	// 			"attributes": {
	// 				"content": "Buy one.",
	// 				"created-at": "2019-11-19T23:17:01-05:00",
	// 				"title": "How to make pizza"
	// 			},
	// 			"id": "article1",
	// 			"links": {
	// 				"self": "/articles/article1"
	// 			},
	// 			"relationships": {
	// 				"author": {
	// 					"links": {
	// 						"related": "/articles/article1/author",
	// 						"self": "/articles/article1/relationships/author"
	// 					}
	// 				}
	// 			},
	// 			"type": "articles"
	// 		}
	// 	],
	// 	"jsonapi": {
	// 		"version": "1.0"
	// 	},
	// 	"links": {
	// 		"self": "/users/user1?fields%5Barticles%5D=author%2Ccontent%2Ccreated-at%2Ctitle\u0026fields%5Busers%5D=articles%2Cregistered-at%2Cusername"
	// 	},
	// 	"meta": {
	// 		"meta": "meta_value"
	// 	}
	// }
}
