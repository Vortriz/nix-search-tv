package nixpkgs

import (
	"cmp"
	"fmt"
	"io"
	"strings"

	"github.com/3timeslazy/nix-search-tv/style"
)

func Preview(out io.Writer, pkg Package) {
	styler := style.StyledText

	// title
	fmt.Fprint(out, styler.Red(styler.Bold(pkg.FullName)))
	fmt.Fprint(out, " ", styler.Dim("("+pkg.GetVersion()+")"))
	if pkg.Meta.Broken {
		fmt.Fprint(out, " ", styler.Red("(broken)"))
	}
	fmt.Fprintln(out)

	fmt.Fprint(out, style.Wrap(pkg.Meta.Description, ""))
	// two new lines instead of one here and after to make `tv` render it as a single new line
	fmt.Fprint(out, "\n\n")

	if pkg.Meta.LongDescription != "" && pkg.Meta.Description != pkg.Meta.LongDescription {
		fmt.Fprint(out, style.StyleLongDescription(styler, pkg.Meta.LongDescription), "\n")
	}

	if len(pkg.Meta.Homepages) > 0 {
		subtitle := "homepage"
		if len(pkg.Meta.Homepages) > 1 {
			subtitle += "s"
		}
		fmt.Fprint(out, styler.Bold(subtitle), "\n")
		fmt.Fprint(out, strings.Join(pkg.Meta.Homepages, "\n"))
		fmt.Fprint(out, "\n\n")
	}

	licenseType := "free"
	if pkg.Meta.Unfree {
		licenseType = "unfree"
	}
	fmt.Fprint(out, styler.Bold("license"))
	fmt.Fprint(out, styler.Dim(" ("+licenseType+")"), "\n")
	fmt.Fprint(out, licensesString(pkg.Meta.Licenses))
	fmt.Fprint(out, "\n\n")

	if pkg.Meta.MainProgram != "" {
		fmt.Fprint(out, styler.Bold("main program"), "\n")
		fmt.Fprint(out, style.PrintCodeBlock("$ "+pkg.Meta.MainProgram))
		fmt.Fprint(out, "\n\n")
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
