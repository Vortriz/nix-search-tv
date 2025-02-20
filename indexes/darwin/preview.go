package darwin

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

	def := pkg.Default
	if def != "" {
		def = textutil.Prop(
			"default", "",
			style.PrintCodeBlock(pkg.Default),
		)
		fmt.Fprintln(out, def)
	}

	example := ""
	if pkg.Example != "" {
		example = textutil.Prop(
			"example", "",
			style.PrintCodeBlock(pkg.Example),
		)
		fmt.Fprintln(out, example)
	}
}

func (pkg *Package) GetSource() string {
	if len(pkg.DeclaredBy) == 1 {
		return pkg.DeclaredBy[0]
	}

	return fmt.Sprintf(
		"https://daiderd.com/nix-darwin/manual/index.html#opt-%s",

		// There are packages with quotes in their names, like
		// system.defaults.".GlobalPreferences"."com.apple.mouse.scaling". For these,
		// the quotes must replaces with "_"
		strings.ReplaceAll(pkg.Name, `"`, "_"),
	)
}
