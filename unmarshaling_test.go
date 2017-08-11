package jsonapi

import (
	"testing"

	"github.com/kkaribu/tchek"
)

func TestUnmarshalResource(t *testing.T) {
	reg := NewMockRegistry()

	res1 := Wrap(&MockType1{})
	url1, err := ParseRawURL(reg, "/mocktypes1/mt1")
	tchek.UnintendedError(err)

	meta1 := map[string]interface{}{
		"str": "a string\\^ç\"",
		"num": float64(42),
		"b":   true,
	}

	doc1 := NewDocument()
	doc1.Data = res1
	doc1.Meta = meta1

	body1, err := Marshal(doc1, url1)
	tchek.UnintendedError(err)

	// buf := &bytes.Buffer{}
	// _ = json.Indent(buf, pl1, "", "\t")
	// pl1 = buf.Bytes()
	// fmt.Printf("PAYLOAD:\n%s\n", pl1)

	dst1 := Wrap(&MockType1{})

	pl1, err := Unmarshal(body1, url1, reg)
	tchek.UnintendedError(err)

	tchek.HaveEqualAttributes(t, -1, res1, dst1)
	tchek.AreEqual(t, -1, meta1, pl1.Meta)
}

// func TestUnmarshalCollection(t *testing.T) {
// 	// loc, _ := time.LoadLocation("")
//
// 	tests := []struct {
// 		payload       string
// 		sample        interface{}
// 		errorExpected bool
// 		expectedData  interface{}
// 	}{
// 		{
// 			// 0
// 			payload:       "payload-1",
// 			sample:        user{},
// 			errorExpected: false,
// 			expectedData:  users,
// 		},
// 	}
//
// 	for n, test := range tests {
// 		content, err := ioutil.ReadFile("tests/" + test.payload + ".json")
// 		tchek.UnintendedError(err)
//
// 		r := Wrap(&user{})
// 		col := WrapCollection(r)
//
// 		err = Unmarshal(content, col)
// 		tchek.ErrorExpected(t, n, test.errorExpected, err)
//
// 		if !test.errorExpected {
// 			tchek.HaveEqualAttributes(t, n, col, test.expectedData)
// 		}
// 	}
// }
