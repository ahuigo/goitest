package interpolate

import (
	"encoding/json"
	"sync"
)

func interpolateBytes(bs []byte, bindings *sync.Map) (interface{}, error) {
	if len(bs) == 0 {
		return nil, nil
	}
	var jsonExpr interface{}
	if err := json.Unmarshal(bs, &jsonExpr); err != nil {
		return nil, err
	}
	return interpolate(jsonExpr, bindings)
}

// interpolate json expression under binding environment
// - jsonExpr interface{}
// - bindings pointer of sync.Map
// logical:
//
//	Number, Boolean, Null: identity with raw value
//	String: two kinds of format
//	  - pretty simple string: "abc"                     =>  "abc"
//	  - interpolation string: "${abc}" under {"abc": 1} => 1
//	Object:
//	  - key: identity with raw key
//	  - value: recursively interpolate value under binding
//	Array:
//	  - elem: recursively interpolate elem under binding
func interpolate(jsonExpr interface{}, bindings *sync.Map) (_ interface{}, err error) {
	switch rawExpr := jsonExpr.(type) {
	case map[string]any:
		args := make(map[string]interface{}, len(rawExpr))
		for k, v := range rawExpr {
			if args[k], err = interpolate(v, bindings); err != nil {
				return nil, err
			}
		}
		return args, nil
	case []interface{}:
		args := make([]interface{}, len(rawExpr))
		for i, item := range rawExpr {
			if args[i], err = interpolate(item, bindings); err != nil {
				return nil, err
			}
		}
		return args, nil
	case string:
		bindingsMap := make(map[string]interface{})
		bindings.Range(func(k, v interface{}) bool {
			bindingsMap[k.(string)] = v
			return true
		})
		data, _ := json.Marshal(bindingsMap)
		return Interpolation(rawExpr, data)
	default:
		return rawExpr, nil
	}
}
