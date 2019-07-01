package jsonapi

import (
	"errors"
	"fmt"
)

// A Schema contains a list of types. It makes sure that each type is
// valid and unique.
//
// Check can be used to validate the relationships between the types.
type Schema struct {
	Types []Type
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
func (s *Schema) RemoveType(typ string) error {
	for i := range s.Types {
		if s.Types[i].Name == typ {
			s.Types = append(s.Types[0:i], s.Types[i+1:]...)
		}
	}

	return nil
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
func (s *Schema) RemoveAttr(typ string, attr string) error {
	for i := range s.Types {
		if s.Types[i].Name == typ {
			return s.Types[i].RemoveAttr(attr)
		}
	}

	return fmt.Errorf("jsonapi: type %s does not exist", typ)
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
func (s *Schema) RemoveRel(typ string, rel string) error {
	for i := range s.Types {
		if s.Types[i].Name == typ {
			return s.Types[i].RemoveRel(rel)
		}
	}

	return fmt.Errorf("jsonapi: type %s does not exist", typ)
}

// HasType returns a boolean indicating whether a type has the specified name
// or not.
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
// A boolean indicates whether a type was found or not.
func (s *Schema) GetType(name string) (Type, bool) {
	for _, typ := range s.Types {
		if typ.Name == name {
			return typ, true
		}
	}
	return Type{}, false
}

// GetResource returns a resource of type SoftResource with the specified
// type. All fields are set to their zero values.
func (s *Schema) GetResource(name string) Resource {
	typ, ok := s.GetType(name)
	if ok {
		return NewSoftResource(typ, nil)
	}
	return nil
}

// Check checks the integrity of all the relationships between the types
// and returns all the errors that were found.
func (s *Schema) Check() []error {
	var (
		ok   bool
		errs = []error{}
	)

	// Check the inverse relationships
	for _, typ := range s.Types {
		// Relationships
		for _, rel := range typ.Rels {
			var targetType Type

			// Does the relationship point to a type that exists?
			if targetType, ok = s.GetType(rel.Type); !ok {
				errs = append(errs, fmt.Errorf(
					"jsonapi: the target type of relationship %s of type %s does not exist",
					rel.Name,
					typ.Name,
				))
			}

			// Inverse relationship (if relevant)
			if rel.InverseName != "" {
				// Is the inverse relationship type the same as its type name?
				if rel.InverseType != typ.Name {
					errs = append(errs, fmt.Errorf(
						"jsonapi: the inverse type of relationship %s should its type's name (%s, not %s)",
						rel.Name,
						typ.Name,
						rel.InverseType,
					))
				}

				// Do both relationships (current and inverse) point to each other?
				var found bool
				for _, invRel := range targetType.Rels {
					if rel.Name == invRel.InverseName && rel.InverseName == invRel.Name {
						found = true
					}
				}
				if !found {
					errs = append(errs, fmt.Errorf(
						"jsonapi: relationship %s of type %s and its inverse do not point each other",
						rel.Name,
						typ.Name,
					))
				}
			}

		}
	}

	return errs
}
