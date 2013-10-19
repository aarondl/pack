package pack

import (
	"strings"
	. "testing"
)

func TestParseDependency(t *T) {
	t.Parallel()

	out, err := ParseDependency(`name`)
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if out.Name != "name" {
		t.Error("Expected name but got:", out.Name)
	}

	out, err = ParseDependency(`name <3.2.5 ~4.2.5`)
	if err != nil {
		t.Fatal("Unexpected error:", err)
	} else if ln := len(out.Constraints); ln != 2 {
		t.Error("Expected 2 constraints, got:", ln)
	}

	if out.Name != "name" {
		t.Error("Expected name to be name but got:", out.Name)
	}

	if op := out.Constraints[0].Operator; op != LessThan {
		t.Error("Expected less than operator, got:", op)
	}

	v, _ := ParseVersion("3.2.5")
	if version := out.Constraints[0].Version; !version.Satisfies(Equal, v) {
		t.Errorf("Expected %v and %v to be equal.", v, version)
	}

	if op := out.Constraints[1].Operator; op != ApproxGreater {
		t.Error("Expected less than operator, got:", op)
	}

	v, _ = ParseVersion("4.2.5")
	if version := out.Constraints[1].Version; !version.Satisfies(Equal, v) {
		t.Errorf("Expected %v and %v to be equal.", v, version)
	}

	out, err = ParseDependency(`github.com/test/test git`)
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	if out.Name != "github.com/test/test" {
		t.Error("Expected name to be name but got:", out.Name)
	}

	if out.URL != "git" {
		t.Error("Expected the repo style but got:", out.URL)
	}

	out, err = ParseDependency(`name git:http://repo.com`)
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	if out.Name != "name" {
		t.Error("Expected name to be name but got:", out.Name)
	}

	if out.URL != "git:http://repo.com" {
		t.Error("Expected the repo url but got:", out.URL)
	}
}

func TestParseDependency_Errors(t *T) {
	t.Parallel()

	_, err := ParseDependency(``)
	if !strings.Contains(err.Error(), `must be in the form`) {
		t.Error("Expected empty error, got:", err)
	}

	_, err = ParseDependency(`name asdf`)
	if err == nil {
		t.Error("Expected an error but got nothing.")
	} else if exp := "form"; !strings.Contains(err.Error(), exp) {
		t.Error("Expected an error matching:", exp, "but got:", err)
	}
}

func TestDependency_String(t *T) {
	t.Parallel()

	var dep Dependency

	if s := dep.String(); s != `` {
		t.Error("Expected empty string, got:", s)
	}

	dep.Constraints = make([]*Constraint, 2)
	dep.Constraints[0] = &Constraint{LessThan, &Version{1, 2, 3, "pre"}}
	dep.Constraints[1] = &Constraint{ApproxGreater, &Version{3, 2, 1, "dev"}}

	if s := dep.String(); s != `` {
		t.Error("Expected empty string, got:", s)
	}

	dep.URL = "hg:hg.io"

	if s := dep.String(); s != `` {
		t.Error("Expected empty string, got:", s)
	}

	dep.Name = "name"

	if s, exp := dep.String(), `name <1.2.3-pre ~3.2.1-dev hg:hg.io`; s != exp {
		t.Error("Expected:", exp, "got:", s)
	}
}

func TestDependency_GetYAML(t *T) {
	t.Parallel()
	d := Dependency{
		"name",
		[]*Constraint{{
			NotEqual,
			&Version{1, 2, 3, `pre`},
		}},
		"git:git+https://repo.com/?hi",
	}
	_, value := d.GetYAML()
	if s, ok := value.(string); !ok {
		t.Error("It should return a string type.")
	} else if s != "name !=1.2.3-pre git:git+https://repo.com/?hi" {
		t.Error("It's not returning the correct string.")
	}
}

func TestDependency_SetYAML(t *T) {
	t.Parallel()
	var d Dependency
	if d.SetYAML("", 10) {
		t.Error("Expecting failure.")
	}
	success := d.SetYAML("", "name >=1.2.3-pre git:git.com")
	if !success {
		t.Error("Expecting success.")
	}
	if d.Name != "name" {
		t.Error("Expected:", d.Name, "to equal: name")
	}
	if exp := "git:git.com"; exp != d.URL {
		t.Error("Expected:", d.URL, "to equal:", exp)
	}
	comp := &Version{1, 2, 3, `pre`}
	if len(d.Constraints) != 1 {
		t.Error("Expected a single constraint.")
	} else if c := d.Constraints[0]; c.Operator != GreaterEqual {
		t.Error("Expected >= operator, got:", c.Operator.String())
	} else if !c.Version.Satisfies(Equal, comp) {
		t.Error("Expected:", c.Version, "to match", comp)
	}
}
