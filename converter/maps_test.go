package converter

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"gopkg.in/mgo.v2/bson"
)

type structWithMapA struct {
	A string
	B int
	C bool
	D *structWithMapB
	E map[string]interface{}
}

type structWithMapB struct {
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

	expected := &structWithMapA{
		A: "something",
		B: 1234,
		C: true,
		D: &structWithMapB{
			X: "abcd",
			Y: map[string]interface{}{
				"O": "qwerty",
			},
		},
		E: map[string]interface{}{
			"P": "qwertz",
		},
	}

	out := &structWithMapA{}
	err := MapToStruct(in, out)
	assert.Nil(t, err)
	assert.Equal(t, expected, out)
	assert.Equal(t, expected.D, out.D)
	fmt.Printf("%v", out)
}

type structWithBsonA struct {
	A string
	B int
	C bool
	D *structWithBsonB
	E bson.M
}

type structWithBsonB struct {
	X string
	Y bson.M
}

func TestBsonToStruct_MatchingTypes(t *testing.T) {
	in := bson.M{
		"A": "something",
		"B": 1234,
		"C": true,
		"D": bson.M{
			"X": "abcd",
			"Y": bson.M{
				"O": "qwerty",
			},
		},
		"E": bson.M{
			"P": "qwertz",
		},
	}

	expected := &structWithBsonA{
		A: "something",
		B: 1234,
		C: true,
		D: &structWithBsonB{
			X: "abcd",
			Y: bson.M{
				"O": "qwerty",
			},
		},
		E: bson.M{
			"P": "qwertz",
		},
	}

	out := &structWithBsonA{}
	err := BsonToStruct(in, out)
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

func TestBsonToMap(t *testing.T) {
	in := bson.M{
		"A": "something",
		"B": 1234,
		"C": true,
		"D": bson.M{
			"X": "abcd",
			"Y": bson.M{
				"O": "qwerty",
			},
		},
		"E": bson.M{
			"P": "qwertz",
		},
	}

	expected := map[string]interface{}{
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

	out := BsonToMap(in)
	assert.Equal(t, expected, out)
}