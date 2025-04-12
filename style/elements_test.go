package style

import (
	"strings"
	"testing"

	"github.com/alecthomas/assert/v2"
)

var s = StyledText

func TestStyleCallouts(t *testing.T) {
	warn := s.Red("> ")
	note := s.Bold("| ")

	cases := []struct {
		Desc     string
		Input    string
		Expected []string
	}{
		{
			"warning callout",
			"::: {.warning}\n  warning\n:::",
			[]string{
				warn,
				warn + "  warning",
				warn,
			},
		},
		{
			"note callout",
			"::: {.note}\n  note\n:::",
			[]string{
				note,
				note + "  note",
				note,
			},
		},
		{
			"unknown callout",
			"::: {.unknown}\n  unknown\n:::",
			[]string{
				note,
				note + "  unknown",
				note,
			},
		},
		{
			":: at the end",
			"::: {.warning}\n  warning\n::",
			[]string{
				warn,
				warn + "  warning",
				warn,
			},
		},
		{
			"multiple",
			"::: {.warning}\n warning\n::\n::: {.note}\n note\n:::",
			[]string{
				warn,
				warn + " warning",
				warn,
				note,
				note + " note",
				note,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			actual := styleCallouts(c.Input)
			expected := strings.Join(c.Expected, "\n")
			assert.Equal(t, expected, actual)
		})
	}
}

// func TestStyleLongDescription0(t *testing.T) {
// 	// 2 :: in the end
// 	// nested {option}
// 	input := "Username for login.\n\n::: {.warning}\n  This option takes precedence over {option}`services.stash.settings.username`\n::\n\n"

// 	styler := StyledText
// 	output := StyleLongDescription(styler, input)

// 	// expected := "Username for login.\n\n::: {.warning}\n  This option takes precedence over {option}`services.stash.settings.username`\n::\n\n"
// 	fmt.Println(output)
// }

// func TestStyleLongDescription2(t *testing.T) {
// 	input := "When using the SLiRP user networking (default), this option allows to\nforward ports to/from the host/guest.\n\n::: {.warning}\nIf the NixOS firewall on the virtual machine is enabled, you also\nhave to open the guest ports to enable the traffic between host and\nguest.\n:::\n\n::: {.note}\nCurrently QEMU supports only IPv4 forwarding.\n:::\n"
// 	// input := "When using the SLiRP user networking (default), this option allows to\nforward ports to/from the host/guest.\n\n::: {.warning}\nIf the NixOS firewall on the virtual machine is enabled, you also\nhave to open the guest ports to enable the traffic between host and\nguest.\n:::\n\n"

// 	styler := StyledText
// 	output := StyleLongDescription(styler, input)

// 	// expected := "Username for login.\n\n::: {.warning}\n  This option takes precedence over {option}`services.stash.settings.username`\n::\n\n"
// 	fmt.Println(output)
// }

// func TestStyleLongDescription3(t *testing.T) {
// 	input := "When using the SLiRP user networking (default), this option allows to\nforward ports to/from the host/guest.\n\n::: {.warning}\nIf the NixOS firewall on the virtual machine is enabled, you also\nhave to open the guest ports to enable the traffic between host and\nguest.\n:::\n\n::: {.note}\nCurrently QEMU supports only IPv4 forwarding.\n:::\n"

// 	styler := StyledText
// 	output := StyleLongDescription(styler, input)

// 	// expected := "Username for login.\n\n::: {.warning}\n  This option takes precedence over {option}`services.stash.settings.username`\n::\n\n"
// 	fmt.Println(output)
// }
