package examples_test

import (
	"fmt"

	"github.com/mfcochauxlaberge/jsonapi"
)

func ExampleSoftResource() {
	// A SoftResource is a struct that implements the Resource
	// interface. It holds a Type object that defines its type and
	// that type is mutable
	sr := &jsonapi.SoftResource{}

	// One use case for a SoftResource is handling a JSON request
	// where not all possible fields have been defined in the payload.
	// The jsonapi library can create a SoftResource with a new type
	// that contains a subset of the fields from the original type.

	// The resource can be modified.
	sr.SetID("user1")

	// When an attribute is added, its value is automatically set to
	// the zero value of the type.
	sr.AddAttr(jsonapi.Attr{
		Name:     "username",
		Type:     jsonapi.AttrTypeString,
		Nullable: false,
	})
	sr.Set("username", "rob")

	fmt.Println(sr.Get("username"))
	// Output:
	// rob
}
