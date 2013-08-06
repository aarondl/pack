package pack

import (
	"bytes"
	"errors"
	"strings"
	. "testing"
)

var testPack = `name: package
importpath: github.com/user/package/import
version: 1.0.0
summary: Package summary
description: Package Description
homepage: www.package.com
repository:
  type: git
  url: github.com/user/package
license: mit
authors:
- name: Author1
  email: author1@email.com
  homepage: blog.author.com
- name: Author2
  email: author2@email.com
  homepage: blog.author.com
contributors:
- name: Contrib1
  email: contrib@email.com
  homepage: github.com/contrib
support:
  website: support.com
  email: email@support.com
  forum: forum.com
  wiki: wiki.com
  issues: github.com/issues
dependencies:
- dep >1.2.3
- dep2 ~1.4.5-pre !=1.5.0
subpackages:
- subpackage
`

type badIO struct {
}
type halfWrite struct {
}

var fakeError = errors.New("Fake error.")

func (f *badIO) Read(in []byte) (int, error) {
	return 0, fakeError
}
func (b *badIO) Write(out []byte) (int, error) {
	return 0, fakeError
}
func (b *halfWrite) Write(out []byte) (int, error) {
	return len(out) / 2, nil
}

func TestParsePack(t *T) {
	t.Parallel()

	buf := bytes.NewBufferString(testPack)
	p, err := ParsePack(buf)
	if err != nil {
		t.Error("Unexpected:", err)
	}
	if p == nil {
		t.Error("It should not return nil.")
	}

	if p.Name != "package" ||
		p.Authors == nil || p.Contributors == nil || p.Version == nil {

		t.Error("Deserializing did not work correctly:", p)
	}

	p, err = ParsePack(&badIO{})
	if err != fakeError {
		t.Error("Should report read errors, got:", err, "want:", fakeError)
	}

	p, err = ParsePack(bytes.NewBufferString("\t"))
	if err == nil {
		t.Error("Should error on parsing failures.")
	}
}

func TestPack_WriteTo(t *T) {
	t.Parallel()

	buf := bytes.NewBufferString(testPack)
	p, err := ParsePack(buf)
	if err != nil {
		t.Error("Unexpected:", err)
	}

	buf = &bytes.Buffer{}
	err = p.WriteTo(buf)
	if err != nil {
		t.Error("Unexpected:", err)
	}
	if str := buf.String(); str != testPack {
		t.Error("Expected:", testPack, "\ngot:", str)
	}
	err = p.WriteTo(&badIO{})
	if err != fakeError {
		t.Error("Should report write errors, got:", err, "want:", fakeError)
	}
	err = p.WriteTo(&halfWrite{})
	if err != errPartialWrite {
		t.Error("Expecting partial write error, got:", err)
	}
}

func TestParseDependency(t *T) {
	t.Parallel()
	var tests = []struct {
		Input  string
		Output Dependency
		Error  string
	}{
		{``, Dependency{}, `empty`},
		{`name <3.2.5`, Dependency{`name`,
			[]*Version{&Version{LessThan, 3, 2, 5, ``}}}, ``},
		{`name <3.2.5 !=1.5.0`, Dependency{`name`,
			[]*Version{
				&Version{LessThan, 3, 2, 5, ``},
				&Version{LessThan, 3, 2, 5, ``},
			}}, ``},
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
		if n, e := len(out.Versions), len(test.Output.Versions); n != e {
			t.Error("Expected:", e, "elements, got:", n)
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
		{Dependency{"", []*Version{&Version{0, 1, 2, 3, ``}}}, ""},
		{Dependency{"name", nil}, "name"},
		{Dependency{"name", []*Version{&Version{GreaterThan, 1, 2, 3, ``}}},
			"name >1.2.3"},
		{Dependency{"name", []*Version{
			{GreaterThan, 1, 2, 3, ``},
			{NotEqual, 1, 5, 0, `pre`},
		}}, "name >1.2.3 !=1.5.0-pre"},
	}

	for _, test := range tests {
		if s := test.Dependency.String(); s != test.Expected {
			t.Error(test, "expected:", s, "to be equal to:", test.Expected)
		}
	}
}

func TestDependency_GetYAML(t *T) {
	t.Parallel()
	d := Dependency{"name", []*Version{&Version{NotEqual, 1, 2, 3, ``}}}
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
	comp := Version{Equal, 1, 2, 3, ``}
	if len(d.Versions) != 1 {
		t.Error("Expected a constraint.")
	} else if !d.Versions[0].Compare(comp) {
		t.Error("Expected:", d.Versions[0], "to match", comp)
	}
}
