package recurse

import (
	"testing"
)

func TestRecurse(t *testing.T) {
	// test the function
	recipes, _, _, err := Recipe("apple pie", []string{"butter", "salt", "pie crust", "cream", "unsalted butter"}, "../../static/data/")
	if err != nil {
		t.Errorf("got error: %s", err)
	}
	_ = recipes
}
