package main

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize[T any](person T) string {
	t := reflect.TypeOf(person)
	if t.Kind() != reflect.Struct {
		return ""
	}
	v := reflect.ValueOf(person)
	serialized := make([]string, 0, t.NumField())
	for i := range t.NumField() {
		fldType := t.Field(i)
		props := strings.Split(fldType.Tag.Get("properties"), ",")
		if props[0] == "" {
			continue
		}
		omitempty := false
		if len(props) == 2 && strings.ToLower(strings.TrimSpace(props[1])) == "omitempty" {
			omitempty = true
		}
		fldVal := v.Field(i)
		if omitempty && (fldVal.IsZero() || (t.Kind() == reflect.Ptr && fldVal.IsNil())) {
			continue
		}

		serialized = append(serialized, fmt.Sprintf("%s=%s", props[0], valueToString(fldVal)))
	}
	return strings.Join(serialized, "\n")
}

func valueToString(fldVal reflect.Value) string {
	var value string
	switch k := fldVal.Kind(); k {
	case reflect.Slice, reflect.Array:
		vals := make([]string, 0, fldVal.Len())
		for y := range fldVal.Len() {
			switch fldVal.Index(y).Kind() {
			case reflect.Array, reflect.Slice:
				vals = append(vals, valueToString(fldVal.Index(y)))
			case reflect.Map:
				vals = append(vals, valueToString(fldVal.Index(y)))
			default:
				vals = append(vals, fmt.Sprintf("%v", fldVal.Index(y)))
			}
		}
		value = strings.Join(vals, ",")
	case reflect.Map:
		vals := make([]string, 0, fldVal.Len())
		iter := fldVal.MapRange()
		for iter.Next() {
			switch iter.Value().Kind() {
			case reflect.Array, reflect.Slice:
				vals = append(vals, fmt.Sprintf("%v:%v", iter.Key(), valueToString(iter.Value())))
			case reflect.Map:
				vals = append(vals, fmt.Sprintf("%v:%v", iter.Key(), valueToString(iter.Value())))
			default:
				vals = append(vals, fmt.Sprintf("%v:%v", iter.Key(), iter.Value()))
			}
		}
		slices.Sort(vals)
		value = strings.Join(vals, ",")
	case reflect.Chan, reflect.Func:
		if fldVal.IsNil() {
			value = "<nil>"
		} else {
			value = fmt.Sprintf("<%s>", k.String())
		}
	case reflect.Pointer:
		if fldVal.IsNil() {
			value = "<nil>"
		} else {
			value = valueToString(fldVal.Elem())
		}
	case reflect.Struct:
		value = fmt.Sprintf("<%s>", fldVal.Type().Name())
	default:
		value = fmt.Sprintf("%v", fldVal)
	}
	return value
}

func TestDifferentTypes(t *testing.T) {
	tests := map[string]struct {
		person any
		result string
	}{
		"test case with channel": {
			person: make(chan int),
			result: "",
		},
		"test case with func": {
			person: func() {},
			result: "",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}

func TestSerializationAggTypes(t *testing.T) {
	type Details struct {
		NonEmpty bool `properties:"non-empty"`
	}
	type Person2 struct {
		Name    string         `properties:"name"`
		Address string         `properties:"address,omitempty"`
		Age     int            `properties:"age"`
		Married bool           `properties:"married"`
		Tags    []string       `properties:"tags"`
		Scores  map[string]int `properties:"scores"`
		Upgrade func()         `properties:"upgrade"`
		Funcs   []func()       `properties:"funcs"`
		Details *Details       `properties:"details,omitempty"`
	}

	serializePerson2 := Serialize[Person2]

	tests := map[string]struct {
		person Person2
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false\ntags=\nscores=\nupgrade=<nil>\nfuncs=",
		},
		"test case with fields": {
			person: Person2{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Upgrade: nil,
			},
			result: "name=John Doe\nage=30\nmarried=true\ntags=\nscores=\nupgrade=<nil>\nfuncs=",
		},
		"test case with slice field": {
			person: Person2{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
				Tags:    []string{"one", "two"},
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true\ntags=one,two\nscores=\nupgrade=<nil>\nfuncs=",
		},
		"test case with map field": {
			person: Person2{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
				Scores:  map[string]int{"code": 100, "theory": 75},
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true\ntags=\nscores=code:100,theory:75\nupgrade=<nil>\nfuncs=",
		},
		"test case with pointer field": {
			person: Person2{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
				Scores:  map[string]int{"code": 100, "theory": 75},
				Details: &Details{},
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true\ntags=\nscores=code:100,theory:75\nupgrade=<nil>\nfuncs=\ndetails=<Details>",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := serializePerson2(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
