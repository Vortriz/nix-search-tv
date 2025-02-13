package homemanager

import (
	"fmt"
	"io"
	"strings"

	"github.com/3timeslazy/nix-search-tv/style"
)

func Preview(out io.Writer, pkg Package) {
	styler := style.StyledText

	fmt.Fprint(out, styleName(pkg.Name))
	fmt.Fprint(out, "\n\n")

	fmt.Fprint(out, styler.Dim(pkg.Description))
	fmt.Fprint(out, "\n\n")

	fmt.Fprint(out, styler.Bold("type"), "\n")
	fmt.Fprint(out, pkg.Type)
	fmt.Fprint(out, "\n\n")

	if pkg.Example.Text != "" {
		fmt.Fprint(out, styler.Bold("example"), "\n")

		example := style.PrintCodeBlock(pkg.Example.Text)
		fmt.Fprint(out, example)

		fmt.Fprint(out, "\n\n")
	}

	// fmt.Fprint(out, "====== Debug ======")
	// fmt.Fprint(out, "\n\n")
	// enc := json.NewEncoder(out)
	// enc.SetIndent("", "  ")
	// enc.Encode(pkg)
}

func styleName(name string) string {
	styler := style.StyledText

	last := strings.LastIndex(name, ".")
	if last == -1 {
		return styler.Red(styler.Bold(name))
	}

	left := styler.Red(name[:last])
	right := styler.Red(styler.Bold(name[last:]))
	return left + right
}
