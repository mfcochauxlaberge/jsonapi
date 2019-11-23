package examples_test

import (
	"fmt"

	"github.com/mfcochauxlaberge/jsonapi"
)

func ExampleWrapper() {
	// The Resource interface is useful for the jsonapi library. It
	// makes manipulation of the resource much easier. But it requires
	// some work and can be annoying when one wants to make a quick
	// project where performance is not an issue.

	// The Wrap function exists to solve this problem. It takes an
	// object or a pointer to an object and return a Wrapper, which
	// which is a struct that implements the Resource interface. When
	// the methods of the interface are used to modify a resource, the
	// Wrapper can use reflection to mutate the original object. That
	// means that if Wrap is given a pointer, that object is modified
	// and the original pointer can still be used to handle the object
	// in a native and type safe way.
	wrap := jsonapi.Wrap(Animal{})

	// The resource can be modified.
	wrap.SetID("animal1")
	wrap.Set("name", "Gopher")

	// Unlike a SoftResource, its type cannot changed. It is defined
	// by the struct definition.
	//
	// A Type object is still generated and can be retrieved.
	_ = wrap.GetType()

	fmt.Println(wrap.Get("name").(string))
	// Output:
	// Gopher
}

// Animal does not implement the Resource interface, but can be
// wrapped if it follows the struct format defined by the jsonapi
// library (ID field, attr tags, rel tags, etc).
type Animal struct {
	ID string `json:"id" api:"animals"`

	Name string `json:"name" api:"attr"`
}
