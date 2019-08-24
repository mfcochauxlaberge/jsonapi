package jsonapi

import (
	"errors"
	"fmt"
)

// A Schema contains a list of types. It makes sure that each type is valid and
// unique.
//
// Check can be used to validate the relationships between the types.
type Schema struct {
	Types []Type

	// Rels stores the relationships found in the schema's types. For
	// two-way relationships, only one is chosen to be part of this
	// map. The chosen one is the one that comes first when sorting
	// both relationships in alphabetical order using the type name
	// first and then the relationship name.
	//
	// For example, a type called Directory has a Parent relationship
	// and a Children relationship. Both relationships have the same
	// type (Directory), so now the name is used for sorting. Children
	// comes before Parent, so the relationship Children from type
	// Directory is stored here. The other one is not stored to avoid
	// duplication (the information is already accessible through the
	// inverse relationship).
	Rels map[string]Rel
}

// AddType adds a type to the schema.
func (s *Schema) AddType(typ Type) error {
	// Validation
	if typ.Name == "" {
		return errors.New("jsonapi: type name is empty")
	}

	// Make sure the name isn't already used
	for i := range s.Types {
		if s.Types[i].Name == typ.Name {
			return fmt.Errorf("jsonapi: type name %s is already used", typ.Name)
		}
	}

	s.Types = append(s.Types, typ)

	return nil
}

// RemoveType removes a type from the schema.
func (s *Schema) RemoveType(typ string) {
	for i := range s.Types {
		if s.Types[i].Name == typ {
			s.Types = append(s.Types[0:i], s.Types[i+1:]...)
		}
	}
}

// AddAttr adds an attribute to the specified type.
func (s *Schema) AddAttr(typ string, attr Attr) error {
	for i := range s.Types {
		if s.Types[i].Name == typ {
			return s.Types[i].AddAttr(attr)
		}
	}

	return fmt.Errorf("jsonapi: type %s does not exist", typ)
}

// RemoveAttr removes an attribute from the specified type.
func (s *Schema) RemoveAttr(typ string, attr string) {
	for i := range s.Types {
		if s.Types[i].Name == typ {
			s.Types[i].RemoveAttr(attr)
		}
	}
}

// AddRel adds a relationship to the specified type.
func (s *Schema) AddRel(typ string, rel Rel) error {
	for i := range s.Types {
		if s.Types[i].Name == typ {
			return s.Types[i].AddRel(rel)
		}
	}

	return fmt.Errorf("jsonapi: type %s does not exist", typ)
}

// RemoveRel removes a relationship from the specified type.
func (s *Schema) RemoveRel(typ string, rel string) {
	for i := range s.Types {
		if s.Types[i].Name == typ {
			s.Types[i].RemoveRel(rel)
		}
	}
}

// HasType returns a boolean indicating whether a type has the specified
// name or not.
func (s *Schema) HasType(name string) bool {
	for i := range s.Types {
		if s.Types[i].Name == name {
			return true
		}
	}
	return false
}

// GetType returns the type associated with the speficied name.
//
// If no type with the given name is found, an zero instance of Type is
// returned. Therefore, checking whether the Name field is empty or not is a
// good way to dertermine whether the type was found or not.
func (s *Schema) GetType(name string) Type {
	for _, typ := range s.Types {
		if typ.Name == name {
			return typ
		}
	}
	return Type{}
}

// Check checks the integrity of all the relationships between the types and
// returns all the errors that were found.
func (s *Schema) Check() []error {
	var (
		errs = []error{}
	)

	// Check the inverse relationships
	for _, typ := range s.Types {
		// Relationships
		for _, rel := range typ.Rels {
			var targetType Type

			// Does the relationship point to a type that exists?
			if targetType = s.GetType(rel.Type); targetType.Name == "" {
				errs = append(errs, fmt.Errorf(
					"jsonapi: the target type of relationship %s of type %s does not exist",
					rel.Name,
					typ.Name,
				))
			}

			// Skip to next relationship here if there's no inverse
			if rel.InverseName == "" {
				continue
			}

			// Is the inverse relationship type the same as its
			// type name?
			if rel.InverseType != typ.Name {
				errs = append(errs, fmt.Errorf(
					"jsonapi: "+
						"the inverse type of relationship %s should its type's name (%s, not %s)",
					rel.Name,
					typ.Name,
					rel.InverseType,
				))
			} else {
				// Do both relationships (current and inverse) point
				// to each other?
				var found bool
				for _, invRel := range targetType.Rels {
					if rel.Name == invRel.InverseName && rel.InverseName == invRel.Name {
						found = true
					}
				}
				if !found {
					errs = append(errs, fmt.Errorf(
						"jsonapi: "+
							"relationship %s of type %s and its inverse do not point each other",
						rel.Name,
						typ.Name,
					))
				}
			}
		}
	}

	return errs
}
