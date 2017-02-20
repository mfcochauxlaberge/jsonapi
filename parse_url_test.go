package jsonapi

import (
	"errors"
	"net/url"
	"testing"
	"time"

	"kkaribu/tchek"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		url                   string
		expectedURLNormalized string
		expectedError         error
	}{

		{
			// 0
			url: ``,
			expectedURLNormalized: ``,
			expectedError:         errors.New("url is invalid"),
		}, {
			// 1
			url: `/`,
			expectedURLNormalized: ``,
			expectedError:         errors.New("url is invalid"),
		}, {
			// 2
			url: `/animals`,
			expectedURLNormalized: `
				/animals
				?page[number]=1
				&page[size]=1000
			`,
			expectedError: nil,
		}, {
			// 3
			url: `
				/animals
				?fields[animals]=name,legs,is-adorable`,
			expectedURLNormalized: `
				/animals
				?fields[animals]=is-adorable,legs,name
				&page[number]=1
				&page[size]=1000
			`,
			expectedError: nil,
		}, {
			// 4
			url: `
				/animals
				?fields[animals]=name,legs,is-adorable,invalid
			`,
			expectedURLNormalized: `
				/animals
				?fields[animals]=is-adorable,legs,name
				&page[number]=1
				&page[size]=1000
			`,
			expectedError: nil,
		}, {
			// 5
			url: `
				/animals
				?include=
					shelter.location,
					shelter.location.country,
					shelter.location,
					shelter.location,
					labels,
					shelter.location,
					shelter.location,
					shelter.location
				`,
			expectedURLNormalized: `
				/animals
				?include=labels,shelter.location.country
				&page[number]=1
				&page[size]=1000`,
			expectedError: nil,
		}, {
			// 6
			url: `
				/animals
				?include=
					shelter.location,
					shelter.location.country,
					shelters,
					labels,
				`,
			expectedURLNormalized: `
				/animals
				?include=labels,shelter.location.country
				&page[number]=1
				&page[size]=1000
			`,
			expectedError: nil,
		}, {
			// 7
			url: `
				/animals
				?include=
					shelter.location,
				&fields[shelters]=location
				&fields[animals]=name,legs,shelter
				`,
			expectedURLNormalized: `
				/animals
				?fields[animals]=legs,name,shelter
				&fields[shelters]=location
				&include=shelter.location
				&page[number]=1
				&page[size]=1000
			`,
			expectedError: nil,
		}, {
			// 8
			url: `
				/animals
				?fields[shelters]=location
				&sort=name,-legs,-
				&include=
					shelter.location,
				&filter[shelter]=s1,s2,s9
				&fields[animals]=name,legs,shelter
				`,
			expectedURLNormalized: `
				/animals
				?fields[animals]=legs,name,shelter
				&fields[shelters]=location
				&filters[shelter]=s1,s2,s9
				&include=shelter.location
				&page[number]=1
				&page[size]=1000
				&sort=name,-legs
			`,
			expectedError: nil,
		}, {
			// 9
			url: `
				/animals
				?fields[shelters]=location,
				&sort=name,,-legs,invalid
				&include=
					shelter.location,
				&filter[shelter]=,s1,s2,s9
				&fields[animals]=,name,,legs,shelter
				`,
			expectedURLNormalized: `
				/animals
				?fields[animals]=legs,name,shelter
				&fields[shelters]=location
				&filters[shelter]=s1,s2,s9
				&include=shelter.location
				&page[number]=1
				&page[size]=1000
				&sort=name,-legs
			`,
			expectedError: nil,
		},
	}

	// App
	reg := NewRegistry()

	reg.RegisterType(animal{})
	reg.RegisterType(shelter{})
	reg.RegisterType(city{})
	reg.RegisterType(country{})

	for n, test := range tests {
		u, _ := url.Parse(tchek.MakeOneLineNoSpaces(test.url))
		url, err := ParseURL(reg, u)

		if test.expectedError == nil {
			tchek.AreEqual(
				t, n,
				tchek.MakeOneLineNoSpaces(test.expectedURLNormalized),
				tchek.MakeOneLineNoSpaces(url.URLNormalized),
			)
		}

		tchek.AreEqual(t, n, test.expectedError, err)
	}
}

type animal struct {
	ID string `json:"id" api:"animals"`

	// Attributes
	Name       string    `json:"name" api:"attr"`
	Legs       int8      `json:"legs" api:"attr"`
	IsAdorable bool      `json:"is-adorable" api:"attr"`
	BornAt     time.Time `json:"born-at" api:"attr"`

	// Relationships
	BirthLocation string   `json:"birth-location" api:"rel,cities"`
	Shelter       string   `json:"shelter" api:"rel,shelters,animals"`
	Toys          []string `json:"toys" api:"rel,toys,owner"`
	Labels        []string `json:"labels" api:"rel,labels"`
}

type shelter struct {
	ID string `json:"id" api:"shelters"`

	// Attributes
	BuiltAt time.Time `json:"built-at" api:"attr"`

	// Relationships
	Location string   `json:"location" api:"rel,cities,shelter"`
	Animals  []string `json:"animals" api:"rel,animals,shelter"`
}

type city struct {
	ID string `json:"id" api:"cities"`

	// Attributes
	Name string `json:"name" api:"attr"`

	// Relationships
	Country string `json:"country" api:"rel,countries"`
	Shelter string `json:"shelter" api:"rel,shelters,location"`
}

type country struct {
	ID string `json:"id" api:"countries"`

	// Attributes
	Name string `json:"name" api:"attr"`

	// Relationships
	Cities []string `json:"cities" api:"rel,cities"`
}
