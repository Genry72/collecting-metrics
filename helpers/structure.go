package helpers

import "reflect"

/*
StructToMap Преобразование стурктуры в map[string]interface{}. Клюс это имя поля, а
значение - значение поля структуры
*/
func StructToMap(data interface{}) map[string]interface{} {

	value := reflect.ValueOf(data)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return nil
	}

	typ := value.Type()

	result := make(map[string]interface{}, value.NumField())

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldName := typ.Field(i).Name
		result[fieldName] = field.Interface()
	}
	return result
}
