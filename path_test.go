package pack

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	. "testing"
)

const (
	fakeGoPath  = "/tmp:/usr/go"
	fakePackset = "default"
)

func Test_NewPaths(t *T) {
	store := os.Getenv(GOPATH)
	defer func() {
		os.Setenv(GOPATH, store)
	}()

	os.Setenv(GOPATH, fakeGoPath)
	p, err := NewPathsFromGopath(fakePackset)
	if err != nil {
		t.Error("Unexpected error:", err)
	}

	expect := fakeGoPath
	if expect != p.Gopath {
		t.Errorf("Expected: %s, got: %s", expect, p.Gopath)
	}
	expect = filepath.Join("/tmp/", GOPACKFOLDER)
	if expect != p.GopackPath {
		t.Errorf("Expected: %s, got: %s", expect, p.GopackPath)
	}
	expect = filepath.Join("/tmp/", GOPACKFOLDER, fakePackset, SRCFOLDER)
	if expect != p.GopacksetPath {
		t.Errorf("Expected: %s, got: %s", expect, p.GopackPath)
	}
	expect = fakeGoPath + string(filepath.ListSeparator) + expect
	if expect != p.CombinedPath {
		t.Errorf("Expected: %s, got: %s", expect, p.CombinedPath)
	}
}

func Test_EnsureDirectory(t *T) {
	if Short() {
		t.SkipNow()
	}
	testdir := filepath.Join(os.TempDir(), "ensuredirectorytest")
	defer os.RemoveAll(testdir)

	created, err := EnsureDirectory(testdir)
	if err != nil {
		t.Error("Unexpected Error:", err)
	}
	if !created {
		t.Error("Expected the folder to be created.")
	}

	_, err = os.Stat(testdir)
	if os.IsNotExist(err) {
		t.Error("Expected the folder to be created.")
	}

	created, err = EnsureDirectory(testdir)
	if err != nil {
		t.Error("Unexpected Error:", err)
	}
	if created {
		t.Error("Expected the folder to exist.")
	}
}

func Test_DirAndFileExists(t *T) {
	if Short() {
		t.SkipNow()
	}
	var exist bool
	var err error
	testdir := filepath.Join(os.TempDir(), "dirandfileexists")
	testfile := filepath.Join(testdir, "testfile.txt")

	_, err = os.Stat(testdir)
	if err == nil || !os.IsNotExist(err) {
		t.Error("Expected the folder to not exist:", err)
	}

	if exist, err = DirExists(testdir); err != nil {
		t.Error("Unexpected error:", err)
	} else if exist {
		t.Error("Expected dir to not exist:", testdir)
	}

	if exist, err = FileExists(testfile); err != nil {
		t.Error("Unexpected error:", err)
	} else if exist {
		t.Error("Expected file to not exist:", testfile)
	}

	err = os.Mkdir(testdir, 0770)
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	defer os.RemoveAll(testdir)

	if exist, err = DirExists(testdir); err != nil {
		t.Error("Unexpected error:", err)
	} else if !exist {
		t.Error("Expected an existing dir:", testdir)
	}

	f, err := os.Create(testfile)
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	err = f.Close()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	if exist, err = FileExists(testfile); err != nil {
		t.Error("Unexpected error:", err)
	} else if !exist {
		t.Error("Expected file to exist:", testfile)
	}

	exist, err = DirExists(testfile)
	if err == nil || !strings.Contains(err.Error(), "dir, but found file") {
		t.Error("Expected an error due to not being dir, but got:", err)
	}

	exist, err = FileExists(testdir)
	if err == nil || !strings.Contains(err.Error(), "file, but found dir") {
		t.Error("Expected an error due to not being file, but got:", err)
	}
}

func Test_PackageExists(t *T) {
	if Short() {
		t.SkipNow()
	}

	tmp := os.TempDir()
	dir := "checkpackageexisttest"
	testdir := filepath.Join(tmp, dir)
	gopath1 := filepath.Join(testdir, "gopath1")
	gopath2 := filepath.Join(testdir, "gopath2")
	gopath := fmt.Sprintf("%c%s%c%s", filepath.ListSeparator, gopath1,
		filepath.ListSeparator, gopath2)
	p, err := NewPaths(gopath, fakePackset)
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	pkg1 := filepath.Join(gopath2, "src", "pkg1")
	pkg2 := filepath.Join(p.GopacksetPath, "github.com", "user", "pkg2")

	err = os.MkdirAll(pkg1, 0770)
	if err != nil {
		t.Error("Error creating dir:", err)
	}
	err = os.MkdirAll(pkg2, 0770)
	if err != nil {
		t.Error("Error creating dir:", err)
	}
	defer os.RemoveAll(testdir)

	path, inGopack, err := p.PackageExists("pkg1")
	if err != nil {
		t.Error("Unexpected Error:", err)
	} else if inGopack {
		t.Error("Expected package to exist outside of gopack")
	} else if path != pkg1 {
		t.Errorf("Expected the path: %s got: %s", pkg1, path)
	}

	path, inGopack, err = p.PackageExists("github.com/user/pkg2")
	if err != nil {
		t.Error("Unexpected Error:", err)
	} else if !inGopack {
		t.Error("Expected the package to exist in gopack")
	} else if path != pkg2 {
		t.Errorf("Expected the path: %s got: %s", pkg2, path)
	}
}

func Test_GopathSetRestore(t *T) {
	// Do this anyways in case the real code fails.
	store := os.Getenv(GOPATH)
	defer os.Setenv(GOPATH, store)

	p, err := NewPaths(fakeGoPath, fakePackset)
	if err != nil {
		t.Error("Unexpected error:", err)
	}

	p.GopathSet()
	if os.Getenv(GOPATH) != p.CombinedPath {
		t.Error("Expected the combined path to have been set.")
	}
	p.GopathRestore()
	if os.Getenv(GOPATH) != p.Gopath {
		t.Error("Expected the original path to have been restored.")
	}
}

func Test_TryUriParse(t *T) {
	uri, err := TryUriParse(`/path/to/file`)
	if uri != nil {
		t.Error("Expected nil url on file path, got:", uri)
	}

	uri, err = TryUriParse(`git://github.com/aarondl/pack`)
	if err != nil {
		t.Error("Expected no error, got:", err)
	}
	if uri == nil {
		t.Error("Uri should not be nil.")
	}

	uri, err = TryUriParse(`ssh+git://github.com/aarondl/pack`)
	if err != nil {
		t.Error("Expected no error, got:", err)
	}
	if uri == nil {
		t.Error("Uri should not be nil.")
	}

	uri, err = TryUriParse(`bad/path`)
	if err == nil {
		t.Error("Expected error, but it was nil.")
	}
}

func Test_SplitAndCullPath(t *T) {
	paths := fmt.Sprintf("%c%c%s%c%c%s%c%s%c%c",
		filepath.ListSeparator,
		filepath.ListSeparator,
		"path",
		filepath.ListSeparator,
		filepath.ListSeparator,
		"path",
		filepath.ListSeparator,
		"path",
		filepath.ListSeparator,
		filepath.ListSeparator)

	result := splitAndCullPath(paths)
	if len(result) != 3 {
		t.Fatal("Expected only 3 valid entries, got:", len(result))
	}
	for _, r := range result {
		if r != "path" {
			t.Fatal("Expected: path,path,path but got:", result)
		}
	}
}
