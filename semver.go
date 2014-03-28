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
var Version = "0.3.0"

/* Struct for semver string comprehension and manipulation.
 * This type and the methods associated are meant only for internal use,
 * and they have been written only with the intention of making the
 * API work easier to comprehend.
 */

type Semver struct {
	major, minor, patch, pre, build string
}

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
var rxNumeric, _ = regexp.Compile("^(0|[1-9])+$")

// Return in the same format as provided, when applicable
func (ver Semver) ConvertToString() string {
	version := strings.Join([]string{ver.major, ver.minor, ver.patch}, ".")
	if len(ver.pre) > 0 {
		version = strings.Join([]string{version, ver.pre}, "-")
	}
	if len(ver.build) > 0 {
		version = strings.Join([]string{version, ver.build}, "+")
	}
	return version
}

// Puntastic function to make a struct from a version string
// This makes it easier to deal with various parts
func ConStructor(version string) *Semver {
	var ver, bld, rel, maj, min, pat string
	ver = version
	if strings.Index(version, "+") > -1 {
		ver, bld = extractor(version, "+")
	}
	if strings.Index(ver, "-") > -1 {
		ver, rel = extractor(ver, "-")
	}
	maj, ver = extractor(ver, ".")
	min, pat = extractor(ver, ".")
	return &Semver{major: maj, minor: min, patch: pat, pre: rel, build: bld}
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
	}
	if a.major == b.major {
		if a.minor < b.minor {
			return true
		}
		if a.minor == b.minor {
			if a.patch < b.patch {
				return true
			}
			if a.patch == b.patch {
				// Compare pre-release
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
	}
	if a.major == b.major {
		if a.minor > b.minor {
			return true
		}
		if a.minor == b.minor {
			if a.patch > b.patch {
				return true
			}
			if a.patch == b.patch {
				return a.edgierThan(b)
			}
		}
	}
	return false
}

// If version and pre-release strings are equivalent, returns true.
func (a Semver) EquivalentTo(b Semver) bool {
	if a.OlderThan(b) || a.NewerThan(b) {
		return false
	}
	return true
}

// Determines if a's pre-release string is higher precedence than b's
func (a Semver) edgierThan(b Semver) bool {
	if len(a.pre) == 0 || len(b.pre) == 0 {
		// Not having a pre-release string signifies precedence
		return a.pre < b.pre
	}
	ed := strings.Split(a.pre, ".")
	gy := strings.Split(b.pre, ".")
	for key := range ed {
		if len(gy) < key+1 {
			return true
		}
		if rxNumeric.MatchString(ed[key]) && rxNumeric.MatchString(gy[key]) {
			if ed[key] != gy[key] {
				if len(ed[key]) > len(gy[key]) {
					return true
				} else {
					if len(ed[key]) < len(gy[key]) {
						return false
					}
					for k := range ed[key] {
						if (ed[key])[k:k+1] != (gy[key])[k:k+1] {
							return (ed[key])[k:k+1] > (gy[key])[k:k+1]
						}
					}
				}
			}
		}
		if ed[key] != gy[key] {
			return ed[key] > gy[key]
		}
	}
	return false
}

/* Public API for interacting with semver strings
 * If you call anything other than these, you're taking ownership of the results
 */

func IsValid(version string) bool {
	return rxMatch.MatchString(version)
}

func Increment(version, bump string) string {
	return ""
}

// Returns true if newer, false if not OR if either input isn't a valid semver
func IsNewer(check, base string) bool {
	if IsValid(check) && IsValid(base) {
		a := ConStructor(check)
		b := ConStructor(base)
		return a.NewerThan(*b)
	}
	return false
}
