package jsonapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Wrapper ...
type Wrapper struct {
	val reflect.Value // Actual value (with content)

	// Structure
	typ   string
	attrs []Attr
	rels  map[string]Rel
}

// Wrap ...
func Wrap(v interface{}) *Wrapper {
	val := reflect.ValueOf(v)

	if val.Kind() != reflect.Ptr {
		panic(errors.New("jsonapi: value has to be a pointer to a struct"))
	}

	if val.Elem().Kind() != reflect.Struct {
		panic(errors.New("jsonapi: value has to be a pointer to a struct"))
	}

	val = val.Elem()

	err := CheckType(val.Interface())
	if err != nil {
		panic(fmt.Sprintf("jsonapi: invalid type: %s", err))
	}

	w := &Wrapper{
		val: val,
	}

	// ID and type
	_, w.typ = IDAndType(v)

	// Attributes
	w.attrs = []Attr{}
	for i := 0; i < w.val.NumField(); i++ {
		fs := w.val.Type().Field(i)
		jsonTag := fs.Tag.Get("json")
		apiTag := fs.Tag.Get("api")

		if apiTag == "attr" {
			w.attrs = append(w.attrs, Attr{
				Name: jsonTag,
				Type: fs.Type.String(),
			})
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
				Name:        jsonTag,
				Type:        relTag[1],
				ToOne:       toOne,
				InverseName: invName,
				InverseType: w.typ,
			}
		}
	}

	return w
}

// IDAndType ...
func (w *Wrapper) IDAndType() (string, string) {
	return IDAndType(w.val.Interface())
}

// Attrs ...
func (w *Wrapper) Attrs() []Attr {
	return w.attrs
}

// Rels ...
func (w *Wrapper) Rels() map[string]Rel {
	return w.rels
}

// Attr ...
func (w *Wrapper) Attr(key string) Attr {
	for _, attr := range w.attrs {
		if attr.Name == key {
			return attr
		}
	}

	panic(fmt.Sprintf("jsonapi: attribute %s does not exist", key))
}

// Rel ...
func (w *Wrapper) Rel(key string) Rel {
	for _, rel := range w.rels {
		if rel.Name == key {
			return rel
		}
	}

	panic(fmt.Sprintf("jsonapi: relationship %s does not exist", key))
}

// // AttrPtrs ...
// func (w *Wrapper) AttrPtrs(attrs []Attr) ([]string, []interface{}) {
// 	keys := []string{}
// 	dsts := []interface{}{}
//
// 	for i := range attrs {
// 		for j := range w.Attrs() {
// 			if attrs[i].Name == w.Attrs()[j].Name {
// 				dest[attrs[i].Name] = &sql.NullString{}
// 			}
// 		}
// 	}
//
// 	return dest
// }
//
// // ApplyAttrPtrs ...
// func (w *Wrapper) ApplyAttrPtrs(ptrs map[string]interface{}) {
// 	// TODO
// }

// New ...
func (w *Wrapper) New() Resource {
	newVal := reflect.New(w.val.Type())

	return Wrap(newVal.Interface())
}

// Get ...
func (w *Wrapper) Get(key string) interface{} {
	return w.getAttr(key, "")
}

// GetID ...
func (w *Wrapper) GetID() string {
	return w.getAttr("", "string").(string)
}

// GetString ...
func (w *Wrapper) GetString(key string) string {
	return w.getAttr(key, "string").(string)
}

// GetStringPtr ...
func (w *Wrapper) GetStringPtr(key string) *string {
	return w.getAttr(key, "*string").(*string)
}

// GetInt ...
func (w *Wrapper) GetInt(key string) int {
	return w.getAttr(key, "int").(int)
}

// GetIntPtr ...
func (w *Wrapper) GetIntPtr(key string) *int {
	return w.getAttr(key, "*int").(*int)
}

// GetInt8 ...
func (w *Wrapper) GetInt8(key string) int8 {
	return w.getAttr(key, "int8").(int8)
}

// GetInt8Ptr ...
func (w *Wrapper) GetInt8Ptr(key string) *int8 {
	return w.getAttr(key, "*int8").(*int8)
}

// GetInt16 ...
func (w *Wrapper) GetInt16(key string) int16 {
	return w.getAttr(key, "int16").(int16)
}

// GetInt16Ptr ...
func (w *Wrapper) GetInt16Ptr(key string) *int16 {
	return w.getAttr(key, "int16").(*int16)
}

// GetInt32 ...
func (w *Wrapper) GetInt32(key string) int32 {
	return w.getAttr(key, "int32").(int32)
}

// GetInt32Ptr ...
func (w *Wrapper) GetInt32Ptr(key string) *int32 {
	return w.getAttr(key, "*int32").(*int32)
}

// GetInt64 ...
func (w *Wrapper) GetInt64(key string) int64 {
	return w.getAttr(key, "int64").(int64)
}

// GetInt64Ptr ...
func (w *Wrapper) GetInt64Ptr(key string) *int64 {
	return w.getAttr(key, "*int64").(*int64)
}

// GetUint ...
func (w *Wrapper) GetUint(key string) uint {
	return w.getAttr(key, "uint").(uint)
}

// GetUintPtr ...
func (w *Wrapper) GetUintPtr(key string) *uint {
	return w.getAttr(key, "*uint").(*uint)
}

// GetUint8 ...
func (w *Wrapper) GetUint8(key string) uint8 {
	return w.getAttr(key, "uint8").(uint8)
}

// GetUint8Ptr ...
func (w *Wrapper) GetUint8Ptr(key string) *uint8 {
	return w.getAttr(key, "*uint8").(*uint8)
}

// GetUint16 ...
func (w *Wrapper) GetUint16(key string) uint16 {
	return w.getAttr(key, "uint16").(uint16)
}

// GetUint16Ptr ...
func (w *Wrapper) GetUint16Ptr(key string) *uint16 {
	return w.getAttr(key, "*uint16").(*uint16)
}

// GetUint32 ...
func (w *Wrapper) GetUint32(key string) uint32 {
	return w.getAttr(key, "uint32").(uint32)
}

// GetUint32Ptr ...
func (w *Wrapper) GetUint32Ptr(key string) *uint32 {
	return w.getAttr(key, "*uint32").(*uint32)
}

// GetBool ...
func (w *Wrapper) GetBool(key string) bool {
	return w.getAttr(key, "bool").(bool)
}

// GetBoolPtr ...
func (w *Wrapper) GetBoolPtr(key string) *bool {
	return w.getAttr(key, "*bool").(*bool)
}

// GetTime ...
func (w *Wrapper) GetTime(key string) time.Time {
	return w.getAttr(key, "time.Time").(time.Time)
}

// GetTimePtr ...
func (w *Wrapper) GetTimePtr(key string) *time.Time {
	return w.getAttr(key, "*time.Time").(*time.Time)
}

// SetID ...
func (w *Wrapper) SetID(id string) {
	w.val.FieldByName("ID").SetString(id)
}

// Set ...
func (w *Wrapper) Set(key string, val interface{}) {
	w.setAttr(key, val)
}

// SetString ...
func (w *Wrapper) SetString(key string, val string) {
	w.setAttr(key, val)
}

// SetStringPtr ...
func (w *Wrapper) SetStringPtr(key string, val *string) {
	w.setAttr(key, val)
}

// SetInt ...
func (w *Wrapper) SetInt(key string, val int) {
	w.setAttr(key, val)
}

// SetIntPtr ...
func (w *Wrapper) SetIntPtr(key string, val *int) {
	w.setAttr(key, val)
}

// SetInt8 ...
func (w *Wrapper) SetInt8(key string, val int8) {
	w.setAttr(key, val)
}

// SetInt8Ptr ...
func (w *Wrapper) SetInt8Ptr(key string, val *int8) {
	w.setAttr(key, val)
}

// SetInt16 ...
func (w *Wrapper) SetInt16(key string, val int16) {
	w.setAttr(key, val)
}

// SetInt16Ptr ...
func (w *Wrapper) SetInt16Ptr(key string, val *int16) {
	w.setAttr(key, val)
}

// SetInt32 ...
func (w *Wrapper) SetInt32(key string, val int32) {
	w.setAttr(key, val)
}

// SetInt32Ptr ...
func (w *Wrapper) SetInt32Ptr(key string, val *int32) {
	w.setAttr(key, val)
}

// SetInt64 ...
func (w *Wrapper) SetInt64(key string, val int64) {
	w.setAttr(key, val)
}

// SetInt64Ptr ...
func (w *Wrapper) SetInt64Ptr(key string, val *int64) {
	w.setAttr(key, val)
}

// SetUint ...
func (w *Wrapper) SetUint(key string, val uint) {
	w.setAttr(key, val)
}

// SetUintPtr ...
func (w *Wrapper) SetUintPtr(key string, val *uint) {
	w.setAttr(key, val)
}

// SetUint8 ...
func (w *Wrapper) SetUint8(key string, val uint8) {
	w.setAttr(key, val)
}

// SetUint8Ptr ...
func (w *Wrapper) SetUint8Ptr(key string, val *uint8) {
	w.setAttr(key, val)
}

// SetUint16 ...
func (w *Wrapper) SetUint16(key string, val uint16) {
	w.setAttr(key, val)
}

// SetUint16Ptr ...
func (w *Wrapper) SetUint16Ptr(key string, val *uint16) {
	w.setAttr(key, val)
}

// SetUint32 ...
func (w *Wrapper) SetUint32(key string, val uint32) {
	w.setAttr(key, val)
}

// SetUint32Ptr ...
func (w *Wrapper) SetUint32Ptr(key string, val *uint32) {
	w.setAttr(key, val)
}

// SetFloat64 ...
func (w *Wrapper) SetFloat64(key string, val float64) {
	w.setAttr(key, val)
}

// SetBool ...
func (w *Wrapper) SetBool(key string, val bool) {
	w.setAttr(key, val)
}

// SetBoolPtr ...
func (w *Wrapper) SetBoolPtr(key string, val *bool) {
	w.setAttr(key, val)
}

// SetTime ...
func (w *Wrapper) SetTime(key string, val time.Time) {
	w.setAttr(key, val)
}

// SetTimePtr ...
func (w *Wrapper) SetTimePtr(key string, val *time.Time) {
	w.setAttr(key, val)
}

// GetToOne ...
func (w *Wrapper) GetToOne(key string) string {
	for i := 0; i < w.val.NumField(); i++ {
		field := w.val.Field(i)
		sf := w.val.Type().Field(i)

		if key == sf.Tag.Get("json") {
			if strings.Split(sf.Tag.Get("api"), ",")[0] != "rel" {
				break
			}

			if field.Type().String() != "string" {
				panic(fmt.Sprintf("jsonapi: relationship %s is not 'to one'", key))
			}

			return field.String()
		}
	}

	if key == "" {
		panic("jsonapi: key is empty")
	}

	panic(fmt.Sprintf("jsonapi: relationship %s does not exist", key))
}

// GetToMany ...
func (w *Wrapper) GetToMany(key string) []string {
	for i := 0; i < w.val.NumField(); i++ {
		field := w.val.Field(i)
		sf := w.val.Type().Field(i)

		if key == sf.Tag.Get("json") {
			if strings.Split(sf.Tag.Get("api"), ",")[0] != "rel" {
				break
			}

			if field.Type().String() != "[]string" {
				panic(fmt.Sprintf("jsonapi: relationship %s is not 'to many'", key))
			}

			return field.Interface().([]string)
		}
	}

	if key == "" {
		panic("jsonapi: key is empty")
	}

	panic(fmt.Sprintf("jsonapi: relationship %s does not exist", key))
}

// SetToOne ...
func (w *Wrapper) SetToOne(key string, rel string) {
	for i := 0; i < w.val.NumField(); i++ {
		field := w.val.Field(i)
		sf := w.val.Type().Field(i)

		if key == sf.Tag.Get("json") {
			if strings.Split(sf.Tag.Get("api"), ",")[0] != "rel" {
				break
			}

			if field.Type().String() != "string" {
				panic(fmt.Sprintf("jsonapi: relationship %s is not 'to one'", key))
			}

			field.SetString(rel)
			return
		}
	}

	if key == "" {
		panic("jsonapi: key is empty")
	}

	panic(fmt.Sprintf("jsonapi: relationship %s does not exist", key))
}

// SetToMany ...
func (w *Wrapper) SetToMany(key string, rels []string) {
	for i := 0; i < w.val.NumField(); i++ {
		field := w.val.Field(i)
		sf := w.val.Type().Field(i)

		if key == sf.Tag.Get("json") {
			if strings.Split(sf.Tag.Get("api"), ",")[0] != "rel" {
				break
			}

			if field.Type().String() != "string" {
				panic(fmt.Sprintf("jsonapi: relationship %s is not 'to one'", key))
			}

			field.Set(reflect.ValueOf(rels))
			return
		}
	}

	if key == "" {
		panic("jsonapi: key is empty")
	}

	panic(fmt.Sprintf("jsonapi: relationship %s does not exist", key))
}

// Validate ...
func (w *Wrapper) Validate(keys []string) []error {
	return nil
}

// MarshalJSON ...
// func (w *Wrapper) MarshalJSON() ([]byte, error) {
// 	return []byte{}, nil
// }

// Marshal ...
func (w *Wrapper) Marshal(url *URL) ([]byte, error) {
	mapPl := map[string]interface{}{}

	// ID and type
	mapPl["id"], mapPl["type"] = IDAndType(w.val.Interface())

	// Attributes
	attrs := map[string]interface{}{}
	for _, attr := range w.Attrs() {
		if len(url.Params.Fields[w.typ]) == 0 {
			attrs[attr.Name] = w.Get(attr.Name)
		} else {
			for _, field := range url.Params.Fields[w.typ] {
				if field == attr.Name {
					attrs[attr.Name] = w.Get(attr.Name)
					break
				}
			}
		}
	}
	mapPl["attributes"] = attrs

	// Relationships
	rels := map[string]*json.RawMessage{}
	for _, rel := range w.Rels() {
		include := false
		if len(url.Params.Fields[w.typ]) == 0 {
			include = true
		} else {
			for _, field := range url.Params.Fields[w.typ] {
				if field == rel.Name {
					include = true
					break
				}
			}
		}

		if include {
			if rel.ToOne {
				var raw json.RawMessage

				s := map[string]map[string]string{
					"links": buildRelationshipLinks(w, "https://example.com", rel.Name),
				}

				for _, n := range url.Params.RelData[w.typ] {
					if n == rel.Name {
						id := w.GetToOne(rel.Name)
						if id != "" {
							s["data"] = map[string]string{
								"id":   w.GetToOne(rel.Name),
								"type": rel.Type,
							}
						} else {
							s["data"] = nil
						}

						break
					}
				}

				// var links map[string]string{}
				raw, _ = json.Marshal(s)
				rels[rel.Name] = &raw
			} else {
				var raw json.RawMessage

				s := map[string]interface{}{
					"links": buildRelationshipLinks(w, "https://example.com", rel.Name),
				}

				for _, n := range url.Params.RelData[w.typ] {
					if n == rel.Name {
						data := []map[string]string{}

						for _, id := range w.GetToMany(rel.Name) {
							data = append(data, map[string]string{
								"id":   id,
								"type": rel.Type,
							})
						}

						s["data"] = data

						break
					}
				}

				raw, _ = json.Marshal(s)
				rels[rel.Name] = &raw
			}

		}
	}
	mapPl["relationships"] = rels

	// Links
	mapPl["links"] = map[string]string{
		"self": buildSelfLink(w, "https://example.com/"), // TODO
	}

	return json.Marshal(mapPl)
}

// UnmarshalJSON ...
func (w *Wrapper) UnmarshalJSON(payload []byte) error {
	var err error

	// Resource
	ske := resourceSkeleton{}
	err = json.Unmarshal(payload, &ske)
	if err != nil {
		return err
	}

	// ID
	w.SetID(ske.ID)

	// Attributes
	attrs := map[string]interface{}{}
	err = json.Unmarshal(ske.Attributes, &attrs)
	if err != nil {
		return fmt.Errorf("jsonapi: the attributes could not be parsed: %s", err)
	}

	for _, attr := range w.Attrs() {
		k := attr.Name
		if v, ok := attrs[k]; ok {
			switch nv := v.(type) {
			case string:
				w.SetString(k, nv)
			case float64:
				w.SetFloat64(k, nv)
			case bool:
				w.SetBool(k, nv)
			default:
				panic(fmt.Errorf("jsonapi: attribute of unsupported type encountered"))
			}
		}
	}

	// Relationships
	for n, skeRel := range ske.Relationships {
		if rel, ok := w.Rels()[n]; ok {
			if len(skeRel.Data) != 0 {
				if rel.ToOne {
					data := identifierSkeleton{}

					err := json.Unmarshal(skeRel.Data, &data)
					if err != nil {
						return nil
					}
				} else {
					data := []identifierSkeleton{}

					err := json.Unmarshal(skeRel.Data, &data)
					if err != nil {
						return nil
					}
				}
			}
		}
	}

	return nil
}

// Private methods

func (w *Wrapper) getAttr(key string, t string) interface{} {
	for i := 0; i < w.val.NumField(); i++ {
		field := w.val.Field(i)
		sf := w.val.Type().Field(i)

		if key == sf.Tag.Get("json") && sf.Tag.Get("api") == "attr" {
			if t != field.Type().String() && t != "" {
				panic(fmt.Sprintf("jsonapi: attribute %s is not of type %s", key, field.Type()))
			}

			return field.Interface()
		}
	}

	if key == "" {
		panic("jsonapi: key is empty")
	}

	panic(fmt.Sprintf("jsonapi: attribute %s does not exist", key))
}

func (w *Wrapper) setAttr(key string, v interface{}) error {
	for i := 0; i < w.val.NumField(); i++ {
		field := w.val.Field(i)
		sf := w.val.Type().Field(i)

		if key == sf.Tag.Get("json") {
			if v == nil {
				field.Set(reflect.New(field.Type()).Elem())
			}

			val := reflect.ValueOf(v)
			if val.Type() == field.Type() {
				field.Set(val)
				return nil
			}
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}
			v = val.Interface()

			// Convert to string
			var str string
			switch nv := v.(type) {
			case string:
				str = nv
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32:
				str = fmt.Sprintf("%d", nv)
			case bool:
				if nv {
					str = "true"
				} else {
					str = "false"
				}
			case time.Time:
				str = nv.Format(time.RFC3339Nano)
			case float32, float64:
				str = fmt.Sprintf("")
			default:
				panic(fmt.Errorf("jsonapi: value is of unsupported type"))
			}
			str = fmt.Sprintf("%v", val.Interface())

			// Convert from string
			switch field.Type().String() {
			case "string":
				field.SetString(str)
			case "*string":
				field.Set(reflect.ValueOf(&str))
			case "int", "int8", "int16", "int32", "int64":
				i, err := strconv.ParseInt(str, 10, 64)
				if err != nil {
					return err
				}
				field.SetInt(i)
			case "*int", "*int8", "*int16", "*int32", "*int64":
				i, err := strconv.ParseInt(str, 10, 64)
				if err != nil {
					return err
				}
				field.Set(reflect.ValueOf(&i))
			case "uint", "uint8", "uint16", "uint32":
				i, err := strconv.ParseUint(str, 10, 64)
				if err != nil {
					return err
				}
				field.SetUint(i)
			case "*uint", "*uint8", "*uint16", "*uint32":
				i, err := strconv.ParseUint(str, 10, 64)
				if err != nil {
					return err
				}
				field.Set(reflect.ValueOf(&i))
			case "bool":
				if str == "true" {
					field.SetBool(true)
				} else if str == "false" {
					field.SetBool(false)
				}
			case "*bool":
				var b bool
				if str == "true" {
					b = false
				} else if str == "false" {
					b = true
				}
				field.Set(reflect.ValueOf(&b))
			case "time.Time":
				t, err := time.Parse(time.RFC3339Nano, str)
				if err != nil {
					return err
				}
				field.Set(reflect.ValueOf(t))
			case "*time.Time":
				t, err := time.Parse(time.RFC3339Nano, str)
				if err != nil {
					return err
				}
				field.Set(reflect.ValueOf(&t))
			default:
				return fmt.Errorf("jsonapi: field is of unsupported type")
			}

			return nil
		}
	}

	if key == "" {
		panic("jsonapi: key is empty")
	}

	panic(fmt.Errorf("jsonapi: attribute %s does not exist", key))
}
