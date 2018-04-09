package converter

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testStructA struct {
	A string
	B int
	C bool
	D *testStructB
	E map[string]interface{}
}

type testStructB struct {
	X string
	Y map[string]interface{}
}

func TestMapToStruct_MatchingTypes(t *testing.T) {
	in := map[string]interface{}{
		"A": "something",
		"B": 1234,
		"C": true,
		"D": map[string]interface{}{
			"X": "abcd",
			"Y": map[string]interface{}{
				"O": "qwerty",
			},
		},
		"E": map[string]interface{}{
			"P": "qwertz",
		},
	}

	expected := &testStructA{
		A: "something",
		B: 1234,
		C: true,
		D: &testStructB{
			X: "abcd",
			Y: map[string]interface{}{
				"O": "qwerty",
			},
		},
		E: map[string]interface{}{
			"P": "qwertz",
		},
	}

	out := &testStructA{}
	err := MapToStruct(in, out)
	assert.Nil(t, err)
	assert.Equal(t, expected, out)
	assert.Equal(t, expected.D, out.D)
	fmt.Printf("%v", out)
}

type testStructC struct {
	A string                 `json:"a,omitempty"`
	B int                    `json:"b"`
	C bool                   `json:"c"`
	D *testStructD           `json:"d"`
	E map[string]interface{} `json:"e"`
}

type testStructD struct {
	X string                 `json:"x"`
	Y map[string]interface{} `json:"y"`
}

func TestMapToStruct_MatchingTypesJsonTags(t *testing.T) {
	in := map[string]interface{}{
		"a": "something",
		"b": 1234,
		"c": true,
		"d": map[string]interface{}{
			"x": "abcd",
			"y": map[string]interface{}{
				"o": "qwerty",
			},
		},
		"e": map[string]interface{}{
			"p": "qwertz",
		},
	}

	expected := &testStructC{
		A: "something",
		B: 1234,
		C: true,
		D: &testStructD{
			X: "abcd",
			Y: map[string]interface{}{
				"o": "qwerty",
			},
		},
		E: map[string]interface{}{
			"p": "qwertz",
		},
	}

	out := &testStructC{}
	err := MapToStruct(in, out)
	assert.Nil(t, err)
	assert.Equal(t, expected, out)
	assert.Equal(t, expected.D, out.D)
	fmt.Printf("%v", out)
}
