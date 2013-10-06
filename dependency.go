package pack

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	errFmtDep = `pack: [%v] invalid name, must start with alphabetic ` +
		`and can only have the following characters: a-z0-9, -, _`
	errFmtConstraint = `pack: [%v] constraints must have the form: ` +
		`[url] (=|!=|>|<|>=|<=|~)version*`
)

var (
	rgxDepName = regexp.MustCompile(`(?i)^([a-z][a-z0-9_\-]+)$`)
	rgxDepUrl  = regexp.MustCompile(
		`^(git|bzr|hg):([a-z0-9\?\-_@\.:/=%&]+)$`)
	rgxConstraint = regexp.MustCompile(
		`(?i)^(=|!=|>|<|>=|<=|~)?([a-z0-9\.-]+)$`)
)

// Dependency is a package dependency.
type Dependency struct {
	Name        string
	URL         string
	Constraints []*Constraint
}

// Constraint is a constraint on a dependency.
type Constraint struct {
	Operator ComparisonOp
	Version  *Version
}

// ParseDependency parses a string into a Dependency.
func ParseDependency(str string) (dep *Dependency, err error) {
	if len(str) == 0 {
		err = errors.New(errMsgEmpty)
		return
	}

	parts := strings.Split(str, " ")
	if !rgxDepName.MatchString(parts[0]) {
		err = fmt.Errorf(errFmtDep, str)
		return
	}

	dep = new(Dependency)
	dep.Name = parts[0]
	if len(parts) == 1 {
		return
	}

	parts = parts[1:]
	if rgxDepUrl.MatchString(parts[0]) {
		dep.URL = parts[0]
		parts = parts[1:]
		if len(parts) == 0 {
			return
		}
	}

	n := len(parts)
	dep.Constraints = make([]*Constraint, n)
	for i := 0; i < n; i++ {
		opVersion := rgxConstraint.FindStringSubmatch(parts[i])
		if opVersion == nil {
			err = fmt.Errorf(errFmtConstraint, parts[i])
			return
		}

		dep.Constraints[i] = new(Constraint)
		con := dep.Constraints[i]
		if len(opVersion[1]) > 0 {
			con.Operator, err = ParseOp(opVersion[1])
			if err != nil {
				return
			}
		} else {
			con.Operator = Equal
		}
		con.Version, err = ParseVersion(opVersion[2])
		if err != nil {
			return
		}
	}
	return
}

// String turns a Dependency into a String.
func (d *Dependency) String() (str string) {
	var buf bytes.Buffer
	if len(d.Name) == 0 {
		return
	}

	buf.WriteString(d.Name)
	if len(d.URL) > 0 {
		buf.WriteByte(' ')
		buf.WriteString(d.URL)
	}
	for _, con := range d.Constraints {
		buf.WriteByte(' ')
		buf.WriteString(con.Operator.String())
		buf.WriteString(con.Version.String())
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
