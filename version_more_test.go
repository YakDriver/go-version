// Copyright IBM Corp. 2019, 2026
// SPDX-License-Identifier: MPL-2.0

package version

import (
	"testing"
	"time"
)

func TestCompare(t *testing.T) {
	t.Parallel()

	type testCase struct {
		v1   string
		v2   string
		want int
	}
	tests := map[string]testCase{
		"semantic less":            {v1: "10.11.4", v2: "10.11.6", want: -1},
		"semantic greater":         {v1: "11.16", v2: "11.5", want: 1},
		"semantic equal":           {v1: "1.0.0", v2: "1.0.0", want: 0},
		"semantic equal normalize": {v1: "1.0", v2: "1.0.0", want: 0},
		"non-semantic less":        {v1: "14.00.3281.6.v1", v2: "14.00.3294.2.v1", want: -1},
		"non-semantic greater":     {v1: "mariadb10.10", v2: "mariadb10.4", want: 1},
		"non-semantic equal":       {v1: "16.00.4085.2.v1", v2: "16.00.4085.2.v1", want: 0},
		"oracle date less":         {v1: "19.0.0.0.ru-2023-10.rur-2023-10.r9", v2: "19.0.0.0.ru-2023-10.rur-2023-10.r10", want: -1},
		"string fallback":          {v1: "db2-ae", v2: "db2-se", want: -1},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := Compare(test.v1, test.v2); got != test.want {
				t.Errorf("Compare(%q, %q) = %d, want %d", test.v1, test.v2, got, test.want)
			}
		})
	}
}

func TestGreaterThan(t *testing.T) {
	t.Parallel()

	type testCase struct {
		v1   string
		v2   string
		want bool
	}
	tests := map[string]testCase{
		"semantic greater":     {v1: "11.16", v2: "11.5", want: true},
		"semantic less":        {v1: "10.4", v2: "10.4.27", want: false},
		"equal":                {v1: "1.0.0", v2: "1.0.0", want: false},
		"non-semantic greater": {v1: "oracle-ee-19", v2: "oracle-ee-9", want: true},
		"non-semantic less":    {v1: "aurora-mysql5.7", v2: "aurora-mysql8.0", want: false},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := GreaterThan(test.v1, test.v2); got != test.want {
				t.Errorf("GreaterThan(%q, %q) = %t, want %t", test.v1, test.v2, got, test.want)
			}
		})
	}
}

func TestEqual(t *testing.T) {
	t.Parallel()

	type testCase struct {
		v1   string
		v2   string
		want bool
	}
	tests := map[string]testCase{
		"identical":              {v1: "1.0.0", v2: "1.0.0", want: true},
		"semantic normalized":    {v1: "1.0", v2: "1.0.0", want: true},
		"different":              {v1: "1.0.0", v2: "1.0.1", want: false},
		"non-semantic identical": {v1: "16.00.4085.2.v1", v2: "16.00.4085.2.v1", want: true},
		"non-semantic different": {v1: "db2-ae", v2: "db2-se", want: false},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := Equal(test.v1, test.v2); got != test.want {
				t.Errorf("Equal(%q, %q) = %t, want %t", test.v1, test.v2, got, test.want)
			}
		})
	}
}

func TestCompareWithTime(t *testing.T) {
	t.Parallel()

	time1 := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	time2 := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	var zeroTime time.Time

	semVer1 := "1.0.0"
	semVer2 := "1.0.1"

	type testCase struct {
		v1CreateTime *time.Time
		v2CreateTime *time.Time
		v1           string
		v2           string
		want         int
	}
	tests := map[string]testCase{
		"distinct times less":     {v1CreateTime: &time1, v2CreateTime: &time2, v1: semVer2, v2: semVer1, want: -1},
		"distinct times greater":  {v1CreateTime: &time2, v2CreateTime: &time1, v1: semVer1, v2: semVer2, want: 1},
		"equal times use version": {v1CreateTime: &time1, v2CreateTime: &time1, v1: semVer1, v2: semVer2, want: -1},
		"nil time uses version":   {v1CreateTime: nil, v2CreateTime: &time1, v1: semVer2, v2: semVer1, want: 1},
		"zero time uses version":  {v1CreateTime: &zeroTime, v2CreateTime: &time1, v1: semVer1, v2: semVer1, want: 0},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := CompareWithTime(test.v1CreateTime, test.v2CreateTime, test.v1, test.v2); got != test.want {
				t.Errorf("CompareWithTime = %d, want %d", got, test.want)
			}
		})
	}
}

func TestGreaterAndEqualWithTime(t *testing.T) {
	t.Parallel()

	time1 := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	time2 := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)

	if !GreaterThanWithTime(&time2, &time1, "1.0.0", "1.0.0") {
		t.Error("GreaterThanWithTime with later time should be true")
	}
	if GreaterThanWithTime(&time1, &time2, "9.9.9", "1.0.0") {
		t.Error("GreaterThanWithTime with earlier time should be false regardless of version")
	}
	if !EqualWithTime(nil, nil, "1.0", "1.0.0") {
		t.Error("EqualWithTime should defer to semantic version equality when times are nil")
	}
	if EqualWithTime(&time1, &time2, "1.0.0", "1.0.0") {
		t.Error("EqualWithTime with distinct times should be false")
	}
}

// TestCompareConsistency verifies that Compare's sign is consistent with the
// boolean helpers derived from it across a representative set of pairs,
// including the reflexive and antisymmetric properties expected by callers
// using these functions as sort comparators.
func TestCompareConsistency(t *testing.T) {
	t.Parallel()

	pairs := [][2]string{
		{"10.11.4", "10.11.6"},
		{"8.0", "8.0.28"},
		{"11.16", "11.5"},
		{"16.00.4085.2.v1", "16.00.4085.2.v1"},
		{"5.7.mysql_aurora.2.11.3", "5.7.mysql_aurora.2.12.1"},
		{"19.0.0.0.ru-2023-10.rur-2023-10.r9", "19.0.0.0.ru-2023-10.rur-2023-10.r10"},
		{"db2-ae", "db2-se"},
		{"oracle-ee-19", "oracle-ee-9"},
	}

	for _, p := range pairs {
		v1, v2 := p[0], p[1]
		c := Compare(v1, v2)

		if got := LessThan(v1, v2); got != (c < 0) {
			t.Errorf("LessThan(%q, %q) = %t, inconsistent with Compare = %d", v1, v2, got, c)
		}
		if got := GreaterThan(v1, v2); got != (c > 0) {
			t.Errorf("GreaterThan(%q, %q) = %t, inconsistent with Compare = %d", v1, v2, got, c)
		}
		if got := Equal(v1, v2); got != (c == 0) {
			t.Errorf("Equal(%q, %q) = %t, inconsistent with Compare = %d", v1, v2, got, c)
		}

		// Antisymmetry: reversing the arguments negates the sign.
		if rev := Compare(v2, v1); rev != -c {
			t.Errorf("Compare(%q, %q) = %d, but Compare(%q, %q) = %d (want %d)", v1, v2, c, v2, v1, rev, -c)
		}

		// Reflexivity: a value compared with itself is equal.
		if self := Compare(v1, v1); self != 0 {
			t.Errorf("Compare(%q, %q) = %d, want 0", v1, v1, self)
		}
	}
}
