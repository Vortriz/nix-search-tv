package optionsfile

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestUnmarshal(t *testing.T) {
	data, err := os.ReadFile("./testdata/options.json")
	assert.NoError(t, err)

	opts := map[string]json.RawMessage{}
	err = json.Unmarshal(data, &opts)
	assert.NoError(t, err)

	cases := []struct {
		Name   string
		Input  string
		Output Package
	}{
		{
			Name:  "declarations as a slice of strings",
			Input: "age.ageBin",
			Output: Package{
				Type:        "string",
				Description: "The age executable to use.\n",
				Declarations: []String{
					"/nix/store/azcvzd254j7wy5fmlbws96a5zfckjw9d-source/modules/age.nix",
				},
				Default: String("\"${pkgs.age}/bin/age\"\n"),
			},
		},
		{
			Name:  "declarations as a slice of objects",
			Input: "nixvim.autoCmd",
			Output: Package{
				Example: String(
					"[\n  {\n    command = \"echo 'Entering a C or C++ file'\";\n    event = [\n      \"BufEnter\"\n      \"BufWinEnter\"\n    ];\n    pattern = [\n      \"*.c\"\n      \"*.h\"\n    ];\n  }\n]",
				),
				Type:        "list of (submodule)",
				Description: "autocmd definitions",
				Declarations: []String{
					"https://github.com/nix-community/nixvim/blob/main/modules/autocmd.nix",
				},
				Default: String("[ ]"),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			actual, ok := opts[c.Input]
			assert.True(t, ok)

			pkg := Package{}
			err := json.Unmarshal(actual, &pkg)
			assert.NoError(t, err)
			assert.True(t, reflect.DeepEqual(pkg, c.Output), "packages are not equal. \nexpected: \n%#v, \ngot: \n%#v", c.Output, pkg)
		})
	}
}
