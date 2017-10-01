package violetear

import "net/http"

// Params string/interface map used with context
type Params map[string]interface{}

// Add param to Params
func (p Params) Add(k, v string) Params {
	if param, ok := p[k]; ok {
		switch param.(type) {
		case string:
			param = []string{param.(string), v}
		case []string:
			param = append(param.([]string), v)
		}
		p[k] = param
	} else {
		p[k] = v
	}
	return p
}

// GetParam returns a value for the parameter set in path
// When having duplicate params pass the index as the last argument to
// retrieve the desired value.
func GetParam(name string, r *http.Request, index ...int) string {
	params := r.Context().Value(ParamsKey).(Params)
	if param := params[":"+name]; param != nil {
		switch param := param.(type) {
		case []string:
			if len(index) > 0 {
				if index[0] > len(param) {
					return ""
				}
				return param[index[0]]
			}
			return param[0]
		default:
			return param.(string)
		}
	}
	return ""
}

// GetParams returns param or params in a []string
func GetParams(name string, r *http.Request) []string {
	params := r.Context().Value(ParamsKey).(Params)
	if param := params[":"+name]; param != nil {
		switch param := param.(type) {
		case []string:
			return param
		default:
			return []string{param.(string)}
		}
	}
	return []string{}
}
