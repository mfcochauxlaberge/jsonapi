package jsonapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Marshal ...
func Marshal(doc *Document, url *URL) ([]byte, error) {
	// Data
	var data json.RawMessage
	var errors json.RawMessage
	var err error

	if res, ok := doc.Data.(Resource); ok {
		// Resource
		data, err = marshalResource(res, doc.PrePath, url.Params.Fields[res.GetType()], doc.RelData)
	} else if col, ok := doc.Data.(Collection); ok {
		// Collection
		data, err = marshalCollection(col, doc.PrePath, url.Params.Fields[col.Type()], doc.RelData)
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
			typ := doc.Included[key].GetType()
			raw, err := marshalResource(doc.Included[key], doc.PrePath, url.Params.Fields[typ], doc.RelData)
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

		if len(inclusions) > 0 {
			plMap["included"] = inclusions
		}
	}

	if len(doc.Meta) > 0 {
		plMap["meta"] = doc.Meta
	}

	if url != nil {
		plMap["links"] = map[string]string{
			"self": doc.PrePath + url.FullURL(),
		}
	}
	plMap["jsonapi"] = map[string]string{"version": "1.0"}

	return json.Marshal(plMap)
}

// Unmarshal ...
func Unmarshal(payload []byte, url *URL, schema *Schema) (*Payload, error) {
	pl := &Payload{}
	ske := &payloadSkeleton{}

	// Unmarshal
	err := json.Unmarshal(payload, ske)
	if err != nil {
		return nil, err
	}

	// Data
	if !url.IsCol && url.RelKind == "" {
		res := schema.GetResource(url.ResType)
		err = json.Unmarshal(ske.Data, res)
		if err != nil {
			return nil, err
		}
		pl.Data = res
	} else if url.RelKind == "self" {
		if !url.IsCol {
			inc := Identifier{}
			err = json.Unmarshal(ske.Data, &inc)
			if err != nil {
				return nil, err
			}
			pl.Data = inc
		} else {
			incs := Identifiers{}
			err = json.Unmarshal(ske.Data, &incs)
			if err != nil {
				return nil, err
			}
			pl.Data = incs
		}
	}

	// Included
	if len(ske.Included) > 0 {
		inc := Identifier{}
		incs := []Identifier{}
		for _, rawInc := range ske.Included {
			err = json.Unmarshal(rawInc, &inc)
			if err != nil {
				return nil, err
			}
			incs = append(incs, inc)
		}

		for i, inc2 := range incs {
			res2 := schema.GetResource(inc2.Type)
			err = json.Unmarshal(ske.Included[i], res2)
			if err != nil {
				return nil, err
			}
			pl.Included[inc2.Type+" "+inc2.ID] = res2
		}
	}

	// Meta
	pl.Meta = ske.Meta

	return pl, nil
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
		return nv.GetID(), nv.GetType()
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
