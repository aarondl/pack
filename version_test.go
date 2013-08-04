package pack

import (
	"strings"
	. "testing"
)

func TestParse(t *T) {
	t.Parallel()
	var tests = []struct {
		Input  string
		Output Version
		Error  string
	}{
		{``, Version{}, `empty`},
		{`2`, Version{}, `form`},
		{`a.2`, Version{}, `form`},
		{`a.2.5`, Version{}, `form`},
		{`!==4.2.1`, Version{}, `form`},
		{`!=>=4.2.1`, Version{}, `form`},
		{`2.1.3`, Version{0, 2, 1, 3}, ``},
		{`4.2.1`, Version{Equal, 4, 2, 1}, ``},
		{`=4.2.1`, Version{Equal, 4, 2, 1}, ``},
		{`!=4.2.1`, Version{NotEqual, 4, 2, 1}, ``},
		{`>4.2.1`, Version{GreaterThan, 4, 2, 1}, ``},
		{`<4.2.1`, Version{LessThan, 4, 2, 1}, ``},
		{`>=4.2.1`, Version{GreaterEqual, 4, 2, 1}, ``},
		{`<=4.2.1`, Version{LessEqual, 4, 2, 1}, ``},
		{`~>4.2.1`, Version{ApproxGreater, 4, 2, 1}, ``},
	}

	for _, test := range tests {
		out, err := ParseVersion(test.Input)
		if err != nil {
			if len(test.Error) == 0 {
				t.Error(test)
				t.Error(test, "had unexpected error:", err)
			} else if !strings.Contains(err.Error(), test.Error) {
				t.Error(test, "expected error message like:",
					test.Error, "got:", err)
			}
		}

		if out != test.Output {
			t.Error("Expected:", out, "to be equal to:", test.Output)
		}
	}
}

func TestCompare(t *T) {
	t.Parallel()
	var tests = []struct {
		Base      Version
		Condition Version
		Result    bool
	}{
		// Base case
		{Version{0, 1, 0, 0}, Version{0, 1, 0, 0}, true},

		// =
		{Version{0, 1, 0, 0}, Version{Equal, 1, 0, 0}, true},
		{Version{0, 1, 0, 0}, Version{Equal, 1, 0, 1}, false},
		{Version{0, 1, 0, 0}, Version{Equal, 1, 1, 0}, false},
		{Version{0, 1, 0, 0}, Version{Equal, 2, 0, 0}, false},

		// !=
		{Version{0, 1, 0, 0}, Version{NotEqual, 1, 0, 0}, false},
		{Version{0, 1, 0, 0}, Version{NotEqual, 1, 0, 1}, true},
		{Version{0, 1, 0, 0}, Version{NotEqual, 1, 1, 0}, true},
		{Version{0, 1, 0, 0}, Version{NotEqual, 2, 0, 0}, true},

		// >
		{Version{0, 1, 0, 1}, Version{GreaterThan, 1, 0, 0}, true},
		{Version{0, 1, 1, 0}, Version{GreaterThan, 1, 0, 0}, true},
		{Version{0, 2, 0, 0}, Version{GreaterThan, 1, 0, 0}, true},
		{Version{0, 1, 0, 0}, Version{GreaterThan, 1, 0, 0}, false},
		{Version{0, 1, 0, 0}, Version{GreaterThan, 1, 0, 1}, false},
		{Version{0, 1, 0, 0}, Version{GreaterThan, 1, 1, 0}, false},
		{Version{0, 1, 0, 0}, Version{GreaterThan, 2, 0, 0}, false},

		// <
		{Version{0, 1, 0, 1}, Version{LessThan, 1, 0, 0}, false},
		{Version{0, 1, 1, 0}, Version{LessThan, 1, 0, 0}, false},
		{Version{0, 2, 0, 0}, Version{LessThan, 1, 0, 0}, false},
		{Version{0, 1, 0, 0}, Version{LessThan, 1, 0, 0}, false},
		{Version{0, 1, 0, 0}, Version{LessThan, 1, 0, 1}, true},
		{Version{0, 1, 0, 0}, Version{LessThan, 1, 1, 0}, true},
		{Version{0, 1, 0, 0}, Version{LessThan, 2, 0, 0}, true},

		// >=
		{Version{0, 1, 0, 1}, Version{GreaterEqual, 1, 0, 0}, true},
		{Version{0, 1, 1, 0}, Version{GreaterEqual, 1, 0, 0}, true},
		{Version{0, 2, 0, 0}, Version{GreaterEqual, 1, 0, 0}, true},
		{Version{0, 1, 0, 0}, Version{GreaterEqual, 1, 0, 0}, true},
		{Version{0, 1, 0, 0}, Version{GreaterEqual, 1, 0, 1}, false},
		{Version{0, 1, 0, 0}, Version{GreaterEqual, 1, 1, 0}, false},
		{Version{0, 1, 0, 0}, Version{GreaterEqual, 2, 0, 0}, false},

		// <=
		{Version{0, 1, 0, 1}, Version{LessEqual, 1, 0, 0}, false},
		{Version{0, 1, 1, 0}, Version{LessEqual, 1, 0, 0}, false},
		{Version{0, 2, 0, 0}, Version{LessEqual, 1, 0, 0}, false},
		{Version{0, 1, 0, 0}, Version{LessEqual, 1, 0, 0}, true},
		{Version{0, 1, 0, 0}, Version{LessEqual, 1, 0, 1}, true},
		{Version{0, 1, 0, 0}, Version{LessEqual, 1, 1, 0}, true},
		{Version{0, 1, 0, 0}, Version{LessEqual, 2, 0, 0}, true},

		// ~>
		{Version{0, 1, 0, 1}, Version{ApproxGreater, 1, 0, 0}, true},
		{Version{0, 1, 1, 0}, Version{ApproxGreater, 1, 0, 0}, true},
		{Version{0, 2, 0, 0}, Version{ApproxGreater, 1, 0, 0}, false},
		{Version{0, 1, 0, 0}, Version{ApproxGreater, 1, 0, 0}, true},
		{Version{0, 1, 0, 0}, Version{ApproxGreater, 1, 0, 1}, false},
		{Version{0, 1, 0, 0}, Version{ApproxGreater, 1, 1, 0}, false},
		{Version{0, 1, 0, 0}, Version{ApproxGreater, 2, 0, 0}, false},
	}

	for _, test := range tests {
		res := test.Base.Compare(test.Condition)

		if res != test.Result {
			t.Error(test, "expected:", res, "to be equal to:", test.Result)
		}
	}
}

func TestVersion_String(t *T) {
	t.Parallel()
	var tests = []struct {
		Version  Version
		Expected string
	}{
		{Version{Equal, 1, 2, 3}, `1.2.3`},
		{Version{NotEqual, 1, 2, 3}, `!=1.2.3`},
		{Version{GreaterThan, 1, 2, 3}, `>1.2.3`},
		{Version{LessThan, 1, 2, 3}, `<1.2.3`},
		{Version{GreaterEqual, 1, 2, 3}, `>=1.2.3`},
		{Version{LessEqual, 1, 2, 3}, `<=1.2.3`},
		{Version{ApproxGreater, 1, 2, 3}, `~>1.2.3`},
	}

	for _, test := range tests {
		if s := test.Version.String(); s != test.Expected {
			t.Error(test, "expected:", s, "to be equal to:", test.Expected)
		}
	}
}

func TestGetYAML(t *T) {
	t.Parallel()
	v := Version{NotEqual, 1, 2, 3}
	_, value := v.GetYAML()
	if s, ok := value.(string); !ok {
		t.Error("It should return a string type.")
	} else if s != "!=1.2.3" {
		t.Error("It's not returning the correct string.")
	}
}

func TestSetYAML(t *T) {
	t.Parallel()
	var v Version
	if v.SetYAML("", "fail") {
		t.Error("Expecting failure.")
	}
	if v.SetYAML("", 10) {
		t.Error("Expecting failure.")
	}
	success := v.SetYAML("", ">=1.2.3")
	if !success {
		t.Error("Expecting success.")
	}
	comp := Version{Equal, 1, 2, 3}
	if !v.Compare(comp) {
		t.Error("Expected:", v, "to match", comp)
	}
}
