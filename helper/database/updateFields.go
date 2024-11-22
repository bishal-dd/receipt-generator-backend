package database

import "reflect"

func UpdateFields(input interface{}) (map[string]interface{}, error ){
    // Convert input to map for selective updates
    updateData := make(map[string]interface{})
    inputValue := reflect.ValueOf(input)
    inputType := inputValue.Type()

    for i := 0; i < inputValue.NumField(); i++ {
        field := inputValue.Field(i)
        
        // Check if field is a pointer and not nil
        if field.Kind() == reflect.Ptr && !field.IsNil() {
            fieldName := inputType.Field(i).Name
            updateData[fieldName] = field.Elem().Interface()
        }
    }

    return updateData, nil
}