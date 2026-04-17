package resolver

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Version struct {
	Major int
	Minor int
	Patch int
	Raw   string
}

func Parse(version string) (Version, error) {
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return Version{}, fmt.Errorf("invalid semver %q", version)
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return Version{}, err
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return Version{}, err
	}
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return Version{}, err
	}
	return Version{Major: major, Minor: minor, Patch: patch, Raw: version}, nil
}

func Resolve(versions []string, rng string) (string, error) {
	parsed := make([]Version, 0, len(versions))
	for _, v := range versions {
		p, err := Parse(v)
		if err != nil {
			continue
		}
		if matches(p, rng) {
			parsed = append(parsed, p)
		}
	}
	if len(parsed) == 0 {
		return "", fmt.Errorf("no version matches range %q", rng)
	}
	sort.Slice(parsed, func(i, j int) bool {
		if parsed[i].Major != parsed[j].Major {
			return parsed[i].Major > parsed[j].Major
		}
		if parsed[i].Minor != parsed[j].Minor {
			return parsed[i].Minor > parsed[j].Minor
		}
		return parsed[i].Patch > parsed[j].Patch
	})
	return parsed[0].Raw, nil
}

func matches(v Version, rng string) bool {
	if strings.HasPrefix(rng, "^") {
		base, err := Parse(strings.TrimPrefix(rng, "^"))
		if err != nil {
			return false
		}
		return v.Major == base.Major && compare(v, base) >= 0
	}
	if strings.HasPrefix(rng, "~") {
		base, err := Parse(strings.TrimPrefix(rng, "~"))
		if err != nil {
			return false
		}
		return v.Major == base.Major && v.Minor == base.Minor && compare(v, base) >= 0
	}
	return v.Raw == rng
}

func compare(a, b Version) int {
	if a.Major != b.Major {
		return a.Major - b.Major
	}
	if a.Minor != b.Minor {
		return a.Minor - b.Minor
	}
	return a.Patch - b.Patch
}
