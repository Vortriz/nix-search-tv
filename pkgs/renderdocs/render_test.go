package renderdocs

import (
	"strings"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestRenderHTML(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    string
		Expected string
	}{
		{
			Name: "Href equals text",
			Input: `
			<a
				href="https://vdirsyncer.pimutils.org/en/stable/config.html#google"
			>
				https://vdirsyncer.pimutils.org/en/stable/config.html#google
			</a>`,
			Expected: "<https://vdirsyncer.pimutils.org/en/stable/config.html#google>",
		},
		{
			Name: "Xref",
			Input: `
			<a
				class="xref"
				href="options.xhtml#opt-accounts.email.accounts._name_.jmap.sessionUrl"
			>
				<code class="option">accounts.email.accounts.&lt;name&gt;.jmap.sessionUrl</code>
			</a>`,
			Expected: "`opt#accounts.email.accounts.<name>.jmap.sessionUrl`",
		},
		{
			Name: "Manpage Link",
			Input: `
			<a
				href="https://www.freedesktop.org/software/systemd/man/systemd.time.html"
			>
				<span class="citerefentry">
					<span class="refentrytitle">systemd.time</span>(7)
				</span>
			</a>`,
			Expected: "`systemd.time(7)`",
		},
		{
			Name: "Pandoc-like callout",
			Input: `
			<div
				class="important"
			>
				<h3 class="title">Important</h3>
				<p>The list must not be empty and must not contain duplicate entries (attrsets which compare equally).</p>
			</div>`,
			Expected: "::: {.important}\nThe list must not be empty and must not contain duplicate entries (attrsets which compare equally).\n:::",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			actual := RenderHTML(tc.Input)
			expected := strings.TrimSpace(tc.Expected)
			actual = strings.TrimSpace(actual)
			assert.Equal(t, expected, actual)
		})
	}
}
