package style

import (
	"fmt"
	"go/doc"
	"go/doc/comment"
	"os"
	"regexp"
	"strings"
	"sync"

	"golang.org/x/term"
)

var (
	reFencedCodeBlock = regexp.MustCompile(`(?ms)\x60\x60\x60+\s*(.*?)\s*\x60\x60\x60+`)
	reInlineHyperlink = regexp.MustCompile(`(?m)\[(.*?)\]\n*\((http.*?)\)`)
	reInlineCode      = regexp.MustCompile(`(?m)\x60\x60?(.*?)\x60\x60?`)
)

var termWidth = sync.OnceValue(func() int {
	termWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	return termWidth
})

func Wrap(text, indent string) string {
	width := min(termWidth(), 80)
	if width < 1 {
		return indent + text
	}

	d := new(doc.Package).Parser().Parse(text)
	pr := &comment.Printer{
		TextCodePrefix: indent + indent,
		TextPrefix:     indent,
		TextWidth:      width,
	}

	return string(pr.Text(d))
}

func StyleLongDescription(styler TextStyler, text string) string {
	linkReplace := styler.Bold("$1") + styler.with(dontEndStyle).Dim(" ($2)")
	codeReplace := styler.Bold("$1") + styler.with(dontEndStyle).Dim("")

	for _, f := range []func(string) string{
		func(text string) string {
			var sb strings.Builder
			sb.Grow(len(text))

			var start int
			for _, is := range reFencedCodeBlock.FindAllStringSubmatchIndex(text, -1) {
				sb.WriteString(text[start:is[0]])
				start = is[1]
				for _, codeLine := range strings.Split(text[is[2]:is[3]], "\n") {
					sb.WriteString("\t")
					sb.WriteString(codeLine)
					sb.WriteString("\n")
				}
			}
			sb.WriteString(text[start:])

			return sb.String()
		},
		func(text string) string { return reInlineHyperlink.ReplaceAllString(text, linkReplace) },
		func(text string) string { return reInlineCode.ReplaceAllString(text, codeReplace) },
		func(text string) string { return Wrap(text, "") }, // this line differs from the original
		styler.Dim,
	} {
		text = f(text)
	}

	return text
}

// PrintCodeBlock prints the given content inside a styled code block
//
// TODO: at the moment of writing, television couldn't correctly
// render box-drawing characters, so use +, | and - for now. Once
// there is a fix, change the lines to those below:
//
//	topBorder := "┌" + strings.Repeat("─", lineWidth-2) + "┐"
//	bottomBorder := "└" + strings.Repeat("─", lineWidth-2) + "┘"
//	leftBorder := "│ "
//	rightBorder := " │"
func PrintCodeBlock(content string) string {
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "FZF_") {
			return PrintCodeBlockV2(content)
		}
	}

	// Determine the maximum width of the content
	lineWidth := findMaxWidth(content) + 4 // Add padding for borders
	topBorder := "+" + strings.Repeat("-", lineWidth-2) + "+"
	bottomBorder := "+" + strings.Repeat("-", lineWidth-2) + "+"
	leftBorder := "| "
	rightBorder := " |"

	block := topBorder + "\n"
	for _, line := range strings.Split(content, "\n") {
		paddedLine := fmt.Sprintf("%-*s", lineWidth-4, line) // Ensure lines fit inside the box
		block += leftBorder + paddedLine + rightBorder + "\n"
	}

	return block + bottomBorder
}

func PrintCodeBlockV2(content string) string {
	// Determine the maximum width of the content
	lineWidth := findMaxWidth(content) + 4 // Add padding for borders
	topBorder := "┌" + strings.Repeat("─", lineWidth-2) + "┐"
	bottomBorder := "└" + strings.Repeat("─", lineWidth-2) + "┘"
	leftBorder := "│ "
	rightBorder := " │"

	block := topBorder + "\n"
	for _, line := range strings.Split(content, "\n") {
		paddedLine := fmt.Sprintf("%-*s", lineWidth-4, line) // Ensure lines fit inside the box
		block += leftBorder + paddedLine + rightBorder + "\n"
	}

	return block + bottomBorder
}

// findMaxWidth calculates the longest line length in the given content
func findMaxWidth(content string) int {
	maxWidth := 0
	for _, line := range strings.Split(content, "\n") {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}
	return maxWidth
}
