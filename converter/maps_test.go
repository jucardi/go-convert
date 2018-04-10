package converter

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefault(t *testing.T) {
	c := Default()
	assert.Equal(t, "json", c.fieldTag)
	assert.False(t, c.useFieldNameOnTagMismatch)
}

func TestMapToStruct_MatchingTypes(t *testing.T) {
	in := sampleMap()
	expected := sampleStruct()
	out := &testStructA{}

	err := NewMapConverter("", false).MapToStruct(in, out)
	assert.Nil(t, err)
	assert.Equal(t, expected, out)
	assert.Equal(t, expected.D, out.D)
	fmt.Printf("%v", out)
}

func TestBsonToStruct_MatchingTypes(t *testing.T) {
	in := sampleBson()
	expected := sampleStructWithBson()
	out := &structWithBsonA{}

	err := NewMapConverter("", false).MapToStruct(in, out)
	assert.Nil(t, err)
	assert.Equal(t, expected, out)
	assert.Equal(t, expected.D, out.D)
	fmt.Printf("%v", out)
}

func TestMapToStruct_MatchingTypesJsonTags(t *testing.T) {
	in := sampleMapMatchingJsonTags()
	expected := sampleStruct()
	out := &testStructA{}

	err := Default().MapToStruct(in, out)
	assert.Nil(t, err)
	assert.Equal(t, expected, out)
	assert.Equal(t, expected.D, out.D)
	fmt.Printf("%v", out)
}

func TestBsonToMap(t *testing.T) {
	in := sampleBson()
	expected := sampleMap()

	out := BsonToMap(in)
	assert.Equal(t, expected, out)
}

func TestMapToMap(t *testing.T) {
	in := map[string]interface{}{
		"a": []string{"a", "b", "c", "b"},
		"b": []string{"1", "2", "3", "4"},
		"c": []string{"x", "y", "z", "w"},
	}

	expected := map[string][]string{
		"a": {"a", "b", "c", "b"},
		"b": {"1", "2", "3", "4"},
		"c": {"x", "y", "z", "w"},
	}

	out := map[string][]string{}
	err := MapToMap(in, out)
	assert.Nil(t, err)
	assert.Equal(t, expected, out)
}

func TestMapConverter_MapToStruct_NestedMapSpecifiedValueType(t *testing.T) {
	in := sampleMapNestedDifferentMap()
	expected := sampleStructNestedMap()
	out := &structNestedMap{}

	err := NewMapConverter("", false).MapToStruct(in, out)
	assert.Nil(t, err)
	assert.Equal(t, expected, out)
}
