package utils

import (
	"reflect"
)

type CollectData struct {
	Arr interface{}
}

// CollectStruct 声明集合
func CollectStruct(arr interface{}) *CollectData {
	return &CollectData{Arr: arr}
}

// interface 转 maps
func interFaceToMaps(v interface{}) []map[string]interface{} {
	val := reflect.ValueOf(v)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		panic("The passed argument is not a slice or array")
	}

	result := make([]map[string]interface{}, val.Len())

	for i := 0; i < val.Len(); i++ {
		result[i] = make(map[string]interface{})
		result[i] = interFaceToMap(val.Index(i).Interface())
	}

	return result
}

// interface 转 map
func interFaceToMap(v interface{}) map[string]interface{} {
	val := reflect.ValueOf(v)
	t := reflect.TypeOf(v)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	result := make(map[string]interface{})

	for j := 0; j < val.NumField(); j++ {
		valInfo := val.Field(j)

		if valInfo.Kind() == reflect.Ptr {
			valInfo = valInfo.Elem()
		}
		if t.Field(j).Anonymous {
			for i, v := range interFaceToMap(valInfo.Interface()) {
				result[i] = v
			}
		} else {
			result[val.Type().Field(j).Name] = valInfo.Interface()
		}
	}

	return result
}

// PluckToArray Pluck 获取一列数据到数组
func (c *CollectData) PluckToArray(key string, result interface{}) {
	arr := interFaceToMaps(c.Arr)

	v := reflect.ValueOf(result)

	if v.Kind() != reflect.Ptr {
		panic("The passed argument is not a ptr")
	}

	typ := reflect.TypeOf(result).Elem()

	if typ.Kind() != reflect.Slice {
		panic("The passed argument is not a *slice")
	}

	val := reflect.MakeSlice(typ, 0, 0)

	for _, i := range arr {
		val = reflect.Append(val, reflect.ValueOf(i[key]))
	}

	v.Elem().Set(val)
}

// PluckToMap Pluck 获取一列数据到数组
func (c *CollectData) PluckToMap(key, value string, result interface{}) {
	arr := interFaceToMaps(c.Arr)

	v := reflect.ValueOf(result)

	if v.Kind() != reflect.Ptr {
		panic("The passed argument is not a ptr")
	}

	typ := reflect.TypeOf(result).Elem()

	if typ.Kind() != reflect.Map {
		panic("The passed argument is not a map")
	}

	for _, i := range arr {
		v.Elem().SetMapIndex(reflect.ValueOf(i[key]), reflect.ValueOf(i[value]))
	}
}
