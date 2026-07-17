# go-version
<!-- Copyright IBM Corp. 2019, 2026 -->
<!-- SPDX-License-Identifier: MPL-2.0 -->

> **I fought the law — and got sorted anyway.**

Rebel versions, gangster versions, versions that never once read the semver spec: all welcome here, all treated with the same respect as their law-abiding, `1.2.3` cousins. **`go-version`** compares version-like strings that don't always play by the rules — the messy, real-world identifiers that strict semantic-versioning libraries choke on — and it **never returns an error**: any two strings can always be ordered.

It was written for the Terraform AWS Provider, where version strings come straight from AWS APIs (for example, RDS database engine versions) and are not something the caller controls. Values like `19.0.0.0.ru-2023-10.rur-2023-10.r10`, `8.0.mysql_aurora.3.05.2`, `SQL Server 2014 12.00.6439.10.v1`, and `oracle-ee-19` are common, and they must still sort sensibly.

## Why not hashicorp/go-version or Masterminds/semver?

Those libraries are excellent — when your versions follow the rules. This package is for when they don't.

| | `YakDriver/go-version` | [`hashicorp/go-version`](https://github.com/hashicorp/go-version) | [`Masterminds/semver`](https://github.com/Masterminds/semver) |
| --- | --- | --- | --- |
| Non-semver strings (e.g. `oracle-ee-19`, `8.0.mysql_aurora.3.05.2`) | Compares them anyway | Returns an error | Returns an error |
| Error handling | Never errors | Errors on parse failure | Errors on parse failure |
| Timestamp-aware ordering | Yes (`*WithTime`) | No | No |
| Semantic-versioning support | Yes (delegates to `hashicorp/go-version`) | Yes (strict) | Yes (strict, constraints) |
| Constraints / ranges (`>=1.2, <2.0`) | No | Yes | Yes |
| Best fit | Untrusted, vendor-supplied version strings | Well-formed versions, constraint checking | Well-formed versions, constraint checking |

In short: reach for `hashicorp/go-version` or `Masterminds/semver` when you control the versions and want strict parsing or constraint ranges. Reach for this package when you **don't** control the versions and just need a total order that won't blow up.

## How comparison works

Comparison uses this precedence, highest to lowest:

1. **Associated timestamps** — only in the `*WithTime` functions, when both times are present, non-zero, and distinct.
2. **Semantic versioning** — parsed via [`hashicorp/go-version`](https://github.com/hashicorp/go-version).
3. **Guerrilla comparison** — for version-like strings that break semver rules, comparing dot-separated parts numerically where possible and by letter/digit runs otherwise.
4. **Plain string comparison** — the final fallback.

## API

Every comparison has a boolean helper and a `*WithTime` variant, all built on `Compare`:

| Function | Returns | Description |
| --- | --- | --- |
| `Compare(v1, v2)` | `int` | `-1`, `0`, or `+1` |
| `LessThan(v1, v2)` | `bool` | `v1` sorts before `v2` |
| `LessThanOrEqual(v1, v2)` | `bool` | `v1` sorts before or equal to `v2` |
| `GreaterThan(v1, v2)` | `bool` | `v1` sorts after `v2` |
| `GreaterThanOrEqual(v1, v2)` | `bool` | `v1` sorts after or equal to `v2` |
| `Equal(v1, v2)` | `bool` | `v1` and `v2` sort equally |
| `CompareWithTime(t1, t2, v1, v2)` | `int` | timestamp-aware `Compare` |
| `LessThanWithTime(t1, t2, v1, v2)` | `bool` | timestamp-aware `LessThan` |
| `LessThanOrEqualWithTime(t1, t2, v1, v2)` | `bool` | timestamp-aware `LessThanOrEqual` |
| `GreaterThanWithTime(t1, t2, v1, v2)` | `bool` | timestamp-aware `GreaterThan` |
| `GreaterThanOrEqualWithTime(t1, t2, v1, v2)` | `bool` | timestamp-aware `GreaterThanOrEqual` |
| `EqualWithTime(t1, t2, v1, v2)` | `bool` | timestamp-aware `Equal` |

The timestamp arguments are `*time.Time`; `nil` or zero times are ignored, in which case the functions fall back to the string comparison.

> **Note:** because semantic versions are compared semantically, `Equal("1.0", "1.0.0")` is `true`.

### Ordering helpers

Convenience helpers for working with collections of version strings, all built on `Compare`:

| Function | Returns | Description |
| --- | --- | --- |
| `Max(versions ...string)` | `string` | greatest version, or `""` if none |
| `Min(versions ...string)` | `string` | least version, or `""` if none |
| `Sort(versions []string)` | | sorts ascending, in place |
| `SortDescending(versions []string)` | | sorts descending, in place |

## Usage

### `LessThan`

For comparing version strings without associated times, use `LessThan()`:

```go
package main

import (
	"fmt"

	"github.com/YakDriver/go-version"
)

func main() {
	// normal semantic-versioning versions
	fmt.Printf("%t\n", version.LessThan("10.11.9", "10.11.10")) // true
	fmt.Printf("%t\n", version.LessThan("10.4", "10.4.27"))      // true
	fmt.Printf("%t\n", version.LessThan("10.6.8", "11"))         // true
	fmt.Printf("%t\n", version.LessThan("1.2rc2", "1.2"))        // true

	// non-semantic-versioning versions
	fmt.Printf("%t\n", version.LessThan("8.0.mysql_aurora.3.1.9", "8.0.mysql_aurora.3.1.10"))                       // true
	fmt.Printf("%t\n", version.LessThan("19.0.0.0.ru-2023-10.rur-2023-10.r9", "19.0.0.0.ru-2023-10.rur-2023-10.r10")) // true
	fmt.Printf("%t\n", version.LessThan("14.00.3281.5.v1", "14.00.3281.6.v1"))                                      // true
	fmt.Printf("%t\n", version.LessThan("oracle-ee-9", "oracle-ee-19"))                                            // true
}
```

### `Compare`, `GreaterThan`, and `Equal`

Use `Compare()` when you need a three-way result (for example, as a `slices.SortFunc` comparator), or the boolean helpers for direct checks:

```go
package main

import (
	"fmt"
	"slices"

	"github.com/YakDriver/go-version"
)

func main() {
	fmt.Println(version.Compare("11.16", "11.5"))     // 1
	fmt.Println(version.Compare("1.0", "1.0.0"))       // 0
	fmt.Println(version.GreaterThan("oracle-ee-19", "oracle-ee-9")) // true
	fmt.Println(version.Equal("1.0", "1.0.0"))         // true

	engines := []string{"8.0.35", "5.7.44", "8.0", "5.7"}
	slices.SortFunc(engines, version.Compare)
	fmt.Println(engines) // [5.7 5.7.44 8.0 8.0.35]
}
```

### `Sort`, `Max`, and `Min`

For working with collections of version strings, use the ordering helpers instead of wiring up `Compare` yourself:

```go
package main

import (
	"fmt"

	"github.com/YakDriver/go-version"
)

func main() {
	engines := []string{"8.0.35", "5.7.44", "8.0", "5.7"}

	version.Sort(engines)
	fmt.Println(engines) // [5.7 5.7.44 8.0 8.0.35]

	version.SortDescending(engines)
	fmt.Println(engines) // [8.0.35 8.0 5.7.44 5.7]

	fmt.Println(version.Max("oracle-ee-9", "oracle-ee-19", "oracle-ee-18")) // oracle-ee-19
	fmt.Println(version.Min("8.0.35", "5.7.44", "8.0", "5.7"))              // 5.7
}
```

### `LessThanWithTime`

To compare versions using their associated times (such as create times), use `LessThanWithTime()`. When both times are present, non-zero, and distinct, the earlier time wins; otherwise the versions are compared directly:

```go
package main

import (
	"fmt"
	"time"

	"github.com/YakDriver/go-version"
)

func main() {
	time1 := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC) // January 1, 2000 00:00:00 UTC
	time2 := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC) // January 1, 2024 00:00:00 UTC
	semVer1 := "1.0.0"
	semVer2 := "1.0.1"

	fmt.Printf("%t\n", version.LessThanWithTime(&time1, &time2, semVer1, semVer2)) // true
	fmt.Printf("%t\n", version.LessThanWithTime(&time1, &time2, semVer2, semVer1)) // true (date only)
	fmt.Printf("%t\n", version.LessThanWithTime(&time2, &time1, semVer1, semVer2)) // false
	fmt.Printf("%t\n", version.LessThanWithTime(&time1, &time1, semVer1, semVer2)) // true (same time, check versions)
}
```

## License

[MPL-2.0](LICENSE)
