package pack

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
)

// DVCSEngine defines behavior required for the DVCS type to function.
type DVCSEngine interface {
	// Close creates a command used to clone with this engine.
	Clone(url, destination string) *exec.Cmd
	// Update creates a command to update the repository from a source.
	Update() *exec.Cmd
	// Checkout creates a command to change the working copy to a specified
	// version.
	Checkout(version string) *exec.Cmd
	// Tags creates a command to retrieve the list of tags.
	Tags() *exec.Cmd
	// FilterTags filters the output of the Tags() command to produce a concise
	// list of validated tags in the repository. Input is given as an array
	// of byte strings.
	FilterTags(input [][]byte) []string
}

// DVCS is the common behavior behind all of the dvcs engines.
type DVCS struct {
	// repo The repository on which all these actions take place.
	Repository string
	// DVCSEngine is the actual dvcs powering these operations.
	engine DVCSEngine
}

// Git is the git implementation of DVCSEngine.
type Git struct {
}

// Hg is the mercurial implementation of DVCSEngine.
type Hg struct {
}

// Bzr is the bazaar implementation of DVCSEngine.
type Bzr struct {
}

// Clone downloads a repository if it doesn't exist on disk.
func (d *DVCS) Clone(url string) error {
	_, err := os.Stat(d.Repository)
	if err == nil {
		return nil
	}

	cmd := d.engine.Clone(url, d.Repository)
	return cmd.Run()
}

// Update updates a repository.
func (d *DVCS) Update() error {
	_, err := os.Stat(d.Repository)
	if os.IsNotExist(err) {
		return fmt.Errorf(`Repo "%s" does not exist.`, d.Repository)
	}

	cmd := d.engine.Update()
	cmd.Dir = d.Repository
	return cmd.Run()
}

// Checkout checks out a version of the repository.
func (d *DVCS) Checkout(version string) error {
	_, err := os.Stat(d.Repository)
	if os.IsNotExist(err) {
		return fmt.Errorf(`Repo "%s" does not exist.`, d.Repository)
	}

	cmd := d.engine.Checkout(version)
	cmd.Dir = d.Repository
	return cmd.Run()
}

// Tags gets the list of all tags for the repository.
func (d *DVCS) Tags() ([]string, error) {
	_, err := os.Stat(d.Repository)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf(`Repo "%s" does not exist.`, d.Repository)
	}

	cmd := d.engine.Tags()
	cmd.Dir = d.Repository
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	tagBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err
	}

	return d.engine.FilterTags(bytes.Split(tagBytes, []byte{'\n'})), nil
}

// Clone does a git clone operation if the path doesn't exist.
func (g Git) Clone(url, destination string) *exec.Cmd {
	return exec.Command("git", "clone", url, destination)
}

// Update does a fetch operation on the repository.
func (g Git) Update() *exec.Cmd {
	return exec.Command("git", "fetch")
}

// Checkout does a git checkout to the version provided.
func (g Git) Checkout(version string) *exec.Cmd {
	return exec.Command("git", "checkout", version)
}

// Tags does a git tag -l to get all the versions in the repository.
func (g Git) Tags() *exec.Cmd {
	return exec.Command("git", "tag", "-l")
}

// FilterTags filters the output of the Tags() command.
func (g Git) FilterTags(byteTags [][]byte) []string {
	tags := make([]string, 0)
	for i := 0; i < len(byteTags); i++ {
		if len(byteTags[i]) == 0 {
			continue
		}
		if rgxVersion.Match(byteTags[i]) {
			tags = append(tags, string(byteTags[i]))
		}
	}
	return tags
}

// Clone does an hg clone operation if the path doesn't exist.
func (h Hg) Clone(url, destination string) *exec.Cmd {
	return exec.Command("hg", "clone", url, destination)
}

// Update does an update operation on the repository.
func (h Hg) Update() *exec.Cmd {
	return exec.Command("hg", "pull")
}

// Checkout does an hg update to the version provided.
func (h Hg) Checkout(version string) *exec.Cmd {
	return exec.Command("hg", "update", version)
}

// Tags does an hg tags to get all the versions in the repository.
func (h Hg) Tags() *exec.Cmd {
	return exec.Command("hg", "tags")
}

// FilterTags filters the output of the Tags() command.
func (h Hg) FilterTags(byteTags [][]byte) []string {
	tags := make([]string, 0)
	for i := 0; i < len(byteTags); i++ {
		if len(byteTags[i]) == 0 {
			continue
		}
		byteTag := bytes.Fields(byteTags[i])[0]
		if rgxVersion.Match(byteTag) {
			tags = append(tags, string(byteTag))
		}
	}
	return tags
}

// Clone does a bzr branch operation if the path doesn't exist.
func (b Bzr) Clone(url, destination string) *exec.Cmd {
	return exec.Command("bzr", "branch", url, destination)
}

// Update does a bzr pull to get the tags.
func (b Bzr) Update() *exec.Cmd {
	return exec.Command("bzr", "pull", "--overwrite-tags", "upstream")
}

// Checkout does a bzr update to the version provided.
func (b Bzr) Checkout(version string) *exec.Cmd {
	return exec.Command("bzr", "update", version)
}

// Tags does a git tag -l to get all the versions that exist in this repository.
func (b Bzr) Tags() *exec.Cmd {
	return exec.Command("bzr", "tags")
}

// FilterTags filters the output of the Tags() command.
func (b Bzr) FilterTags(byteTags [][]byte) []string {
	tags := make([]string, 0)
	for i := 0; i < len(byteTags); i++ {
		if len(byteTags[i]) == 0 {
			continue
		}
		byteTag := bytes.Fields(byteTags[i])[0]
		if rgxVersion.Match(byteTag) {
			tags = append(tags, string(byteTag))
		}
	}
	return tags
}

// tryUriParse tries to parse the given string into a uri.
func tryUriParse(pathOrUrl string) (*url.URL, error) {
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
