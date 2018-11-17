package jsonapi

// Document ...
type Document struct {
	// Data
	Data interface{}

	// Included
	Included map[string]Resource

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

// NewDocument ...
func NewDocument() *Document {
	return &Document{
		Included:  map[string]Resource{},
		Resources: map[string]map[string]struct{}{},
		Links:     map[string]Link{},
		RelData:   map[string][]string{},
		Meta:      map[string]interface{}{},
	}
}

// Include ...
func (d *Document) Include(res Resource) {
	key := res.GetType() + " " + res.GetID()

	if len(d.Included) == 0 {
		d.Included = map[string]Resource{}
	}

	if dres, ok := d.Data.(Resource); ok {
		// Check resource
		rkey := dres.GetID() + " " + dres.GetType()

		if rkey == key {
			return
		}
	} else if col, ok := d.Data.(Collection); ok {
		// Check Collection
		ctyp := col.Sample().GetType()
		if ctyp == res.GetType() {
			for i := 0; i < col.Len(); i++ {
				rkey := col.Elem(i).GetID() + " " + col.Elem(i).GetType()

				if rkey == key {
					return
				}
			}
		}
	}

	// Check already included resources
	if _, ok := d.Included[key]; ok {
		return
	}

	d.Included[key] = res
}

// MarshalJSON ...
// func (d *Document) Marshal() ([]byte, error) {
// 	// Data
// 	var data json.RawMessage
// 	var errors json.RawMessage
// 	var err error
// 	if d.Errors != nil {
// 		errors, err = json.Marshal(d.Errors)
// 	} else if res, ok := d.Data.(Resource); ok {
// 		_, typ := res.IDAndType()
// 		data, err = marshalResource(res, d.URL.Host, d.URL.Params.Fields[typ], d.RelData)
// 	} else if col, ok := d.Data.(Collection); ok {
// 		data, err = marshalCollection(col, d.URL.Host, d.URL.Params.Fields[col.Type()], d.RelData)
// 	} else if id, ok := d.Data.(Identifier); ok {
// 		data, err = json.Marshal(id)
// 	} else if ids, ok := d.Data.(Identifiers); ok {
// 		data, err = json.Marshal(ids)
// 	} else {
// 		data = []byte("null")
// 	}
//
// 	if err != nil {
// 		return []byte{}, err
// 	}
//
// 	// Included
// 	inclusions := []*json.RawMessage{}
// 	if len(data) > 0 {
// 		for key := range d.Included {
// 			_, typ := d.Included[key].IDAndType()
// 			raw, err := marshalResource(d.Included[key], d.URL.Host, d.URL.Params.Fields[typ], d.RelData)
// 			if err != nil {
// 				return []byte{}, err
// 			}
// 			rawm := json.RawMessage(raw)
// 			inclusions = append(inclusions, &rawm)
// 		}
// 	}
//
// 	// Marshaling
// 	plMap := map[string]interface{}{}
//
// 	if len(data) > 0 {
// 		plMap["data"] = data
// 	}
//
// 	if len(d.Links) > 0 {
// 		plMap["links"] = d.Links
// 	}
//
// 	if len(errors) > 0 {
// 		plMap["errors"] = errors
// 	}
//
// 	if len(inclusions) > 0 {
// 		plMap["included"] = inclusions
// 	}
//
// 	if len(d.Meta) > 0 {
// 		plMap["meta"] = d.Meta
// 	}
//
// 	if len(d.JSONAPI) > 0 {
// 		plMap["jsonapi"] = d.JSONAPI
// 	}
//
// 	return json.Marshal(plMap)
// }
