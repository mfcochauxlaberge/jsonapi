package examples_test

import (
	"fmt"

	"github.com/mfcochauxlaberge/jsonapi"
)

func ExampleSoftCollection() {
	// A SoftCollection implements the Collection interface. It can
	// only contain SoftResources which all share the same common type.
	sc := &jsonapi.SoftCollection{}
	sc.Type = &jsonapi.Type{
		Name: "users",
	}
	sc.AddAttr(jsonapi.Attr{
		Name:     "username",
		Type:     jsonapi.AttrTypeString,
		Nullable: false,
	})

	// A SoftResource is added to the collection. Its type will be
	// set to the SoftCollection's type.
	sr := &jsonapi.SoftResource{}
	sc.Add(sr)

	// Normally, the following line would not work if username was
	// not alrady defined.
	sr.Set("username", "rob")

	// An attribute is added to the type through the SoftCollection.
	sc.AddAttr(jsonapi.Attr{
		Name:     "admin",
		Type:     jsonapi.AttrTypeBool,
		Nullable: false,
	})

	// Now, all resources inside the collection have the new field.
	fmt.Println(sc.At(0).Get("admin"))
	// Output: false
}
