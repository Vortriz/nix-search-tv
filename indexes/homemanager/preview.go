package homemanager

import (
	"fmt"
	"io"
	"strings"

	"github.com/3timeslazy/nix-search-tv/indexes/textutil"
	"github.com/3timeslazy/nix-search-tv/style"
)

func Preview(out io.Writer, pkg Package) {
	styler := style.StyledText

	pkgTitle := textutil.PkgName(pkg.Name) + "\n"
	fmt.Fprint(out, pkgTitle)

	desc := strings.TrimSpace(pkg.Description)
	desc = styler.Dim(desc)
	fmt.Fprintln(out, desc+"\n")

	typ := textutil.Prop("type", "", pkg.Type)
	fmt.Fprintln(out, typ)

	def := pkg.Default.Text
	if def != "" {
		def = textutil.Prop(
			"default", "",
			style.PrintCodeBlock(pkg.Default.Text),
		)
		fmt.Fprintln(out, def)
	}

	example := ""
	if pkg.Example.Text != "" {
		example = textutil.Prop(
			"example", "",
			style.PrintCodeBlock(pkg.Example.Text),
		)
		fmt.Fprintln(out, example)
	}
}
