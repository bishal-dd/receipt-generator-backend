package database

import (
	"reflect"
	"time"

	"github.com/bishal-dd/receipt-generator-backend/helper/ids"
)


func CreateFields[T any](input interface{}) *T {
	model := new(T)

	// Use reflection to set common fields
	v := reflect.ValueOf(model).Elem()

	// Set ID if the struct has an ID field
	if idField := v.FieldByName("ID"); idField.IsValid() && idField.CanSet() {
		idField.SetString(ids.UILD())
	}

	// Set CreatedAt if the struct has a CreatedAt field
	if createdAtField := v.FieldByName("CreatedAt"); createdAtField.IsValid() && createdAtField.CanSet() {
		createdAtField.Set(reflect.ValueOf(time.Now()))
	}

	// Copy input fields to the model
	inputV := reflect.ValueOf(input)
	inputT := inputV.Type()

	for i := 0; i < inputV.NumField(); i++ {
		fieldName := inputT.Field(i).Name
		field := v.FieldByName(fieldName)

		if field.IsValid() && field.CanSet() {
			field.Set(inputV.Field(i))
		}
	}

	return model
}
