package pack

import (
	"strings"
	. "testing"
)

func TestParseDependency(t *T) {
	t.Parallel()

	out, err := ParseDependency(``)
	if !strings.Contains(err.Error(), `empty`) {
		t.Error("Expected empty error, got:", err)
	}

	out, err = ParseDependency(`1234 <3.2.5`)
	if !strings.Contains(err.Error(), `name`) {
		t.Error("Expected name error, got:", err)
	}

	out, err = ParseDependency(`name ~=3.2.5`)
	if !strings.Contains(err.Error(), `constraints must have`) {
		t.Error("Expected constraints error, got:", err)
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

	dep.Name = "name"

	if s, exp := dep.String(), `name <1.2.3-pre ~3.2.1-dev`; s != exp {
		t.Error("Expected:", exp, "got:", s)
	}
}

func TestDependency_GetYAML(t *T) {
	t.Parallel()
	d := Dependency{"name",
		[]*Constraint{{
			NotEqual,
			&Version{1, 2, 3, `pre`},
		}},
	}
	_, value := d.GetYAML()
	if s, ok := value.(string); !ok {
		t.Error("It should return a string type.")
	} else if s != "name !=1.2.3-pre" {
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
	success := d.SetYAML("", "name >=1.2.3-pre")
	if !success {
		t.Error("Expecting success.")
	}
	if d.Name != "name" {
		t.Error("Expected:", d.Name, "to equal: name")
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
