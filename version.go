// Copyright IBM Corp. 2019, 2026
// SPDX-License-Identifier: MPL-2.0

// Package version compares version-like strings that may or may not obey
// semantic-versioning rules.
//
// It is built for the messy, real-world version identifiers returned by cloud
// APIs (for example, AWS RDS database engine versions) where strict semantic
// versioning parsers either error out or produce surprising results. Comparison
// never returns an error: when a value cannot be parsed as a semantic version,
// the package falls back to progressively looser strategies so that any two
// strings can always be ordered.
//
// Precedence of comparison, highest to lowest:
//
//  1. Associated timestamps (only in the WithTime variants)
//  2. Semantic versioning, via github.com/hashicorp/go-version
//  3. A "guerrilla" comparison for version-like strings that break semver rules
//  4. A plain lexicographic string comparison
package version

import (
	"cmp"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
)

// re splits a dot-separated part into runs of letters and runs of digits so
// that, for example, "ru-2023-10" compares sensibly against "ru-2024-01".
var re = regexp.MustCompile(`([a-zA-Z]+|\d+)`)

// Compare compares two version strings.
//
// It returns -1 if v1 sorts before v2, 0 if they are considered equal, and +1
// if v1 sorts after v2. Comparison never errors: values that are not valid
// semantic versions fall back to the guerrilla and string strategies described
// in the package documentation.
func Compare(v1, v2 string) int {
	a, err := version.NewVersion(v1)
	if err != nil {
		return compareStringsGuerrilla(v1, v2)
	}

	b, err := version.NewVersion(v2)
	if err != nil {
		return compareStringsGuerrilla(v1, v2)
	}

	return a.Compare(b)
}

// CompareWithTime compares two versions, giving precedence to their associated
// timestamps (such as create times) when both are present, non-zero, and
// distinct. Otherwise it falls back to [Compare].
//
// It returns -1, 0, or +1, following the same convention as [Compare].
func CompareWithTime(v1CreateTime, v2CreateTime *time.Time, v1, v2 string) int {
	var zero time.Time
	if v1CreateTime != nil && v2CreateTime != nil && !v1CreateTime.Equal(zero) && !v2CreateTime.Equal(zero) && !v1CreateTime.Equal(*v2CreateTime) {
		if v1CreateTime.Before(*v2CreateTime) {
			return -1
		}
		return 1
	}

	return Compare(v1, v2)
}

// LessThan reports whether v1 sorts before v2. See the package documentation
// for the comparison precedence.
func LessThan(v1, v2 string) bool {
	return Compare(v1, v2) < 0
}

// LessThanWithTime reports whether v1 sorts before v2, giving precedence to the
// associated timestamps when both are present, non-zero, and distinct. See the
// package documentation for the comparison precedence.
func LessThanWithTime(v1CreateTime, v2CreateTime *time.Time, v1, v2 string) bool {
	return CompareWithTime(v1CreateTime, v2CreateTime, v1, v2) < 0
}

// GreaterThan reports whether v1 sorts after v2. See the package documentation
// for the comparison precedence.
func GreaterThan(v1, v2 string) bool {
	return Compare(v1, v2) > 0
}

// GreaterThanWithTime reports whether v1 sorts after v2, giving precedence to
// the associated timestamps when both are present, non-zero, and distinct.
func GreaterThanWithTime(v1CreateTime, v2CreateTime *time.Time, v1, v2 string) bool {
	return CompareWithTime(v1CreateTime, v2CreateTime, v1, v2) > 0
}

// Equal reports whether v1 and v2 are considered equal. Note that semantic
// versions are compared semantically, so "1.0" and "1.0.0" are equal.
func Equal(v1, v2 string) bool {
	return Compare(v1, v2) == 0
}

// EqualWithTime reports whether v1 and v2 are considered equal, giving
// precedence to the associated timestamps when both are present, non-zero, and
// distinct.
func EqualWithTime(v1CreateTime, v2CreateTime *time.Time, v1, v2 string) bool {
	return CompareWithTime(v1CreateTime, v2CreateTime, v1, v2) == 0
}

// compareStringsGuerrilla compares two version-like strings that are not valid
// semantic versions. It walks the dot-separated parts, comparing numerically
// where both parts are integers and falling back to a sub-part comparison
// otherwise, then finally to a plain string comparison.
func compareStringsGuerrilla(v1, v2 string) int {
	if v1 == v2 { // save some time if they are equal
		return 0
	}

	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	for i := range min(len(parts1), len(parts2)) {
		num1, err1 := strconv.Atoi(parts1[i])
		num2, err2 := strconv.Atoi(parts2[i])

		if (err1 != nil || err2 != nil) && parts1[i] != parts2[i] {
			if c := compareSubparts(parts1[i], parts2[i]); c != 0 {
				return c
			}

			continue
		}

		if num1 != num2 {
			return cmp.Compare(num1, num2)
		}
	}

	return strings.Compare(v1, v2)
}

// compareSubparts compares two dot-separated parts by splitting each into runs
// of letters and digits, comparing numerically where both sub-parts are
// integers and lexicographically otherwise.
func compareSubparts(p1, p2 string) int {
	subp1 := re.FindAllString(p1, -1)
	subp2 := re.FindAllString(p2, -1)

	for i := range min(len(subp1), len(subp2)) {
		num1, err1 := strconv.Atoi(subp1[i])
		num2, err2 := strconv.Atoi(subp2[i])

		if (err1 != nil || err2 != nil) && subp1[i] != subp2[i] {
			if c := strings.Compare(subp1[i], subp2[i]); c != 0 {
				return c
			}
		}

		if num1 != num2 {
			return cmp.Compare(num1, num2)
		}
	}

	return 0
}
