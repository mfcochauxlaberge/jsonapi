package jsonapi_test

import (
	"testing"

	. "github.com/mfcochauxlaberge/jsonapi"
	"github.com/stretchr/testify/assert"
)

func TestRelTempName(t *testing.T) {
	tests := []struct {
		name     string
		relTemp  RelTemp
		expected string
	}{
		{
			name: "regular relationship",
			relTemp: RelTemp{
				Type1: "t1",
				Name1: "n1",
				Type2: "t2",
				Name2: "n2",
			},
			expected: "t1_n1_t2_n2",
		}, {
			name: "regular relationship reverse",
			relTemp: RelTemp{
				Type1: "t2",
				Name1: "n2",
				Type2: "t1",
				Name2: "n1",
			},
			expected: "t1_n1_t2_n2",
		}, {
			name: "same type",
			relTemp: RelTemp{
				Type1: "t1",
				Name1: "n1",
				Type2: "t1",
				Name2: "n2",
			},
			expected: "t1_n1_t1_n2",
		}, {
			name: "same type reverse",
			relTemp: RelTemp{
				Type1: "t1",
				Name1: "n2",
				Type2: "t1",
				Name2: "n1",
			},
			expected: "t1_n1_t1_n2",
		}, {
			name: "same relationship name",
			relTemp: RelTemp{
				Type1: "t1",
				Name1: "n1",
				Type2: "t2",
				Name2: "n1",
			},
			expected: "t1_n1_t2_n1",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			// if test.expected != test.relTemp.Name() {
			// 	t.Errorf(
			// 		"expected %q, got %q",
			// 		test.expected,
			// 		test.relTemp.Name(),
			// 	)
			// }
			assert.Equal(test.expected, test.relTemp.Name())
		})
	}
}
