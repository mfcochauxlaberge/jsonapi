package jsonapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func Marshal(doc *Document, url *URL) ([]byte, error) {
	// Data
	var data json.RawMessage
	var errors json.RawMessage
	var err error

	if res, ok := doc.Data.(Resource); ok {
		// Resource
		_, typ := res.IDAndType()
		data, err = marshalResource(res, url.Host, url.Params.Fields[typ], doc.RelData)
	} else if col, ok := doc.Data.(Collection); ok {
		// Collection
		data, err = marshalCollection(col, url.Host, url.Params.Fields[col.Type()], doc.RelData)
	} else if id, ok := doc.Data.(Identifier); ok {
		// Identifer
		data, err = json.Marshal(id)
	} else if ids, ok := doc.Data.(Identifiers); ok {
		// Identifiers
		data, err = json.Marshal(ids)
	} else if e, ok := doc.Data.(Error); ok {
		// Error
		errors, err = json.Marshal([]Error{e})
	} else if es, ok := doc.Data.([]Error); ok {
		// Errors
		errors, err = json.Marshal(es)
	} else {
		data = []byte("null")
	}

	if err != nil {
		return []byte{}, err
	}

	// Included
	inclusions := []*json.RawMessage{}
	if len(data) > 0 {
		for key := range doc.Included {
			_, typ := doc.Included[key].IDAndType()
			raw, err := marshalResource(doc.Included[key], url.Host, url.Params.Fields[typ], doc.RelData)
			if err != nil {
				return []byte{}, err
			}
			rawm := json.RawMessage(raw)
			inclusions = append(inclusions, &rawm)
		}
	}

	// Marshaling
	plMap := map[string]interface{}{}

	if len(errors) > 0 {
		plMap["errors"] = errors
	} else if len(data) > 0 {
		plMap["data"] = data
	}

	if len(doc.Links) > 0 {
		plMap["links"] = doc.Links
	}

	if len(inclusions) > 0 {
		plMap["included"] = inclusions
	}

	if len(doc.Meta) > 0 {
		plMap["meta"] = doc.Meta
	}

	if len(doc.JSONAPI) > 0 {
		plMap["jsonapi"] = doc.JSONAPI
	}

	return json.Marshal(plMap)
}

// Unmarshal ...
func Unmarshal(payload []byte, v interface{}) (*Document, error) {
	doc := &Document{}
	ske := &documentSkeleton{}

	// Unmarshal
	err := json.Unmarshal(payload, ske)
	if err != nil {
		return nil, err
	}

	// Resource or collection
	if v == nil {
	} else if res, ok := v.(Resource); ok {
		err = json.Unmarshal(ske.Data, v)
		if err != nil {
			return nil, err
		}
		doc.Data = res
	} else if col, ok := v.(Collection); ok {
		err = json.Unmarshal(ske.Data, v)
		if err != nil {
			return nil, err
		}
		doc.Data = col
	} else if id, ok := v.(Identifier); ok {
		err = json.Unmarshal(ske.Data, v)
		if err != nil {
			return nil, err
		}
		doc.Data = id
	} else if ids, ok := v.(Identifiers); ok {
		err = json.Unmarshal(ske.Data, v)
		if err != nil {
			return nil, err
		}
		doc.Data = ids
	} else {
		panic("v in Unmarshal doest not implement Resource or Collection")
	}

	doc.Meta = ske.Meta
	doc.JSONAPI = ske.JSONAPI

	return doc, nil
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
