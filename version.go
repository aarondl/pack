package pack

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

const (
	intBase     = 10
	intSize     = 32
	errMsgEmpty = `pack: Version string must not be empty.`
	errFmtMatch = `pack: [%v] must be in the form (~|>|<|!)=major.minor.patch`
)

var (
	rgxVersion = regexp.MustCompile(
		`^(=|!=|>|<|>=|<=|~>)?([0-9]+)\.([0-9]+)\.([0-9]+)$`)
)

// A comparison operator type.
type ComparisonOp int

// Defines the comparison operator types.
const (
	// Equal is the = operator.
	Equal ComparisonOp = iota
	// NotEqual is the != operator.
	NotEqual
	// GreaterThan is the > operator.
	GreaterThan
	// LessThan is the < operator.
	LessThan
	// GreaterEqual is the >= operator.
	GreaterEqual
	// LessEqual is the <= operator.
	LessEqual
	// ApproxGreater is the ~> operator.
	// This operator means "greater than or equal to so long as the major
	// version is not incremented".
	ApproxGreater
)

// Version is a semantic version number with an optional comparison operator.
// For example: 2.1.0
// 2 = Major, 1 = Minor, 0 = Patch
// For a more thorough explanation see: http://semver.org/
type Version struct {
	// Operator is the operator included in this version.
	Operator ComparisonOp
	// Major version of the package.
	Major uint
	// Minor version of the package.
	Minor uint
	// Patch version of the package.
	Patch uint
}

// Parse a string into a version struct.
func Parse(str string) (version Version, err error) {
	if len(str) == 0 {
		err = errors.New(errMsgEmpty)
		return
	}
	parts := rgxVersion.FindStringSubmatch(str)

	if parts == nil {
		err = fmt.Errorf(errFmtMatch, str)
		return
	}

	switch parts[1] {
	case `!=`:
		version.Operator = NotEqual
	case `>`:
		version.Operator = GreaterThan
	case `<`:
		version.Operator = LessThan
	case `>=`:
		version.Operator = GreaterEqual
	case `<=`:
		version.Operator = LessEqual
	case `~>`:
		version.Operator = ApproxGreater
	}

	var n uint64
	if n, err = strconv.ParseUint(parts[2], intBase, intSize); err != nil {
		return
	} else {
		version.Major = uint(n)
	}

	if n, err = strconv.ParseUint(parts[3], intBase, intSize); err != nil {
		return
	} else {
		version.Minor = uint(n)
	}

	if n, err = strconv.ParseUint(parts[4], intBase, intSize); err != nil {
		return
	} else {
		version.Patch = uint(n)
	}

	return
}

// Checks that the base version (lhs) satisfies the condition version on the rhs
// Example: 2.0.0 is the base version, and <=2.1.3 is the condition version
// will return true.
func (b Version) Compare(c Version) bool {
	switch c.Operator {
	case Equal:
		return b.Major == c.Major && b.Minor == c.Minor && b.Patch == c.Patch
	case NotEqual:
		return b.Major != c.Major || b.Minor != c.Minor || b.Patch != c.Patch
	case GreaterThan:
		return b.Major > c.Major ||
			b.Major == c.Major && b.Minor > c.Minor ||
			b.Major == c.Major && b.Minor == c.Minor && b.Patch > c.Patch
	case LessThan:
		return b.Major < c.Major ||
			b.Major == c.Major && b.Minor < c.Minor ||
			b.Major == c.Major && b.Minor == c.Minor && b.Patch < c.Patch
	case GreaterEqual:
		return b.Major == c.Major && b.Minor == c.Minor && b.Patch == c.Patch ||
			b.Major > c.Major ||
			b.Major == c.Major && b.Minor > c.Minor ||
			b.Major == c.Major && b.Minor == c.Minor && b.Patch > c.Patch
	case LessEqual:
		return b.Major == c.Major && b.Minor == c.Minor && b.Patch == c.Patch ||
			b.Major < c.Major ||
			b.Major == c.Major && b.Minor < c.Minor ||
			b.Major == c.Major && b.Minor == c.Minor && b.Patch < c.Patch
	case ApproxGreater:
		return b.Major == c.Major && b.Minor == c.Minor && b.Patch == c.Patch ||
			b.Major == c.Major && b.Minor > c.Minor ||
			b.Major == c.Major && b.Minor == c.Minor && b.Patch > c.Patch
	}
	return false
}

// String changes the version into a string representation.
func (v Version) String() string {
	var sym string
	switch v.Operator {
	case NotEqual:
		sym = `!=`
	case GreaterThan:
		sym = `>`
	case LessThan:
		sym = `<`
	case GreaterEqual:
		sym = `>=`
	case LessEqual:
		sym = `<=`
	case ApproxGreater:
		sym = `~>`
	}
	return fmt.Sprintf(`%s%d.%d.%d`, sym, v.Major, v.Minor, v.Patch)
}

// GetYAML implements the goyaml Getter interface.
func (v *Version) GetYAML() (_ string, value interface{}) {
	return "", v.String()
}

// SetYAML implements the goyaml Setter interface.
func (v *Version) SetYAML(_ string, value interface{}) (ok bool) {
	var s string
	var err error
	if s, ok = value.(string); ok {
		*v, err = Parse(s)
		ok = err == nil
	}
	return
}
