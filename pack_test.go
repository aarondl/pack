package pack

import (
	"strings"
	. "testing"
)

func TestPack(t *T) {
	t.Parallel()
}

func TestParseDependency(t *T) {
	t.Parallel()
	var tests = []struct {
		Input  string
		Output Dependency
		Error  string
	}{
		{``, Dependency{}, `empty`},
		{`name <3.2.5`, Dependency{`name`, &Version{LessThan, 3, 2, 5}}, ``},
		{`1234 a.2.5`, Dependency{}, `name`},
		{`name1234 a.2`, Dependency{}, `form`},
	}

	for _, test := range tests {
		out, err := ParseDependency(test.Input)
		if err != nil {
			if len(test.Error) == 0 {
				t.Error(test)
				t.Error(test, "had unexpected error:", err)
			} else if !strings.Contains(err.Error(), test.Error) {
				t.Error(test, "expected error message like:",
					test.Error, "got:", err)
			}
			continue
		}

		if out.Name != test.Output.Name {
			t.Error("Expected:", out.Name, "to be equal to:", test.Output.Name)
		}
		if out.Version != nil && *out.Version != *test.Output.Version {
			t.Error("Expected:", *out.Version, "to be equal to:",
				*test.Output.Version)
		}
	}
}

func TestDependency_String(t *T) {
	t.Parallel()
	var tests = []struct {
		Dependency Dependency
		Expected   string
	}{
		{Dependency{}, ""},
		{Dependency{"", &Version{0, 1, 2, 3}}, ""},
		{Dependency{"name", nil}, "name"},
		{Dependency{"name", &Version{GreaterThan, 1, 2, 3}}, "name >1.2.3"},
	}

	for _, test := range tests {
		if s := test.Dependency.String(); s != test.Expected {
			t.Error(test, "expected:", s, "to be equal to:", test.Expected)
		}
	}
}

func TestDependency_GetYAML(t *T) {
	t.Parallel()
	d := Dependency{"name", &Version{NotEqual, 1, 2, 3}}
	_, value := d.GetYAML()
	if s, ok := value.(string); !ok {
		t.Error("It should return a string type.")
	} else if s != "name !=1.2.3" {
		t.Error("It's not returning the correct string.")
	}
}

func TestDependency_SetYAML(t *T) {
	t.Parallel()
	var d Dependency
	if d.SetYAML("", "1234fail") {
		t.Error("Expecting failure.")
	}
	if d.SetYAML("", 10) {
		t.Error("Expecting failure.")
	}
	success := d.SetYAML("", "name >=1.2.3")
	if !success {
		t.Error("Expecting success.")
	}
	if d.Name != "name" {
		t.Error("Expected:", d.Name, "to equal: name")
	}
	comp := Version{Equal, 1, 2, 3}
	if !d.Version.Compare(comp) {
		t.Error("Expected:", d.Version, "to match", comp)
	}
}
