// Package kitsync applies the HAWP kit + provider-overlay update, driven
// entirely by the existing core/providers/manifest.yaml — no new mapping
// format is invented. `hawp update` uses this to refresh a downstream
// repo's .hawp/kit/ and installed provider overlay in the same release
// as the CLI binary.
package kitsync

import (
	"os"

	"gopkg.in/yaml.v3"
)

// InstallRule is one source->destination mapping for a provider, matching
// manifest.yaml's installs_to entries exactly.
type InstallRule struct {
	Dest    string `yaml:"dest"`
	From    string `yaml:"from"`
	Pattern string `yaml:"pattern"`
	Install string `yaml:"install"` // "refresh" | "seed-if-missing" | "seed_if_missing" | ""
	Update  string `yaml:"update"`  // "refresh" | "seed-if-missing" | "seed_if_missing" | "skip" | "" (missing = always refresh)
}

// Provider is one provider's manifest entry.
type Provider struct {
	Source     string        `yaml:"source"`
	InstallsTo []InstallRule `yaml:"installs_to"`
	DocRef     string        `yaml:"doc_ref"`
}

// Manifest is the parsed core/providers/manifest.yaml.
type Manifest struct {
	Providers map[string]Provider `yaml:"providers"`
}

// UpdateMode reports how this rule behaves during an update pass:
//   - "refresh": always overwrite from the provider pack
//   - "seed-if-missing": write only when the destination is absent
//   - "skip": never touch the destination on update
// Missing Update defaults to "refresh" for ordinary synced content.
func (r InstallRule) UpdateMode() string {
	switch r.Update {
	case "skip":
		return "skip"
	case "seed-if-missing", "seed_if_missing":
		return "seed-if-missing"
	default:
		return "refresh"
	}
}

// IsSeedIfMissing reports whether this rule is a seed-once placement:
// the file/directory is written only if the destination does not yet
// exist (e.g. CLAUDE.md, AGENTS.md — user is expected to customize
// after initial install).
func (r InstallRule) IsSeedIfMissing() bool {
	return r.Install == "seed-if-missing" || r.Install == "seed_if_missing"
}

// AllProviderNames returns a sorted list of every provider name in
// the manifest, for use when --provider all is requested.
func (m *Manifest) AllProviderNames() []string {
	names := make([]string, 0, len(m.Providers))
	for name := range m.Providers {
		names = append(names, name)
	}
	return names
}

// ParseManifest reads and parses manifest.yaml at path.
func ParseManifest(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var manifest Manifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}
	return &manifest, nil
}
