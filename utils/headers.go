package utils

type Headers map[string]string

func (h Headers) Add(k string, v string) Headers {
	_, present := h[k]
	if !present {
		h[k] = v
	}
	return h
}

func (h Headers) Get(k string) string {
	v, present := h[k]
	if present {
		return v
	}
	return ""
}
