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
	Name     string
	Email    string
	Homepage string
}

// Support contains the locations at which to find support for the package.
type Support struct {
	Website string
	Email   string
	Forum   string
	Wiki    string
	Issues  string
}

// Repository is a repository.
type Repository struct {
	// Type can be one of: git/mercurial/bazaar
	Type string
	Url  string
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
	// Support
	Support Support
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
