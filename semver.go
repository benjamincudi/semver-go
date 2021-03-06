package semver

import (
	"regexp"
	"strconv"
	"strings"
)

/* A minimalist interface for semantic versioning, based on http://semver.org
 * This is based on Semver v2.0.0, and will adhere to those standards.
 * Updates will be documented at the bottom of readme.md,
 */

// This is the current version of this package
var Version = "2.0.0-alpha.0.1.1"

type PrereleaseTag struct {
	tag string
}

func makePrerelease(tag string) (PrereleaseTag, string) {
	if rxPre.MatchString(tag) {
		return PrereleaseTag{tag: tag}, ""
	} else {
		return PrereleaseTag{}, "Invalid prerelease tag"
	}
}

type BuildTag struct {
	tag string
}

func makeBuild(tag string) (BuildTag, string) {
	if rxBuild.MatchString(tag) {
		return BuildTag{tag: tag}, ""
	} else {
		return BuildTag{}, "Invalid build tag"
	}
}

/* Struct for semver string comprehension and manipulation.
 * This type and the methods associated are meant only for internal use,
 * and they have been written only with the intention of making the
 * API work easier to comprehend.
 */

type Semver struct {
	major, minor, patch uint64
	pre                 PrereleaseTag
	build               BuildTag
}

const (
	PATCH string = "patch"
	MINOR string = "minor"
	MAJOR string = "major"
)

/* Regex for matching semantic version strings, explained:
 * '^'
 * '(0|[1-9]\\d*)'	 								// major
 * '\\.(0|[1-9]\\d*)'	 							// minor
 * '\\.(0|[1-9]\\d*)'	 							// patch
 * '(?:-'	 										// start prerelease
 * '('	 											// capture
 * '(?:' 											// first identifier
 * '0|' 											// 0, or
 * '[1-9]\\d*|' 									// numeric identifier, or
 * '\\d*[a-zA-Z-][a-zA-Z0-9-]*' 					// id with at least one non-number
 * ')' 												// end first identifier
 * '(?:\\.'  										// dot-separated
 * '(?:0|[1-9]\\d*|\\d*[a-zA-Z-][a-zA-Z0-9-]*)'  	// identifier
 * ')*'  											// zero or more of those
 * ')'  											// end prerelease capture
 * ')?'  											// prerelease is optional
 * '(?:'  											// build tag (non-capturing)
 * '\\+[0-9A-Za-z-]+(?:\\.[0-9A-Za-z-]+)*'  		// pretty much anything goes
 * ')?' 											// build tag is optional
 * '$'
 *
 * Credit: https://github.com/mojombo/semver/issues/110#issuecomment-19433829
 */
var rxMatch, _ = regexp.Compile("^(0|[1-9]\\d*)\\.(0|[1-9]\\d*)\\.(0|[1-9]\\d*)(?:-((?:0|[1-9]\\d*|\\d*[a-zA-Z-][a-zA-Z0-9-]*)(?:\\.(?:0|[1-9]\\d*|\\d*[a-zA-Z-][a-zA-Z0-9-]*))*))?(?:\\+[0-9A-Za-z-]+(?:\\.[0-9A-Za-z-]+)*)?$")
var rxNumeric, _ = regexp.Compile("^(0|[1-9])+$")                                                                                 // For checking pre-release identifiers to see if they are pure numeric
var rxPre, _ = regexp.Compile("^((?:0|[1-9]\\d*|\\d*[a-zA-Z-][a-zA-Z0-9-]*)(?:\\.(?:0|[1-9]\\d*|\\d*[a-zA-Z-][a-zA-Z0-9-]*))*)$") // The rules for a pre-release string, omitting the leading dash
var rxBuild, _ = regexp.Compile("^(?:[0-9A-Za-z-]+(?:\\.[0-9A-Za-z-]+)*)$")                                                       // The rules for a build string, omitting the leading plus

// Return in the same format as provided, when applicable
func (ver Semver) ConvertToString() string {
	version := strings.Join([]string{strconv.FormatUint(ver.major, 10), strconv.FormatUint(ver.minor, 10), strconv.FormatUint(ver.patch, 10)}, ".")
	if len(ver.pre.tag) > 0 {
		version = strings.Join([]string{version, ver.pre.tag}, "-")
	}
	if len(ver.build.tag) > 0 {
		version = strings.Join([]string{version, ver.build.tag}, "+")
	}
	return version
}

// Puntastic function to make a struct from a version string
// This makes it easier to deal with various parts
func ConStructor(version string) (*Semver, string) {
	if !IsValid(version) {
		return &Semver{}, "Not a valid version string"
	}
	var ver, err, bd, rl, a, b, c string
	var bld BuildTag
	var rel PrereleaseTag
	ver = version
	if strings.Index(version, "+") > -1 {
		ver, bd = extractor(version, "+")
		bld, err = makeBuild(bd)
		if len(err) > 0 {
			return &Semver{}, "Not a valid build string"
		}
	}
	if strings.Index(ver, "-") > -1 {
		ver, rl = extractor(ver, "-")
		rel, err = makePrerelease(rl)
		if len(err) > 0 {
			return &Semver{}, "Not a valid prerelease string"
		}
	}
	a, ver = extractor(ver, ".")
	b, c = extractor(ver, ".")
	maj, _ := strconv.ParseUint(a, 10, 0)
	min, _ := strconv.ParseUint(b, 10, 0)
	pat, _ := strconv.ParseUint(c, 10, 0)
	return &Semver{major: maj, minor: min, patch: pat, pre: rel, build: bld}, ""
}

// Helper to do tediously repetitive slicing of strings
func extractor(base, mark string) (string, string) {
	lead := base[:strings.Index(base, mark)]
	tail := base[strings.Index(base, mark)+1:]
	return lead, tail
}

// Compare normal version string to see if 'a' is older than 'b'
// If normal version is entirely the same, compare pre-release strings
func (a Semver) OlderThan(b Semver) bool {
	if a.major < b.major {
		return true
	} else if a.major == b.major {
		if a.minor < b.minor {
			return true
		} else if a.minor == b.minor {
			if a.patch < b.patch {
				return true
			} else if a.patch == b.patch {
				return b.edgierThan(a)
			}
		}
	}
	return false
}

// Compare normal version string to see if 'a' is newer than 'b'
// If normal version is entirely the same, compare pre-release strings
func (a Semver) NewerThan(b Semver) bool {
	if a.major > b.major {
		return true
	} else if a.major == b.major {
		if a.minor > b.minor {
			return true
		} else if a.minor == b.minor {
			if a.patch > b.patch {
				return true
			} else if a.patch == b.patch {
				return a.edgierThan(b)
			}
		}
	}
	return false
}

// If version and pre-release strings are equivalent, returns true.
func (a Semver) EquivalentTo(b Semver) bool {
	return a.OlderThan(b) == a.NewerThan(b)
}

// Determines if a's pre-release string is higher precedence than b's
func (a Semver) edgierThan(b Semver) bool {
	if len(a.pre.tag) == 0 || len(b.pre.tag) == 0 {
		// Not having a pre-release string signifies precedence
		return a.pre.tag < b.pre.tag
	}
	ed := strings.Split(a.pre.tag, ".")
	gy := strings.Split(b.pre.tag, ".")
	for key := range ed {
		if len(gy) < key+1 {
			return true
		}
		if rxNumeric.MatchString(ed[key]) && rxNumeric.MatchString(gy[key]) {
			if ed[key] != gy[key] {
				left, _ := strconv.ParseInt(ed[key], 10, 0)
				right, _ := strconv.ParseInt(gy[key], 10, 0)
				return left > right
			}
		}
		if ed[key] != gy[key] {
			return ed[key] > gy[key]
		}
	}
	return false
}

// Normal version components can be bumped simply
func (a *Semver) incrementVersion(enum string) {
	switch enum {
	case "major":
		a.major += 1
		a.minor, a.patch = 0, 0
	case "minor":
		a.minor += 1
		a.patch = 0
	case "patch":
		a.patch += 1
	}

}

/* Public API for interacting with semver strings
 * If you call anything other than these, you're taking ownership of the results
 */

func IsValid(version string) bool {
	return rxMatch.MatchString(version)
}

func Increment(version, enum string) string {
	if IsValid(version) {
		a, _ := ConStructor(version)
		a.incrementVersion(enum)
		return a.ConvertToString()
	}
	return "Invalid Version"
}

// Returns true if newer, false if not OR if either input isn't a valid semver
func IsNewer(s, v string) bool {
	if IsValid(s) && IsValid(v) {
		a, _ := ConStructor(s)
		b, _ := ConStructor(v)
		return a.NewerThan(*b)
	}
	return false
}
