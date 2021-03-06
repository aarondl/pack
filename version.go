package pack

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	intBase       = 10
	intSize       = 32
	errMsgEmpty   = `pack: String must not be empty.`
	errFmtVersion = `pack: [%v] must be in the form: major.minor.patch-release`
	errFmtOp      = `pack: [%v] must be one of: = != > < >= <= ~`
)

var (
	// rgxVersion ensures:
	// 1. Major, minor, patch versions exist and are numeric with no leading 0s
	// 2. Release is preceeded by a dash
	// 3. Release's tokens are sepearated by .
	// 4. Release's tokens must be: numeric or alphanumeric starting with alpha.
	rgxVersion = regexp.MustCompile(
		`(?i)^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)` +
			`(?:-((?:[a-z][a-z0-9]*|[1-9][0-9]*)` +
			`(?:\.(?:[a-z][a-z0-9]*|[1-9][0-9]*))*))?$`)
)

// ComparisonOp represents a boolean operator.
type ComparisonOp int

// Defines the comparison operator types.
const (
	// Equal is the = operator.
	Equal ComparisonOp = iota + 1
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
	// ApproxGreater is the ~ operator.
	// This operator means "greater than or equal to so long as the major
	// version is not incremented".
	ApproxGreater
)

// Version is a semantic version number with an optional comparison operator.
// For example: 2.1.0-alpha.1
// 2 = Major, 1 = Minor, 0 = Patch, alpha.1 = Release
// For a more thorough explanation see: http://semver.org/
type Version struct {
	// Major version of the package.
	Major uint
	// Minor version of the package.
	Minor uint
	// Patch version of the package.
	Patch uint
	// Release version of the package.
	Release string
}

// ParseVersion parses a string into a version.
func ParseVersion(str string) (version *Version, err error) {
	if len(str) == 0 {
		err = errors.New(errMsgEmpty)
		return
	}
	parts := rgxVersion.FindStringSubmatch(str)

	if parts == nil {
		err = fmt.Errorf(errFmtVersion, str)
		return
	}

	version = new(Version)
	var n uint64
	if n, err = strconv.ParseUint(parts[1], intBase, intSize); err != nil {
		return
	}
	version.Major = uint(n)

	if n, err = strconv.ParseUint(parts[2], intBase, intSize); err != nil {
		return
	}
	version.Minor = uint(n)

	if n, err = strconv.ParseUint(parts[3], intBase, intSize); err != nil {
		return
	}
	version.Patch = uint(n)

	version.Release = parts[4]

	return
}

// ParseOp parses an operation string into a comparison operator type.
func ParseOp(str string) (op ComparisonOp, err error) {
	switch str {
	case `=`:
		op = Equal
	case `!=`:
		op = NotEqual
	case `>`:
		op = GreaterThan
	case `<`:
		op = LessThan
	case `>=`:
		op = GreaterEqual
	case `<=`:
		op = LessEqual
	case `~`:
		op = ApproxGreater
	default:
		err = fmt.Errorf(errFmtOp, str)
	}
	return
}

// String turns a comparison operator back into a string.
func (op ComparisonOp) String() (str string) {
	switch op {
	case Equal:
		str = `=`
	case NotEqual:
		str = `!=`
	case GreaterThan:
		str = `>`
	case LessThan:
		str = `<`
	case GreaterEqual:
		str = `>=`
	case LessEqual:
		str = `<=`
	case ApproxGreater:
		str = `~`
	}
	return
}

// Satisfies checks that the base version (lhs) satisfies the condition version
// (rhs).
// Example: 2.0.0 is the base version, and <=2.1.3 is the condition version
// will return true. Comparison is according to http://semver.org/
func (b *Version) Satisfies(op ComparisonOp, c *Version) (ok bool) {
	switch op {
	case Equal:
		ok = b.Major == c.Major && b.Minor == c.Minor && b.Patch == c.Patch &&
			b.Release == c.Release
	case NotEqual:
		ok = b.Major != c.Major || b.Minor != c.Minor || b.Patch != c.Patch ||
			b.Release != c.Release
	case GreaterThan:
		ok = b.Major > c.Major ||
			b.Major == c.Major && b.Minor > c.Minor ||
			b.Major == c.Major && b.Minor == c.Minor && b.Patch > c.Patch ||
			b.Major == c.Major && b.Minor == c.Minor && b.Patch == c.Patch &&
				compareReleases(b.Release, c.Release) > 0
	case LessThan:
		ok = b.Major < c.Major ||
			b.Major == c.Major && b.Minor < c.Minor ||
			b.Major == c.Major && b.Minor == c.Minor && b.Patch < c.Patch ||
			b.Major == c.Major && b.Minor == c.Minor && b.Patch == c.Patch &&
				compareReleases(b.Release, c.Release) < 0
	case GreaterEqual:
		ok = b.Major == c.Major && b.Minor == c.Minor && b.Patch == c.Patch &&
			b.Release == c.Release ||
			b.Major > c.Major ||
			b.Major == c.Major && b.Minor > c.Minor ||
			b.Major == c.Major && b.Minor == c.Minor && b.Patch > c.Patch ||
			b.Major == c.Major && b.Minor == c.Minor && b.Patch == c.Patch &&
				compareReleases(b.Release, c.Release) >= 0
	case LessEqual:
		ok = b.Major == c.Major && b.Minor == c.Minor && b.Patch == c.Patch &&
			b.Release == c.Release ||
			b.Major < c.Major ||
			b.Major == c.Major && b.Minor < c.Minor ||
			b.Major == c.Major && b.Minor == c.Minor && b.Patch < c.Patch ||
			b.Major == c.Major && b.Minor == c.Minor && b.Patch == c.Patch &&
				compareReleases(b.Release, c.Release) <= 0
	case ApproxGreater:
		ok = b.Major == c.Major && b.Minor == c.Minor && b.Patch == c.Patch &&
			b.Release == c.Release ||
			b.Major == c.Major && b.Minor > c.Minor ||
			b.Major == c.Major && b.Minor == c.Minor && b.Patch > c.Patch ||
			b.Major == c.Major && b.Minor == c.Minor && b.Patch == c.Patch &&
				compareReleases(b.Release, c.Release) >= 0
	}
	return
}

// compareReleases returns an integer depicting the relationship between
// release strings. Comparison is according to http://semver.org/
func compareReleases(base, compare string) int {
	if len(base) == 0 && len(compare) == 0 {
		return 0
	}
	b := strings.Split(base, ".")
	c := strings.Split(compare, ".")
	i, lb, lc := 0, len(b), len(c)
	for ; i < lb && i < lc; i++ {
		bnum, errb := strconv.ParseInt(b[i], 10, 64)
		cnum, errc := strconv.ParseInt(c[i], 10, 64)
		bIsNum, cIsNum := errb == nil, errc == nil
		switch {
		case bIsNum && !cIsNum:
			return 1
		case !bIsNum && cIsNum:
			return -1
		case bIsNum && cIsNum:
			if val := bnum - cnum; val > 0 {
				return -1
			} else if val < 0 {
				return 1
			}
		case !bIsNum && !cIsNum:
			if val := compareStrings(b[i], c[i]); val != 0 {
				return val
			}
		}
	}

	if i < lb {
		return -1
	} else if i < lc {
		return 1
	}

	return 0
}

// compareStrings is a c-style string comparison.
func compareStrings(lhs, rhs string) int {
	var i = 0
	l, r := len(lhs), len(rhs)
	for ; i < l && i < r; i++ {
		if val := int(lhs[i]) - int(rhs[i]); val > 0 {
			return -1
		} else if val < 0 {
			return 1
		}
	}

	if i < l {
		return -1
	} else if i < r {
		return 1
	}

	return 0
}

// Zero checks to see if this is a completely zero'd Version.
func (v *Version) Zero() bool {
	return v.Major == 0 && v.Minor == 0 && v.Patch == 0
}

// String changes the version into a string representation.
func (v Version) String() string {
	var release string
	if len(v.Release) > 0 {
		release = "-" + v.Release
	}
	return fmt.Sprintf(
		`%d.%d.%d%s`, v.Major, v.Minor, v.Patch, release)
}

// GetYAML implements the goyaml Getter interface.
func (v *Version) GetYAML() (_ string, value interface{}) {
	return "", v.String()
}

// SetYAML implements the goyaml Setter interface.
func (v *Version) SetYAML(_ string, value interface{}) (ok bool) {
	var s string
	var err error
	var tmp *Version
	if s, ok = value.(string); ok {
		tmp, err = ParseVersion(s)
		if ok = tmp != nil && err == nil; !ok {
			return
		}
		*v = *tmp
	}
	return
}
