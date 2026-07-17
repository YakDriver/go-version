// Copyright IBM Corp. 2019, 2026
// SPDX-License-Identifier: MPL-2.0

package version

import (
	"slices"
	"testing"
	"time"
)

func TestLessThanOrEqual(t *testing.T) {
	t.Parallel()

	type testCase struct {
		v1   string
		v2   string
		want bool
	}
	tests := map[string]testCase{
		"less":                 {v1: "10.11.4", v2: "10.11.6", want: true},
		"equal":                {v1: "1.0.0", v2: "1.0.0", want: true},
		"equal normalized":     {v1: "1.0", v2: "1.0.0", want: true},
		"greater":              {v1: "11.16", v2: "11.5", want: false},
		"non-semantic equal":   {v1: "oracle-ee-19", v2: "oracle-ee-19", want: true},
		"non-semantic less":    {v1: "oracle-ee-9", v2: "oracle-ee-19", want: true},
		"non-semantic greater": {v1: "oracle-ee-19", v2: "oracle-ee-9", want: false},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := LessThanOrEqual(test.v1, test.v2); got != test.want {
				t.Errorf("LessThanOrEqual(%q, %q) = %t, want %t", test.v1, test.v2, got, test.want)
			}
		})
	}
}

func TestGreaterThanOrEqual(t *testing.T) {
	t.Parallel()

	type testCase struct {
		v1   string
		v2   string
		want bool
	}
	tests := map[string]testCase{
		"greater":              {v1: "11.16", v2: "11.5", want: true},
		"equal":                {v1: "1.0.0", v2: "1.0.0", want: true},
		"equal normalized":     {v1: "1.0", v2: "1.0.0", want: true},
		"less":                 {v1: "10.4", v2: "10.4.27", want: false},
		"non-semantic equal":   {v1: "oracle-ee-19", v2: "oracle-ee-19", want: true},
		"non-semantic greater": {v1: "oracle-ee-19", v2: "oracle-ee-9", want: true},
		"non-semantic less":    {v1: "oracle-ee-9", v2: "oracle-ee-19", want: false},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := GreaterThanOrEqual(test.v1, test.v2); got != test.want {
				t.Errorf("GreaterThanOrEqual(%q, %q) = %t, want %t", test.v1, test.v2, got, test.want)
			}
		})
	}
}

func TestOrEqualWithTime(t *testing.T) {
	t.Parallel()

	time1 := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	time2 := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)

	// Distinct times take precedence over the version strings.
	if !LessThanOrEqualWithTime(&time1, &time2, "9.9.9", "1.0.0") {
		t.Error("LessThanOrEqualWithTime should honor the earlier time regardless of version")
	}
	if !GreaterThanOrEqualWithTime(&time2, &time1, "1.0.0", "9.9.9") {
		t.Error("GreaterThanOrEqualWithTime should honor the later time regardless of version")
	}

	// Nil times fall back to version comparison, including equality.
	if !LessThanOrEqualWithTime(nil, nil, "1.0", "1.0.0") {
		t.Error("LessThanOrEqualWithTime should be true for semantically equal versions")
	}
	if !GreaterThanOrEqualWithTime(nil, nil, "1.0", "1.0.0") {
		t.Error("GreaterThanOrEqualWithTime should be true for semantically equal versions")
	}
}

func TestMax(t *testing.T) {
	t.Parallel()

	type testCase struct {
		versions []string
		want     string
	}
	tests := map[string]testCase{
		"empty":             {versions: nil, want: ""},
		"single":            {versions: []string{"1.2.3"}, want: "1.2.3"},
		"semantic":          {versions: []string{"8.0.35", "5.7.44", "8.0", "5.7"}, want: "8.0.35"},
		"non-semantic":      {versions: []string{"oracle-ee-9", "oracle-ee-19", "oracle-ee-18"}, want: "oracle-ee-19"},
		"mixed sql server":  {versions: []string{"12.00.6439.10.v1", "14.00.3451.2.v1", "13.00.6435.1.v1"}, want: "14.00.3451.2.v1"},
		"equal keeps first": {versions: []string{"1.0", "1.0.0"}, want: "1.0"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := Max(test.versions...); got != test.want {
				t.Errorf("Max(%v) = %q, want %q", test.versions, got, test.want)
			}
		})
	}
}

func TestMin(t *testing.T) {
	t.Parallel()

	type testCase struct {
		versions []string
		want     string
	}
	tests := map[string]testCase{
		"empty":             {versions: nil, want: ""},
		"single":            {versions: []string{"1.2.3"}, want: "1.2.3"},
		"semantic":          {versions: []string{"8.0.35", "5.7.44", "8.0", "5.7"}, want: "5.7"},
		"non-semantic":      {versions: []string{"oracle-ee-9", "oracle-ee-19", "oracle-ee-18"}, want: "oracle-ee-9"},
		"equal keeps first": {versions: []string{"1.0.0", "1.0"}, want: "1.0.0"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := Min(test.versions...); got != test.want {
				t.Errorf("Min(%v) = %q, want %q", test.versions, got, test.want)
			}
		})
	}
}

func TestSort(t *testing.T) {
	t.Parallel()

	versions := []string{"8.0.35", "5.7.44", "8.0", "5.7"}
	Sort(versions)

	want := []string{"5.7", "5.7.44", "8.0", "8.0.35"}
	if !slices.Equal(versions, want) {
		t.Errorf("Sort = %v, want %v", versions, want)
	}
}

func TestSortNonSemantic(t *testing.T) {
	t.Parallel()

	versions := []string{
		"19.0.0.0.ru-2023-10.rur-2023-10.r10",
		"19.0.0.0.ru-2020-01.rur-2020-01.r1",
		"19.0.0.0.ru-2023-10.rur-2023-10.r9",
	}
	Sort(versions)

	want := []string{
		"19.0.0.0.ru-2020-01.rur-2020-01.r1",
		"19.0.0.0.ru-2023-10.rur-2023-10.r9",
		"19.0.0.0.ru-2023-10.rur-2023-10.r10",
	}
	if !slices.Equal(versions, want) {
		t.Errorf("Sort = %v, want %v", versions, want)
	}
}

func TestSortDescending(t *testing.T) {
	t.Parallel()

	versions := []string{"8.0", "5.7.44", "8.0.35", "5.7"}
	SortDescending(versions)

	want := []string{"8.0.35", "8.0", "5.7.44", "5.7"}
	if !slices.Equal(versions, want) {
		t.Errorf("SortDescending = %v, want %v", versions, want)
	}
}
