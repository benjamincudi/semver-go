package semver

import (
	"regexp"
	"strings"
)

/* A minimalist interface for semantic versioning, based on http://semver.org
 * This is based on Semver v2.0.0, and will adhere to those standards.
 * Updates will be documented at the bottom of readme.md,
 */

// This is the current version of this package
var Version = "0.1.0"

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

// Return in the same format as provided, when applicable
func (ver Semver) convertToString() string {
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
func conStructor(version string) *Semver {
	var ver, bld, rel, maj, min, pat string
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

/* Public API for interacting with semver strings
 * If you call anything other than these, you're taking ownership of the results
 */

func IsValid(version string) bool {
	return rxMatch.MatchString(version)
}

func Increment(version, bump string) string {
	return ""
}

func IsNewer(newer, current string) bool {
	return true
}
