package main

import (
	"fmt"
	"reflect"
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
		if omitempty && fldVal.IsZero() {
			continue
		}
		serialized = append(serialized, fmt.Sprintf("%s=%v", props[0], fldVal))
	}
	return strings.Join(serialized, "\n")
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
