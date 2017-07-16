package jsonapi

import (
	"encoding/json"
	"testing"

	"github.com/kkaribu/tchek"
)

func TestUnmarshalResource(t *testing.T) {
	reg := NewMockRegistry()

	res1 := Wrap(&MockType1{})
	url1, err := ParseRawURL(reg, "/mocktypes/mt1")
	tchek.UnintendedError(err)

	meta1 := map[string]interface{}{
		"str": "a string\\^ç\"",
		"num": float64(42),
		"b":   true,
	}

	jsonapi1 := map[string]interface{}{
		"version": "1.0",
	}

	doc1 := NewDocument()
	doc1.URL = url1
	doc1.Meta = meta1
	doc1.JSONAPI = jsonapi1

	pl1, err := json.Marshal(doc1)
	tchek.UnintendedError(err)

	// buf := &bytes.Buffer{}
	// _ = json.Indent(buf, pl1, "", "\t")
	// pl1 = buf.Bytes()
	// fmt.Printf("PAYLOAD:\n%s\n", pl1)

	dst1 := Wrap(&MockType1{})

	doc2, err := Unmarshal(pl1, nil)
	tchek.UnintendedError(err)

	tchek.HaveEqualAttributes(t, -1, res1, dst1)
	tchek.AreEqual(t, -1, meta1, doc2.Meta)
	tchek.AreEqual(t, -1, jsonapi1, doc2.JSONAPI)
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
