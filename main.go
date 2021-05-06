package main

import (
	"errors"
	"fmt"
	"reflect"
)

// 1. Написать функцию, которая принимает на вход структуру in
// (struct или кастомную struct) и values map[string]interface{}
// (key - название поля структуры, которому нужно присвоить value этой мапы).
// Необходимо по значениям из мапы изменить входящую структуру in с помощью пакета reflect.

type Client struct {
	Username string
	Number   string
}

func main() {
	u := &Client{}
	m := map[string]interface{}{
		// "Username": "client",
		// "Number":   1,
		// "Hello":    "good",
	}

	err := myFunc(u, m)

	if err != nil {
		fmt.Printf("%v+", err)
	}

	fmt.Println(u)
}

func myFunc(str interface{}, m map[string]interface{}) error {

	if m == nil {
		return errors.New("arg is nil")
	}

	strValue := reflect.ValueOf(str)
	strElem := strValue.Elem()

	for key, _ := range m {

		if strElem.FieldByName(key).IsValid() {
			mapValue := reflect.ValueOf(m[key])
			strElemField := strElem.FieldByName(key)

			if strElemField.Kind() == mapValue.Kind() {
				strElemField.Set(mapValue)
			}

		}
	}

	return nil
}
