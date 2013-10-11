/*
Package pack provides low level types, parsing, and semantic versioning tools
to facilitate the reading, writing, and comparison of package metadata.
*/
package pack

import (
	"errors"
	"io"
	"io/ioutil"
	"launchpad.net/goyaml"
)

var (
	errPartialWrite = errors.New(`pack: Partial write on pack serialization.`)
)

// Author is metadata about an author.
type Author struct {
	Name     string   `yaml:",omitempty"`
	Emails   []string `yaml:",omitempty"`
	Homepage string   `yaml:",omitempty"`
}

// Support contains the locations at which to find support for the package.
type Support struct {
	Website string `yaml:",omitempty"`
	Email   string `yaml:",omitempty"`
	Forum   string `yaml:",omitempty"`
	Wiki    string `yaml:",omitempty"`
	Issues  string `yaml:",omitempty"`
}

// Repository is a version control repository endpoint.
type Repository struct {
	// Type can be one of: git/mercurial/bazaar
	Type string `yaml:",omitempty"`
	URL  string `yaml:",omitempty"`
}

// Pack is the metadata of a package.
type Pack struct {
	// Display name for the package, ImportPath's trailing name if not provided.
	Name string `yaml:",omitempty"`
	// The import path of the package.
	ImportPath string `yaml:",omitempty"`
	// Version
	Version *Version `yaml:",omitempty"`
	// Short description of the package.
	Summary string `yaml:",omitempty"`
	// Longer description of the package.
	Description string `yaml:",omitempty"`
	// Homepage
	Homepage string `yaml:",omitempty"`
	// Repository
	Repository *Repository `yaml:",omitempty"`
	// License type ie. MIT, LGPL-3.0+, GPL-3.0+, Apache-2.0
	License string `yaml:",omitempty"`
	// Authors
	Authors []*Author `yaml:",omitempty"`
	// Contributors
	Contributors []*Author `yaml:",omitempty"`
	// Support
	Support *Support `yaml:",omitempty"`
	// Dependencies of the package.
	Dependencies []*Dependency `yaml:",omitempty"`
	// Environments of the package.
	Environments map[string][]*Dependency `yaml:",omitempty"`
	// Subpackages are used to mark packages that should be tagged with this
	// same metadata. They must be subdirectories. This is useful for
	// having subpackages within the same vcs repository.
	Subpackages []string `yaml:",omitempty"`
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
