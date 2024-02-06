# go-version

**`go-version`** attempts to compare versions. It uses this precedence in comparing:

1. The earliest time associated with a version (only when using `LessThanWithTime()`)
2. Semantic-versioning parsing based on [go-version](https://github.com/hashicorp/go-version)
3. If the version breaks semantic versioning rules, a smart comparison looking at the version as a version-like thing
4. A plain string comparison

_Version comparison will not return errors in any situation._

`go-version` is optimized for the versions used by AWS RDS database engine versions. It will likely work with many other version-like situations.

# Usage examples


## `LessThan`

If you're just comparing version strings without associated times, such as the version create times, use `LessThan()`:

```go
package main

import (
    "time"

    "github.com/YakDriver/go-version"
)

func main() {
    // normal semantic-versioning versions
    fmt.Printf("%t\n", version.LessThan("10.11.9", "10.11.10")) // true
    fmt.Printf("%t\n", version.LessThan("10.4", "10.4.27")) // true
    fmt.Printf("%t\n", version.LessThan("10.6.8", "11")) // true
    fmt.Printf("%t\n", version.LessThan("1.2rc2", "1.2")) // true

    // non-semantic-versioning versions
    fmt.Printf("%t\n", version.LessThan("8.0.mysql_aurora.3.1.9", "8.0.mysql_aurora.3.1.10")) // false
    fmt.Printf("%t\n", version.LessThan("19.0.0.0.ru-2023-10.rur-2023-10.r9", "19.0.0.0.ru-2023-10.rur-2023-10.r10")) // true 
    fmt.Printf("%t\n", version.LessThan("14.00.3281.5.v1", "14.00.3281.6.v1")) // true   
    fmt.Printf("%t\n", version.LessThan("oracle-ee-9", "oracle-ee-19")) // true            
}
```

## `LessThanWithTime`
To compare versions with associated times, such as create times, use `LessThanWithTime()`:

```go
package main

import (
    "time"

    "github.com/YakDriver/go-version"
)

func main() {
    time1 := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC) // January 1, 2000 00:00:00UTC
    time2 := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC) // January 1, 2024 00:00:00UTC

	semVer1 := "1.0.0"
	semVer2 := "1.0.1"

    fmt.Printf("%t\n", version.LessThanWithTime(time1, time2, semVer1, semVer2)) // true
    fmt.Printf("%t\n", version.LessThanWithTime(time1, time2, semVer2, semVer1)) // true (date only)
    fmt.Printf("%t\n", version.LessThanWithTime(time2, time1, semVer1, semVer2)) // false
    fmt.Printf("%t\n", version.LessThanWithTime(time1, time1, semVer1, semVer2)) // true (same time, check versions)
}
```
