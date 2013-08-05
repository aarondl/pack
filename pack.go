/*
Package pack provides low level types, parsing, and semantic versioning tools
to facilitate the reading, writing, and comparison of package metadata.
*/
package pack

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"launchpad.net/goyaml"
	"regexp"
)

const (
	errFmtDep = `pack: [%v] must be in the form: name( versionconstraint)?`
)

var (
	errPartialWrite = errors.New(`pack: Partial write on pack serialization.`)
	rgxDependency   = regexp.MustCompile(
		`(?i)^([a-z][a-z0-9_\-]+)( [a-z0-9<>=!~\.]+)?$`)
)

// Author is metadata about an author.
type Author struct {
	Name     string
	Email    string
	Homepage string
}

// Repository is a repository.
type Repository struct {
	// Type can be one of: git/mercurial/bazaar
	Type string
	Url  string
}

// Dependency is a package dependency.
type Dependency struct {
	Name    string
	Version *Version
}

// Pack is the metadata of a package.
type Pack struct {
	// Display name for the package, ImportPath's trailing name if not provided.
	Name string
	// The import path of the package.
	ImportPath string
	// Version
	Version *Version
	// Short description of the package.
	Summary string
	// Longer description of the package.
	Description string
	// Homepage
	Homepage string
	// Repository
	Repository Repository
	// License type
	License string
	// Authors
	Authors []*Author
	// Contributors
	Contributors []*Author
	// Dependencies of the package.
	Dependencies []*Dependency
	// Subpackages are used to mark packages that should be tagged with this
	// same metadata. They must be subdirectories. This is useful for
	// having subpackages within the same vcs repository.
	Subpackages []string
}

// ParsePack reads yaml from a reader and parses it into a pack object.
func ParsePack(reader io.Reader) (*Pack, error) {
	var p *Pack

	read, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	p = new(Pack)
	err = goyaml.Unmarshal(read, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// WriteTo writes the pack object to the passed in writer.
func (p *Pack) WriteTo(writer io.Writer) error {
	written, err := goyaml.Marshal(p)
	if err != nil {
		return err
	}

	n, err := writer.Write(written)
	if err != nil {
		return err
	}
	if n != len(written) {
		return errPartialWrite
	}

	return nil
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
		var v Version
		v, err = ParseVersion(parts[2][1:])
		if err != nil {
			return
		}
		dep.Version = &v
	}
	return
}

// String turns a Dependency into a String.
func (d *Dependency) String() (str string) {
	str = d.Name
	if len(str) > 0 && d.Version != nil {
		str += " " + d.Version.String()
	}
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
