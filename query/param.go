package query

import (
	"reflect"
)

// Given a list of names and arbitrary number of arguments (primitive values,
// maps or structs), extract the arguments named in the list, place them in the
// an interface list and return a pointer of it. Nil will be returned if nowhere
// can a named argument be found.
// - nameList can contain empty strings "" to indicate that its value should
// 	 come from a positional argument. The first argument will be mapped to the
//   first empty string, and so on. The last arguments should be pointers to
//   struct instances or maps. Their order don't matter, but if more then one
//   struct field / map entry have the same name / key, the first one's value
//   will be used. Any arguments that are not struct / map / pointer to them
//   will be ignored once all the spots corresponding to the empty strings in
//   nameList are filled.
//
// For example:
// ```
// var k = &struct{Hello int, World int}{1, 2}
// var l = &struct{Hello int, World int}{3, 4}
// PrepareParameters([]string{"Hello", "", "World"}, "str", &k, &l)
// ```
// The call will return the list `[]interface{}{1, "str", 2}`
func PrepareParameters(nameList *[]string, args ... interface{}) *[]interface{} {
	l := len(*nameList)
	var res = make([]interface{}, l)
	var namePos = make(map[string]int, l)

	j := 0
	filled := l
	argLen := len(args)
	for i, name := range *nameList {
		if name == "" {
			if j >= argLen {
				return nil
			}
			res[i] = args[j]
			j++
			filled--
		} else {
			namePos[name] = i
		}
	}
	for i := j; i < argLen; i++ {
		arg := args[i]
		tmp := reflect.ValueOf(arg)
		if tmp.Kind() == reflect.Ptr {
			tmp = tmp.Elem() // Should be a struct / map
		}
		kind := tmp.Kind()
		// Inspect the struct and map, decrement 'filled' as 'res' gets filled
		if kind == reflect.Struct {
			structType := tmp.Type()
			numFields := structType.NumField()
			for i := 0; i < numFields; i++ {
				name := structType.Field(i).Name
				if pos, in := namePos[name]; in {
					delete(namePos, name)
					res[pos] = tmp.Field(i).Interface()
					filled--
				}
			}
		} else if kind == reflect.Map {
			names := tmp.MapKeys()
			for _, name := range names {
				// All the non-string keys will be ignored.
				if key, ok := name.Interface().(string); ok {
					if pos, in := namePos[key]; in {
						delete(namePos, key)
						res[pos] = tmp.MapIndex(name).Interface()
						filled--
					}
				}
			}
		}
		if filled == 0 {
			break
		}
	}
	if filled > 0 {
		return nil
	}
	return &res
}