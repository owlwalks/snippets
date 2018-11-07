package snippets

import (
	"html/template"
	"io"
)

func escapeHTML(w io.Writer, b []byte) {
	template.HTMLEscape(w, b)
}

func escapeJS(w io.Writer, b []byte) {
	template.JSEscape(w, b)
}
