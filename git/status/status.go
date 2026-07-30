package status

import (
	"strings"
)

// Counts is the backup-style multi-bucket aggregate of porcelain status lines.
type Counts struct {
	Modified, Added, Deleted, Untracked, Renamed, Copied, Unmerged int
}

// ChangeCounts is the four-bucket aggregate used by inspect-style summaries.
// Untracked (??) counts as Added.
type ChangeCounts struct {
	Added, Changed, Renamed, Deleted int
}

// ParsePorcelain maps `git status --porcelain` text to backup-style Counts.
// Rules (first match wins per line):
//   - ?? → Untracked
//   - R (index) → Renamed
//   - C (index) → Copied
//   - A (index) → Added
//   - U in either column → Unmerged (may also count D/M below)
//   - D in either column → Deleted
//   - M in either column → Modified
func ParsePorcelain(porcelain string) Counts {
	var counts Counts
	for _, line := range strings.Split(porcelain, "\n") {
		if line == "" {
			continue
		}
		// Keep leading XY as-is; require at least two status chars.
		if len(line) < 2 {
			continue
		}
		x, y := line[0], line[1]
		if x == '?' && y == '?' {
			counts.Untracked++
			continue
		}
		if x == 'R' {
			counts.Renamed++
			continue
		}
		if x == 'C' {
			counts.Copied++
			continue
		}
		if x == 'A' {
			counts.Added++
			continue
		}
		if x == 'U' || y == 'U' {
			counts.Unmerged++
		}
		if x == 'D' || y == 'D' {
			counts.Deleted++
			continue
		}
		if x == 'M' || y == 'M' {
			counts.Modified++
		}
	}
	return counts
}

// ParseChangeCounts maps porcelain lines to four ChangeCounts buckets.
// Rules (first match wins):
//   - ?? → Added
//   - R in either column → Renamed
//   - A in either column → Added
//   - D in either column → Deleted
//   - else (including M) → Changed
func ParseChangeCounts(porcelain string) ChangeCounts {
	var counts ChangeCounts
	for _, line := range strings.Split(porcelain, "\n") {
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "??") {
			counts.Added++
			continue
		}
		if len(line) < 2 {
			counts.Changed++
			continue
		}
		x, y := line[0], line[1]
		switch {
		case x == 'R' || y == 'R':
			counts.Renamed++
		case x == 'A' || y == 'A':
			counts.Added++
		case x == 'D' || y == 'D':
			counts.Deleted++
		default:
			counts.Changed++
		}
	}
	return counts
}
