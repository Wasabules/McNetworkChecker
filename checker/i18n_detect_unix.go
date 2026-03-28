//go:build !windows

package checker

import (
	"os"
	"strings"
)

var supported = map[string]bool{"fr": true, "en": true, "es": true, "de": true, "pt": true}

func DetectOSLocale() string {
	for _, env := range []string{"LC_MESSAGES", "LC_ALL", "LANG"} {
		if v := os.Getenv(env); len(v) >= 2 {
			lang := strings.ToLower(v[:2])
			if supported[lang] {
				return lang
			}
		}
	}
	return "en"
}
