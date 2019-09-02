package jsonapi_test

import (
	"fmt"
	"sort"
	"strconv"
	"testing"
	"time"

	. "github.com/mfcochauxlaberge/jsonapi"

	"github.com/stretchr/testify/assert"
)

func TestRange(t *testing.T) {
	assert := assert.New(t)

	// Collection
	col := Resources{}
	typ := &Type{}
	_ = typ.AddAttr(Attr{
		Name:     "attr1",
		Type:     AttrTypeString,
		Nullable: false,
	})
	_ = typ.AddAttr(Attr{
		Name:     "attr2",
		Type:     AttrTypeInt,
		Nullable: false,
	})

	resources := []struct {
		id     string
		fields map[string]interface{}
	}{
		{
			id: "res1",
			fields: map[string]interface{}{
				"attr1": "string1",
				"attr2": 2,
			},
		}, {
			id: "res2",
			fields: map[string]interface{}{
				"attr1": "string2",
				"attr2": 2,
			},
		}, {
			id: "res3",
			fields: map[string]interface{}{
				"attr1": "string2",
				"attr2": 0,
			},
		}, {
			id: "res4",
			fields: map[string]interface{}{
				"attr1": "string2",
				"attr2": 2,
			},
		}, {
			id: "res5",
			fields: map[string]interface{}{
				"attr1": "string3",
				"attr2": 2,
			},
		}, {
			id: "res6",
			fields: map[string]interface{}{
				"attr1": "string3",
				"attr2": 4,
			},
		}, {
			id: "res7",
			fields: map[string]interface{}{
				"attr1": "string4",
				"attr2": 0,
			},
		}, {
			id: "res8",
			fields: map[string]interface{}{
				"attr1": "string5",
				"attr2": 2,
			},
		},
	}

	for _, res := range resources {
		sr := &SoftResource{}
		sr.SetType(typ)
		sr.SetID(res.id)
		for field, val := range res.fields {
			sr.Set(field, val)
		}
		col.Add(sr)
	}

	// Range test 1
	ranged := Range(
		// Collection
		&col,
		// IDs
		[]string{},
		// Filter
		nil,
		// Sort
		[]string{},
		// PageSize
		10,
		// PageNumber
		0,
	)

	expectedIDs := []string{
		"res1", "res2", "res3", "res4", "res5", "res6", "res7", "res8",
	}
	ids := []string{}
	for i := 0; i < ranged.Len(); i++ {
		ids = append(ids, ranged.At(i).GetID())
	}
	assert.Equal(expectedIDs, ids, "range of IDs (1)")

	// Range test 2
	ranged = Range(
		// Collection
		&col,
		// IDs
		[]string{"res1", "res2", "res3", "res4", "res5", "res6"},
		// Filter
		&Filter{Field: "attr2", Op: "=", Val: 2},
		// Sort
		[]string{"-attr1"},
		// PageSize
		2,
		// PageNumber
		0,
	)

	expectedIDs = []string{"res5", "res2"}
	ids = []string{}
	for i := 0; i < ranged.Len(); i++ {
		ids = append(ids, ranged.At(i).GetID())
	}
	assert.Equal(expectedIDs, ids, "range of IDs (2)")

	// Range test 3
	assert.Equal(
		0,
		Range(&Resources{}, nil, nil, nil, 1, 100).Len(),
		"range of IDs (3)",
	)
}

func TestSortResources(t *testing.T) {
	assert := assert.New(t)

	var (
		now            = time.Now()
		col Collection = &Resources{}
	)

	// A collection of resources will be created and
	// one attribute will be added for each entry from
	// the following slice.
	// The point is to provoke all possible scenarios
	// for each attribute type.
	attrs := []struct {
		vals [2]interface{}
	}{
		// non-nullable
		{vals: [2]interface{}{"", "a"}},
		{vals: [2]interface{}{int(-1), int(0)}},
		{vals: [2]interface{}{int8(-1), int8(0)}},
		{vals: [2]interface{}{int16(-1), int16(0)}},
		{vals: [2]interface{}{int32(-1), int32(0)}},
		{vals: [2]interface{}{int64(-1), int64(0)}},
		{vals: [2]interface{}{uint(0), uint(1)}},
		{vals: [2]interface{}{uint8(0), uint8(1)}},
		{vals: [2]interface{}{uint16(0), uint16(1)}},
		{vals: [2]interface{}{uint32(0), uint32(1)}},
		{vals: [2]interface{}{uint64(0), uint64(1)}},
		{vals: [2]interface{}{false, true}},
		{vals: [2]interface{}{now, now.Add(time.Second)}},
		// nullable
		{vals: [2]interface{}{nilptr("string"), nilptr("string")}},
		{vals: [2]interface{}{nilptr("string"), ptr("a")}},
		{vals: [2]interface{}{ptr(""), nilptr("string")}},
		{vals: [2]interface{}{ptr(""), ptr("")}},
		{vals: [2]interface{}{ptr(""), ptr("a")}},
		{vals: [2]interface{}{nilptr("int"), nilptr("int")}},
		{vals: [2]interface{}{nilptr("int"), ptr(int(0))}},
		{vals: [2]interface{}{ptr(int(-1)), nilptr("int")}},
		{vals: [2]interface{}{ptr(int(-1)), ptr(int(-1))}},
		{vals: [2]interface{}{ptr(int(-1)), ptr(int(0))}},
		{vals: [2]interface{}{nilptr("int8"), nilptr("int8")}},
		{vals: [2]interface{}{nilptr("int8"), ptr(int8(0))}},
		{vals: [2]interface{}{ptr(int8(-1)), nilptr("int8")}},
		{vals: [2]interface{}{ptr(int8(-1)), ptr(int8(-1))}},
		{vals: [2]interface{}{ptr(int8(-1)), ptr(int8(0))}},
		{vals: [2]interface{}{nilptr("int16"), nilptr("int16")}},
		{vals: [2]interface{}{nilptr("int16"), ptr(int16(0))}},
		{vals: [2]interface{}{ptr(int16(-1)), nilptr("int16")}},
		{vals: [2]interface{}{ptr(int16(-1)), ptr(int16(-1))}},
		{vals: [2]interface{}{ptr(int16(-1)), ptr(int16(0))}},
		{vals: [2]interface{}{nilptr("int32"), nilptr("int32")}},
		{vals: [2]interface{}{nilptr("int32"), ptr(int32(0))}},
		{vals: [2]interface{}{ptr(int32(-1)), nilptr("int32")}},
		{vals: [2]interface{}{ptr(int32(-1)), ptr(int32(-1))}},
		{vals: [2]interface{}{ptr(int32(-1)), ptr(int32(0))}},
		{vals: [2]interface{}{nilptr("int64"), nilptr("int64")}},
		{vals: [2]interface{}{nilptr("int64"), ptr(int64(0))}},
		{vals: [2]interface{}{ptr(int64(-1)), nilptr("int64")}},
		{vals: [2]interface{}{ptr(int64(-1)), ptr(int64(-1))}},
		{vals: [2]interface{}{ptr(int64(-1)), ptr(int64(0))}},
		{vals: [2]interface{}{nilptr("uint"), nilptr("uint")}},
		{vals: [2]interface{}{nilptr("uint"), ptr(uint(0))}},
		{vals: [2]interface{}{ptr(uint(0)), nilptr("uint")}},
		{vals: [2]interface{}{ptr(uint(0)), ptr(uint(0))}},
		{vals: [2]interface{}{ptr(uint(0)), ptr(uint(1))}},
		{vals: [2]interface{}{nilptr("uint8"), nilptr("uint8")}},
		{vals: [2]interface{}{nilptr("uint8"), ptr(uint8(0))}},
		{vals: [2]interface{}{ptr(uint8(0)), nilptr("uint8")}},
		{vals: [2]interface{}{ptr(uint8(0)), ptr(uint8(0))}},
		{vals: [2]interface{}{ptr(uint8(0)), ptr(uint8(1))}},
		{vals: [2]interface{}{nilptr("uint16"), nilptr("uint16")}},
		{vals: [2]interface{}{nilptr("uint16"), ptr(uint16(0))}},
		{vals: [2]interface{}{ptr(uint16(0)), nilptr("uint16")}},
		{vals: [2]interface{}{ptr(uint16(0)), ptr(uint16(0))}},
		{vals: [2]interface{}{ptr(uint16(0)), ptr(uint16(1))}},
		{vals: [2]interface{}{nilptr("uint32"), nilptr("uint32")}},
		{vals: [2]interface{}{nilptr("uint32"), ptr(uint32(0))}},
		{vals: [2]interface{}{ptr(uint32(0)), nilptr("uint32")}},
		{vals: [2]interface{}{ptr(uint32(0)), ptr(uint32(0))}},
		{vals: [2]interface{}{ptr(uint32(0)), ptr(uint32(1))}},
		{vals: [2]interface{}{nilptr("uint64"), nilptr("uint64")}},
		{vals: [2]interface{}{nilptr("uint64"), ptr(uint64(0))}},
		{vals: [2]interface{}{ptr(uint64(0)), nilptr("uint64")}},
		{vals: [2]interface{}{ptr(uint64(0)), ptr(uint64(0))}},
		{vals: [2]interface{}{ptr(uint64(0)), ptr(uint64(1))}},
		{vals: [2]interface{}{nilptr("bool"), nilptr("bool")}},
		{vals: [2]interface{}{nilptr("bool"), ptr(false)}},
		{vals: [2]interface{}{ptr(false), nilptr("bool")}},
		{vals: [2]interface{}{ptr(false), ptr(false)}},
		{vals: [2]interface{}{ptr(false), ptr(true)}},
		{vals: [2]interface{}{nilptr("time.Time"), nilptr("time.Time")}},
		{vals: [2]interface{}{nilptr("time.Time"), ptr(now)}},
		{vals: [2]interface{}{ptr(now), ptr(now)}},
		{vals: [2]interface{}{ptr(now), ptr(now.Add(time.Second))}},
	}

	// Add attributes to type
	typ := &Type{Name: "type"}
	for i, t := range attrs {
		ti, null := GetAttrType(fmt.Sprintf("%T", t.vals[0]))
		_ = typ.AddAttr(Attr{
			Name:     "attr" + strconv.Itoa(i),
			Type:     ti,
			Nullable: null,
		})
	}

	// Add resources and attributes
	for i := range attrs {
		sr := &SoftResource{
			Type: typ,
		}
		sr.SetID("id" + strconv.Itoa(i))
		for j := range attrs {
			if i != j {
				sr.Set("attr"+strconv.Itoa(j), attrs[j].vals[0])
			} else {
				sr.Set("attr"+strconv.Itoa(j), attrs[j].vals[1])
			}
		}
		col.Add(sr)
	}

	// Sort collection
	rules := []string{}
	for i := 0; i < col.Len(); i++ {
		reverse := ""
		if i%3 == 0 {
			reverse = "-"
		}
		rules = append(rules, reverse+"attr"+strconv.Itoa(i))
	}
	rules = append(rules, "id")
	page := Range(
		col,
		nil,
		nil,
		rules,
		1000,
		0,
	)

	// Sorted IDs from the collection
	ids := []string{}
	for i := 0; i < page.Len(); i++ {
		ids = append(ids, page.At(i).GetID())
	}

	expectedIDs := []string{
		"id0", "id3", "id6", "id9", "id12", "id20", "id24", "id25", "id27",
		"id35", "id39", "id40", "id42", "id50", "id54", "id55", "id57",
		"id69", "id70", "id72", "id10", "id13", "id16", "id18", "id21",
		"id23", "id26", "id28", "id31", "id33", "id36", "id38", "id41",
		"id43", "id46", "id48", "id51", "id53", "id56", "id58", "id61",
		"id63", "id64", "id65", "id66", "id67", "id68", "id71", "id73",
		"id75", "id76", "id74", "id62", "id60", "id59", "id52", "id49",
		"id47", "id45", "id44", "id37", "id34", "id32", "id30", "id29",
		"id22", "id19", "id17", "id15", "id14", "id11", "id8", "id7", "id5",
		"id4", "id2", "id1",
	}
	assert.Equal(expectedIDs, ids, fmt.Sprintf("sort with rules: %v", rules))

	// Sort with an empty list of sorting rules.
	// col.Sort([]string{})
	page = Range(
		col,
		nil,
		nil,
		[]string{},
		1000,
		0,
	)

	ids = []string{}
	for i := 0; i < page.Len(); i++ {
		ids = append(ids, page.At(i).GetID())
	}

	sort.Strings(expectedIDs)
	assert.Equal(expectedIDs, ids, "sort by ID")

	// Sort collection with different types
	sr1 := &SoftResource{}
	sr1.SetID("sr1")
	col1 := &Resources{Wrap(mocktype{}), sr1}
	assert.Panics(func() {
		_ = Range(col1, nil, nil, []string{"field", "id"}, 100, 0)
	})

	// Sort collection with unknown attribute
	col1 = &Resources{
		Wrap(mocktype{}),
		Wrap(mocktype{}),
	}
	assert.Panics(func() {
		_ = Range(col1, nil, nil, []string{"unknown", "id"}, 100, 0)
	})

	// Sort collection with attribute of different type
	col1 = &Resources{
		&SoftResource{
			Type: &Type{
				Name: "type",
				Attrs: map[string]Attr{
					"samename": Attr{
						Name:     "samename",
						Type:     AttrTypeString,
						Nullable: false,
					},
				},
			},
		},
		&SoftResource{
			Type: &Type{
				Name: "type",
				Attrs: map[string]Attr{
					"samename": Attr{
						Name:     "samename",
						Type:     AttrTypeString,
						Nullable: true,
					},
				},
			},
		},
	}
	assert.Panics(func() {
		_ = Range(col1, nil, nil, []string{"samename", "id"}, 100, 0)
	})
}
