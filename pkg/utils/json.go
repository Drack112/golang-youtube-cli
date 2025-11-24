package utils

import "strconv"

func GetText(obj map[string]any, keys ...string) string {
	current := any(obj)
	for _, key := range keys {
		switch v := current.(type) {
		case map[string]any:
			current = v[key]
		case []any:
			if len(v) > 0 {
				if m, ok := v[0].(map[string]any); ok {
					current = m[key]
				} else {
					return ""
				}
			} else {
				return ""
			}
		default:
			return ""
		}
		if current == nil {
			return ""
		}
	}

	switch v := current.(type) {
	case string:
		return v
	case map[string]any:
		// Handle text objects like {"text": "value"}
		if text, ok := v["text"].(string); ok {
			return text
		}
		if simpleText, ok := v["simpleText"].(string); ok {
			return simpleText
		}
	}

	return ""
}

// Interface{} -> string
func Str(v any) string {
	if v == nil {
		return ""
	}

	switch s := v.(type) {
	case string:
		return s
	default:
		return ""
	}
}

func DeepGet(obj map[string]any, keys ...string) any {
	current := any(obj)
	for _, key := range keys {
		switch v := current.(type) {
		case map[string]any:
			current = v[key]
		case []any:
			if key == "" {
				continue
			}
			// Try to parse key as index
			if idx, err := strconv.Atoi(key); err == nil && idx >= 0 && idx < len(v) {
				current = v[idx]
			} else {
				return nil
			}
		default:
			return nil
		}

		if current == nil {
			return nil
		}
	}

	return current
}
