package version

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
)

var (
	re = regexp.MustCompile(`([a-zA-Z]+|\d+)`)
)

// CompareVersions compares two versions, returning true if v1 is less than v2.
// Precedence of comparison is:
//  1. CreateTime
//  2. go-version (semantic versioning)
//  3. Guerilla comparison (versions that do not conform to semantic versioning rules)
//  4. String comparison
func CompareVersions(v1CreateTime1, v2CreateTime2 *time.Time, v1, v2 string) bool {
	var zero time.Time
	if v1CreateTime1 != nil && v2CreateTime2 != nil && !v1CreateTime1.Equal(zero) && !v2CreateTime2.Equal(zero) && !v1CreateTime1.Equal(*v2CreateTime2) {
		return v1CreateTime1.Before(*v2CreateTime2)
	}

	return CompareVersionStrings(v1, v2)
}

// CompareVersionStrings compares two version strings, returning true if v1 is less than v2.
// Precedence of comparison is:
//  1. go-version (semantic versioning)
//  2. Guerilla comparison (versions that do not conform to semantic versioning rules)
//  3. String comparison
func CompareVersionStrings(v1, v2 string) bool {
	a, err := version.NewVersion(v1)
	if err != nil {
		return compareVersionStringsGuerrilla(v1, v2)
	}

	b, err := version.NewVersion(v2)
	if err != nil {
		return compareVersionStringsGuerrilla(v1, v2)
	}

	return a.LessThan(b)
}

func compareVersionStringsGuerrilla(v1, v2 string) bool {
	if v1 == v2 { // save some time if they are equal
		return false
	}

	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	for i := 0; i < len(parts1) && i < len(parts2); i++ {
		num1, err1 := strconv.Atoi(parts1[i])
		num2, err2 := strconv.Atoi(parts2[i])

		if (err1 != nil || err2 != nil) && parts1[i] != parts2[i] {
			// string comparison
			switch compareSubparts(parts1[i], parts2[i]) {
			case -1:
				return true
			case 1:
				return false
			}

			continue
		}

		if num1 != num2 {
			// number comparison
			return num1 < num2
		}
	}

	// string comparison
	return v1 < v2
}

func compareSubparts(p1, p2 string) int {
	subp1 := re.FindAllString(p1, -1)
	subp2 := re.FindAllString(p2, -1)

	for i := 0; i < len(subp1) && i < len(subp2); i++ {
		num1, err1 := strconv.Atoi(subp1[i])
		num2, err2 := strconv.Atoi(subp2[i])

		if (err1 != nil || err2 != nil) && subp1[i] != subp2[i] {
			// string comparison
			if subp1[i] < subp2[i] {
				return -1
			}

			if subp1[i] > subp2[i] {
				return 1
			}
		}

		if num1 != num2 {
			// number comparison
			if num1 < num2 {
				return -1
			}

			if num1 > num2 {
				return 1
			}
		}
	}

	return 0
}
