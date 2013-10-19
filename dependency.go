package pack

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

const (
	errFmtName = `pack: [%v] must be in the form: ` +
		`importpath [constraints]* [url]?`
	errFmtConstraint = `pack: [%v] constraints must have the form: ` +
		`(=|!=|>|<|>=|<=|~)version`
	errFmtUrl = `pack: [%v] urls must have the form: (git|hg|bzr)(:url)?`
)

var (
	rgxDepUrl = regexp.MustCompile(
		`(?i)^(git|bzr|hg)(?::([a-z0-9\?\-_@\.:/=%&]+))?$`)
	rgxConstraint = regexp.MustCompile(
		`(?i)^(=|!=|>|<|>=|<=|~)?([0-9]\.[0-9]+\.[0-9]+(?:-[a-z0-9\-\.]+)?)$`)
)

// Dependency is a package dependency.
type Dependency struct {
	Name        string
	Constraints []*Constraint
	URL         string
}

// Constraint is a constraint on a dependency.
type Constraint struct {
	Operator ComparisonOp
	Version  *Version
}

// ParseDependency parses a string into a Dependency.
func ParseDependency(str string) (*Dependency, error) {
	var dep *Dependency
	var n, i int
	var err error

	var parts = strings.Split(str, " ")
	if len(str) == 0 || len(parts[0]) == 0 {
		return nil, fmt.Errorf(errFmtName, str)
	}

	dep = new(Dependency)
	dep.Name = parts[0]
	parts = parts[1:]
	if n = len(parts); n == 0 {
		return dep, nil
	}

	for i = 0; i < n; i++ {
		opVersion := rgxConstraint.FindStringSubmatch(parts[i])
		if opVersion == nil {
			if i+1 == n {
				if dep.Constraints != nil {
					dep.Constraints = dep.Constraints[:n-1]
				}
				break // Give a chance for url parsing too.
			}
			return nil, fmt.Errorf(errFmtConstraint, parts[i])
		}

		if dep.Constraints == nil {
			dep.Constraints = make([]*Constraint, n)
		}
		dep.Constraints[i] = new(Constraint)
		con := dep.Constraints[i]
		if len(opVersion[1]) > 0 {
			con.Operator, err = ParseOp(opVersion[1])
			if err != nil {
				return nil, err
			}
		} else {
			con.Operator = Equal
		}
		con.Version, err = ParseVersion(opVersion[2])
		if err != nil {
			panic(str + " -- " + parts[i] + " -- " + opVersion[0])
			return nil, err
		}
	}

	parts = parts[i:]
	if len(parts) == 0 {
		return dep, nil
	}

	if rgxDepUrl.MatchString(parts[0]) {
		dep.URL = parts[0]
	} else {
		return nil, fmt.Errorf(errFmtUrl, parts[0])
	}

	return dep, nil
}

// String turns a Dependency into a String.
func (d *Dependency) String() (str string) {
	var buf bytes.Buffer
	if len(d.Name) == 0 {
		return
	}

	buf.WriteString(d.Name)
	for _, con := range d.Constraints {
		buf.WriteByte(' ')
		buf.WriteString(con.Operator.String())
		buf.WriteString(con.Version.String())
	}
	if len(d.URL) > 0 {
		buf.WriteByte(' ')
		buf.WriteString(d.URL)
	}
	str = buf.String()
	return
}

// GetYAML implements the goyaml Getter interface.
func (v *Dependency) GetYAML() (_ string, value interface{}) {
	return "", v.String()
}

// SetYAML implements the goyaml Setter interface.
func (d *Dependency) SetYAML(_ string, value interface{}) (ok bool) {
	var s string
	var err error
	var tmp *Dependency
	if s, ok = value.(string); ok {
		tmp, err = ParseDependency(s)
		if ok = tmp != nil && err == nil; !ok {
			return
		}
		*d = *tmp
	}
	return
}
