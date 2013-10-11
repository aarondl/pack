package pack

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	. "testing"
)

// unzipArchive is here to ease porting to windows later. (Instead of just exec
// unzip for example).
func unzipArchive(zipfile, targetdir string) error {
	file, err := zip.OpenReader(zipfile)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, f := range file.File {
		filename := filepath.Join(targetdir, f.Name)
		if f.FileInfo().IsDir() {
			if err = os.MkdirAll(filename, 0770); err != nil {
				return err
			}
			continue
		}

		zipped, err := f.Open()
		if err != nil {
			return err
		}
		err = copyFileTo(zipped, filename)
		if err != nil {
			return err
		}
	}

	return nil
}

// copyFileTo copies from a reader to the file given. Writing this as a
// separate function let's us defer the closes.
func copyFileTo(zipped io.ReadCloser, filepath string) error {
	defer zipped.Close()
	writeTo, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer writeTo.Close()
	_, err = io.Copy(writeTo, zipped)
	if err != nil {
		return err
	}

	return nil
}

func TestGit(t *T) {
	if Short() {
		t.SkipNow()
	}

	var err error
	tmpDir := os.TempDir()
	gopackTestDir := filepath.Join(tmpDir, "gopacktest")
	gitOrigin := filepath.Join(gopackTestDir, "gitOrigin")
	gitClone := filepath.Join(gopackTestDir, "gitClone")
	var deleteTestDir = func() {
		err = os.RemoveAll(gopackTestDir)
		if err != nil && !os.IsNotExist(err) {
			t.Fatal("Error removing temporary dir:", err)
		}
	}

	deleteTestDir()
	if err = os.MkdirAll(gitOrigin, 0770); err != nil {
		t.Fatal("Could not create directory:", err)
	}
	defer deleteTestDir()

	if err = unzipArchive("testgit.zip", gitOrigin); err != nil {
		t.Fatal("Failed to unzip git archive:", err)
	}

	var git = DVCS{gitOrigin, Git{}}
	var tags []string
	if tags, err = git.Tags(); err != nil {
		t.Fatal("Failed to retrieve tags:", err)
	} else if len(tags) != 2 {
		t.Error("Expected 2 tags, got:", len(tags), tags)
	}
	for _, tag := range tags {
		if err = git.Checkout(tag); err != nil {
			t.Error("Failed to checkout tag:", err)
		}
	}
	git.Repository = gitClone
	if err = git.Clone(gitOrigin); err != nil {
		t.Error("Failed to clone repository:", err)
	}
	if err = git.Clone(gitOrigin); err != nil {
		t.Error("Expected no error on useless clone but got:", err)
	}
	if err = git.Update(); err != nil {
		t.Error("Failed to update repository:", err)
	}
}

func TestHg(t *T) {
	if Short() {
		t.SkipNow()
	}

	var err error
	tmpDir := os.TempDir()
	gopackTestDir := filepath.Join(tmpDir, "gopacktest")
	hgOrigin := filepath.Join(gopackTestDir, "hgOrigin")
	hgClone := filepath.Join(gopackTestDir, "hgClone")
	var deleteTestDir = func() {
		err = os.RemoveAll(gopackTestDir)
		if err != nil && !os.IsNotExist(err) {
			t.Fatal("Error removing temporary dir:", err)
		}
	}

	deleteTestDir()
	if err = os.MkdirAll(hgOrigin, 0770); err != nil {
		t.Fatal("Could not create directory:", err)
	}
	defer deleteTestDir()

	if err = unzipArchive("testhg.zip", hgOrigin); err != nil {
		t.Fatal("Failed to unzip hg archive:", err)
	}

	var hg = DVCS{hgOrigin, Hg{}}
	var tags []string
	if tags, err = hg.Tags(); err != nil {
		t.Fatal("Failed to retrieve tags:", err)
	} else if len(tags) != 2 {
		t.Error("Expected 2 tags, got:", len(tags), tags)
	}
	for _, tag := range tags {
		if err = hg.Checkout(tag); err != nil {
			t.Error("Failed to checkout tag:", err)
		}
	}
	hg.Repository = hgClone
	if err = hg.Clone(hgOrigin); err != nil {
		t.Error("Failed to clone repository:", err)
	}
	if err = hg.Clone(hgOrigin); err != nil {
		t.Error("Expected no error on useless clone but got:", err)
	}
	if err = hg.Update(); err != nil {
		t.Error("Failed to update repository:", err)
	}
}

func TestBzr(t *T) {
	// When this becomes an actual issue, deal with it.
	t.Log("Is bzr actually a dvcs?")
	t.SkipNow()

	if Short() {
		t.SkipNow()
	}

	var err error
	tmpDir := os.TempDir()
	gopackTestDir := filepath.Join(tmpDir, "gopacktest")
	bzrOrigin := filepath.Join(gopackTestDir, "bzrOrigin")
	bzrClone := filepath.Join(gopackTestDir, "bzrClone")
	var deleteTestDir = func() {
		err = os.RemoveAll(gopackTestDir)
		if err != nil && !os.IsNotExist(err) {
			t.Fatal("Error removing temporary dir:", err)
		}
	}

	deleteTestDir()
	if err = os.MkdirAll(bzrOrigin, 0770); err != nil {
		t.Fatal("Could not create directory:", err)
	}
	defer deleteTestDir()

	if err = unzipArchive("testbzr.zip", bzrOrigin); err != nil {
		t.Fatal("Failed to unzip bzr archive:", err)
	}

	var bzr = DVCS{bzrOrigin, Bzr{}}
	var tags []string
	if tags, err = bzr.Tags(); err != nil {
		t.Fatal("Failed to retrieve tags:", err)
	} else if len(tags) != 2 {
		t.Error("Expected 2 tags, got:", len(tags), tags)
	}
	for _, tag := range tags {
		if err = bzr.Checkout(tag); err != nil {
			t.Error("Failed to checkout tag:", err)
		}
	}
	bzr.Repository = bzrClone
	if err = bzr.Clone(bzrOrigin); err != nil {
		t.Error("Failed to clone repository:", err)
	}
	if err = bzr.Clone(bzrOrigin); err != nil {
		t.Error("Expected no error on useless clone but got:", err)
	}
	if err = bzr.Update(); err != nil {
		t.Error("Failed to update repository:", err)
	}
}

func TestDVCS_UriParse(t *T) {
	t.Parallel()
	uri, err := tryUriParse(`/path/to/file`)
	if uri != nil {
		t.Error("Expected nil url on file path, got:", uri)
	}

	uri, err = tryUriParse(`git://github.com/aarondl/pack`)
	if err != nil {
		t.Error("Expected no error, got:", err)
	}
	if uri == nil {
		t.Error("Uri should not be nil.")
	}

	uri, err = tryUriParse(`ssh+git://github.com/aarondl/pack`)
	if err != nil {
		t.Error("Expected no error, got:", err)
	}
	if uri == nil {
		t.Error("Uri should not be nil.")
	}

	uri, err = tryUriParse(`bad/path`)
	if err == nil {
		t.Error("Expected error, but it was nil.")
	}
}
