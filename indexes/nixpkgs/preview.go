package nixpkgs

import (
	"cmp"
	"fmt"
	"io"
	"runtime"
	"slices"
	"strings"

	"github.com/3timeslazy/nix-search-tv/indexes/textutil"
	"github.com/3timeslazy/nix-search-tv/style"
)

func Preview(out io.Writer, pkg Package) {
	styler := style.StyledText

	pkgTitle := textutil.PkgName(pkg.Name) + " " + styler.Dim("("+pkg.GetVersion()+")")
	if pkg.Meta.Broken {
		pkgTitle += " " + styler.Red("(broken)")
	}
	fmt.Fprintln(out, pkgTitle)

	desc := ""
	if pkg.Meta.Description != "" {
		desc = style.Wrap(pkg.Meta.Description, "") + "\n"
	}
	fmt.Fprintln(out, desc)

	longDesc := ""
	if pkg.Meta.LongDescription != "" && pkg.Meta.Description != pkg.Meta.LongDescription {
		longDesc = style.StyleLongDescription(styler, pkg.Meta.LongDescription)
		fmt.Fprintln(out, longDesc)
	}

	homepages := ""
	if hmpgs := len(pkg.Meta.Homepages); hmpgs > 0 {
		homepages = textutil.Prop(
			textutil.IfElse(hmpgs > 1, "homepage", "homepages"), "",
			strings.Join(pkg.Meta.Homepages, "\n"),
		)
		fmt.Fprintln(out, homepages)
	}

	licenseType := textutil.IfElse(pkg.Meta.Unfree, "unfree", "free")
	license := textutil.Prop(
		"license", styler.Dim("("+licenseType+")"),
		licensesString(pkg.Meta.Licenses),
	)
	fmt.Fprintln(out, license)

	mainProg := ""
	if pkg.Meta.MainProgram != "" {
		mainProg = textutil.Prop(
			"main program", "",
			style.PrintCodeBlock("$ "+pkg.Meta.MainProgram),
		)
		fmt.Fprintln(out, mainProg)
	}

	platforms := ""
	if len(pkg.Meta.Platforms) > 0 {
		platforms = textutil.Prop(
			"platforms", "",
			platformsString(pkg.Meta.Platforms),
		)
		fmt.Fprintln(out, platforms)
	}

}

func licensesString(ls []License) string {
	if len(ls) == 0 {
		return "No License"
	}

	ss := []string{}
	for _, l := range ls {
		ss = append(ss, cmp.Or(l.SpdxID, l.FullName))
	}

	return strings.Join(ss, "\n")
}

var printablePlatforms = []string{
	"x86_64-linux",
	"aarch64-linux",
	"i686-linux",
	"x86_64-darwin",
	"aarch64-darwin",
}

func platformsString(ps []string) string {
	toPrint := []string{}
	cp := currentPlatform()

	// iterate in order printablePlatforms->ps and not
	// vise versa to always print platforms in the same order
	// without sorting
	for _, printable := range printablePlatforms {
		if !slices.Contains(ps, printable) {
			continue
		}
		if printable != cp {
			printable = style.StyledText.Dim(printable)
		}
		toPrint = append(toPrint, printable)
	}

	return strings.Join(toPrint, "\n")
}

var go2nixArch = map[string]string{
	"arm64": "aarch64",
	"amd64": "x86_64",
}

func currentPlatform() string {
	arch := go2nixArch[runtime.GOARCH]
	kern := runtime.GOOS
	return arch + "-" + kern
}




