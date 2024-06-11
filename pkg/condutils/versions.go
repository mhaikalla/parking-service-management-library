package condutils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	// MajorVersionParseError error when parse major version from string to int.
	MajorVersionParseError = `major version contains non numeric, got an error %v`

	// MinorVersionParseError error when parse minor version from string to int.
	MinorVersionParseError = `minor version contains non numeric, got an error %v`

	// PatchVersionsParseError error when parse patch version from string to int.
	PatchVersionsParseError = `patch version contains non numeric, got an error %v`

	// SemanticVersionLayoutError error when semantic version layout isn't <major>.<minor>.<patch>.
	SemanticVersionLayoutError = `cannot parse semantic version, semantic version layout must be <major>.<minor>.<patch> with numeric string`
)

// SemanticVersion is hold attribute for semantic versioning
type SemanticVersion struct {
	// major mean have breaking change
	major int

	// minor mean have backward compability
	minor int

	// patch mean have patch or hot fix for issues
	patch int

	// versionLong is semantic version join not using dot
	versionLong int
}

// NewSemanticVersion create a new Semantic Versioning using string version.
// `version` is string with pattern <major>.<minor>.<patch> which is only support numeric string
// return error if contains alphabet.
func NewSemanticVersion(version string) (SemanticVersion, error) {
	res := SemanticVersion{}
	splitedVersion := strings.Split(version, ".")

	if len(splitedVersion) < 3 {
		return res, errors.New(SemanticVersionLayoutError)
	}

	majorStr := splitedVersion[0]
	minorStr := splitedVersion[1]
	patchStr := splitedVersion[2]

	majorVer, majorVerErr := strconv.Atoi(majorStr)
	if majorVerErr != nil {
		return res, fmt.Errorf(MajorVersionParseError, majorVerErr)
	}

	minorVer, minorVerErr := strconv.Atoi(minorStr)
	if minorVerErr != nil {
		return res, fmt.Errorf(MinorVersionParseError, minorVerErr)
	}

	patchVer, patchVerErr := strconv.Atoi(patchStr)
	if patchVerErr != nil {
		return res, fmt.Errorf(PatchVersionsParseError, patchVerErr)
	}

	versionLong, versionLongErr := strconv.Atoi(majorStr + minorStr + patchStr)
	if versionLongErr != nil {
		return res, errors.New(SemanticVersionLayoutError)
	}

	res.major = majorVer
	res.minor = minorVer
	res.patch = patchVer
	res.versionLong = versionLong

	return res, nil
}

// After check if this version after `ver`
func (sm *SemanticVersion) After(ver SemanticVersion) bool {
	if sm.major > ver.major {
		return true
	}

	if sm.major < ver.major {
		return false
	}

	if sm.minor > ver.minor {
		return true
	}

	if sm.minor < ver.minor {
		return false
	}

	if sm.patch > ver.patch {
		return true
	}

	return false
}

// Before check if this version berfore `ver`
func (sm *SemanticVersion) Before(ver SemanticVersion) bool {
	if sm.major < ver.major {
		return true
	}

	if sm.major > ver.major {
		return false
	}

	if sm.minor < ver.minor {
		return true
	}

	if sm.minor > ver.minor {
		return false
	}

	if sm.patch < ver.patch {
		return true
	}

	return false
}

// Major get this major version
func (sm *SemanticVersion) Major() int {
	return sm.major
}

// Minor get this Minor version
func (sm *SemanticVersion) Minor() int {
	return sm.minor
}

// Patch get this patch version
func (sm *SemanticVersion) Patch() int {
	return sm.patch
}

// ToString return string of this semantic version
func (sm *SemanticVersion) ToString() string {
	return fmt.Sprintf("%v.%v.%v", sm.major, sm.minor, sm.patch)
}
