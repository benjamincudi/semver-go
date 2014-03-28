package semver

import (
	"testing"
)

func TestConStructor(t *testing.T) {
	var ver1 = "3.6.1"
	v1 := ConStructor(ver1)
	var result = (v1.major == "3" && v1.minor == "6" && v1.patch == "1")
	if !result {
		t.Errorf("Constructor failed on %s, returned %+v", ver1, v1)
	}

	var ver2 = "1.0.0-alpha.3"
	v2 := ConStructor(ver2)
	result = (v2.pre == "alpha.3")
	if !result {
		t.Errorf("Constructor failed on %s, returned %+v", ver2, v2)
	}

	var ver3 = "8.7.9+exp.sha.111334"
	v3 := ConStructor(ver3)
	result = (v3.build == "exp.sha.111334")
	if !result {
		t.Errorf("Constructor failed on %s, returned %+v", ver3, v3)
	}

	var ver4 = "1.3.0-rc.1+exp.sha.5114f85"
	v4 := ConStructor(ver4)
	result = (v4.major == "1" && v4.minor == "3" && v4.patch == "0" && v4.pre == "rc.1" && v4.build == "exp.sha.5114f85")
	if !result {
		t.Errorf("Constructor failed on %s, returned %+v", ver4, v4)
	}

}

func TestConvertToString(t *testing.T) {
	var ver1 = "3.6.1"
	v1 := ConStructor(ver1)
	var result = (ver1 == v1.ConvertToString())
	if !result {
		t.Errorf("String Conversion failed on %s, returned %+v", ver1, v1.ConvertToString())
	}

	var ver2 = "1.0.0-alpha.3"
	v2 := ConStructor(ver2)
	result = (ver2 == v2.ConvertToString())
	if !result {
		t.Errorf("Constructor failed on %s, returned %+v", ver2, v2.ConvertToString())
	}

	var ver3 = "8.7.9+exp.sha.111334"
	v3 := ConStructor(ver3)
	result = (ver3 == v3.ConvertToString())
	if !result {
		t.Errorf("Constructor failed on %s, returned %+v", ver3, v3.ConvertToString())
	}

	var ver4 = "1.3.0-rc.1+exp.sha.5114f85"
	v4 := ConStructor(ver4)
	result = (ver4 == v4.ConvertToString())
	if !result {
		t.Errorf("Constructor failed on %s, returned %+v", ver4, v4.ConvertToString())
	}
}

func TestComparisons(t *testing.T) {
	ver1, ver2, ver3, ver4, ver5 := "1.3.0", "1.3.2", "1.4.0", "2.0.1", "2.0.1+build.125124"
	v1 := ConStructor(ver1)
	v2 := ConStructor(ver2)
	v3 := ConStructor(ver3)
	v4 := ConStructor(ver4)
	v5 := ConStructor(ver5)
	if v1.NewerThan(*v2) {
		t.Errorf("%s is not newer than %s", ver1, ver2)
	}
	if v5.OlderThan(*v1) {
		t.Errorf("%s is not older than %s", ver5, ver1)
	}
	if v3.OlderThan(*v2) {
		t.Errorf("%s is not older than %s", ver3, ver2)
	}
	if !v4.EquivalentTo(*v5) {
		t.Errorf("%s is actually equivalent to %s", ver4, ver5)
	}
}

/* Semantic Versioning 2.0.0
 * http://semver.org/
 *
 * The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT",
 * "RECOMMENDED", "MAY", and "OPTIONAL" in this document are to be interpreted as described in
 * RFC 2119(http://tools.ietf.org/html/rfc2119).
 *
 *	1) 	Software using Semantic Versioning MUST declare a public API. This API could be
 *		declared in the code itself or exist strictly in documentation. However it is done,
 *		it should be precise and comprehensive.
 *
 *	2) 	A normal version number MUST take the form X.Y.Z where X, Y, and Z are non-negative
 *		integers, and MUST NOT contain leading zeroes.
 *		X is the major version, Y is the minor version, and Z is the patch version.
 *		Each element MUST increase numerically. For instance: 1.9.0 -> 1.10.0 -> 1.11.0.
 */
func TestIsValid(t *testing.T) {
	var validVersions = [...]string{"1.9.0", "0.9.0", "435.2458.518334", "0.0.0"}
	var invalidVersions = [...]string{"1.09.0", "1.a.0", "1.2.3b", "a.0.5", "5.4.2.10", "1.9", "1.9.", "1..", "1.-9.0", "-1.9.0", "1.9.-25"}
	for _, version := range validVersions {
		result := IsValid(version)
		if !result {
			t.Errorf("Should return true for %s, was false", version)
		}
	}
	for _, version := range invalidVersions {
		result := IsValid(version)
		if result {
			t.Errorf("Should return false for %s, was true", version)
		}
	}
}

/*	3) 	Once a versioned package has been released, the contents of that version MUST NOT be
 * 		modified. Any modifications MUST be released as a new version.
 *
 *	4) 	Major version zero (0.y.z) is for initial development. Anything may change at any time.
 *		The public API should not be considered stable.
 *
 *	5) 	Version 1.0.0 defines the public API. The way in which the version number is incremented
 *		after this release is dependent on this public API and how it changes.
 */

func TestIncrement(t *testing.T) {
	var ver = "1.3.9"

	/*	6) 	Patch version Z (x.y.Z | x > 0) MUST be incremented if only backwards compatible bug fixes
	 *		are introduced. A bug fix is defined as an internal change that fixes incorrect behavior.
	 */
	ver = Increment(ver, PATCH)
	if ver != "1.3.10" {
		t.Errorf("New version should be 1.3.10, was: ", ver)
	}
	/*	7) 	Minor version Y (x.Y.z | x > 0) MUST be incremented if new, backwards compatible
	 * 	functionality is introduced to the public API.
	 * 	It MUST be incremented if any public API functionality is marked as deprecated.
	 *		It MAY be incremented if substantial new functionality or improvements are introduced
	 *		within the private code.
	 *		It MAY include patch level changes.
	 *		Patch version MUST be reset to 0 when minor version is incremented.
	 */
	ver = Increment(ver, PATCH)
	ver = Increment(ver, MINOR)
	if ver != "1.4.0" {
		t.Errorf("New version should be 1.4.0, was: ", ver)
	}
	/*	8)	Major version X (X.y.z | X > 0) MUST be incremented if any backwards incompatible changes
	 *	 	are introduced to the public API.
	 *		It MAY include minor and patch level changes.
	 *		Patch and minor version MUST be reset to 0 when major version is incremented.
	 */
	ver = Increment(ver, PATCH)
	ver = Increment(ver, MINOR)
	ver = Increment(ver, PATCH)
	ver = Increment(ver, MAJOR)
	if ver != "2.0.0" {
		t.Errorf("New version should be 2.0.0, was: ", ver)
	}
}

//assertEquals('10.0.0', IncrementMajor('9.0.0'));
//assertEquals('1.10.0', IncrementMinor('1.9.0'));
//assertEquals('1.10.100', IncrementPatch('1.10.99'));
//assertEquals('1.0.1', IncrementPatch('1.0.0'));
//assertEquals('1.0.10', IncrementPatch('1.0.9'));
//assertEquals('1.0.100', IncrementPatch('1.0.99'));
//assertEquals('1.0.1000', IncrementPatch('1.0.999'));

//assertEquals('1.10.0', IncrementMinor('1.9.0'));
//assertEquals('1.100.0', IncrementMinor('1.99.5'));
//assertEquals('1.1000.0', IncrementMinor('1.999.100'));

//assertEquals('10.0.0', IncrementMajor('9.0.0'));
//assertEquals('10.0.0', IncrementMajor('9.1.0'));
//assertEquals('10.0.0', IncrementMajor('9.0.10'));
//assertEquals('10.0.0', IncrementMajor('9.8.534'));

/*	9) 	A pre-release version MAY be denoted by appending a hyphen and a series of dot separated
 *	 	identifiers immediately following the patch version.
 *		Identifiers MUST comprise only ASCII alphanumerics and hyphen [0-9A-Za-z-].
 *		Identifiers MUST NOT be empty.
 *		Numeric identifiers MUST NOT include leading zeroes.
 *		Pre-release versions have a lower precedence than the associated normal version.
 *		A pre-release version indicates that the version is unstable and might not satisfy the
 *		intended compatibility requirements as denoted by its associated normal version.
 *		Examples: 1.0.0-alpha, 1.0.0-alpha.1, 1.0.0-0.3.7, 1.0.0-x.7.z.92.
 *
 *	10) Build metadata MAY be denoted by appending a plus sign and a series of dot separated
 * 	identifiers immediately following the patch or pre-release version.
 *		Identifiers MUST comprise only ASCII alphanumerics and hyphen [0-9A-Za-z-].
 *		Identifiers MUST NOT be empty.
 *		Build metadata SHOULD be ignored when determining version precedence.
 *		Thus two versions that differ only in the build metadata, have the same precedence.
 *		Examples: 1.0.0-alpha+001, 1.0.0+20130313144700, 1.0.0-beta+exp.sha.5114f85.
 */

func TestValidPostfix(t *testing.T) {
	var validVersions = [...]string{"1.0.0-alpha", "1.0.0-alpha.1", "1.0.0-0.3.7", "1.0.0-x.7.z.92"}
	var result = false
	for _, version := range validVersions {
		result = IsValid(version)
		if !result {
			t.Errorf("Version %s should be valid, returned false", version)
		}
	}

	var ver = "1.0.0-01a"
	result = IsValid(ver)
	if !result {
		t.Errorf("Version %s is valid, returned false", ver)
	}
	ver = "1.0.0-01"
	result = IsValid(ver)
	if result {
		t.Errorf("Version %s is invalid, returned true", ver)
	}

	var invalidVersions = [...]string{"1.0.0-beta.02.1", "1.0.0-beta.$.13", "1.0.0-beta.2..1"}
	for _, version := range invalidVersions {
		result = IsValid(version)
		if result {
			t.Errorf("Version %s is invalid, but was reported as valid", version)
		}
	}

	validVersions = [...]string{"1.0.0-alpha+001", "1.0.0+20130313144700", "1.0.0-beta+exp.sha.5114f85", "1.0.0-alpha+2014-03"}
	for _, version := range validVersions {
		result = IsValid(version)
		if !result {
			t.Errorf("Version %s should be valid, returned false", version)
		}
	}
}

/*	11) Precedence refers to how versions are compared to each other when ordered.
 *		Precedence MUST be calculated by separating the version into major, minor, patch and
 *		pre-release identifiers in that order (Build metadata does not figure into precedence).
 *		Precedence is determined by the first difference when comparing each of these identifiers
 *		from left to right as follows:
 *
 *		Major, minor, and patch versions are always compared numerically.
 *			Example: 1.0.0 < 2.0.0 < 2.1.0 < 2.1.1.
 *		When major, minor, and patch are equal, a pre-release version has lower precedence than a
 *		normal version.
 *			Example: 1.0.0-alpha < 1.0.0.
 *
 *		Precedence for two pre-release versions with the same major, minor, and patch version MUST
 *		be determined by comparing each dot	separated identifier from left to right until a
 *		difference is found as follows:
 *
 *		Identifiers consisting of only digits are compared numerically and identifiers with letters
 *		or hyphens are compared lexically in ASCII sort order.
 *		Numeric identifiers always have lower precedence than non-numeric identifiers.
 *		A larger set of pre-release fields has a higher precedence than a smaller set, if all of the
 *		preceding identifiers are equal.
 *			Example: 1.0.0-alpha < 1.0.0-alpha.1 < 1.0.0-alpha.beta < 1.0.0-beta < 1.0.0-beta.2 < 1.0.0-beta.11 < 1.0.0-rc.1 < 1.0.0.
 */

func TestIsNewer(t *testing.T) {
	var ver = "3.1.2"
	var result = IsNewer(ver, "3.1.1")
	if !result {
		t.Errorf("Version %s is newer than 3.1.1, but was reported otherwise", ver)
	}
	result = IsNewer("3.2.0", ver)
	if !result {
		t.Errorf("Version 3.2.0 is newer than %s, but was reported otherwise", ver)
	}
	result = IsNewer("4.0.0", ver)
	if !result {
		t.Errorf("Version 4.0.0 is newer than %s, but was reported otherwise", ver)
	}
	result = IsNewer(ver, "3.1.2-alpha")
	if !result {
		t.Errorf("3.1.2-alpha is a prerelease, %s should report as newer", ver)
	}

	var prereleaseVersions = [...]string{"3.1.2-alpha", "3.1.2-alpha.1", "3.1.2-alpha.beta", "3.1.2-beta", "3.1.2-beta.2", "3.1.2-beta.11", "3.1.2-rc.1"}
	var count = len(prereleaseVersions) - 1
	for key, version := range prereleaseVersions {
		if key < count {
			result = IsNewer(prereleaseVersions[key+1], version)
			if !result {
				t.Errorf("Version %s is newer than %s, but was reported otherwise", prereleaseVersions[key+1], version)
			}
		} else {
			result = IsNewer(ver, version)
			if !result {
				t.Errorf("Version %s is newer than %s, but was reported otherwise", ver, version)
			}
		}
	}
}
