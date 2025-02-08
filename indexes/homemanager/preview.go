package homemanager

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/3timeslazy/nix-search-tv/style"
)

func Preview(out io.Writer, pkg Package) {
	styler := style.StyledText

	// fmt.Fprint(out, styler.Red(styler.Bold(pkg.Name)))
	fmt.Fprint(out, styleName(pkg.Key))
	fmt.Fprint(out, "\n\n")

	fmt.Fprint(out, styler.Dim(pkg.Description))
	fmt.Fprint(out, "\n\n")

	fmt.Fprint(out, styler.Bold("type"), "\n")
	fmt.Fprint(out, pkg.Type)
	fmt.Fprint(out, "\n\n")

	// if len(pkg.Subs) > 0 {
	// 	fmt.Fprint(out, styler.Bold("options"), "\n")

	// 	pkg.Subs[0] = "  " + pkg.Subs[0]
	// 	subs := strings.Join(pkg.Subs, "\n  ")

	// 	fmt.Fprint(out, subs)
	// 	fmt.Fprint(out, "\n\n")
	// }
	if len(pkg.Example) > 0 {
		fmt.Fprint(out, styler.Bold("example"), "\n")

		example := style.PrintCodeBlock(pkg.Example["text"].(string))
		fmt.Fprint(out, example)

		fmt.Fprint(out, "\n\n")
	}

	fmt.Fprint(out, "===========")
	fmt.Fprint(out, "\n\n")
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	enc.Encode(pkg)
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
