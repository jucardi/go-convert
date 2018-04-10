package converter

import "gopkg.in/mgo.v2/bson"

type testStructA struct {
	A string                 `json:"a,omitempty"`
	B int                    `json:"b"`
	C bool                   `json:"c"`
	D *testStructB           `json:"d"`
	E map[string]interface{} `json:"e"`
}

type testStructB struct {
	X string                 `json:"x"`
	Y map[string]interface{} `json:"y"`
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

type structNestedMap struct {
	A string
	B int
	C map[string][]string
}

func sampleStruct() *testStructA {
	return &testStructA{
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
}

func sampleStructWithBson() *structWithBsonA {
	return &structWithBsonA{
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
}

func sampleBson() bson.M {
	return bson.M{
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
}

func sampleMap() map[string]interface{} {
	return map[string]interface{}{
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
}

func sampleMapMatchingJsonTags() map[string]interface{} {
	return map[string]interface{}{
		"a": "something",
		"b": 1234,
		"c": true,
		"d": map[string]interface{}{
			"x": "abcd",
			"y": map[string]interface{}{
				"O": "qwerty",
			},
		},
		"e": map[string]interface{}{
			"P": "qwertz",
		},
	}
}

func sampleStructNestedMap() *structNestedMap {
	return &structNestedMap{
		A: "abcd",
		B: 1234,
		C: map[string][]string{
			"a": {"a", "b", "c", "b"},
			"b": {"1", "2", "3", "4"},
			"c": {"x", "y", "z", "w"},
		},
	}
}

func sampleMapNestedDifferentMap() map[string]interface{} {
	return map[string]interface{}{
		"A": "abcd",
		"B": 1234,
		"C": map[string]interface{}{
			"a": []string{"a", "b", "c", "b"},
			"b": []string{"1", "2", "3", "4"},
			"c": []string{"x", "y", "z", "w"},
		},
	}
}
