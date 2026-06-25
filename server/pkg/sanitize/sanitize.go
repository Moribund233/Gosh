package sanitize

import "html"

func String(s string) string {
	return html.EscapeString(s)
}

func Strings(ss []string) []string {
	out := make([]string, len(ss))
	for i, s := range ss {
		out[i] = String(s)
	}
	return out
}

func MapStrings(m map[string]interface{}) {
	for k, v := range m {
		switch val := v.(type) {
		case string:
			m[k] = String(val)
		case map[string]interface{}:
			MapStrings(val)
		case []interface{}:
			for i, item := range val {
				if s, ok := item.(string); ok {
					val[i] = String(s)
				} else if m2, ok := item.(map[string]interface{}); ok {
					MapStrings(m2)
				}
			}
		}
	}
}
