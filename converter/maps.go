package converter

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	ErrorNoPtr = errors.New("output must be a pointer")
)

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
func MapToStruct(in map[string]interface{}, out interface{}, omitErrors ...bool) error {
	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Ptr {
		return ErrorNoPtr
	}

	if v.IsNil() {
		return errors.New("out pointer must not be nil")
	}

	for k, v := range in {
		if err := SetField(out, k, v); err != nil {
			fmt.Println(err)
			if len(omitErrors) == 0 || !omitErrors[0] {
				return fmt.Errorf("error assigning value to field '%s', %v", k, err)
			}
		}
	}

	return nil
}

// SetField: Sets the field value of a struct pointer.
func SetField(obj interface{}, name string, value interface{}) error {
	v := reflect.ValueOf(obj)

	if v.Kind() != reflect.Ptr {
		return errors.New("must receive a pointer to a struct")
	}

	structValue := v.Elem()

	if structValue.Kind() != reflect.Struct {
		return errors.New("pointer does not point to a struct")
	}

	var structFieldValue reflect.Value

	if jsonName := findFieldByJsonTag(structValue.Type(), name); jsonName != "" {
		structFieldValue = structValue.FieldByName(jsonName)
	}

	if !structFieldValue.IsValid() {
		structFieldValue = structValue.FieldByName(name)
	}

	if !structFieldValue.IsValid() {
		return fmt.Errorf("no such field by name or json tag '%s' in %v", name, structValue.Type())
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		if val.Kind() != reflect.Map {
			return fmt.Errorf("provided value type did not match obj field type and it is not a map, unable to convert. Name %v, expected %v, actual %v", name, structFieldType, val.Type())
		}

		isNil, newVal := ensureOrCreatePtrToStruct(structFieldValue)

		if !isNil {
			newVal = structFieldValue
		}

		if err := MapToStruct(val.Interface().(map[string]interface{}), newVal.Interface()); err != nil {
			return fmt.Errorf("unable to convert mapt to struct pointer for internal field %s, %v", name, err)
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

func findFieldByJsonTag(structType reflect.Type, tag string) string {
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		t := field.Tag.Get("json")

		// To make sure a match happens when using metadata in tags such as `json:"name,omitempty"`
		split := strings.Split(t, ",")

		if tag == split[0] {
			return field.Name
		}
	}
	return ""
}
