package textutil

import (
	"strings"

	"github.com/3timeslazy/nix-search-tv/style"
)

var s = style.StyledText

func PkgName(pkg string) string {
	styler := style.StyledText

	last := strings.LastIndex(pkg, ".")
	if last == -1 {
		return styler.Red(styler.Bold(pkg))
	}

	left := styler.Red(pkg[:last])
	right := styler.Red(styler.Bold(pkg[last:]))
	return left + right
}

func Prop(name string, mods string, text string) string {
	name = s.Bold(name)
	if mods != "" {
		name += " " + mods
	}
	return name + "\n" + text + "\n"
}

func IfElse(cond bool, ok, notok string) string {
	if cond {
		return ok
	}
	return notok
}

