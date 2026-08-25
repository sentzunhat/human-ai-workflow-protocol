package update

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	domainupdate "github.com/sentzunhat/hawp/librarian/src/internal/domain/update"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/download"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/filesystem"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/githubrelease"
)

const (
	cacheTTL      = 48 * time.Hour
	phase1        = 15 * time.Minute // first notice → "will auto-update in 5 min"
	phase2        = 5 * time.Minute  // → "will auto-update in 1 min"
	phase3        = 1 * time.Minute  // → auto-update
	autoUpdateLag = 3 * time.Second
)

// updateCache is the transient check result written to ~/.hawp/cache/update-check.json.
type updateCache struct {
	CheckedAt     time.Time  `json:"checked_at"`
	Current       string     `json:"current"`
	Latest        string     `json:"latest"`
	HasUpdate     bool       `json:"has_update"`
	FirstNoticeAt *time.Time `json:"first_notice_at,omitempty"`
}

// UpdateConfig is the user preference file at ~/.hawp/config/update.json.
type UpdateConfig struct {
	AutoUpdate bool `json:"auto_update"` // default true — set false to suppress Phase 4
}

func hawpHome() (filesystem.HawpHome, bool) {
	home, err := os.UserHomeDir()
	if err != nil {
		return filesystem.HawpHome{}, false
	}
	return filesystem.ResolveHawpHome(home), true
}

// LoadUpdateConfig reads ~/.hawp/config/update.json.
// Returns defaults (auto_update: true) when the file is absent or unreadable.
func LoadUpdateConfig() UpdateConfig {
	h, ok := hawpHome()
	if !ok {
		return UpdateConfig{AutoUpdate: true}
	}
	data, err := os.ReadFile(h.UpdateConfigFile)
	if err != nil {
		return UpdateConfig{AutoUpdate: true}
	}
	var c UpdateConfig
	if err := json.Unmarshal(data, &c); err != nil {
		return UpdateConfig{AutoUpdate: true}
	}
	return c
}

// SetAutoUpdate writes ~/.hawp/config/update.json with the given auto_update value.
func SetAutoUpdate(enabled bool) error {
	h, ok := hawpHome()
	if !ok {
		return fmt.Errorf("could not resolve home directory")
	}
	if err := os.MkdirAll(h.Config, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(UpdateConfig{AutoUpdate: enabled}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(h.UpdateConfigFile, data, 0644)
}

func loadCache() *updateCache {
	h, ok := hawpHome()
	if !ok {
		return nil
	}
	data, err := os.ReadFile(h.UpdateCacheFile)
	if err != nil {
		return nil
	}
	var c updateCache
	if err := json.Unmarshal(data, &c); err != nil {
		return nil
	}
	return &c
}

func saveCache(c *updateCache) {
	h, ok := hawpHome()
	if !ok {
		return
	}
	if err := os.MkdirAll(h.Cache, 0755); err != nil {
		return
	}
	data, err := json.Marshal(c)
	if err != nil {
		return
	}
	_ = os.WriteFile(h.UpdateCacheFile, data, 0644)
}

// refreshCacheBackground fires a goroutine that hits GitHub releases and
// updates the cache. Errors are silently discarded — a stale cache is
// always preferable to blocking the user's command.
func refreshCacheBackground(currentVersion string) {
	go func() {
		client := githubrelease.NewClient()
		status, err := Check(client, domainupdate.Repo, currentVersion)
		if err != nil {
			return
		}
		c := &updateCache{
			CheckedAt: time.Now(),
			Current:   currentVersion,
			Latest:    status.Latest,
			HasUpdate: status.UpdateAvailable,
		}
		// Preserve first_notice_at so the countdown is not reset on refresh.
		if existing := loadCache(); existing != nil {
			c.FirstNoticeAt = existing.FirstNoticeAt
		}
		saveCache(c)
	}()
}

// CheckAndNotify is called (via defer) at the end of most commands. It:
//  1. Refreshes the update cache in the background if stale (>48 h).
//  2. Prints a countdown notice to stderr when an update is pending.
//  3. In Phase 4 (≥21 min after first notice): announces, pauses 3 s, then
//     self-replaces the binary and syncs the kit.
//
// Suppressed when: command is mcp/update/version, --no-update-check is
// passed, or ~/.hawp/config/update.json has "auto_update": false (Phase 4
// only — notices still print so the user knows an update is waiting).
func CheckAndNotify(currentVersion string, skipCheck bool) {
	if skipCheck {
		return
	}

	cache := loadCache()
	if cache == nil || time.Since(cache.CheckedAt) > cacheTTL {
		refreshCacheBackground(currentVersion)
	}

	if cache == nil || !cache.HasUpdate || cache.Latest == currentVersion {
		return
	}

	// Record when we first surfaced the notice for this version.
	if cache.FirstNoticeAt == nil {
		now := time.Now()
		cache.FirstNoticeAt = &now
		saveCache(cache)
	}

	elapsed := time.Since(*cache.FirstNoticeAt)
	total := phase1 + phase2 + phase3

	switch {
	case elapsed < phase1:
		remaining := phase1 - elapsed
		fmt.Fprintf(os.Stderr,
			"\nhawp %s is available (current: %s) — auto-update in %d min. Run `hawp update` to install now.\n",
			cache.Latest, currentVersion, int(remaining.Minutes())+1)

	case elapsed < phase1+phase2:
		remaining := phase1 + phase2 - elapsed
		fmt.Fprintf(os.Stderr,
			"\nhawp %s available — auto-update in %d min. Run `hawp update` to install now.\n",
			cache.Latest, int(remaining.Minutes())+1)

	case elapsed < total:
		remaining := total - elapsed
		fmt.Fprintf(os.Stderr,
			"\nhawp %s available — auto-update in %d min. Last chance: `hawp update` or `hawp update --disable-auto`.\n",
			cache.Latest, int(remaining.Minutes())+1)

	default:
		cfg := LoadUpdateConfig()
		if !cfg.AutoUpdate {
			fmt.Fprintf(os.Stderr,
				"\nhawp %s is available — auto-update is disabled. Run `hawp update` to install.\n",
				cache.Latest)
			return
		}

		fmt.Fprintf(os.Stderr,
			"\nhawp %s is available — auto-updating in 3 seconds. Press Ctrl-C to cancel.\n",
			cache.Latest)
		time.Sleep(autoUpdateLag)
		fmt.Fprintf(os.Stderr, "Updating now...\n")

		execPath, err := os.Executable()
		if err != nil {
			fmt.Fprintf(os.Stderr, "auto-update: could not determine executable path: %v\n", err)
			return
		}
		applied, err := Apply(download.NewHTTPFetcher(), githubrelease.NewClient(), domainupdate.Repo, execPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "auto-update failed: %v\n", err)
			return
		}
		// Clear the cache so notices stop.
		saveCache(&updateCache{
			CheckedAt: time.Now(),
			Current:   applied,
			Latest:    applied,
			HasUpdate: false,
		})
		fmt.Fprintf(os.Stderr, "Updated to hawp %s. Re-run your command.\n", applied)
	}
}
