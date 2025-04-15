package renderdocs

import (
	"github.com/JohannesKaufmann/dom"
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	"golang.org/x/net/html"
)

func RenderHTML(text string) string {
	conv := converter.NewConverter(
		converter.WithEscapeMode(converter.EscapeModeDisabled),
		converter.WithPlugins(
			base.NewBasePlugin(),
			commonmark.NewCommonmarkPlugin(),
		),
	)

	conv.Register.RendererFor(
		"a", converter.TagTypeInline,
		func(ctx converter.Context, w converter.Writer, n *html.Node) converter.RenderStatus {
			href, _ := dom.GetAttribute(n, "href")
			text := dom.CollectText(n)
			if href == text {
				w.WriteString("<" + href + ">")
				return converter.RenderSuccess
			}

			if dom.HasClass(n, "xref") {
				w.WriteString("`opt#" + text + "`")
				return converter.RenderSuccess
			}

			child := dom.FirstChildNode(n)
			if child != nil && child.Data == "span" {
				if dom.HasClass(child, "citerefentry") {
					w.WriteString("`" + text + "`")
					return converter.RenderSuccess
				}
			}

			return converter.RenderTryNext
		},
		converter.PriorityEarly,
	)

	conv.Register.RendererFor(
		"div", converter.TagTypeInline,
		func(ctx converter.Context, w converter.Writer, n *html.Node) converter.RenderStatus {
			for _, class := range dom.GetClasses(n) {
				switch class {
				default:
					continue

				case "important", "warning", "caution", "note":
				}

				if n.FirstChild != nil && n.FirstChild.Data == "h3" {
					n.RemoveChild(n.FirstChild)
				}
				w.WriteString("::: {." + class + "}\n")
				w.WriteString(dom.CollectText(n))
				w.WriteString("\n:::")
				return converter.RenderSuccess
			}

			return converter.RenderTryNext
		},
		converter.PriorityEarly,
	)

	text, _ = conv.ConvertString(text)
	return text
}
