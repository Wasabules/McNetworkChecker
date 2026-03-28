package checker

import (
	"encoding/json"
	"strings"
)

// MOTDSpan is a single styled text fragment for frontend rendering.
type MOTDSpan struct {
	Text          string `json:"text"`
	Color         string `json:"color"`
	Bold          bool   `json:"bold,omitempty"`
	Italic        bool   `json:"italic,omitempty"`
	Underlined    bool   `json:"underlined,omitempty"`
	Strikethrough bool   `json:"strikethrough,omitempty"`
	Obfuscated    bool   `json:"obfuscated,omitempty"`
}

// Minecraft named color codes → hex
var mcNamedColors = map[string]string{
	"black":        "#000000",
	"dark_blue":    "#0000AA",
	"dark_green":   "#00AA00",
	"dark_aqua":    "#00AAAA",
	"dark_red":     "#AA0000",
	"dark_purple":  "#AA00AA",
	"gold":         "#FFAA00",
	"gray":         "#AAAAAA",
	"dark_gray":    "#555555",
	"blue":         "#5555FF",
	"green":        "#55FF55",
	"aqua":         "#55FFFF",
	"red":          "#FF5555",
	"light_purple": "#FF55FF",
	"yellow":       "#FFFF55",
	"white":        "#FFFFFF",
}

func resolveColor(c string) string {
	if strings.HasPrefix(c, "#") {
		return c
	}
	if hex, ok := mcNamedColors[c]; ok {
		return hex
	}
	return "#FFFFFF"
}

type motdStyle struct {
	Color         string
	Bold          bool
	Italic        bool
	Underlined    bool
	Strikethrough bool
	Obfuscated    bool
}

// fullChatComponent captures all Minecraft chat formatting fields.
// Pointer bools distinguish "not set" (nil → inherit) from "explicitly false" (override).
type fullChatComponent struct {
	Text          string            `json:"text"`
	Color         string            `json:"color"`
	Bold          *bool             `json:"bold"`
	Italic        *bool             `json:"italic"`
	Underlined    *bool             `json:"underlined"`
	Strikethrough *bool             `json:"strikethrough"`
	Obfuscated    *bool             `json:"obfuscated"`
	Extra         []json.RawMessage `json:"extra"`
}

// parseDescription returns plain text (for logs/report) and rich spans (for UI).
func parseDescription(raw json.RawMessage) (string, []MOTDSpan) {
	if len(raw) == 0 {
		return "", nil
	}

	var str string
	if err := json.Unmarshal(raw, &str); err == nil {
		return str, []MOTDSpan{{Text: str, Color: "#FFFFFF"}}
	}

	spans := flattenComponent(raw, motdStyle{Color: "#FFFFFF"})
	var sb strings.Builder
	for _, span := range spans {
		sb.WriteString(span.Text)
	}
	return sb.String(), spans
}

// flattenComponent recursively converts a chat component tree into a flat
// list of styled spans. Handles mixed arrays of strings and objects.
func flattenComponent(raw json.RawMessage, parent motdStyle) []MOTDSpan {
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		if s == "" {
			return nil
		}
		return []MOTDSpan{{
			Text: s, Color: parent.Color, Bold: parent.Bold,
			Italic: parent.Italic, Underlined: parent.Underlined,
			Strikethrough: parent.Strikethrough, Obfuscated: parent.Obfuscated,
		}}
	}

	var comp fullChatComponent
	if err := json.Unmarshal(raw, &comp); err != nil {
		return nil
	}

	style := parent
	if comp.Color != "" {
		style.Color = resolveColor(comp.Color)
	}
	if comp.Bold != nil {
		style.Bold = *comp.Bold
	}
	if comp.Italic != nil {
		style.Italic = *comp.Italic
	}
	if comp.Underlined != nil {
		style.Underlined = *comp.Underlined
	}
	if comp.Strikethrough != nil {
		style.Strikethrough = *comp.Strikethrough
	}
	if comp.Obfuscated != nil {
		style.Obfuscated = *comp.Obfuscated
	}

	var spans []MOTDSpan
	if comp.Text != "" {
		spans = append(spans, MOTDSpan{
			Text: comp.Text, Color: style.Color, Bold: style.Bold,
			Italic: style.Italic, Underlined: style.Underlined,
			Strikethrough: style.Strikethrough, Obfuscated: style.Obfuscated,
		})
	}
	for _, child := range comp.Extra {
		spans = append(spans, flattenComponent(child, style)...)
	}
	return spans
}
