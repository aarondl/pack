/*
Package pack provides low level structs, parsing, and semantic versioning tools
to facilitate the reading, writing, and comparison of package meta data.
*/
package pack

// Pack is the meta data type of a package.
type Pack struct {
	// Display name for the package, ImportPath's trailing name if not provided.
	Name string
	// The import path of the package.
	ImportPath string
	// Version
	Version Version
}
