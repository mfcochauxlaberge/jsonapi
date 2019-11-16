package jsonapi

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Wrapper wraps a reflect.Value that represents a struct.
//
// The Wrap function can be used to wrap a struct and make a Wrapper object.
//
// It implements the Resource interface, so the value can be handled as if it
// were a Resource.
type Wrapper struct {
	val reflect.Value // Actual value (with content)

	// Structure
	typ   string
	attrs map[string]Attr
	rels  map[string]Rel
}

// Wrap wraps v (a struct or a pointer to a struct) and returns a Wrapper that
// can be used as a Resource to handle the given value.
//
// If v is not a pointer, the changes applied to the Wrapper object won't affect
// the underlying object (which will be a new instance of v's type).
func Wrap(v interface{}) *Wrapper {
	val := reflect.ValueOf(v)

	if val.Kind() != reflect.Ptr {
		if val.Kind() != reflect.Struct {
			panic(errors.New("jsonapi: value has to be a pointer to a struct"))
		}

		val = reflect.New(val.Type())
	} else if val.Elem().Kind() != reflect.Struct {
		panic(errors.New("jsonapi: value has to be a pointer to a struct"))
	}

	val = val.Elem()

	err := Check(val.Interface())
	if err != nil {
		panic(fmt.Sprintf("jsonapi: invalid type: %s", err))
	}

	w := &Wrapper{
		val: val,
	}

	// ID and type
	_, w.typ = IDAndType(v)

	// Attributes
	w.attrs = map[string]Attr{}
	for i := 0; i < w.val.NumField(); i++ {
		fs := w.val.Type().Field(i)
		jsonTag := fs.Tag.Get("json")
		apiTag := fs.Tag.Get("api")

		if apiTag == "attr" {
			typ, null := GetAttrType(fs.Type.String())
			w.attrs[jsonTag] = Attr{
				Name:     jsonTag,
				Type:     typ,
				Nullable: null,
			}
		}
	}

	// Relationships
	w.rels = map[string]Rel{}
	for i := 0; i < w.val.NumField(); i++ {
		fs := w.val.Type().Field(i)
		jsonTag := fs.Tag.Get("json")
		relTag := strings.Split(fs.Tag.Get("api"), ",")
		invName := ""

		if len(relTag) == 3 {
			invName = relTag[2]
		}

		toOne := true
		if fs.Type.String() == "[]string" {
			toOne = false
		}

		if relTag[0] == "rel" {
			w.rels[jsonTag] = Rel{
				FromName: jsonTag,
				ToType:   relTag[1],
				ToOne:    toOne,
				ToName:   invName,
				FromType: w.typ,
			}
		}
	}

	return w
}

// IDAndType returns the ID and the type of the Wrapper.
func (w *Wrapper) IDAndType() (string, string) {
	return IDAndType(w.val.Interface())
}

// Attrs returns the attributes of the Wrapper.
func (w *Wrapper) Attrs() map[string]Attr {
	return w.attrs
}

// Rels returns the relationships of the Wrapper.
func (w *Wrapper) Rels() map[string]Rel {
	return w.rels
}

// Attr returns the attribute that corresponds to the given key.
func (w *Wrapper) Attr(key string) Attr {
	for _, attr := range w.attrs {
		if attr.Name == key {
			return attr
		}
	}

	return Attr{}
}

// Rel returns the relationship that corresponds to the given key.
func (w *Wrapper) Rel(key string) Rel {
	for _, rel := range w.rels {
		if rel.FromName == key {
			return rel
		}
	}

	return Rel{}
}

// New returns a copy of the resource under the wrapper.
func (w *Wrapper) New() Resource {
	newVal := reflect.New(w.val.Type())

	return Wrap(newVal.Interface())
}

// GetID returns the wrapped resource's ID.
func (w *Wrapper) GetID() string {
	id, _ := IDAndType(w.val.Interface())
	return id
}

// GetType returns the wrapped resource's type.
func (w *Wrapper) GetType() Type {
	return Type{
		Name:  w.typ,
		Attrs: w.attrs,
		Rels:  w.rels,
	}
}

// Get returns the value associated to the attribute named after key.
func (w *Wrapper) Get(key string) interface{} {
	return w.getAttr(key)
}

// SetID sets the ID of the wrapped resource.
func (w *Wrapper) SetID(id string) {
	w.val.FieldByName("ID").SetString(id)
}

// Set sets the value associated to the attribute named after key.
func (w *Wrapper) Set(key string, val interface{}) {
	w.setAttr(key, val)
}

// GetToOne returns the value associated with the relationship named after key.
func (w *Wrapper) GetToOne(key string) string {
	for i := 0; i < w.val.NumField(); i++ {
		field := w.val.Field(i)
		sf := w.val.Type().Field(i)

		if key == sf.Tag.Get("json") {
			if strings.Split(sf.Tag.Get("api"), ",")[0] != "rel" {
				panic(fmt.Sprintf("jsonapi: field %q is not a relationship", key))
			}

			if field.Type().String() != "string" {
				panic(fmt.Sprintf("jsonapi: relationship %q is not 'to one'", key))
			}

			return field.String()
		}
	}

	if key == "" {
		panic("jsonapi: key is empty")
	}

	panic(fmt.Sprintf("jsonapi: relationship %q does not exist", key))
}

// GetToMany returns the value associated with the relationship named after key.
func (w *Wrapper) GetToMany(key string) []string {
	for i := 0; i < w.val.NumField(); i++ {
		field := w.val.Field(i)
		sf := w.val.Type().Field(i)

		if key == sf.Tag.Get("json") {
			if strings.Split(sf.Tag.Get("api"), ",")[0] != "rel" {
				panic(fmt.Sprintf("jsonapi: field %q is not a relationship", key))
			}

			if field.Type().String() != "[]string" {
				panic(fmt.Sprintf("jsonapi: relationship %q is not 'to many'", key))
			}

			return field.Interface().([]string)
		}
	}

	if key == "" {
		panic("jsonapi: key is empty")
	}

	panic(fmt.Sprintf("jsonapi: relationship %q does not exist", key))
}

// SetToOne sets the value associated to the relationship named after key.
func (w *Wrapper) SetToOne(key string, rel string) {
	for i := 0; i < w.val.NumField(); i++ {
		field := w.val.Field(i)
		sf := w.val.Type().Field(i)

		if key == sf.Tag.Get("json") {
			if strings.Split(sf.Tag.Get("api"), ",")[0] != "rel" {
				panic(fmt.Sprintf("jsonapi: field %q is not a relationship", key))
			}

			if field.Type().String() != "string" {
				panic(fmt.Sprintf("jsonapi: relationship %q is not 'to one'", key))
			}

			field.SetString(rel)

			return
		}
	}

	if key == "" {
		panic("jsonapi: key is empty")
	}

	panic(fmt.Sprintf("jsonapi: relationship %q does not exist", key))
}

// SetToMany sets the value associated to the relationship named after key.
func (w *Wrapper) SetToMany(key string, rels []string) {
	for i := 0; i < w.val.NumField(); i++ {
		field := w.val.Field(i)
		sf := w.val.Type().Field(i)

		if key == sf.Tag.Get("json") {
			if strings.Split(sf.Tag.Get("api"), ",")[0] != "rel" {
				panic(fmt.Sprintf("jsonapi: field %q is not a relationship", key))
			}

			if field.Type().String() != "[]string" {
				panic(fmt.Sprintf("jsonapi: relationship %q is not 'to many'", key))
			}

			field.Set(reflect.ValueOf(rels))

			return
		}
	}

	if key == "" {
		panic("jsonapi: key is empty")
	}

	panic(fmt.Sprintf("jsonapi: relationship %q does not exist", key))
}

// Copy makes a copy of the wrapped resource and returns it.
//
// The returned value's concrete type is also a Wrapper.
func (w *Wrapper) Copy() Resource {
	nw := Wrap(reflect.New(w.val.Type()).Interface())

	// Attributes
	for _, attr := range w.Attrs() {
		nw.Set(attr.Name, w.Get(attr.Name))
	}

	// Relationships
	for _, rel := range w.Rels() {
		if rel.ToOne {
			nw.SetToOne(rel.FromName, w.GetToOne(rel.FromName))
		} else {
			nw.SetToMany(rel.FromName, w.GetToMany(rel.FromName))
		}
	}

	return nw
}

// Private methods

func (w *Wrapper) getAttr(key string) interface{} {
	for i := 0; i < w.val.NumField(); i++ {
		field := w.val.Field(i)
		sf := w.val.Type().Field(i)

		if key == sf.Tag.Get("json") && sf.Tag.Get("api") == "attr" {
			if strings.HasPrefix(field.Type().String(), "*") && field.IsNil() {
				return nil
			}

			return field.Interface()
		}
	}

	if key == "" {
		panic("jsonapi: key is empty")
	}

	panic(fmt.Sprintf("jsonapi: attribute %q does not exist", key))
}

func (w *Wrapper) setAttr(key string, v interface{}) {
	for i := 0; i < w.val.NumField(); i++ {
		field := w.val.Field(i)
		sf := w.val.Type().Field(i)

		if key == sf.Tag.Get("json") {
			if v == nil {
				field.Set(reflect.New(field.Type()).Elem())
				return
			}

			val := reflect.ValueOf(v)
			if val.Type() == field.Type() {
				field.Set(val)
				return
			}

			panic(fmt.Sprintf("jsonapi: value is of wrong type (expected %q, got %q)",
				field.Type(),
				val.Type(),
			))
		}
	}

	if key == "" {
		panic("jsonapi: key is empty")
	}

	panic(fmt.Errorf("jsonapi: attribute %q does not exist", key))
}
