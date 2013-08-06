package pack

import (
	"os"
)

// ParsePackFile opens a file for reading and parses it into a Pack.
func ParsePackFile(filename string) (p *Pack, err error) {
	var file *os.File
	file, err = os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()
	p, err = ParsePack(file)
	return
}

// WritePackFile opens a file for writing and writes the Pack to it.
func (p *Pack) WritePackFile(filename string) (err error) {
	var file *os.File
	file, err = os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()
	err = p.WriteTo(file)
	return
}
