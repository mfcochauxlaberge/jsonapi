package jsonapi

import (
	"errors"
	"sort"
	"strings"
	"time"
)

var _ Collection = (*SoftCollection)(nil)

// SoftCollection is a collection of SoftResources where the type can be changed
// for all elements at once by modifying the Type field.
type SoftCollection struct {
	Type *Type

	col  []*SoftResource
	sort []string
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

// Range returns a subset of the collection arranged according to the
// given parameters.
func (s *SoftCollection) Range(ids []string, filter *Filter, sort []string, fields []string, pageSize uint, pageNumber uint) []Resource {
	rangeCol := &SoftCollection{}
	rangeCol.SetType(s.Type)

	// Filter IDs
	if len(ids) > 0 {
		for _, rec := range s.col {
			for _, id := range ids {
				if rec.id == id {
					rangeCol.Add(rec)
				}
			}
		}
	} else {
		for _, rec := range s.col {
			rangeCol.Add(rec)
		}
	}

	// Filter
	if filter != nil {
		i := 0
		for i < len(rangeCol.col) {
			if !filter.IsAllowed(rangeCol.col[i]) {
				rangeCol.col = append(rangeCol.col[:i], rangeCol.col[i+1:]...)
			} else {
				i++
			}
		}
	}

	// Sort
	rangeCol.Sort(sort)

	// Pagination
	skip := int(pageNumber * pageSize)
	if skip >= len(rangeCol.col) {
		rangeCol = &SoftCollection{}
	} else {
		page := &SoftCollection{}
		page.SetType(s.Type)
		for i := skip; i < len(rangeCol.col) && i < skip+int(pageSize); i++ {
			page.Add(rangeCol.col[i])
		}
		rangeCol = page
	}

	// Make the collection
	col := []Resource{}
	for _, rec := range rangeCol.col {
		col = append(col, rec)
	}

	return col
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
			sr.SetToOne(rel.Name, r.GetToOne(rel.Name))
		} else {
			sr.SetToMany(rel.Name, r.GetToMany(rel.Name))
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

// UnmarshalJSON populates a SoftCollection from the given payload.
//
// Only the attributes and relationships defined in the SoftCollection's Type
// field will be considered.
func (s *SoftCollection) UnmarshalJSON(payload []byte) error {
	// TODO Implement this method
	return errors.New("jsonapi: SoftCollection.UnmarshalJSON not yet implemented")
}

// Sort rearranges the order of the collection according the rules.
func (s *SoftCollection) Sort(rules []string) {
	s.sort = rules

	if len(s.sort) == 0 {
		s.sort = []string{"id"}
	}

	sort.Sort(s)
}

// Swap implements sort.Interface's Swap method.
func (s *SoftCollection) Swap(i, j int) {
	s.col[i], s.col[j] = s.col[j], s.col[i]
}

// Less implements sort.Interface's Less method.
func (s *SoftCollection) Less(i, j int) bool {
	for _, r := range s.sort {
		inverse := false
		if strings.HasPrefix(r, "-") {
			r = r[1:]
			inverse = true
		}

		if r == "id" {
			return s.col[i].GetID() < s.col[j].GetID() != inverse
		}

		// Here we return true if v < v2.
		// The "!= inverse" part acts as a XOR operation so that
		// the opposite boolean is returned when inverse sorting
		// is required.
		switch v := s.col[i].data[r].(type) {
		case string:
			v2 := s.col[j].data[r].(string)
			if v == v2 {
				continue
			}
			return v < v2 != inverse
		case int:
			v2 := s.col[j].data[r].(int)
			if v == v2 {
				continue
			}
			return v < v2 != inverse
		case int8:
			v2 := s.col[j].data[r].(int8)
			if v == v2 {
				continue
			}
			return v < v2 != inverse
		case int16:
			v2 := s.col[j].data[r].(int16)
			if v == v2 {
				continue
			}
			return v < v2 != inverse
		case int32:
			v2 := s.col[j].data[r].(int32)
			if v == v2 {
				continue
			}
			return v < v2 != inverse
		case int64:
			v2 := s.col[j].data[r].(int64)
			if v == v2 {
				continue
			}
			return v < v2 != inverse
		case uint:
			v2 := s.col[j].data[r].(uint)
			if v == v2 {
				continue
			}
			return v < v2 != inverse
		case uint8:
			v2 := s.col[j].data[r].(uint8)
			if v == v2 {
				continue
			}
			return v < v2 != inverse
		case uint16:
			v2 := s.col[j].data[r].(uint16)
			if v == v2 {
				continue
			}
			return v < v2 != inverse
		case uint32:
			v2 := s.col[j].data[r].(uint32)
			if v == v2 {
				continue
			}
			return v < v2 != inverse
		case bool:
			v2 := s.col[j].data[r].(bool)
			if v == v2 {
				continue
			}
			return !v != inverse
		case time.Time:
			if v.Equal(s.col[j].data[r].(time.Time)) {
				continue
			}
			return v.Before(s.col[j].data[r].(time.Time)) != inverse
		case *string:
			p := s.col[j].data[r].(*string)
			if v == p {
				continue
			}
			if v == nil {
				return !inverse
			}
			if p == nil {
				return inverse
			}
			if *v == *p {
				continue
			}
			return *v < *p != inverse
		case *int:
			p := s.col[j].data[r].(*int)
			if v == p {
				continue
			}
			if v == nil {
				return !inverse
			}
			if p == nil {
				return inverse
			}
			if *v == *p {
				continue
			}
			return *v < *p != inverse
		case *int8:
			p := s.col[j].data[r].(*int8)
			if v == p {
				continue
			}
			if v == nil {
				return !inverse
			}
			if p == nil {
				return inverse
			}
			if *v == *p {
				continue
			}
			return *v < *p != inverse
		case *int16:
			p := s.col[j].data[r].(*int16)
			if v == p {
				continue
			}
			if v == nil {
				return !inverse
			}
			if p == nil {
				return inverse
			}
			if *v == *p {
				continue
			}
			return *v < *p != inverse
		case *int32:
			p := s.col[j].data[r].(*int32)
			if v == p {
				continue
			}
			if v == nil {
				return !inverse
			}
			if p == nil {
				return inverse
			}
			if *v == *p {
				continue
			}
			return *v < *p != inverse
		case *int64:
			p := s.col[j].data[r].(*int64)
			if v == p {
				continue
			}
			if v == nil {
				return !inverse
			}
			if p == nil {
				return inverse
			}
			if *v == *p {
				continue
			}
			return *v < *p != inverse
		case *uint:
			p := s.col[j].data[r].(*uint)
			if v == p {
				continue
			}
			if v == nil {
				return !inverse
			}
			if p == nil {
				return inverse
			}
			if *v == *p {
				continue
			}
			return *v < *p != inverse
		case *uint8:
			p := s.col[j].data[r].(*uint8)
			if v == p {
				continue
			}
			if v == nil {
				return !inverse
			}
			if p == nil {
				return inverse
			}
			if *v == *p {
				continue
			}
			return *v < *p != inverse
		case *uint16:
			p := s.col[j].data[r].(*uint16)
			if v == p {
				continue
			}
			if v == nil {
				return !inverse
			}
			if p == nil {
				return inverse
			}
			if *v == *p {
				continue
			}
			return *v < *p != inverse
		case *uint32:
			p := s.col[j].data[r].(*uint32)
			if v == p {
				continue
			}
			if v == nil {
				return !inverse
			}
			if p == nil {
				return inverse
			}
			if *v == *p {
				continue
			}
			return *v < *p != inverse
		case *bool:
			p := s.col[j].data[r].(*bool)
			if v == p {
				continue
			}
			if v == nil {
				return !inverse
			}
			if p == nil {
				return inverse
			}
			if *v == *p {
				continue
			}
			return !*v != inverse
		case *time.Time:
			p := s.col[j].data[r].(*time.Time)
			if v == p {
				continue
			}
			if v == nil {
				return !inverse
			}
			if p == nil {
				return inverse
			}
			if v.Equal(*p) {
				continue
			}
			return v.Before(*p) != inverse
		}
	}

	return false
}
