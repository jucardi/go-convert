package converter

import (
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"strings"
)

var (
	bsonType = reflect.TypeOf(bson.M{})
	instance *MapConverter
)

type MapConverter struct {
	fieldTag                  string
	useFieldNameOnTagMismatch bool
}

// Default:  returns the default instance of the MapConverter
func Default() *MapConverter {
	if instance == nil {
		instance = NewMapConverter("json", false)
	}
	return instance
}

// SetDefault: Replaces the default instance with the given instance.
func SetDefault(converter *MapConverter) *MapConverter {
	instance = converter
	return instance
}

// NewMapConverter: Creates a new instance of MapConverter that may have a different configuration than the default instance. The arg 'tag' specifies if a tag should be used to
// match the fields. Eg. `json`, `bson`, `xml`. The default instance is configured to use `json`
func NewMapConverter(tag string, useFieldNameOnTagMismatch bool) *MapConverter {
	return &MapConverter{
		fieldTag:                  tag,
		useFieldNameOnTagMismatch: useFieldNameOnTagMismatch,
	}
}

// MapToStruct: Attempts to convert a map[string]interface{} to a provided pointer to a struct. The conversion happens recursively, meaning that if a struct reference is defined
// in a parent struct ref, it will be automatically created and mapped as well.
//
//        Eg: Given the following struct definitions:
//
//                type StructA struct {                       type StructB struct {
//                    A string                                    X string
//                    B int                                       Y map[string]interface{}
//                    C bool                                  }
//                    D *StructB
//                    E map[string]interface{}
//                }
//
//        By doing the following:
//
//                out := &StructA{}
//                converter.MapToStruct(in, out)
//
//
//        The `in` map defined below will produce the `out` also defined below
//
//            in := map[string]interface{}{                      out := &StructA{
//                "A": "something",                                  A: "something",
//                "B": 1234,                                         B: 1234,
//                "C": true,                                         C: true,
//                "D": map[string]interface{}{                       D: &StructB{
//                    "X": "abcd",                                       X: "abcd",
//                    "Y": map[string]interface{}{                       Y: map[string]interface{}{
//                        "O": "qwerty",                                     "O": "qwerty",
//                    },                                                 },
//                },                                                 },
//                "E": map[string]interface{}{                       E: map[string]interface{}{
//                    "P": "qwertz",                                     "P": "qwertz",
//                },                                                 },
//            }                                                  }
//
// NOTE: If the struct has `json` tags defined, the converter will attempt to match by json tags before attempting to match by field name. See the use cases
// defined in the test file.S
//
func (m *MapConverter) MapToStruct(in map[string]interface{}, out interface{}, omitErrors ...bool) error {
	return m.BsonToStruct(bson.M(in), out, omitErrors...)
}

// BsonToStruct:  Attempts to convert a bson.M to a provided pointer to a struct. The conversion happens recursively, meaning that if a struct reference is defined
// in a parent struct ref, it will be automatically created and mapped as well. Works similar to MapToStruct. Read MapToStruct docs for more information.
func (m *MapConverter) BsonToStruct(in bson.M, out interface{}, omitErrors ...bool) error {
	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Ptr {
		return errors.New("output must be a pointer")
	}

	if v.IsNil() {
		return errors.New("out pointer must not be nil")
	}

	for k, v := range in {
		if v == nil {
			continue
		}
		if err := m.SetField(out, k, v); err != nil {
			fmt.Println(err)
			if len(omitErrors) == 0 || !omitErrors[0] {
				return fmt.Errorf("error assigning value to field '%s', %v", k, err)
			}
		}
	}

	return nil
}

// SetField: Sets the field value of a struct pointer.
func (m *MapConverter) SetField(obj interface{}, name string, value interface{}) error {
	v := reflect.ValueOf(obj)

	if v.Kind() != reflect.Ptr {
		return errors.New("must receive a pointer to a struct")
	}

	structValue := v.Elem()

	if structValue.Kind() != reflect.Struct {
		return errors.New("pointer does not point to a struct")
	}

	var structFieldValue reflect.Value

	if m.fieldTag != "" {
		if jsonName := m.findFieldByJsonTag(structValue.Type(), name); jsonName != "" {
			structFieldValue = structValue.FieldByName(jsonName)
		}

		if !structFieldValue.IsValid() && !m.useFieldNameOnTagMismatch {
			return fmt.Errorf("no field matches the '%s' tag with value '%s' in %v", m.fieldTag, name, structValue.Type())
		}
	}

	if !structFieldValue.IsValid() {
		structFieldValue = structValue.FieldByName(name)
	}

	if !structFieldValue.IsValid() {
		return fmt.Errorf("no field by name '%s' in %v", name, structValue.Type())
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)

	if structFieldType != val.Type() {
		if val.Kind() != reflect.Map && val.Type() != bsonType {
			return fmt.Errorf("provided value type did not match obj field type and it is not a map, unable to convert. Name %v, expected %v, actual %v", name, structFieldType, val.Type())
		}

		isNil, newVal := ensureOrCreatePtrToStruct(structFieldValue)

		if !isNil {
			newVal = structFieldValue
		}

		if val.Type() == reflect.TypeOf(bson.M{}) {
			if err := m.BsonToStruct(val.Interface().(bson.M), newVal.Interface()); err != nil {
				return fmt.Errorf("unable to convert bson to struct pointer for internal field %s, %v", name, err)
			}
		} else if newVal.Kind() == reflect.Map && newVal.Type().Elem().String() != "interface {}" {
			newVal = reflect.MakeMap(newVal.Type())
			mapToMap(val.Interface().(map[string]interface{}), newVal)
		} else if err := m.MapToStruct(val.Interface().(map[string]interface{}), newVal.Interface()); err != nil {
			return fmt.Errorf("unable to convert map to struct pointer for internal field %s, %v", name, err)
		}

		val = newVal
	}

	structFieldValue.Set(val)
	return nil
}

func ensureOrCreatePtrToStruct(v reflect.Value) (isNil bool, ret reflect.Value) {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			v = reflect.New(v.Type().Elem())
			isNil = true
			continue
		}
		v = v.Elem()
	}
	return isNil, reflect.New(v.Type())
}

func (m *MapConverter) findFieldByJsonTag(structType reflect.Type, tagValue string) string {
	if m.fieldTag == "" {
		return ""
	}
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		t := field.Tag.Get(m.fieldTag)

		// To make sure a match happens when using metadata in tags such as `json:"name,omitempty"`
		split := strings.Split(t, ",")

		if tagValue == split[0] {
			return field.Name
		}
	}
	return ""
}

// BsonToMap: Easily converts from the mgo bson.M type to a map[string]interface{} recursively
func BsonToMap(in bson.M) map[string]interface{} {
	ret := map[string]interface{}{}

	for k, v := range in {
		if reflect.TypeOf(v) == bsonType {
			ret[k] = BsonToMap(v.(bson.M))
		} else {
			ret[k] = v
		}
	}

	return ret
}

// MapToMap: Is a map converter helper to easily convert a map[string]interface{} to a map[string]*someType, assuming the values contained by map[string]interface{} are in fact
// the same type for the output.
func MapToMap(in map[string]interface{}, out interface{}, omitErrors ...bool) error {
	outVal := reflect.ValueOf(out)
	return mapToMap(in, outVal, omitErrors...)
}

func mapToMap(in map[string]interface{}, outVal reflect.Value, omitErrors ...bool) error {
	if outVal.Kind() != reflect.Map {
		return errors.New("'out' must be a map")
	}

	if outVal.IsNil() {
		return errors.New("the output map must be initialized")
	}

	mapType := outVal.Type().Elem()

	for k, v := range in {
		kVal := reflect.ValueOf(k)
		vVal := reflect.ValueOf(v)

		if vVal.Type() == mapType {
			outVal.SetMapIndex(kVal, vVal)
			continue
		}

		err := fmt.Errorf("error assigning value to field '%s', type mismatch %v != %v, ", k, vVal.Type(), mapType)
		fmt.Println(err)

		if len(omitErrors) == 0 || !omitErrors[0] {
			return err
		}
	}

	return nil
}
