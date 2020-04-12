package jsonapi

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Check checks that the given value can be used with this library and returns
// the first error it finds.
//
// It makes sure that the struct has an ID field of type string and that the api
// key of the field tags are properly formatted.
//
// If nil is returned, then the value can be safely used with this library.
func Check(v interface{}) error {
	value := reflect.ValueOf(v)
	kind := value.Kind()

	// Check wether it's a struct
	if kind != reflect.Struct {
		return errors.New("jsonapi: not a struct")
	}

	// Check ID field
	var (
		idField reflect.StructField
		ok      bool
	)

	if idField, ok = value.Type().FieldByName("ID"); !ok {
		return errors.New("jsonapi: struct doesn't have an ID field")
	}

	resType := idField.Tag.Get("api")
	if resType == "" {
		return errors.New("jsonapi: ID field's api tag is empty")
	}

	// Check attributes
	for i := 0; i < value.NumField(); i++ {
		sf := value.Type().Field(i)

		if sf.Tag.Get("api") == "attr" {
			isValid := false

			switch sf.Type.String() {
			case
				"string",
				"int", "int8", "int16", "int32", "int64",
				"uint", "uint8", "uint16", "uint32", "uint64",
				"bool",
				"time.Time",
				"[]uint8",
				"*string",
				"*int", "*int8", "*int16", "*int32", "*int64",
				"*uint", "*uint8", "*uint16", "*uint32", "*uint64",
				"*bool",
				"*time.Time",
				"*[]uint8":
				isValid = true
			}

			if !isValid {
				return fmt.Errorf(
					"jsonapi: attribute %q of type %q is of unsupported type",
					sf.Name,
					resType,
				)
			}
		}
	}

	// Check relationships
	for i := 0; i < value.NumField(); i++ {
		sf := value.Type().Field(i)

		if strings.HasPrefix(sf.Tag.Get("api"), "rel,") {
			s := strings.Split(sf.Tag.Get("api"), ",")

			if len(s) < 2 || len(s) > 3 {
				return fmt.Errorf(
					"jsonapi: api tag of relationship %q of struct %q is invalid",
					sf.Name,
					value.Type().Name(),
				)
			}

			if sf.Type.String() != "string" && sf.Type.String() != "[]string" {
				return fmt.Errorf(
					"jsonapi: relationship %q of type %q is not string or []string",
					sf.Name,
					resType,
				)
			}
		}
	}

	return nil
}

// BuildType takes a struct or a pointer to a struct to analyse and builds a
// Type object that is returned.
//
// If an error is returned, the Type object will be empty.
func BuildType(v interface{}) (Type, error) {
	typ := Type{}

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return typ, errors.New("jsonapi: value must represent a struct")
	}

	err := Check(val.Interface())
	if err != nil {
		return typ, fmt.Errorf("jsonapi: invalid type: %q", err)
	}

	// ID and type
	_, typ.Name = IDAndType(v)

	// Attributes
	typ.Attrs = map[string]Attr{}

	for i := 0; i < val.NumField(); i++ {
		fs := val.Type().Field(i)
		jsonTag := fs.Tag.Get("json")
		apiTag := fs.Tag.Get("api")

		if apiTag == "attr" {
			fieldType, null := GetAttrType(fs.Type.String())
			typ.Attrs[jsonTag] = Attr{
				Name:     jsonTag,
				Type:     fieldType,
				Nullable: null,
			}
		}
	}

	// Relationships
	typ.Rels = map[string]Rel{}

	for i := 0; i < val.NumField(); i++ {
		fs := val.Type().Field(i)
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
			typ.Rels[jsonTag] = Rel{
				FromName: jsonTag,
				ToOne:    toOne,
				ToType:   relTag[1],
				ToName:   invName,
				FromType: typ.Name,
			}
		}
	}

	// NewFunc
	res := Wrap(reflect.New(val.Type()).Interface())
	typ.NewFunc = res.Copy

	return typ, nil
}

// MustBuildType calls BuildType and returns the result.
//
// It panics if the error returned by BuildType is not nil.
func MustBuildType(v interface{}) Type {
	typ, err := BuildType(v)
	if err != nil {
		panic(err)
	}

	return typ
}

// IDAndType returns the ID and the type of the resource represented by v.
//
// Two empty strings are returned if v is not recognized as a resource.
// CheckType can be used to check the validity of a struct.
func IDAndType(v interface{}) (string, string) {
	if res, ok := v.(Resource); ok {
		return res.GetID(), res.GetType().Name
	}

	val := reflect.ValueOf(v)

	// Allows pointers to structs
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() == reflect.Struct {
		idF := val.FieldByName("ID")

		if !idF.IsValid() {
			return "", ""
		}

		idSF, _ := val.Type().FieldByName("ID")

		if idF.Kind() == reflect.String {
			return idF.String(), idSF.Tag.Get("api")
		}
	}

	return "", ""
}
