package version

import (
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
		name, test := name, test
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
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := LessThan(test.v1, test.v2)

			if got != test.want {
				t.Errorf("unexpected result: %t, wanted: %t; [(%s vs %s)]", got, test.want, test.v1, test.v2)
			}
		})
	}
}
