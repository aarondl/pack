package pack

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	errFmtDep = `pack: [%v] must be in the form: name( versionconstraint)*`
)

var (
	rgxDependency = regexp.MustCompile(
		`(?i)^([a-z][a-z0-9_\-]+)((:? [a-z0-9<>=!~\.-]+)+)?$`)
)

// Dependency is a package dependency.
type Dependency struct {
	Name     string
	Versions []*Version
}

// ParseDependency parses a string into a Dependency.
func ParseDependency(str string) (dep Dependency, err error) {
	if len(str) == 0 {
		err = errors.New(errMsgEmpty)
		return
	}

	parts := rgxDependency.FindStringSubmatch(str)
	if parts == nil {
		err = fmt.Errorf(errFmtDep, str)
		return
	}

	dep.Name = parts[1]
	if len(parts[2]) > 0 {
		splits := strings.Split(parts[2], " ")[1:]
		ln := len(splits)

		dep.Versions = make([]*Version, ln)
		for i := 0; i < ln; i++ {
			dep.Versions[i] = new(Version)
			*dep.Versions[i], err = ParseVersion(splits[i])
			if err != nil {
				return
			}
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
	if d.Versions != nil {
		for i := 0; i < len(d.Versions); i++ {
			buf.WriteByte(' ')
			buf.WriteString(d.Versions[i].String())
		}
	}
	str = buf.String()
	return
}

// GetYAML implements the goyaml Getter interface.
func (v *Dependency) GetYAML() (_ string, value interface{}) {
	return "", v.String()
}

// SetYAML implements the goyaml Setter interface.
func (v *Dependency) SetYAML(_ string, value interface{}) (ok bool) {
	var s string
	var err error
	if s, ok = value.(string); ok {
		*v, err = ParseDependency(s)
		ok = err == nil
	}
	return
}
