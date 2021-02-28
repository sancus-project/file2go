package static

type FuncMap map[string]interface{}

func (c Collection) NewFuncMap(hashify bool) FuncMap {

	m := make(map[string]interface{}, 2)
	m["Filename"] = c.getFilenameFunc(hashify)
	m["Filetype"] = c.getFiletypeFunc(hashify)

	return m
}

func (c Collection) getFilenameFunc(hashify bool) interface{} {

	if hashify {
		return func(fn0 string) string {
			if fn1, ok := c.Redirects[fn0]; ok {
				return fn1
			}
			return fn0
		}
	}

	return func(fn0 string) string {
		return fn0
	}
}

func (c Collection) getFiletypeFunc(hashify bool) interface{} {

	return func(fn0 string) string {
		if v, ok := c.Files[fn0]; ok {
			return v.ContentType
		} else {
			return "application/octet-stream"
		}
	}
}
