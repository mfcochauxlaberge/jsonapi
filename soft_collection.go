package jsonapi

// SoftCollection is a collection of SoftResources where the type can be changed
// for all elements at once by modifying the Type field.
type SoftCollection struct {
	Type *Type

	col []*SoftResource
}

// SetType sets the collection's type.
func (s *SoftCollection) SetType(typ *Type) {
	s.Type = typ
}

// GetType returns the collection's type.
func (s *SoftCollection) GetType() Type {
	return *s.Type
}

// AddAttr adds an attribute to all of the resources in the collection.
func (s *SoftCollection) AddAttr(attr Attr) error {
	return s.Type.AddAttr(attr)
}

// AddRel adds a relationship to all of the resources in the collection.
func (s *SoftCollection) AddRel(rel Rel) error {
	return s.Type.AddRel(rel)
}

// Len returns the length of the collection.
func (s *SoftCollection) Len() int {
	return len(s.col)
}

// At returns the element at index i.
func (s *SoftCollection) At(i int) Resource {
	if i >= 0 && i < len(s.col) {
		return s.col[i]
	}

	return nil
}

// Resource returns the element with an ID equal to id.
//
// It builds and returns a SoftResource with only the specified fields.
func (s *SoftCollection) Resource(id string, fields []string) Resource {
	for i := range s.col {
		if s.col[i].GetID() == id {
			return s.col[i]
		}
	}

	return nil
}

// Add creates a SoftResource and adds it to the collection.
func (s *SoftCollection) Add(r Resource) {
	// A SoftResource is built from the Resource and
	// then it is added to the collection.
	sr := &SoftResource{}
	sr.id = r.GetID()
	sr.Type = s.Type

	for _, attr := range r.Attrs() {
		sr.AddAttr(attr)
		sr.Set(attr.Name, r.Get(attr.Name))
	}

	for _, rel := range r.Rels() {
		sr.AddRel(rel)

		if rel.ToOne {
			sr.SetToOne(rel.FromName, r.GetToOne(rel.FromName))
		} else {
			sr.SetToMany(rel.FromName, r.GetToMany(rel.FromName))
		}
	}

	s.col = append(s.col, sr)
}

// Remove removes the resource with an ID equal to id.
//
// Nothing happens if no resource has such an ID.
func (s *SoftCollection) Remove(id string) {
	for i := range s.col {
		if s.col[i].GetID() == id {
			s.col = append(s.col[:i], s.col[i+1:]...)
			return
		}
	}
}
