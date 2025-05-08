package main

import (
	"errors"
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

func Serialize(object any) string {
	return serialize(reflect.ValueOf(object))
}

func serialize(val reflect.Value) string {
	t := val.Type()
	if t.Kind() != reflect.Struct {
		return ""
	}
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
		fldVal := val.Field(i)
		if omitempty && fldVal.IsZero() {
			continue
		}

		serialized = append(serialized, fmt.Sprintf("%s=%s", props[0], valueToString(fldVal)))
	}
	return strings.Join(serialized, "\n")
}

func valueToString(fldVal reflect.Value) string {
	var value string
	switch k := fldVal.Kind(); k {
	case reflect.Invalid:
		return "<invalid>"
	case reflect.Slice, reflect.Array:
		vals := make([]string, 0, fldVal.Len())
		for y := range fldVal.Len() {
			vals = append(vals, valueToString(fldVal.Index(y)))
		}
		value = strings.Join(vals, ",")
	case reflect.Map:
		vals := make([]string, 0, fldVal.Len())
		iter := fldVal.MapRange()
		for iter.Next() {
			vals = append(vals, fmt.Sprintf("%v:%v", iter.Key(), iter.Value()))
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
		value = fmt.Sprintf("[\n%s\n]", serialize(fldVal))
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
		Description string `properties:"desc"`
		Salary      int    `properties:"salary"`
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
		Error   error          `properties:"error,omitempty"`
	}

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
				Details: &Details{Description: "description", Salary: 100},
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true\ntags=\nscores=code:100,theory:75\nupgrade=<nil>\nfuncs=\ndetails=[\ndesc=description\nsalary=100\n]",
		},
		"test case with interface field": {
			person: Person2{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
				Scores:  map[string]int{"code": 100, "theory": 75},
				Error:   errors.New("error message"),
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true\ntags=\nscores=code:100,theory:75\nupgrade=<nil>\nfuncs=\nerror=error message",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
