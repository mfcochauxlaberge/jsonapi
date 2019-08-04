package jsonapi

// A Document represents a JSON:API document.
type Document struct {
	// Data
	Data interface{}

	// Included
	Included []Resource

	// References
	Resources map[string]map[string]struct{}
	Links     map[string]Link

	// Relationships where data has to be included in payload
	RelData map[string][]string

	// Top-level members
	Meta map[string]interface{}

	// Errors
	Errors []Error

	// Internal
	PrePath string
}

// NewDocument returns a pointer to a new Document.
func NewDocument() *Document {
	return &Document{
		Included:  []Resource{},
		Resources: map[string]map[string]struct{}{},
		Links:     map[string]Link{},
		RelData:   map[string][]string{},
		Meta:      map[string]interface{}{},
	}
}

// Include adds res to the set of resources to be included under the
// included top-level field.
//
// It also makes sure that resources are not added twice.
func (d *Document) Include(res Resource) {
	key := res.GetID() + " " + res.GetType().Name

	if len(d.Included) == 0 {
		d.Included = []Resource{}
	}

	if dres, ok := d.Data.(Resource); ok {
		// Check resource
		rkey := dres.GetID() + " " + dres.GetType().Name

		if rkey == key {
			return
		}
	} else if col, ok := d.Data.(Collection); ok {
		// Check Collection
		ctyp := col.GetType()
		if ctyp.Name == res.GetType().Name {
			for i := 0; i < col.Len(); i++ {
				rkey := col.At(i).GetID() + " " + col.At(i).GetType().Name

				if rkey == key {
					return
				}
			}
		}
	}

	// Check already included resources
	for _, res := range d.Included {
		if key == res.GetID()+" "+res.GetType().Name {
			return
		}
	}

	d.Included = append(d.Included, res)
}
