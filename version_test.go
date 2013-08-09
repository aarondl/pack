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
		{`2.1.3`, Version{2, 1, 3, ``}, ``},
		{`4.2.1`, Version{4, 2, 1, ``}, ``},

		// Release
		{`4.2.1-.pre`, Version{}, `form`},
		{`4.2.1-`, Version{}, `form`},
		{`4.2.1-=`, Version{}, `form`},
		{`4.2.1-1pre`, Version{}, `form`},
		{`4.2.1-01`, Version{}, `form`},
		{`4.2.1-pre`, Version{4, 2, 1, `pre`}, ``},
		{`4.2.1-pre1`, Version{4, 2, 1, `pre1`}, ``},
		{`4.2.1-pre.1`, Version{4, 2, 1, `pre.1`}, ``},
		{`4.2.1-pre.1.alpha`, Version{4, 2, 1, `pre.1.alpha`}, ``},
	}

	for _, test := range tests {
		out, err := ParseVersion(test.Input)
		if err != nil {
			if len(test.Error) == 0 {
				t.Error(test, "had unexpected error:", err)
			} else if !strings.Contains(err.Error(), test.Error) {
				t.Error(test, "expected error message like:",
					test.Error, "got:", err)
			}
		} else if out != nil && *out != test.Output {
			t.Error("Output:", out, "to be equal to:", test.Output)
		}
	}
}

func TestSatisfies(t *T) {
	t.Parallel()
	var tests = []struct {
		Base      string
		Op        string
		Condition string
		Result    bool
	}{
		// Base case
		{"1.0.0", "=", "1.0.0", true},

		// =
		{"1.0.0-a", "=", "1.0.0-a", true},
		{"1.0.0", "=", "1.0.1", false},
		{"1.0.0", "=", "1.1.0", false},
		{"1.0.0", "=", "2.0.0", false},
		{"1.0.0-a", "=", "1.0.0", false},

		// !=
		{"1.0.0-a", "!=", "1.0.0-a", false},
		{"1.0.0", "!=", "1.0.0", false},
		{"1.0.0", "!=", "1.0.1", true},
		{"1.0.0", "!=", "1.1.0", true},
		{"1.0.0", "!=", "2.0.0", true},
		{"1.0.0-a", "!=", "1.0.0", true},

		// >
		{"1.0.1", ">", "1.0.0", true},
		{"1.1.0", ">", "1.0.0", true},
		{"2.0.0", ">", "1.0.0", true},
		{"1.0.0", ">", "1.0.0", false},
		{"1.0.0", ">", "1.0.1", false},
		{"1.0.0", ">", "1.1.0", false},
		{"1.0.0", ">", "2.0.0", false},
		{"1.0.0-a", ">", "1.0.0-a", false},
		{"1.0.0-a", ">", "1.0.1-a", false},
		{"1.0.0-a", ">", "1.0.0", false},
		{"1.0.0", ">", "1.0.0-a", true},
		{"1.0.0-a", ">", "2.0.0", false},
		{"2.0.0", ">", "1.0.0-a", true},

		// <
		{"1.0.1", "<", "1.0.0", false},
		{"1.1.0", "<", "1.0.0", false},
		{"2.0.0", "<", "1.0.0", false},
		{"1.0.0", "<", "1.0.0", false},
		{"1.0.0", "<", "1.0.1", true},
		{"1.0.0", "<", "1.1.0", true},
		{"1.0.0", "<", "2.0.0", true},
		{"1.0.0-a", "<", "1.0.0-a", false},
		{"1.0.0-a", "<", "1.0.1-a", true},
		{"1.0.0-a", "<", "1.0.0", true},
		{"1.0.0", "<", "1.0.0-a", false},
		{"1.0.0-a", "<", "2.0.0", true},
		{"2.0.0", "<", "1.0.0-a", false},

		// >=
		{"1.0.1", ">=", "1.0.0", true},
		{"1.1.0", ">=", "1.0.0", true},
		{"2.0.0", ">=", "1.0.0", true},
		{"1.0.0", ">=", "1.0.0", true},
		{"1.0.0", ">=", "1.0.1", false},
		{"1.0.0", ">=", "1.1.0", false},
		{"1.0.0", ">=", "2.0.0", false},
		{"1.0.0-a", ">=", "1.0.0-a", true},
		{"1.0.0-a", ">=", "1.0.1-a", false},
		{"1.0.0-a", ">=", "1.0.0", false},
		{"1.0.0", ">=", "1.0.0-a", true},
		{"1.0.0-a", ">=", "2.0.0", false},
		{"2.0.0", ">=", "1.0.0-a", true},

		// <=
		{"1.0.1", "<=", "1.0.0", false},
		{"1.1.0", "<=", "1.0.0", false},
		{"2.0.0", "<=", "1.0.0", false},
		{"1.0.0", "<=", "1.0.0", true},
		{"1.0.0", "<=", "1.0.1", true},
		{"1.0.0", "<=", "1.1.0", true},
		{"1.0.0", "<=", "2.0.0", true},
		{"1.0.0-a", "<=", "1.0.0-a", true},
		{"1.0.0-a", "<=", "1.0.1-a", true},
		{"1.0.0-a", "<=", "1.0.0", true},
		{"1.0.0", "<=", "1.0.0-a", false},
		{"1.0.0-a", "<=", "2.0.0", true},
		{"2.0.0", "<=", "1.0.0-a", false},

		// ~
		{"1.0.1", "~", "1.0.0", true},
		{"1.1.0", "~", "1.0.0", true},
		{"2.0.0", "~", "1.0.0", false},
		{"1.0.0", "~", "1.0.0", true},
		{"1.0.0", "~", "1.0.1", false},
		{"1.0.0", "~", "1.1.0", false},
		{"1.0.0", "~", "2.0.0", false},
		{"1.0.0-a", "~", "1.0.0-a", true},
		{"1.0.0-a", "~", "1.0.1-a", false},
		{"1.0.0-a", "~", "1.0.0", false},
		{"1.0.0", "~", "1.0.0-a", true},
		{"1.0.0-a", "~", "2.0.0", false},
		{"2.0.0", "~", "1.0.0-a", false},
		{"1.0.0-a", "~", "2.0.0-a", false},
		{"2.0.0-a", "~", "1.0.0-a", false},
	}

	for _, test := range tests {
		base, err := ParseVersion(test.Base)
		if err != nil {
			t.Error("Error parsing base version:", err)
		}
		op, err := ParseOp(test.Op)
		if err != nil {
			t.Error("Error parsing operator:", err)
		}
		cond, err := ParseVersion(test.Condition)
		if err != nil {
			t.Error("Error parsing base version:", err)
		}

		res := base.Satisfies(op, cond)

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
		Version Version
		Output  string
	}{
		{Version{0, 0, 0, ``}, `0.0.0`},
		{Version{1, 2, 3, ``}, `1.2.3`},
		{Version{1, 2, 3, `1.3.patch`}, `1.2.3-1.3.patch`},
	}

	for _, test := range tests {
		if s := test.Version.String(); s != test.Output {
			t.Error(test, "expected:", s, "to be equal to:", test.Output)
		}
	}
}

func TestVersion_GetYAML(t *T) {
	t.Parallel()
	v := Version{1, 2, 3, ``}
	_, value := v.GetYAML()
	if s, ok := value.(string); !ok {
		t.Error("It should return a string type.")
	} else if s != "1.2.3" {
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
	success := v.SetYAML("", "1.2.3-pre")
	if !success {
		t.Error("Expecting success.")
	}
	comp := &Version{1, 2, 3, `pre`}
	if !v.Satisfies(Equal, comp) {
		t.Error("Output:", v, "to match", comp)
	}
}

func TestCompareOp_Parse(t *T) {
	var tests = []struct {
		Input  string
		Output ComparisonOp
		Error  string
	}{
		{`0`, 0, `must be one of`},
		{`=`, Equal, ``},
		{`!=`, NotEqual, ``},
		{`>`, GreaterThan, ``},
		{`<`, LessThan, ``},
		{`>=`, GreaterEqual, ``},
		{`<=`, LessEqual, ``},
		{`~`, ApproxGreater, ``},
	}

	for _, test := range tests {
		op, err := ParseOp(test.Input)
		if err != nil {
			if len(test.Error) == 0 {
				t.Error(test, "had unexpected error:", err)
			} else if !strings.Contains(err.Error(), test.Error) {
				t.Error(test, "expected error message like:",
					test.Error, "got:", err)
			}
		}

		if op != test.Output {
			t.Errorf("Expected: %v got: %v", test.Output, op)
		}
	}
}

func TestCompareOp_String(t *T) {
	var tests = []struct {
		Input  ComparisonOp
		Output string
		Empty  bool
	}{
		{0, ``, true},
		{Equal, `=`, false},
		{NotEqual, `!=`, false},
		{GreaterThan, `>`, false},
		{LessThan, `<`, false},
		{GreaterEqual, `>=`, false},
		{LessEqual, `<=`, false},
		{ApproxGreater, `~`, false},
	}

	for _, test := range tests {
		s := test.Input.String()
		if test.Empty && len(s) > 0 {
			t.Error("Expected empty string but got:", s)
		} else if !test.Empty && s != test.Output {
			t.Error(test, "expected:", s, "to be equal to:", test.Output)
		}
	}
}
