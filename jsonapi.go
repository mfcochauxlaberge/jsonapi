package jsonapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Marshal ...
func Marshal(v interface{}, url *URL, extra Extra) ([]byte, error) {
	// Document
	doc := &Document{
		URL:     url,
		Meta:    extra.Meta,
		JSONAPI: extra.JSONAPI,
	}

	if res, ok := v.(Resource); ok {
		doc.Resource = res
	} else if col, ok := v.(Collection); ok {
		doc.Collection = col
	} else if ident, ok := v.(Identifier); ok {
		doc.Identifier = ident
	} else if idents, ok := v.(Identifiers); ok {
		doc.Identifiers = idents
	} else if err, ok := v.(Error); ok {
		doc.Errors = Errors{err}
	} else if errs, ok := v.(Errors); ok {
		doc.Errors = errs
	} else {
		panic(fmt.Errorf("jsonapi: cannot marshal unsupported type %s", reflect.ValueOf(v).Type().String()))
	}

	return json.Marshal(doc)
}

// Unmarshal ...
func Unmarshal(payload []byte, v interface{}) error {
	ske := &documentSkeleton{}

	// Unmarshal
	err := json.Unmarshal(payload, ske)
	if err != nil {
		return err
	}

	var res Resource
	var col Collection
	var ok bool
	if res, ok = v.(Resource); ok {
		err = json.Unmarshal(ske.Data, &res)
		if err != nil {
			return err
		}
	} else if col, ok = v.(Collection); ok {
		err = json.Unmarshal(ske.Data, &col)
		if err != nil {
			return err
		}
	}

	return nil
}

// CheckType ...
func CheckType(v interface{}) error {
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
			case "string", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "bool", "time.Time", "*string", "*int", "*int8", "*int16", "*int32", "*int64", "*uint", "*uint8", "*uint16", "*uint32", "*uint64", "*bool", "*time.Time":
				isValid = true
			}

			if !isValid {
				return fmt.Errorf("jsonapi: attribute %s of type %s is of unsupported type", sf.Name, resType)
			}
		}
	}

	// Check relationships
	for i := 0; i < value.NumField(); i++ {
		sf := value.Type().Field(i)

		if strings.HasPrefix(sf.Tag.Get("api"), "rel,") {
			s := strings.Split(sf.Tag.Get("api"), ",")

			if len(s) < 2 || len(s) > 3 {
				return fmt.Errorf("jsonapi: api tag of relationship %s of struct %s is invalid", sf.Name, value.Type().Name())
			}

			if sf.Type.String() != "string" && sf.Type.String() != "[]string" {
				return fmt.Errorf("jsonapi: relationship %s of type %s is not string or []string", sf.Name, resType)
			}
		}
	}

	return nil
}

// IDAndType returns the ID and the type of the resource represented by v.
//
// Two empty strings are returned if v is not recognized as a resource.
// CheckType can be used to check the validity of a struct.
func IDAndType(v interface{}) (string, string) {
	switch nv := v.(type) {
	case Resource:
		return nv.IDAndType()
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
