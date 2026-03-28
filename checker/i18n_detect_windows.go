//go:build windows

package checker

import "syscall"

var (
	kernel32                   = syscall.NewLazyDLL("kernel32.dll")
	procGetUserDefaultUILang   = kernel32.NewProc("GetUserDefaultUILanguage")
)

func DetectOSLocale() string {
	ret, _, _ := procGetUserDefaultUILang.Call()
	primary := uint16(ret) & 0x3FF
	switch primary {
	case 0x0C:
		return "fr"
	case 0x09:
		return "en"
	case 0x0A:
		return "es"
	case 0x07:
		return "de"
	case 0x16:
		return "pt"
	default:
		return "en"
	}
}
