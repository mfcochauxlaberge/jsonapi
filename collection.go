package jsonapi

import "encoding/json"

// A Collection defines the interface of a structure that can manage a set of
// ordered resources of the same type.
type Collection interface {
	// Type returns the name of the resources' type.
	GetType() Type

	// Len returns the number of resources in the collection.
	Len() int

	// At returns the resource at index i.
	At(int) Resource

	// Add adds a resource in the collection.
	Add(Resource)
}

// MarshalCollection marshals a Collection into a JSON-encoded payload.
func MarshalCollection(c Collection, prepath string, fields map[string][]string, relData map[string][]string) []byte {
	var raws []*json.RawMessage

	if c.Len() == 0 {
		return []byte("[]")
	}

	for i := 0; i < c.Len(); i++ {
		r := c.At(i)
		raw := json.RawMessage(
			MarshalResource(r, prepath, fields[r.GetType().Name], relData),
		)
		raws = append(raws, &raw)
	}

	// NOTE An error should not happen.
	pl, _ := json.Marshal(raws)

	return pl
}

// UnmarshalCollection unmarshals a JSON-encoded payload into a Collection.
func UnmarshalCollection(data []byte, schema *Schema) (Collection, error) {
	var cske []json.RawMessage

	err := json.Unmarshal(data, &cske)
	if err != nil {
		return nil, err
	}

	col := &Resources{}

	for i := range cske {
		res, err := UnmarshalResource(cske[i], schema)
		if err != nil {
			return nil, err
		}

		col.Add(res)
	}

	return col, nil
}

// Resources is a slice of objects that implement the Resource interface. They
// do not necessarily have the same type.
type Resources []Resource

// GetType returns a zero Type object because the collection does not represent
// a particular type.
func (r *Resources) GetType() Type {
	return Type{}
}

// Len returns the number of elements in r.
func (r *Resources) Len() int {
	return len(*r)
}

// At returns the number of elements in r.
func (r *Resources) At(i int) Resource {
	if i >= 0 && i < r.Len() {
		return (*r)[i]
	}

	return nil
}

// Add adds a Resource object to r.
func (r *Resources) Add(res Resource) {
	*r = append(*r, res)
}
