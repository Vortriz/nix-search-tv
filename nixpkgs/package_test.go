package nixpkgs

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestUnmarshalElem(t *testing.T) {
	input := `
	{
	  "Elem": 1
	}
	`
	es := struct {
		Elem ElemOrSlice[int]
	}{}

	err := json.Unmarshal([]byte(input), &es)
	if err != nil {
		t.Fatal(err)
	}
	if es.Elem[0] != 1 {
		t.Fatalf("expected 1, but got %v", es.Elem[0])
	}
}

func TestUnmarshalSlice(t *testing.T) {
	input := `
	{
	  "Slice": [1, 2]
	}
	`
	es := struct {
		Slice ElemOrSlice[int]
	}{}

	err := json.Unmarshal([]byte(input), &es)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual([]int(es.Slice), []int{1, 2}) {
		t.Fatalf("expected [1 2], but got %v", es.Slice)
	}
}
