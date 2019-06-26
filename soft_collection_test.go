package jsonapi_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	. "github.com/mfcochauxlaberge/jsonapi"

	"github.com/mitchellh/copystructure"
	"github.com/stretchr/testify/assert"
)

func TestSoftCollection(t *testing.T) {
	sc := &SoftCollection{}

	// Add type
	typ := Type{Name: "thistype"}
	typ.AddAttr(Attr{
		Name: "attr1",
		Type: AttrTypeInt,
		Null: false,
	})
	typ.AddAttr(Attr{
		Name: "attr2",
		Type: AttrTypeString,
		Null: true,
	})
	typ.AddRel(Rel{
		Name:         "rel1",
		Type:         "othertype",
		ToOne:        true,
		InverseName:  "rel2",
		InverseType:  "thistype",
		InverseToOne: true,
	})
	typ.AddRel(Rel{
		Name:         "rel3",
		Type:         "othertype",
		ToOne:        false,
		InverseName:  "rel4",
		InverseType:  "thistype",
		InverseToOne: true,
	})

	// Make a copy so that modifying the original typ
	// does not modify the SoftCollection's type.
	typcopy := copystructure.Must(copystructure.Copy(typ)).(Type)
	sc.SetType(&typcopy)

	assert.Equal(t, sc.GetType(), &typ)

	// Modify the SoftCollection's type and the local type
	// at the same time and check whether they still are
	// the same.
	attr3 := Attr{
		Name: "attr3",
		Type: AttrTypeBool,
		Null: false,
	}
	rel5 := Rel{
		Name:         "rel5",
		Type:         "othertype",
		ToOne:        true,
		InverseName:  "rel6",
		InverseType:  "thistype",
		InverseToOne: false,
	}
	typ.AddAttr(attr3)
	sc.AddAttr(attr3)
	typ.AddRel(rel5)
	sc.AddRel(rel5)

	assert.Equal(t, sc.GetType(), &typ)

	// Add a SoftResource with more fields than those
	// specified in the SoftCollection.
	sr := NewSoftResource(Type{Name: "thirdtype"}, nil)
	attr4 := Attr{
		Name: "attr4",
		Type: AttrTypeUint16,
		Null: true,
	}
	sr.AddAttr(attr4)
	typ.AddAttr(attr4)
	rel7 := Rel{
		Name:         "rel7",
		Type:         "othertype",
		ToOne:        true,
		InverseName:  "rel8",
		InverseType:  "thirdtype",
		InverseToOne: true,
	}
	sr.AddRel(rel7)
	typ.AddRel(rel7)

	sc.Add(sr)

	assert.Equal(t, sc.GetType(), &typ)

	// Add more elements to the SoftCollection.
	sr = &SoftResource{}
	sr.SetID("res1")
	sc.Add(sr)
	sr = &SoftResource{}
	sr.SetID("res2")
	sc.Add(sr)

	assert.Equal(t, 3, sc.Len())

	// Remove an element.
	sc.Remove("res1")
	sc.Remove("res99")

	assert.Equal(t, 2, sc.Len())
}

func TestSoftCollectionResource(t *testing.T) {
	sc := &SoftCollection{}
	sc.SetType(&Type{})

	sc.GetType().Name = "type1"
	sc.GetType().AddAttr(Attr{
		Name: "attr1",
		Type: AttrTypeString,
		Null: false,
	})
	sc.GetType().AddAttr(Attr{
		Name: "attr2",
		Type: AttrTypeInt,
		Null: true,
	})
	sc.GetType().AddRel(Rel{
		Name:  "rel1",
		Type:  "type2",
		ToOne: true,
	})

	sr := &SoftResource{}
	sr.SetType(sc.GetType())
	sr.SetID("res1")
	sr.Set("attr", "value1")
	sc.Add(sr)

	// Resource with all fields
	assert.Equal(t, sr, sc.Resource("res1", nil))

	// Resource with some fields
	// TODO Fix this test. It seems like defining any set of
	// fields will make the assert pass.
	assert.Equal(t, sr, sc.Resource("res1", []string{"attr2", "rel1"}))

	// Resource not found
	assert.Equal(t, nil, sc.Resource("notfound", nil))
}

func TestSoftCollectionSort(t *testing.T) {
	var (
		now = time.Now()
		sc  = &SoftCollection{}
	)

	attrs := []struct {
		typ  string
		vals [2]interface{}
	}{
		{
			typ:  "string",
			vals: [2]interface{}{"", "a"},
		}, {
			typ:  "int",
			vals: [2]interface{}{int(-1), int(0)},
		}, {
			typ:  "int8",
			vals: [2]interface{}{int8(-1), int8(0)},
		}, {
			typ:  "int16",
			vals: [2]interface{}{int16(-1), int16(0)},
		}, {
			typ:  "int32",
			vals: [2]interface{}{int32(-1), int32(0)},
		}, {
			typ:  "int64",
			vals: [2]interface{}{int64(-1), int64(0)},
		}, {
			typ:  "uint",
			vals: [2]interface{}{uint(0), uint(1)},
		}, {
			typ:  "uint8",
			vals: [2]interface{}{uint8(0), uint8(1)},
		}, {
			typ:  "uint16",
			vals: [2]interface{}{uint16(0), uint16(1)},
		}, {
			typ:  "uint32",
			vals: [2]interface{}{uint32(0), uint32(1)},
		}, {
			typ:  "uint64",
			vals: [2]interface{}{uint64(0), uint64(1)},
		}, {
			typ:  "bool",
			vals: [2]interface{}{false, true},
		}, {
			typ: "time.Time",
			vals: [2]interface{}{
				now,
				now.Add(time.Second),
			},
		}, {
			typ:  "*string",
			vals: [2]interface{}{nilptr("string"), nilptr("string")},
		}, {
			typ:  "*string",
			vals: [2]interface{}{nilptr("string"), ptr("a")},
		}, {
			typ:  "*int",
			vals: [2]interface{}{nilptr("int"), nilptr("int")},
		}, {
			typ:  "*int",
			vals: [2]interface{}{nilptr("int"), ptr(int(0))},
		}, {
			typ:  "*int8",
			vals: [2]interface{}{nilptr("int8"), nilptr("int8")},
		}, {
			typ:  "*int8",
			vals: [2]interface{}{nilptr("int8"), ptr(int8(0))},
		}, {
			typ:  "*int16",
			vals: [2]interface{}{nilptr("int16"), nilptr("int16")},
		}, {
			typ:  "*int16",
			vals: [2]interface{}{nilptr("int16"), ptr(int16(0))},
		}, {
			typ:  "*int32",
			vals: [2]interface{}{nilptr("int32"), nilptr("int32")},
		}, {
			typ:  "*int32",
			vals: [2]interface{}{nilptr("int32"), ptr(int32(0))},
		}, {
			typ:  "*int64",
			vals: [2]interface{}{nilptr("int64"), nilptr("int64")},
		}, {
			typ:  "*int64",
			vals: [2]interface{}{nilptr("int64"), ptr(int64(0))},
		}, {
			typ:  "*uint",
			vals: [2]interface{}{nilptr("uint"), nilptr("uint")},
		}, {
			typ:  "*uint",
			vals: [2]interface{}{nilptr("uint"), ptr(uint(0))},
		}, {
			typ:  "*uint8",
			vals: [2]interface{}{nilptr("uint8"), nilptr("uint8")},
		}, {
			typ:  "*uint8",
			vals: [2]interface{}{nilptr("uint8"), ptr(uint8(0))},
		}, {
			typ:  "*uint16",
			vals: [2]interface{}{nilptr("uint16"), nilptr("uint16")},
		}, {
			typ:  "*uint16",
			vals: [2]interface{}{nilptr("uint16"), ptr(uint16(0))},
		}, {
			typ:  "*uint32",
			vals: [2]interface{}{nilptr("uint32"), nilptr("uint32")},
		}, {
			typ:  "*uint32",
			vals: [2]interface{}{nilptr("uint32"), ptr(uint32(0))},
		}, {
			typ:  "*uint64",
			vals: [2]interface{}{nilptr("uint64"), nilptr("uint64")},
		}, {
			typ:  "*uint64",
			vals: [2]interface{}{nilptr("uint64"), ptr(uint64(0))},
		}, {
			typ:  "*bool",
			vals: [2]interface{}{nilptr("bool"), ptr(true)},
		}, {
			typ:  "*bool",
			vals: [2]interface{}{nilptr("bool"), nilptr("bool")},
		}, {
			typ: "*time.Time",
			vals: [2]interface{}{
				nilptr("time.Time"),
				now,
			},
		}, {
			typ: "*time.Time",
			vals: [2]interface{}{
				nilptr("time.Time"),
				nilptr("time.Time"),
			},
		},
	}

	// 1, 3 => 3
	// 2, 6 => 12
	// 3, 8 => 24
	// 4, 10 => 40
	// Formula: number of types * (number of types + 2)

	// 1
	// 1
	// 2

	// 21
	// 21
	// 12
	// 12
	// 11
	// 11

	// 211
	// 211
	// 121
	// 121
	// 112
	// 112
	// 111
	// 111

	// 2111
	// 2111
	// 1211
	// 1211
	// 1121
	// 1121
	// 1112
	// 1112
	// 1111
	// 1111

	// Add attributes to type
	typ := &Type{Name: "type"}
	for i, t := range attrs {
		ti, null := GetAttrType(t.typ)
		typ.AddAttr(Attr{
			Name: "attr" + strconv.Itoa(i),
			Type: ti,
			Null: null,
		})
	}
	sc.SetType(typ)

	// Add resources
	num := len(attrs)*2 + 2
	for n := 0; n < num; n++ {
		sr := &SoftResource{}
		sr.SetType(typ)
		sr.SetID("id" + strconv.Itoa(n))
		i2 := (n - (n % 2)) / 2
		fmt.Printf("n: %d, i2: %d\n", n, i2)
		for i := 0; i < len(attrs); i++ {
			// i2 := len(attrs) + 1 - n
			if i != i2 {
				// fmt.Printf("first value\n")
				if i == 2 {
					fmt.Printf("setting %s/%s to %v (first)\n", sr.GetID(), "attr"+strconv.Itoa(i), attrs[i].vals[0])
				}
				sr.Set("attr"+strconv.Itoa(i), attrs[i].vals[0])
			} else {
				// fmt.Printf("second value\n")
				if i == 2 {
					fmt.Printf("setting %s/%s to %v (second)\n", sr.GetID(), "attr"+strconv.Itoa(i), attrs[i].vals[1])
				}
				sr.Set("attr"+strconv.Itoa(i), attrs[i].vals[1])
			}
		}
		sc.Add(sr)
	}

	for j := 0; j < sc.Len(); j++ {
		res := sc.Elem(j)
		fmt.Printf("Resource: %s (%s)\n", res.GetID(), res.GetType().Name)
		typ := res.GetType()
		for _, field := range typ.Fields() {
			fmt.Printf("  %s: '%v' (%T)\n", field, res.Get(field), res.Get(field))
		}
	}

	// Sort collection
	rules := []string{}
	for i := 0; i < sc.Len(); i++ {
		reverse := ""
		if i%3 == 0 {
			reverse = "-"
		}
		rules = append(rules, reverse+"attr"+strconv.Itoa(i))
	}
	rules = append(rules, "id")
	sc.Sort(rules)

	// Sorted IDs from the collection
	ids := []string{}
	for i := 0; i < sc.Len(); i++ {
		ids = append(ids, sc.Elem(i).GetID())
	}

	// expectedIDs := []string{
	// 	"id14", "id2", "id0", "id193", "id4", "id5", "id6", "id7", "id8",
	// 	"id9", "id10", "id11", "id12", "id192", "id1", "id15", "id16", "id17",
	// 	"id18", "id19", "id20", "id21", "id22", "id23", "id24", "id25", "id26",
	// 	"id27", "id28", "id29", "id30", "id31", "id32", "id33", "id34", "id35",
	// 	"id36", "id37", "id38", "id39", "id40", "id41", "id42", "id43", "id44",
	// 	"id45", "id46", "id47", "id48", "id49", "id50", "id51", "id52", "id53",
	// 	"id54", "id55", "id56", "id57", "id58", "id59", "id60", "id61", "id62",
	// 	"id63", "id64", "id65", "id66", "id67", "id68", "id69", "id70", "id71",
	// 	"id72", "id73", "id74", "id75", "id76", "id77", "id78", "id79", "id80",
	// 	"id81", "id82", "id83", "id84", "id85", "id86", "id87", "id88", "id89",
	// 	"id90", "id91", "id92", "id93", "id94", "id95", "id96", "id97", "id98",
	// 	"id99", "id100", "id101", "id102", "id103", "id104", "id105", "id106",
	// 	"id107", "id108", "id109", "id110", "id111", "id112", "id113", "id114",
	// 	"id115", "id116", "id117", "id118", "id119", "id120", "id121", "id122",
	// 	"id123", "id124", "id125", "id126", "id127", "id128", "id129", "id130",
	// 	"id131", "id132", "id133", "id134", "id135", "id136", "id137", "id138",
	// 	"id139", "id140", "id141", "id142", "id143", "id144", "id145", "id146",
	// 	"id147", "id148", "id149", "id150", "id151", "id152", "id153", "id154",
	// 	"id155", "id156", "id157", "id158", "id159", "id160", "id161", "id162",
	// 	"id163", "id164", "id165", "id166", "id167", "id168", "id169", "id170",
	// 	"id171", "id172", "id173", "id174", "id175", "id176", "id177", "id178",
	// 	"id179", "id180", "id181", "id182", "id183", "id184", "id185", "id186",
	// 	"id187", "id188", "id189", "id190", "id191", "id194", "id3", "id13",
	// }
	// assert.Equal(t, expectedIDs, ids, fmt.Sprintf("rules: %v", rules))

	// // Sort with an empty list of sorting rules.
	// sc.Sort([]string{})

	// ids = []string{}
	// for i := 0; i < sc.Len(); i++ {
	// 	ids = append(ids, sc.Elem(i).GetID())
	// }

	// expectedIDs = []string{
	// 	"id1", "id10", "id2", "id3", "id4", "id5", "id6", "id7", "id8", "id9",
	// }
	// assert.Equal(t, expectedIDs, ids)
}
