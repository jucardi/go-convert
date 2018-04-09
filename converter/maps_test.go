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

	//out := &testStructA{}
	out := &testStructA{}
	err := MapToStruct(in, out)
	assert.Nil(t, err)
	assert.Equal(t, expected, out)
	assert.Equal(t, expected.D, out.D)
	fmt.Printf("%v", out)
}
