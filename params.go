package violetear

// Params string/interface map used with context
type Params []param

func (p *Params) Set(name, value string) {
	for _, param := range *p {
		if param.name == name {
			switch param.value.(type) {
			case string:
				param.value = []string{param.value.(string), value}
			case []string:
				param.value = append(param.value.([]string), value)
			}
			return
		}
	}
	*p = append(*p, param{name, value})
}

func (p *Params) Get(name string) *param {
	return nil
}

type param struct {
	name  string
	value interface{}
}

/*
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
*/
