package utils

// import (
// 	"reflect"
// 	"fmt"
// )

// type Json map[interface{}]interface {}
type Json map[string]interface{}

// func (j *Json) AsObject() map[string]interface {} {
// 	res := map[string]interface{}{}
// 	for i, v := range *j {
// 		res[i.(string)] = v
// 	}
// 	return res
// }

func (j *Json) Contains(test interface{}) bool {
	for i, _ := range *j {
		if i == test {
			return true
		}
	}
	return false
}

func (j *Json) get(index interface{}) interface{} {
	return (*j)[index.(string)]
}

// func (j *Json) indexOf(test interface{}) interface{} {
// 	for i, _  := range *j {
// 		if(i == test) {
// 			return true
// 		}
// 	}
// 	return false
// }
