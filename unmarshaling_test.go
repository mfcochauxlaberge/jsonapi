package jsonapi

import (
	"io/ioutil"
	"testing"
	"time"

	"kkaribu/tchek"
)

func TestUnmarshalResource(t *testing.T) {
	loc, _ := time.LoadLocation("")

	tests := []struct {
		payload       string
		dst           interface{}
		errorExpected bool
		expectedData  interface{}
		// expectedOpts  Options
	}{
		{
			// 0
			payload:       "payload-0",
			dst:           &user{},
			errorExpected: false,
			expectedData: user{
				ID:        "1",
				Name:      "Bob",
				Age:       36,
				CreatedAt: time.Date(2017, 1, 2, 3, 4, 5, 6, loc),
			},
		},
	}

	for n, test := range tests {
		content, err := ioutil.ReadFile("tests/" + test.payload + ".json")
		tchek.UnintendedError(err)

		res := Wrap(test.dst)

		err = Unmarshal(content, res)
		tchek.ErrorExpected(t, n, test.errorExpected, err)

		if !test.errorExpected {
			tchek.HaveEqualAttributes(t, n, test.dst, test.expectedData)
		}
	}
}

func TestUnmarshalCollection(t *testing.T) {
	// loc, _ := time.LoadLocation("")

	tests := []struct {
		payload       string
		sample        interface{}
		errorExpected bool
		expectedData  interface{}
	}{
		{
			// 0
			payload:       "payload-1",
			sample:        user{},
			errorExpected: false,
			expectedData:  users,
		},
	}

	for n, test := range tests {
		content, err := ioutil.ReadFile("tests/" + test.payload + ".json")
		tchek.UnintendedError(err)

		r := Wrap(user{})
		col := WrapCollection(r)

		err = Unmarshal(content, col)
		tchek.ErrorExpected(t, n, test.errorExpected, err)

		if !test.errorExpected {
			tchek.HaveEqualAttributes(t, n, col, test.expectedData)
		}
	}
}
