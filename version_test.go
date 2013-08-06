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
		// Plain error cases
		{``, Version{}, `empty`},
		{`2`, Version{}, `form`},
		{`a.2`, Version{}, `form`},
		{`a.2.5`, Version{}, `form`},

		// Regex sanity
		{`04.2.1`, Version{}, `form`},
		{`4.2.01`, Version{}, `form`},
		{`!==4.2.1`, Version{}, `form`},
		{`!=>=4.2.1`, Version{}, `form`},

		// Nice cases
		{`2.1.3`, Version{0, 2, 1, 3, ``}, ``},
		{`4.2.1`, Version{Equal, 4, 2, 1, ``}, ``},

		// Operators
		{`=4.2.1`, Version{Equal, 4, 2, 1, ``}, ``},
		{`!=4.2.1`, Version{NotEqual, 4, 2, 1, ``}, ``},
		{`>4.2.1`, Version{GreaterThan, 4, 2, 1, ``}, ``},
		{`<4.2.1`, Version{LessThan, 4, 2, 1, ``}, ``},
		{`>=4.2.1`, Version{GreaterEqual, 4, 2, 1, ``}, ``},
		{`<=4.2.1`, Version{LessEqual, 4, 2, 1, ``}, ``},
		{`~4.2.1`, Version{ApproxGreater, 4, 2, 1, ``}, ``},

		// Release
		{`~4.2.1-.pre`, Version{}, `form`},
		{`~4.2.1-`, Version{}, `form`},
		{`~4.2.1-=`, Version{}, `form`},
		{`~4.2.1-1pre`, Version{}, `form`},
		{`~4.2.1-01`, Version{}, `form`},
		{`>=4.2.1-pre`, Version{GreaterEqual, 4, 2, 1, `pre`}, ``},
		{`<=4.2.1-pre1`, Version{LessEqual, 4, 2, 1, `pre1`}, ``},
		{`<=4.2.1-pre.1`, Version{LessEqual, 4, 2, 1, `pre.1`}, ``},
		{`<=4.2.1-pre.1.alpha`, Version{LessEqual, 4, 2, 1, `pre.1.alpha`}, ``},
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
		Base      string
		Condition string
		Result    bool
	}{
		// Base case
		{"1.0.0", "1.0.0", true},

		// =
		{"1.0.0-a", "=1.0.0-a", true},
		{"1.0.0", "=1.0.1", false},
		{"1.0.0", "=1.1.0", false},
		{"1.0.0", "=2.0.0", false},
		{"1.0.0-a", "=1.0.0", false},

		// !=
		{"1.0.0-a", "!=1.0.0-a", false},
		{"1.0.0", "!=1.0.0", false},
		{"1.0.0", "!=1.0.1", true},
		{"1.0.0", "!=1.1.0", true},
		{"1.0.0", "!=2.0.0", true},
		{"1.0.0-a", "!=1.0.0", true},

		// >
		{"1.0.1", ">1.0.0", true},
		{"1.1.0", ">1.0.0", true},
		{"2.0.0", ">1.0.0", true},
		{"1.0.0", ">1.0.0", false},
		{"1.0.0", ">1.0.1", false},
		{"1.0.0", ">1.1.0", false},
		{"1.0.0", ">2.0.0", false},
		{"1.0.0-a", ">1.0.0-a", false},
		{"1.0.0-a", ">1.0.1-a", false},
		{"1.0.0-a", ">1.0.0", false},
		{"1.0.0", ">1.0.0-a", true},
		{"1.0.0-a", ">2.0.0", false},
		{"2.0.0", ">1.0.0-a", true},

		// <
		{"1.0.1", "<1.0.0", false},
		{"1.1.0", "<1.0.0", false},
		{"2.0.0", "<1.0.0", false},
		{"1.0.0", "<1.0.0", false},
		{"1.0.0", "<1.0.1", true},
		{"1.0.0", "<1.1.0", true},
		{"1.0.0", "<2.0.0", true},
		{"1.0.0-a", "<1.0.0-a", false},
		{"1.0.0-a", "<1.0.1-a", true},
		{"1.0.0-a", "<1.0.0", true},
		{"1.0.0", "<1.0.0-a", false},
		{"1.0.0-a", "<2.0.0", true},
		{"2.0.0", "<1.0.0-a", false},

		// >=
		{"1.0.1", ">=1.0.0", true},
		{"1.1.0", ">=1.0.0", true},
		{"2.0.0", ">=1.0.0", true},
		{"1.0.0", ">=1.0.0", true},
		{"1.0.0", ">=1.0.1", false},
		{"1.0.0", ">=1.1.0", false},
		{"1.0.0", ">=2.0.0", false},
		{"1.0.0-a", ">=1.0.0-a", true},
		{"1.0.0-a", ">=1.0.1-a", false},
		{"1.0.0-a", ">=1.0.0", false},
		{"1.0.0", ">=1.0.0-a", true},
		{"1.0.0-a", ">=2.0.0", false},
		{"2.0.0", ">=1.0.0-a", true},

		// <=
		{"1.0.1", "<=1.0.0", false},
		{"1.1.0", "<=1.0.0", false},
		{"2.0.0", "<=1.0.0", false},
		{"1.0.0", "<=1.0.0", true},
		{"1.0.0", "<=1.0.1", true},
		{"1.0.0", "<=1.1.0", true},
		{"1.0.0", "<=2.0.0", true},
		{"1.0.0-a", "<=1.0.0-a", true},
		{"1.0.0-a", "<=1.0.1-a", true},
		{"1.0.0-a", "<=1.0.0", true},
		{"1.0.0", "<=1.0.0-a", false},
		{"1.0.0-a", "<=2.0.0", true},
		{"2.0.0", "<=1.0.0-a", false},

		// ~
		{"1.0.1", "~1.0.0", true},
		{"1.1.0", "~1.0.0", true},
		{"2.0.0", "~1.0.0", false},
		{"1.0.0", "~1.0.0", true},
		{"1.0.0", "~1.0.1", false},
		{"1.0.0", "~1.1.0", false},
		{"1.0.0", "~2.0.0", false},
		{"1.0.0-a", "~1.0.0-a", true},
		{"1.0.0-a", "~1.0.1-a", false},
		{"1.0.0-a", "~1.0.0", false},
		{"1.0.0", "~1.0.0-a", true},
		{"1.0.0-a", "~2.0.0", false},
		{"2.0.0", "~1.0.0-a", false},
		{"1.0.0-a", "~2.0.0-a", false},
		{"2.0.0-a", "~1.0.0-a", false},
	}

	for _, test := range tests {
		base, err := ParseVersion(test.Base)
		if err != nil {
			t.Error("Error parsing base version:", err)
		}
		cond, err := ParseVersion(test.Condition)
		if err != nil {
			t.Error("Error parsing base version:", err)
		}

		res := base.Compare(cond)

		if res != test.Result {
			t.Errorf("%v %v || expected: %v got: %v", test.Base, test.Condition,
				test.Result, res)
		}
	}
}

func TestCompareReleases(t *T) {
	t.Parallel()
	var tests = []struct {
		Base    string
		Compare string
		Result  int
	}{
		{``, ``, 0},
		{``, `a`, 1},
		{`a`, ``, -1},
		{`a`, `a`, 0},
		{`1`, `a`, 1},
		{`a`, `1`, -1},
		{`a`, `a.b`, 1},
		{`a.b`, `a`, -1},
		{`a1`, `a2`, 1},
		{`a2`, `a1`, -1},
		{`ab`, `abc`, 1},
		{`abc`, `ab`, -1},
		{`a.1`, `a.2`, 1},
		{`a.2`, `a.1`, -1},
		{`1.a`, `2.a`, 1},
		{`2.a`, `1.a`, -1},
	}

	for _, test := range tests {
		res := compareReleases(test.Base, test.Compare)

		if res != test.Result {
			t.Error(test, "expected:", res, "to be equal to:", test.Result)
		}
	}
}

func TestCompareStrings(t *T) {
	t.Parallel()
	var tests = []struct {
		Base    string
		Compare string
		Result  int
	}{
		{``, ``, 0},
		{``, `a`, 1},
		{`a`, ``, -1},
		{`a`, `ab`, 1},
		{`ab`, `a`, -1},
		{`ab`, `ab`, 0},
	}

	for _, test := range tests {
		res := compareStrings(test.Base, test.Compare)

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
		{Version{Equal, 1, 2, 3, ``}, `1.2.3`},
		{Version{NotEqual, 1, 2, 3, ``}, `!=1.2.3`},
		{Version{GreaterThan, 1, 2, 3, ``}, `>1.2.3`},
		{Version{LessThan, 1, 2, 3, ``}, `<1.2.3`},
		{Version{GreaterEqual, 1, 2, 3, ``}, `>=1.2.3`},
		{Version{LessEqual, 1, 2, 3, ``}, `<=1.2.3`},
		{Version{ApproxGreater, 1, 2, 3, ``}, `~1.2.3`},
		{Version{ApproxGreater, 1, 2, 3, `1.3.patch`}, `~1.2.3-1.3.patch`},
	}

	for _, test := range tests {
		if s := test.Version.String(); s != test.Expected {
			t.Error(test, "expected:", s, "to be equal to:", test.Expected)
		}
	}
}

func TestVersion_GetYAML(t *T) {
	t.Parallel()
	v := Version{NotEqual, 1, 2, 3, ``}
	_, value := v.GetYAML()
	if s, ok := value.(string); !ok {
		t.Error("It should return a string type.")
	} else if s != "!=1.2.3" {
		t.Error("It's not returning the correct string.")
	}
}

func TestVersion_SetYAML(t *T) {
	t.Parallel()
	var v Version
	if v.SetYAML("", "fail") {
		t.Error("Expecting failure.")
	}
	if v.SetYAML("", 10) {
		t.Error("Expecting failure.")
	}
	success := v.SetYAML("", ">=1.2.3-pre")
	if !success {
		t.Error("Expecting success.")
	}
	comp := Version{Equal, 1, 2, 3, `pre`}
	if !v.Compare(comp) {
		t.Error("Expected:", v, "to match", comp)
	}
}
