//+------------------------------------------------------------------+
//|                                             MetaTrader 5 Web API |
//|                             Copyright 2000-2025, MetaQuotes Ltd. |
//|                                               www.metaquotes.net |
//+------------------------------------------------------------------+
package mt5webapi

import (
	"fmt"
	"strings"
)

// JSONWriter builds JSON output using a strings.Builder.
type JSONWriter struct {
	builder        strings.Builder
	isFirstElement bool
}

// NewJSONWriter creates a new JSONWriter.
func NewJSONWriter() *JSONWriter {
	return &JSONWriter{}
}

// String returns the accumulated JSON string.
func (w *JSONWriter) String() string {
	return w.builder.String()
}

// WriteBeginObject writes '{' and resets the first-element flag.
func (w *JSONWriter) WriteBeginObject() {
	w.isFirstElement = true
	w.builder.WriteByte('{')
}

// WriteEndObject writes '}'.
func (w *JSONWriter) WriteEndObject() {
	w.builder.WriteByte('}')
}

// writeBeginAttribute writes the attribute name prefix: ,"name":
func (w *JSONWriter) writeBeginAttribute(name string) {
	if !w.isFirstElement {
		w.builder.WriteByte(',')
	}
	fmt.Fprintf(&w.builder, "\"%s\":", name)
}

// writeBeginCollection writes '['.
func (w *JSONWriter) writeBeginCollection() {
	w.isFirstElement = true
	w.builder.WriteByte('[')
}

// writeEndCollection writes ']'.
func (w *JSONWriter) writeEndCollection() {
	w.builder.WriteByte(']')
}

// WriteAttributeString writes a JSON string attribute: "name":"value"
func (w *JSONWriter) WriteAttributeString(name, value string) {
	w.WriteAttributeStringQuote(name, value, true)
}

// WriteAttributeStringQuote writes a string attribute with optional quoting.
func (w *JSONWriter) WriteAttributeStringQuote(name, value string, quoteString bool) {
	w.writeBeginAttribute(name)
	w.writeObjectString(value, quoteString)
	w.isFirstElement = false
}

// WriteAttributeInt writes an int attribute.
func (w *JSONWriter) WriteAttributeInt(name string, value int) {
	w.writeBeginAttribute(name)
	fmt.Fprintf(&w.builder, "%d", value)
	w.isFirstElement = false
}

// WriteAttributeUint writes a uint attribute.
func (w *JSONWriter) WriteAttributeUint(name string, value uint) {
	w.writeBeginAttribute(name)
	fmt.Fprintf(&w.builder, "%d", value)
	w.isFirstElement = false
}

// WriteAttributeInt64 writes an int64 attribute.
func (w *JSONWriter) WriteAttributeInt64(name string, value int64) {
	w.writeBeginAttribute(name)
	fmt.Fprintf(&w.builder, "%d", value)
	w.isFirstElement = false
}

// WriteAttributeUint32 writes a uint32 attribute.
func (w *JSONWriter) WriteAttributeUint32(name string, value uint32) {
	w.writeBeginAttribute(name)
	fmt.Fprintf(&w.builder, "%d", value)
	w.isFirstElement = false
}

// WriteAttributeInt32 writes an int32 attribute.
func (w *JSONWriter) WriteAttributeInt32(name string, value int32) {
	w.writeBeginAttribute(name)
	fmt.Fprintf(&w.builder, "%d", value)
	w.isFirstElement = false
}

// WriteAttributeUint64 writes a uint64 attribute.
func (w *JSONWriter) WriteAttributeUint64(name string, value uint64) {
	w.writeBeginAttribute(name)
	fmt.Fprintf(&w.builder, "%d", value)
	w.isFirstElement = false
}

// WriteAttributeRaw writes a pre-formatted raw JSON value (no quoting).
func (w *JSONWriter) WriteAttributeRaw(name string, rawJSON string) {
	w.writeBeginAttribute(name)
	w.builder.WriteString(rawJSON)
	w.isFirstElement = false
}

// WriteAttributeFloat64 writes a float64 attribute with invariant decimal separator.
func (w *JSONWriter) WriteAttributeFloat64(name string, value float64) {
	w.writeBeginAttribute(name)
	w.writeObjectFloat64(value)
	w.isFirstElement = false
}

// WriteAttributeFloat32 writes a float32 attribute with invariant decimal separator.
func (w *JSONWriter) WriteAttributeFloat32(name string, value float32) {
	w.writeBeginAttribute(name)
	w.writeObjectFloat32(value)
	w.isFirstElement = false
}

// WriteAttributeIEnumerable writes a collection attribute (array of strings or numbers).
func (w *JSONWriter) WriteAttributeIEnumerable(name string, values []interface{}) {
	w.writeBeginAttribute(name)
	w.isFirstElement = true
	w.WriteCollection(values)
	w.isFirstElement = false
}

// writeObjectString writes a JSON string value with optional quoting.
func (w *JSONWriter) writeObjectString(value string, quoteString bool) {
	if quoteString {
		w.builder.WriteByte('"')
		w.builder.WriteString(QuoteString(value))
		w.builder.WriteByte('"')
	} else {
		w.builder.WriteString(value)
	}
}

// writeObjectFloat64 writes a float64 value using Go's %G format (matches C# "G").
func (w *JSONWriter) writeObjectFloat64(value float64) {
	fmt.Fprintf(&w.builder, "%G", value)
}

// writeObjectFloat32 writes a float32 value using Go's %G format (matches C# "G").
func (w *JSONWriter) writeObjectFloat32(value float32) {
	fmt.Fprintf(&w.builder, "%G", value)
}

// WriteCollection writes an array of values: [val1,val2,...].
// Supported element types: string, int, int64, uint, uint64, float64, float32, []interface{}.
func (w *JSONWriter) WriteCollection(collection []interface{}) {
	w.writeBeginCollection()
	for _, item := range collection {
		if !w.isFirstElement {
			w.builder.WriteByte(',')
		}
		w.isFirstElement = false
		switch v := item.(type) {
		case string:
			w.writeObjectString(v, true)
		case int:
			fmt.Fprintf(&w.builder, "%d", v)
		case int64:
			fmt.Fprintf(&w.builder, "%d", v)
		case uint:
			fmt.Fprintf(&w.builder, "%d", v)
		case uint64:
			fmt.Fprintf(&w.builder, "%d", v)
		case float64:
			w.writeObjectFloat64(v)
		case float32:
			w.writeObjectFloat32(v)
		case []interface{}:
			w.isFirstElement = true
			w.WriteCollection(v)
		default:
			fmt.Fprintf(&w.builder, "%v", v)
		}
	}
	w.writeEndCollection()
}

// QuoteString escapes special characters in a string for JSON output.
// Handles: \r, \t, ", ', <, >, \\, \n, \b, \f, and control chars < 0x20.
func QuoteString(value string) string {
	if len(value) == 0 {
		return ""
	}
	var builder *strings.Builder
	startIndex := 0
	count := 0
	for i, c := range value {
		runeLen := len(string(c))
		needsEscape := c == '\r' || c == '\t' || c == '"' || c == '\'' ||
			c == '<' || c == '>' || c == '\\' || c == '\n' ||
			c == '\b' || c == '\f' || c < ' '
		if needsEscape {
			if builder == nil {
				b := strings.Builder{}
				b.Grow(len(value) + 5)
				builder = &b
			}
			if count > 0 {
				builder.WriteString(value[startIndex : startIndex+count])
			}
			startIndex = i + runeLen
			count = 0
		}
		switch c {
		case '<', '>', '\'':
			if builder != nil {
				fmt.Fprintf(builder, "\\u%04x", c)
			}
			continue
		case '\\':
			if builder != nil {
				builder.WriteString(`\\`)
			}
			continue
		case '\b':
			if builder != nil {
				builder.WriteString(`\b`)
			}
			continue
		case '\t':
			if builder != nil {
				builder.WriteString(`\t`)
			}
			continue
		case '\n':
			if builder != nil {
				builder.WriteString(`\n`)
			}
			continue
		case '\f':
			if builder != nil {
				builder.WriteString(`\f`)
			}
			continue
		case '\r':
			if builder != nil {
				builder.WriteString(`\r`)
			}
			continue
		case '"':
			if builder != nil {
				builder.WriteString(`\"`)
			}
			continue
		}
		if c < ' ' {
			if builder != nil {
				fmt.Fprintf(builder, "\\u%04x", c)
			}
		} else {
			count += runeLen
		}
	}
	if builder == nil {
		return value
	}
	if count > 0 {
		builder.WriteString(value[startIndex : startIndex+count])
	}
	return builder.String()
}
