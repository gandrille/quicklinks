package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gandrille/go-commons/filesystem"
	"github.com/gandrille/go-commons/strpair"
)

func New(filePath string) ([]Configuration, error) {
	content, fileerr := filesystem.ReadFileAsString(filePath)
	if fileerr != nil {
		return nil, errors.New("Can't read file " + filePath)
	}

	entries, err := parseEntries(content)
	if err != nil {
		return nil, err
	}

	return buildConfigurations(entries)
}

func parseEntries(content string) ([]entry, error) {
	var entries []entry
	for _, line := range strings.Split(content, "\n") {
		l := strings.TrimSpace(line)
		if strings.HasPrefix(l, "[") && strings.HasSuffix(l, "]") {
			inner := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(l, "["), "]"))
			idx := strings.Index(inner, " ")
			var e entry
			if idx == -1 {
				return nil, fmt.Errorf("key is missing after %s on line %s", inner, line)
			} else {
				kind := strings.TrimSpace(inner[:idx])
				key := strings.TrimSpace(inner[idx+1:])
				if strings.Contains(key, " ") {
					return nil, fmt.Errorf("key \"%s\" must NOT contain spaces on line %s", key, line)
				}
				e = entry{kind, key, []strpair.StrPair{}}
			}
			if e.kind == "" {
				return nil, fmt.Errorf("Can't parse command kind on line %s", line)
			}
			entries = append(entries, e)
		} else if strings.HasPrefix(l, "#") {
			// this is a comment... skip it
		} else if strings.Contains(l, "=") {
			if len(entries) == 0 {
				return nil, fmt.Errorf("No entity started with [kind description] before line %s", line)
			}
			idx := strings.Index(l, "=")
			key := strings.TrimSpace(l[:idx])
			value := strings.TrimSpace(l[idx+1:])
			pair := strpair.New(key, value)
			e := &entries[len(entries)-1]
			e.parameters = append(e.parameters, pair)
		} else if l != "" {
			return nil, fmt.Errorf("Can't parse line %s", line)
		}
	}

	return entries, nil
}

func buildConfigurations(entries []entry) ([]Configuration, error) {
	var configs []Configuration
	builders := builders()

	for _, entry := range entries {
		builder := findBuilder(builders, entry.kind)
		if builder == nil {
			return nil, fmt.Errorf("Don't know how to build entry with type %s", entry.kind)
		}
		conf, err := (*builder).build(entry)
		if err != nil {
			return nil, err
		}
		configs = append(configs, *conf)
	}

	return configs, nil
}

func findBuilder(builders []configurationBuilder, kind string) *configurationBuilder {
	for _, b := range builders {
		if b.kind() == kind {
			return &b
		}
	}
	return nil
}
