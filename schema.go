package jsonapi

import "database/sql"

// Type ...
type Type struct {
	Name    string
	Fields  []string
	Attrs   map[string]Attr
	Rels    map[string]Rel
	Default Resource
}

// Attr ...
type Attr struct {
	Name    string
	Type    string
	Null    bool
	Default sql.NullString
}

// Rel ...
type Rel struct {
	Name         string
	Type         string
	ToOne        bool
	InverseName  string
	InverseType  string
	InverseToOne bool
}

// Inverse ...
func (r *Rel) Inverse() Rel {
	return Rel{
		Name:         r.InverseName,
		Type:         r.InverseType,
		ToOne:        r.InverseToOne,
		InverseName:  r.Name,
		InverseType:  r.Type,
		InverseToOne: r.ToOne,
	}
}
