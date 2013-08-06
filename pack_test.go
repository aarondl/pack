package pack

import (
	"bytes"
	"errors"
	. "testing"
)

var testPack = `name: package
importpath: github.com/user/package/import
version: 1.0.0
summary: Package summary
description: Package Description
homepage: www.package.com
repository:
  type: git
  url: github.com/user/package
license: mit
authors:
- name: Author1
  email: author1@email.com
  homepage: blog.author.com
- name: Author2
  email: author2@email.com
  homepage: blog.author.com
contributors:
- name: Contrib1
  email: contrib@email.com
  homepage: github.com/contrib
support:
  website: support.com
  email: email@support.com
  forum: forum.com
  wiki: wiki.com
  issues: github.com/issues
dependencies:
- dep >1.2.3
- dep2 ~1.4.5-pre !=1.5.0
subpackages:
- subpackage
`

type badIO struct {
}
type halfWrite struct {
}

var fakeError = errors.New("Fake error.")

func (f *badIO) Read(in []byte) (int, error) {
	return 0, fakeError
}
func (b *badIO) Write(out []byte) (int, error) {
	return 0, fakeError
}
func (b *halfWrite) Write(out []byte) (int, error) {
	return len(out) / 2, nil
}

func TestParsePack(t *T) {
	t.Parallel()

	buf := bytes.NewBufferString(testPack)
	p, err := ParsePack(buf)
	if err != nil {
		t.Error("Unexpected:", err)
	}
	if p == nil {
		t.Error("It should not return nil.")
	}

	if p.Name != "package" ||
		p.Authors == nil || p.Contributors == nil || p.Version == nil {

		t.Error("Deserializing did not work correctly:", p)
	}

	p, err = ParsePack(&badIO{})
	if err != fakeError {
		t.Error("Should report read errors, got:", err, "want:", fakeError)
	}

	p, err = ParsePack(bytes.NewBufferString("\t"))
	if err == nil {
		t.Error("Should error on parsing failures.")
	}
}

func TestPack_WriteTo(t *T) {
	t.Parallel()

	buf := bytes.NewBufferString(testPack)
	p, err := ParsePack(buf)
	if err != nil {
		t.Error("Unexpected:", err)
	}

	buf = &bytes.Buffer{}
	err = p.WriteTo(buf)
	if err != nil {
		t.Error("Unexpected:", err)
	}
	if str := buf.String(); str != testPack {
		t.Error("Expected:", testPack, "\ngot:", str)
	}
	err = p.WriteTo(&badIO{})
	if err != fakeError {
		t.Error("Should report write errors, got:", err, "want:", fakeError)
	}
	err = p.WriteTo(&halfWrite{})
	if err != errPartialWrite {
		t.Error("Expecting partial write error, got:", err)
	}
}
