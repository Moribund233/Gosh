package sanitize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"plain text", "hello", "hello"},
		{"html tags", "<script>alert(1)</script>", "&lt;script&gt;alert(1)&lt;/script&gt;"},
		{"quotes", `"hello"`, "&#34;hello&#34;"},
		{"ampersand", "a&b", "a&amp;b"},
		{"empty", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, String(tt.input))
		})
	}
}

func TestStrings(t *testing.T) {
	input := []string{"<a>", "b", "<script>alert(1)</script>"}
	want := []string{"&lt;a&gt;", "b", "&lt;script&gt;alert(1)&lt;/script&gt;"}
	assert.Equal(t, want, Strings(input))
}

func TestStrings_Empty(t *testing.T) {
	assert.Empty(t, Strings(nil))
	assert.Empty(t, Strings([]string{}))
}

func TestMapStrings(t *testing.T) {
	m := map[string]interface{}{
		"name":       "<script>",
		"desc":       "safe text",
		"nested":     map[string]interface{}{"inner": "<img>"},
		"list":       []interface{}{"<a>", "b"},
		"number":     42,
		"bool":       true,
	}
	MapStrings(m)

	assert.Equal(t, "&lt;script&gt;", m["name"])
	assert.Equal(t, "safe text", m["desc"])
	assert.Equal(t, "&lt;img&gt;", m["nested"].(map[string]interface{})["inner"])
	assert.Equal(t, "&lt;a&gt;", m["list"].([]interface{})[0])
	assert.Equal(t, "b", m["list"].([]interface{})[1])
	assert.Equal(t, 42, m["number"])
	assert.Equal(t, true, m["bool"])
}

func TestMapStrings_NestedList(t *testing.T) {
	m := map[string]interface{}{
		"items": []interface{}{
			map[string]interface{}{"title": "<b>bold</b>"},
			map[string]interface{}{"title": "normal"},
		},
	}
	MapStrings(m)

	items := m["items"].([]interface{})
	assert.Equal(t, "&lt;b&gt;bold&lt;/b&gt;", items[0].(map[string]interface{})["title"])
	assert.Equal(t, "normal", items[1].(map[string]interface{})["title"])
}

func TestMapStrings_Empty(t *testing.T) {
	MapStrings(nil)
	MapStrings(map[string]interface{}{})
}
