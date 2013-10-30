package pack

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
)

const (
	GOPATH       = "GOPATH"
	GOPACKFOLDER = "gopack"
	SRCFOLDER    = "src"
	CONFIGFILE   = "config.yaml"
)

var (
	errGoPathNotSet  = errors.New("GOPATH must be set to use this tool.")
)

// Paths contains all the paths used by gopack.
type Paths struct {
	Gopath           string
	Gopaths          []string
	GopackPath       string
	GopacksetPath    string
	CombinedPath     string
	GopackConfigPath string
	packset          string
}

// NewPaths uses the environment to locate all the paths to be used and returns
// them in a paths variable.
func NewPaths(gopath, packset string) (*Paths, error) {
	if len(gopath) == 0 {
		return nil, errGoPathNotSet
	}
	p := &Paths{Gopath: gopath}
	p.Gopaths = filepath.SplitList(gopath)
	p.GopackPath = filepath.Join(p.Gopaths[0], GOPACKFOLDER)
	p.GopackConfigPath = filepath.Join(p.GopackPath, CONFIGFILE)
	p.packset = packset
	p.GopacksetPath = filepath.Join(p.GopackPath, p.packset, SRCFOLDER)
	p.CombinedPath = p.Gopath + string(filepath.ListSeparator) + p.GopacksetPath
	return p, nil
}

// NewPathsFromGopath creates a new paths based on the gopath from the env.
func NewPathsFromGopath(packset string) (*Paths, error) {
	return NewPaths(os.Getenv(GOPATH), packset)
}

// SetPackset updates the packset and all paths that include packset.
func (p *Paths) SetPackset(packset string) {
	p.packset = packset
	p.GopacksetPath = filepath.Join(p.GopackPath, p.packset, SRCFOLDER)
	p.CombinedPath = p.Gopath + string(filepath.ListSeparator) + p.GopacksetPath
}

// Packset returns the current packset.
func (p *Paths) Packset() string {
	return p.Packset()
}

// GopathRestore restores the original gopath variable.
func (p *Paths) GopathRestore() {
	os.Setenv(GOPATH, p.Gopath)
}

// GopathAppend adds the combined path to the current gopath.
func (p *Paths) GopathSet() {
	os.Setenv(GOPATH, p.CombinedPath)
}

// PackageExists checks a packages existence. If it exists it will return
// a path, a boolean indicating whether or not it was found in the GOPACKSETPATH
func (p *Paths) PackageExists(imp string) (string, bool, error) {
	for _, gopath := range p.Gopaths {
		packagepath := filepath.Join(gopath, SRCFOLDER, imp)
		exist, err := DirExists(packagepath)
		if err != nil {
			return "", false, err
		} else if exist {
			return packagepath, false, nil
		}
	}

	packagepath := filepath.Join(p.GopacksetPath, imp)
	exist, err := DirExists(packagepath)
	if err != nil {
		return "", false, err
	} else if exist {
		panic("LOL")
		return packagepath, true, nil
	}

	return "", false, nil
}

// EnsureDirectory ensures a directory exists, or it creates it. Returns
// true if the directory had to be created.
func EnsureDirectory(dir string) (bool, error) {
	if exists, err := DirExists(dir); err != nil {
		return false, err
	} else if exists {
		return false, nil
	}

	err := os.MkdirAll(dir, 0770)
	if err != nil {
		return false, err
	}
	return true, nil
}

// DirExists checks to see if a directory exists.
func DirExists(dir string) (bool, error) {
	f, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return false, err
	}
	if !f.IsDir() {
		return false, fmt.Errorf("Expected %s to be dir, but found file.", dir)
	}
	return true, nil
}

// TryUriParse tries to parse the given string into a uri.
func TryUriParse(pathOrUrl string) (*url.URL, error) {
	if filepath.IsAbs(pathOrUrl) {
		return nil, nil
	}
	url, err := url.ParseRequestURI(pathOrUrl)
	if err != nil {
		return nil, err
	}
	if !url.IsAbs() {
		return nil, fmt.Errorf(`Expected "%s", to be an absolute path or url.`,
			pathOrUrl)
	}
	return url, nil
}
