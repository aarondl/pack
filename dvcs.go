package pack

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"regexp"
)

const (
	gitTagErr = "fatal: No names found, cannot describe anything.\n"
)

var (
	rgxGitDescribe = regexp.MustCompile(`-[0-9]+-g[a-f0-9]\s?$`)
)

// DVCS represents a distributed version control system.
type DVCS interface {
	// Status runs a status command to see if there's actually a usable dvcs
	// at this location.
	Status() error
	// Close creates a command used to clone with this engine.
	Clone(url string) error
	// Update creates a command to update the repository from a source.
	Update() error
	// Checkout creates a command to change the working copy to a specified
	// version.
	Checkout(version string) error
	// Tags creates a command to retrieve the list of tags.
	Tags() ([]string, error)
	// CurrentTag retrieves the current tag if there is one.
	CurrentTag() (string, error)
	// SetRepoPath allows overriding of the path that was set on creation.
	SetRepoPath(path string)
}

// dvcsHelper provides various helper functions for the dvcs implementations.
type dvcsHelper struct {
	// Repository is the location of the repository.
	Repository string
}

// SetRepoPath allows overriding of the path that was set on creation.
func (d *dvcsHelper) SetRepoPath(path string) {
	d.Repository = path
}

// getCmdOutput wraps all the crazy error handling required to get input
// from a command.
func (_ dvcsHelper) getCmdOutput(cmd *exec.Cmd) ([]byte, []byte, error) {
	var stdout, stderr []byte
	var stdoutPipe, stderrPipe io.ReadCloser
	var err error
	if stdoutPipe, err = cmd.StdoutPipe(); err != nil {
		return nil, nil, err
	}
	if stderrPipe, err = cmd.StderrPipe(); err != nil {
		return nil, nil, err
	}
	if err = cmd.Start(); err != nil {
		return nil, nil, err
	}
	if stdout, err = ioutil.ReadAll(stdoutPipe); err != nil {
		return stdout, stderr, nil
	}
	if stderr, err = ioutil.ReadAll(stderrPipe); err != nil {
		return stdout, stderr, nil
	}
	if err = cmd.Wait(); err != nil {
		return stdout, stderr, nil
	}

	return stdout, stderr, nil
}

// Git uses the git toolset to implement the dvcs interface.
type Git struct {
	dvcsHelper
}

// NewGit returns a new instance of the git dvcs.
func NewGit(repo string) DVCS {
	return &Git{dvcsHelper{repo}}
}

// Hg uses the mercurial toolset to implement the dvcs interface.
type Hg struct {
	dvcsHelper
}

// NewHg returns a new instance of the hg dvcs.
func NewHg(repo string) DVCS {
	return &Hg{dvcsHelper{repo}}
}

// Bzr uses the bazaar toolset to implement the dvcs interface.
type Bzr struct {
	dvcsHelper
}

// NewBzr returns a new instance of the bzr dvcs.
func NewBzr(repo string) DVCS {
	return &Bzr{dvcsHelper{repo}}
}

// repoExists checks to see if a repo exists, returns an error if it does not.
func (d dvcsHelper) repoExists() error {
	if exists, err := DirExists(d.Repository); err != nil {
		return err
	} else if !exists {
		return fmt.Errorf(`Repo "%s" does not exist.`, d.Repository)
	}
	return nil
}

// Status performs a status check on the repository to see if it's actually
// a git repository.
func (g *Git) Status() error {
	if err := g.repoExists(); err != nil {
		return err
	}

	cmd := exec.Command("git", "status")
	cmd.Dir = g.Repository
	return cmd.Run()
}

// Clone downloads a repository if it doesn't exist on disk.
func (g *Git) Clone(url string) error {
	if err := g.repoExists(); err == nil {
		return nil
	}

	cmd := exec.Command("git", "clone", url, g.Repository)
	return cmd.Run()
}

// Update updates a repository from the default remote.
func (g *Git) Update() error {
	if err := g.repoExists(); err != nil {
		return err
	}

	cmd := exec.Command("git", "fetch")
	cmd.Dir = g.Repository
	return cmd.Run()
}

// Checkout checks out a version of the repository.
func (g *Git) Checkout(version string) error {
	if err := g.repoExists(); err != nil {
		return err
	}

	cmd := exec.Command("git", "checkout", version)
	cmd.Dir = g.Repository
	return cmd.Run()
}

// Tags gets the list of all tags for the repository.
func (g *Git) Tags() ([]string, error) {
	if err := g.repoExists(); err != nil {
		return nil, err
	}

	cmd := exec.Command("git", "tag", "-l")
	cmd.Dir = g.Repository
	stdout, _, err := g.getCmdOutput(cmd)
	if err != nil {
		return nil, err
	}

	if len(stdout) == 0 {
		return nil, nil
	}

	tags := make([]string, 0)
	tagBytes := bytes.Split(stdout, []byte{'\n'})
	for i := 0; i < len(tagBytes); i++ {
		if len(tagBytes[i]) == 0 {
			continue
		}
		if rgxVersion.Match(tagBytes[i]) {
			tags = append(tags, string(tagBytes[i]))
		}
	}
	return tags, nil
}

// CurrentTag retrieves the current tag of the repository, or empty string if
// no tag exists.
func (g *Git) CurrentTag() (string, error) {
	var tag string
	if err := g.repoExists(); err != nil {
		return tag, err
	}

	cmd := exec.Command("git", "describe", "--tags")
	cmd.Dir = g.Repository
	stdout, stderr, err := g.getCmdOutput(cmd)
	if err != nil {
		if stderr != nil && string(stderr) == gitTagErr {
			return tag, nil
		}
		return tag, err
	}

	if len(stdout) == 0 || rgxGitDescribe.Match(stdout) {
		return tag, nil
	}

	return string(bytes.TrimSpace(stdout)), nil
}

// Status performs a status check on the repository to see if it's actually
// an hg repository.
func (h *Hg) Status() error {
	if err := h.repoExists(); err != nil {
		return err
	}

	cmd := exec.Command("hg", "status")
	cmd.Dir = h.Repository
	return cmd.Run()
}

// Clone downloads a repository if it doesn't exist on disk.
func (h *Hg) Clone(url string) error {
	if err := h.repoExists(); err == nil {
		return nil
	}

	cmd := exec.Command("hg", "clone", url, h.Repository)
	return cmd.Run()
}

// Update updates a repository from the default remote.
func (h *Hg) Update() error {
	if err := h.repoExists(); err != nil {
		return err
	}

	cmd := exec.Command("hg", "pull")
	cmd.Dir = h.Repository
	return cmd.Run()
}

// Checkout checks out a version of the repository.
func (h *Hg) Checkout(version string) error {
	if err := h.repoExists(); err != nil {
		return err
	}

	cmd := exec.Command("hg", "checkout", version)
	cmd.Dir = h.Repository
	return cmd.Run()
}

// Tags gets the list of all tags for the repository.
func (h *Hg) Tags() ([]string, error) {
	if err := h.repoExists(); err != nil {
		return nil, err
	}

	cmd := exec.Command("hg", "tags")
	cmd.Dir = h.Repository
	stdout, _, err := h.getCmdOutput(cmd)
	if err != nil {
		return nil, err
	}

	if len(stdout) == 0 {
		return nil, nil
	}

	tagBytes := bytes.Split(stdout, []byte{'\n'})
	tags := make([]string, 0)
	for i := 0; i < len(tagBytes); i++ {
		if len(tagBytes[i]) == 0 {
			continue
		}
		tagByte := bytes.Fields(tagBytes[i])[0]
		if rgxVersion.Match(tagByte) {
			tags = append(tags, string(tagByte))
		}
	}
	return tags, nil
}

// CurrentTag retrieves the current tag of the repository, or empty string if
// no tag exists.
func (h *Hg) CurrentTag() (string, error) {
	var tag string
	if err := h.repoExists(); err != nil {
		return tag, err
	}

	cmd := exec.Command("hg", "identify")
	cmd.Dir = h.Repository
	stdout, _, err := h.getCmdOutput(cmd)
	if err != nil {
		return tag, err
	}

	if len(stdout) == 0 {
		return tag, nil
	}

	parts := bytes.Fields(stdout)
	if len(parts) < 2 {
		return tag, nil
	}

	tagBytes := bytes.Split(parts[1], []byte{'/'})
	for i := 0; i < len(tagBytes); i++ {
		if len(tagBytes[i]) == 0 {
			continue
		}
		if rgxVersion.Match(tagBytes[i]) {
			tag = string(tagBytes[i])
			break
		}
	}

	return tag, nil
}

// Status performs a status check on the repository to see if it's actually
// a bzr repository.
func (b *Bzr) Status() error {
	if err := b.repoExists(); err != nil {
		return err
	}

	cmd := exec.Command("bzr", "status")
	return cmd.Run()
}

// Clone downloads a repository if it doesn't exist on disk.
func (b *Bzr) Clone(url string) error {
	if err := b.repoExists(); err == nil {
		return nil
	}

	return fmt.Errorf("Not implemented yet!")
}

// Update updates a repository from the default remote.
func (b *Bzr) Update() error {
	if err := b.repoExists(); err != nil {
		return err
	}

	return fmt.Errorf("Not implemented yet!")
}

// Checkout checks out a version of the repository.
func (b *Bzr) Checkout(version string) error {
	if err := b.repoExists(); err != nil {
		return err
	}

	return fmt.Errorf("Not implemented yet!")
}

// Tags gets the list of all tags for the repository.
func (b *Bzr) Tags() ([]string, error) {
	if err := b.repoExists(); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("Not implemented yet!")
}

// CurrentTag retrieves the current tag of the repository, or empty string if
// no tag exists.
func (b *Bzr) CurrentTag() (string, error) {
	var tag string
	if err := b.repoExists(); err != nil {
		return tag, err
	}

	return "", fmt.Errorf("Not implemented yet!")
}
