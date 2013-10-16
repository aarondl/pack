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

func testDvcsHelper(t *T, zipfile string, dvcs DVCS) {
	var err error
	tmpDir := os.TempDir()
	gopackTestDir := filepath.Join(tmpDir, "gopacktest")
	dvcsOrigin := filepath.Join(gopackTestDir, "dvcsOrigin")
	dvcsClone := filepath.Join(gopackTestDir, "dvcsClone")
	var deleteTestDir = func() {
		err = os.RemoveAll(gopackTestDir)
		if err != nil && !os.IsNotExist(err) {
			t.Fatal("Error removing temporary dir:", err)
		}
	}

	deleteTestDir()
	if err = os.MkdirAll(dvcsOrigin, 0770); err != nil {
		t.Fatal("Could not create directory:", err)
	}
	defer deleteTestDir()

	if err = unzipArchive(zipfile, dvcsOrigin); err != nil {
		t.Fatal("Failed to unzip dvcs archive:", err)
	}

	var tags []string
	dvcs.SetRepoPath(dvcsOrigin)
	if err = dvcs.Status(); err != nil {
		t.Fatal("Status should not error if it is a repo:", err)
	}
	if tags, err = dvcs.Tags(); err != nil {
		t.Fatal("Failed to retrieve tags:", err)
	} else if len(tags) != 2 {
		t.Error("Expected 2 tags, got:", len(tags), tags)
	}
	for _, tag := range tags {
		if err = dvcs.Checkout(tag); err != nil {
			t.Error("Failed to checkout tag:", err)
		}
		if ctag, err := dvcs.CurrentTag(); err != nil {
			t.Error("Failed to retrieve current tag:", err)
		} else if ctag != tag {
			t.Errorf("Expected tag: %s, got: %s", tag, ctag)
		}
	}
	dvcs.SetRepoPath(dvcsClone)
	if err = dvcs.Clone(dvcsOrigin); err != nil {
		t.Error("Failed to clone repository:", err)
	}
	if err = dvcs.Clone(dvcsOrigin); err != nil {
		t.Error("Expected no error on useless clone but got:", err)
	}
	if err = dvcs.Update(); err != nil {
		t.Error("Failed to update repository:", err)
	}
}

func TestGit(t *T) {
	if Short() {
		t.SkipNow()
	}

	testDvcsHelper(t, "testgit.zip", &Git{})
}

func TestHg(t *T) {
	if Short() {
		t.SkipNow()
	}

	testDvcsHelper(t, "testhg.zip", &Hg{})
}

func TestBzr(t *T) {
	// When this becomes an actual issue, deal with it.
	t.Log("Is bzr actually a dvcs?")
	t.SkipNow()

	if Short() {
		t.SkipNow()
	}

	testDvcsHelper(t, "testbzr.zip", &Bzr{})
}
