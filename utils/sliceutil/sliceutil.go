package sliceutil

import "reflect"

// Contains find elem in a slice and returns if it exists
func Contains(slice, elem interface{}) bool {
	arrV := reflect.ValueOf(slice)
	if arrV.Kind() == reflect.Slice {
		for i := 0; i < arrV.Len(); i++ {
			if arrV.Index(i).Interface() == elem {
				return true
			}
		}
	}
	return false
}

// IndexOf find elem index in a slice and returns it
func IndexOf(slice, elem interface{}) int {
	arrV := reflect.ValueOf(slice)
	if arrV.Kind() == reflect.Slice {
		for i := 0; i < arrV.Len(); i++ {
			if arrV.Index(i).Interface() == elem {
				return i
			}
		}
	}
	return -1
}
