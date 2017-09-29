package violetear

import (
	"net/http"
)

// Params string/interface map used with context
type (
	Params []param

	param struct {
		name  string
		value interface{}
	}
)

// Set add value to params
func (p Params) Add(name, value string) Params {
	for i, param := range p {
		if param.name == name {
			switch param.value.(type) {
			case string:
				param.value = []string{param.value.(string), value}
			case []string:
				param.value = append(param.value.([]string), value)
			}
			p[i] = param
			return p
		}
	}
	p = append(p, param{name, value})
	return p
}

// GetParam returns a value for the parameter set in path
// When having duplicate params pass the index as the last argument to
// retrieve the desired value.
func GetParam(name string, r *http.Request, index ...int) string {
	params := r.Context().Value(ParamsKey).(Params)
	for _, param := range params {
		if param.name == ":"+name {
			switch p := param.value.(type) {
			case []string:
				if len(index) > 0 {
					if index[0] > len(p) {
						return ""
					}
					return p[index[0]]
				}
				return p[0]
			default:
				return p.(string)
			}
		}
	}
	return ""
}

// GetParams returns param or params in a []string
func GetParams(name string, r *http.Request) []string {
	params := r.Context().Value(ParamsKey).(Params)
	for _, param := range params {
		if param.name == ":"+name {
			switch p := param.value.(type) {
			case []string:
				return p
			default:
				return []string{p.(string)}
			}
		}
	}
	return []string{}
}
