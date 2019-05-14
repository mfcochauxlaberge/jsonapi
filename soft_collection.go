package jsonapi

import (
	"errors"
	"sort"
	"strings"
	"sync"
	"time"
)

// SoftCollection ...
type SoftCollection struct {
	typ  *Type
	col  []*SoftResource
	sort []string

	m sync.Mutex
}

// Type ...
func (s *SoftCollection) Type() Type {
	return *s.typ
}

// SetType ...
func (s *SoftCollection) SetType(typ *Type) {
	s.typ = typ
}

// AddAttr ...
func (s *SoftCollection) AddAttr(attr Attr) error {
	return s.typ.AddAttr(attr)
}

// AddRel ...
func (s *SoftCollection) AddRel(rel Rel) error {
	return s.typ.AddRel(rel)

}

// Len ...
func (s *SoftCollection) Len() int {
	return len(s.col)
}

// Elem ...
func (s *SoftCollection) Elem(i int) Resource {
	if i >= 0 && i < len(s.col) {
		return s.col[i]
	}
	return nil
}

// Resource ...
func (s *SoftCollection) Resource(id string, fields []string) Resource {
	for i := range s.col {
		if s.col[i].GetID() == id {
			sr := s.col[i].Copy()
			return sr
		}
	}
	return nil
}

// Add ...
func (s *SoftCollection) Add(r Resource) {
	// A SoftResource is built from the Resource and
	// then it is added to the collection.
	sr := &SoftResource{}
	sr.id = r.GetID()
	sr.typ = s.typ

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

// UnmarshalJSON ...
func (s *SoftCollection) UnmarshalJSON(payload []byte) error {
	return errors.New("jsonapi: SoftCollection.UnmarshalJSON unimplemented")
}

// Sort ...
func (s *SoftCollection) Sort(rules []string) {
	s.sort = rules

	if len(s.sort) == 0 {
		s.sort = []string{"id"}
	}

	sort.Sort(s)
}

// Swap ...
func (s *SoftCollection) Swap(i, j int) {
	s.col[i], s.col[j] = s.col[j], s.col[i]
}

// Less ...
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
				return true != inverse
			}
			if p == nil {
				return false != inverse
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
				return true != inverse
			}
			if p == nil {
				return false != inverse
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
				return true != inverse
			}
			if p == nil {
				return false != inverse
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
				return true != inverse
			}
			if p == nil {
				return false != inverse
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
				return true != inverse
			}
			if p == nil {
				return false != inverse
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
				return true != inverse
			}
			if p == nil {
				return false != inverse
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
				return true != inverse
			}
			if p == nil {
				return false != inverse
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
				return true != inverse
			}
			if p == nil {
				return false != inverse
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
				return true != inverse
			}
			if p == nil {
				return false != inverse
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
				return true != inverse
			}
			if p == nil {
				return false != inverse
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
				return true != inverse
			}
			if p == nil {
				return false != inverse
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
				return true != inverse
			}
			if p == nil {
				return false != inverse
			}
			if v.Equal(*p) {
				continue
			}
			return v.Before(*p) != inverse
		}
	}

	return false
}
