// Copyright IBM Corp. 2019, 2026
// SPDX-License-Identifier: MPL-2.0

package version

import (
	"slices"
	"testing"
	"time"
)

func TestLessThanWithTime(t *testing.T) {
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
		want         bool
	}
	tests := map[string]testCase{
		"nil create time 1": {
			v1CreateTime: nil,
			v2CreateTime: nil,
			v1:           semVer1,
			v2:           semVer2,
			want:         true,
		},
		"nil create time 2": {
			v1CreateTime: &time1,
			v2CreateTime: nil,
			v1:           semVer1,
			v2:           semVer2,
			want:         true,
		},
		"nil create time 3": {
			v1CreateTime: nil,
			v2CreateTime: &time1,
			v1:           semVer1,
			v2:           semVer2,
			want:         true,
		},
		"empty v1 v2": {
			v1CreateTime: nil,
			v2CreateTime: nil,
			v1:           "",
			v2:           "",
			want:         false,
		},
		"empty v1": {
			v1CreateTime: nil,
			v2CreateTime: nil,
			v1:           "",
			v2:           semVer2,
			want:         true,
		},
		"empty v2": {
			v1CreateTime: nil,
			v2CreateTime: nil,
			v1:           semVer1,
			v2:           "",
			want:         false,
		},
		"equal create time": {
			v1CreateTime: &time1,
			v2CreateTime: &time1,
			v1:           semVer1,
			v2:           semVer2,
			want:         true,
		},
		"zero create time 1": {
			v1CreateTime: &zeroTime,
			v2CreateTime: &time1,
			v1:           semVer1,
			v2:           semVer2,
			want:         true,
		},
		"zero create time 2": {
			v1CreateTime: &time1,
			v2CreateTime: &zeroTime,
			v1:           semVer1,
			v2:           semVer2,
			want:         true,
		},
		"create time 1": {
			v1CreateTime: &time1,
			v2CreateTime: &time2,
			v1:           semVer1,
			v2:           semVer2,
			want:         true,
		},
		"create time 2": {
			v1CreateTime: &time2,
			v2CreateTime: &time1,
			v1:           semVer1,
			v2:           semVer2,
			want:         false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := LessThanWithTime(test.v1CreateTime, test.v2CreateTime, test.v1, test.v2)

			if got != test.want {
				t.Errorf("unexpected result: %t, wanted: %t; [(%s vs %s), (%s vs %s)]", got, test.want, test.v1CreateTime, test.v2CreateTime, test.v1, test.v2)
			}
		})
	}
}

func TestLessThan(t *testing.T) {
	t.Parallel()

	type testCase struct {
		v1   string
		v2   string
		want bool
	}
	tests := map[string]testCase{
		"semantic 1": {
			v1:   "10.11.4",
			v2:   "10.11.6",
			want: true,
		},
		"semantic 2": {
			v1:   "10.11.4",
			v2:   "10.11.6",
			want: true,
		},
		"semantic 3": {
			v1:   "10.4",
			v2:   "10.4.27",
			want: true,
		},
		"semantic 4": {
			v1:   "10.4.27",
			v2:   "10.4.32",
			want: true,
		},
		"semantic 5": {
			v1:   "10.4.32",
			v2:   "10.5.16",
			want: true,
		},
		"semantic 6": {
			v1:   "10.5.16",
			v2:   "10.6.8",
			want: true,
		},
		"semantic 7": {
			v1:   "10.6.8",
			v2:   "11",
			want: true,
		},
		"semantic 8": {
			v1:   "11",
			v2:   "11.16",
			want: true,
		},
		"semantic 9": {
			v1:   "11.16",
			v2:   "11.5",
			want: false,
		},
		"semantic 10": {
			v1:   "11.5",
			v2:   "11.9",
			want: true,
		},
		"semantic 11": {
			v1:   "11.9",
			v2:   "12",
			want: true,
		},
		"semantic 14": {
			v1:   "19",
			v2:   "18",
			want: false,
		},
		"semantic 15": {
			v1:   "5.7.44",
			v2:   "8.0",
			want: true,
		},
		"semantic 16": {
			v1:   "8.0",
			v2:   "8.0.28",
			want: true,
		},
		"semantic 17": {
			v1:   "12.00",
			v2:   "12.17",
			want: true,
		},
		"semantic 18": {
			v1:   "13",
			v2:   "12.00",
			want: false,
		},
		"semantic and non 1": {
			v1:   "8.0.35",
			v2:   "11.5.9.0.sb00000000.r1",
			want: true,
		},
		"non-semantic 4": {
			v1:   "11.5.9.0.sb00000000.r1",
			v2:   "12.00.6439.10.v1",
			want: true,
		},
		"non-semantic 5": {
			v1:   "12.00.6439.10.v1",
			v2:   "12.00.6444.4.v1",
			want: true,
		},
		"non-semantic 6": {
			v1:   "12.00.6444.4.v1",
			v2:   "12.00.6449.1.v1",
			want: true,
		},
		"non-semantic 7": {
			v1:   "12.00.6449.1.v1",
			v2:   "13.00.6430.49.v1",
			want: true,
		},
		"non-semantic 8": {
			v1:   "13.00.6430.49.v1",
			v2:   "13.00.6435.1.v1",
			want: true,
		},
		"non-semantic 9": {
			v1:   "13.00.6435.1.v1",
			v2:   "14.00.3281.6.v1",
			want: true,
		},
		"non-semantic 10": {
			v1:   "14.00.3281.6.v1",
			v2:   "14.00.3294.2.v1",
			want: true,
		},
		"non-semantic 11": {
			v1:   "14.00.3294.2.v1",
			v2:   "14.00.3460.9.v1",
			want: true,
		},
		"non-semantic 12": {
			v1:   "14.00.3460.9.v1",
			v2:   "14.00.3465.1.v1",
			want: true,
		},
		"non-semantic 13": {
			v1:   "14.00.3465.1.v1",
			v2:   "15.00.4322.2.v1",
			want: true,
		},
		"non-semantic 14": {
			v1:   "15.00.4335.1.v9",
			v2:   "15.00.4335.1.v10",
			want: true,
		},
		"non-semantic 15": {
			v1:   "15.00.4335.1.v1",
			v2:   "15.00.4345.5.v1",
			want: true,
		},
		"non-semantic 16": {
			v1:   "15.00.4345.5.v1",
			v2:   "16.00.4085.2.v1",
			want: true,
		},
		"non-semantic 17": {
			v1:   "16.00.4085.2.v1",
			v2:   "16.00.4095.4.v1",
			want: true,
		},
		"non-semantic 18": {
			v1:   "16.00.4095.4.v1",
			v2:   "19.0.0.0.ru-2020-01.rur-2020-01.r1",
			want: true,
		},
		"non-semantic 19": {
			v1:   "19.0.0.0.ru-2020-01.rur-2020-01.r1",
			v2:   "19.0.0.0.ru-2020-04.rur-2020-04.r1",
			want: true,
		},
		"non-semantic 20": {
			v1:   "19.0.0.0.ru-2020-04.rur-2020-04.r1",
			v2:   "19.0.0.0.ru-2021-10.rur-2021-10.r1",
			want: true,
		},
		"non-semantic 21": {
			v1:   "19.0.0.0.ru-2021-10.rur-2021-10.r1",
			v2:   "19.0.0.0.ru-2022-01.rur-2022-01.r1",
			want: true,
		},
		"non-semantic 22": {
			v1:   "19.0.0.0.ru-2022-01.rur-2022-01.r1",
			v2:   "19.0.0.0.ru-2023-10.rur-2023-10.r1",
			want: true,
		},
		"non-semantic 23": {
			v1:   "19.0.0.0.ru-2023-10.rur-2023-10.r1",
			v2:   "19.0.0.0.ru-2024-01.rur-2024-01.r1",
			want: true,
		},
		"non-semantic 23a": {
			v1:   "19.0.0.0.ru-2023-10.rur-2023-10.r9",
			v2:   "19.0.0.0.ru-2023-10.rur-2023-10.r10",
			want: true,
		},
		"non-semantic 25": {
			v1:   "5.7.mysql_aurora.2.11.3",
			v2:   "5.7.mysql_aurora.2.12.1",
			want: true,
		},
		"non-semantic 26": {
			v1:   "5.7.mysql_aurora.2.12.1",
			v2:   "8.0.mysql_aurora.3.01.0",
			want: true,
		},
		"non-semantic 27": {
			v1:   "8.0.mysql_aurora.3.01.0",
			v2:   "8.0.mysql_aurora.3.05.2",
			want: true,
		},
		"non-semantic 29": {
			v1:   "aurora-mysql5.7",
			v2:   "aurora-mysql8.0",
			want: true,
		},
		"non-semantic 30": {
			v1:   "aurora-mysql8.0",
			v2:   "aurora-mysql15.1.1",
			want: true,
		},
		"non-semantic 31": {
			v1:   "db2-ae",
			v2:   "db2-se",
			want: true,
		},
		"non-semantic 33": {
			v1:   "MariaDB 10.4.27",
			v2:   "MariaDB 10.4.28",
			want: true,
		},
		"non-semantic 34": {
			v1:   "mariadb10.10",
			v2:   "mariadb10.4",
			want: false,
		},
		"non-semantic 36": {
			v1:   "Oracle 19.0.0.0.ru-2019-07.rur-2019-07.r1",
			v2:   "Oracle 19.0.0.0.ru-2020-04.rur-2020-04.r1",
			want: true,
		},
		"non-semantic 37": {
			v1:   "Oracle 19.0.0.0.ru-2020-04.rur-2020-04.r1",
			v2:   "Oracle 19.0.0.0.ru-2021-07.rur-2021-07.r1",
			want: true,
		},
		"non-semantic 38": {
			v1:   "Oracle 19.0.0.0.ru-2021-07.rur-2021-07.r1",
			v2:   "Oracle 19.0.0.0.ru-2021-10.rur-2021-10.r1",
			want: true,
		},
		"non-semantic 39": {
			v1:   "Oracle 19.0.0.0.ru-2021-10.rur-2021-10.r1",
			v2:   "Oracle 19.0.0.0.ru-2022-01.rur-2022-01.r1",
			want: true,
		},
		"non-semantic 40": {
			v1:   "Oracle 19.0.0.0.ru-2022-01.rur-2022-01.r1",
			v2:   "Oracle 19.0.0.0.ru-2023-10.rur-2023-10.r1",
			want: true,
		},
		"non-semantic 41": {
			v1:   "Oracle 19.0.0.0.ru-2023-10.rur-2023-10.r1",
			v2:   "Oracle 19.0.0.0.ru-2024-01.rur-2024-01.r1",
			want: true,
		},
		"non-semantic 42": {
			v1:   "oracle-ee-18",
			v2:   "oracle-ee-19",
			want: true,
		},
		"non-semantic 43": {
			v1:   "oracle-ee-19",
			v2:   "oracle-ee-9",
			want: false,
		},
		"non-semantic 44": {
			v1:   "SQL Server 2014 12.00.6293.0.v1",
			v2:   "SQL Server 2014 12.00.6329.1.v1",
			want: true,
		},
		"non-semantic 45": {
			v1:   "SQL Server 2014 12.00.6329.1.v1",
			v2:   "SQL Server 2014 12.00.6433.1.v1",
			want: true,
		},
		"non-semantic 46": {
			v1:   "SQL Server 2014 12.00.6433.1.v1",
			v2:   "SQL Server 2014 12.00.6439.10.v1",
			want: true,
		},
		"non-semantic 47": {
			v1:   "SQL Server 2014 12.00.6439.10.v1",
			v2:   "SQL Server 2017 14.00.3451.2.v1",
			want: true,
		},
		"non-semantic 49": {
			v1:   "sqlserver-ee-13.0",
			v2:   "sqlserver-ee-14.0",
			want: true,
		},
		"non-semantic 50 equal": {
			v1:   "16.00.4085.2.v1",
			v2:   "16.00.4085.2.v1",
			want: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := LessThan(test.v1, test.v2)

			if got != test.want {
				t.Errorf("unexpected result: %t, wanted: %t; [(%s vs %s)]", got, test.want, test.v1, test.v2)
			}
		})
	}
}

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

func TestSortStable(t *testing.T) {
	t.Parallel()

	// "1.0", "1.0.0", and "1.0.0.0" all compare equal; a stable sort must keep
	// their original relative order.
	versions := []string{"1.0.0", "2.0", "1.0", "1.0.0.0"}
	Sort(versions)

	want := []string{"1.0.0", "1.0", "1.0.0.0", "2.0"}
	if !slices.Equal(versions, want) {
		t.Errorf("Sort = %v, want %v (equal elements should keep input order)", versions, want)
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
