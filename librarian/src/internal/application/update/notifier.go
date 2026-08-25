package update

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/download"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/githubrelease"
	domainupdate "github.com/sentzunhat/hawp/librarian/src/internal/domain/update"
)

const (
	cacheTTL      = 48 * time.Hour
	phase1        = 15 * time.Minute // first notice → "5 min" notice
	phase2        = 5 * time.Minute  // "5 min" → "1 min" notice
	phase3        = 1 * time.Minute  // "1 min" → auto-update
	autoUpdateLag = 3 * time.Second  // pause before replacing binary
)

type updateCache struct {
	CheckedAt     time.Time  `json:"checked_at"`
	Current       string     `json:"current"`
	Latest        string     `json:"latest"`
	HasUpdate     bool       `json:"has_update"`
	FirstNoticeAt *time.Time `json:"first_notice_at,omitempty"`
}

func cachePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".hawp", "update-check.json"), nil
}

func loadCache() *updateCache {
	p, err := cachePath()
	if err != nil {
		return nil
	}
	data, err := os.ReadFile(p)
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
	p, err := cachePath()
	if err != nil {
		return
	}
	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		return
	}
	data, err := json.Marshal(c)
	if err != nil {
		return
	}
	_ = os.WriteFile(p, data, 0644)
}

// refreshCacheBackground fires a goroutine that hits GitHub releases and
// updates the cache file. Errors are silently discarded — a stale cache is
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
		// Preserve first_notice_at from any existing cache so the countdown
		// is not reset when the background check refreshes the record.
		if existing := loadCache(); existing != nil {
			c.FirstNoticeAt = existing.FirstNoticeAt
		}
		saveCache(c)
	}()
}

// CheckAndNotify is called at the end of most commands. It:
//  1. Refreshes the update cache in the background if it is stale (>48 h).
//  2. If an update is pending, prints a countdown notice to stderr.
//  3. After 21 minutes of notices, auto-installs the update in-process.
//
// Skip commands: mcp (long-running), update (already handling updates),
// version (informational), and any invocation that passes --no-update-check.
func CheckAndNotify(currentVersion string, skipCheck bool) {
	if skipCheck {
		return
	}

	cache := loadCache()
	stale := cache == nil || time.Since(cache.CheckedAt) > cacheTTL
	if stale {
		refreshCacheBackground(currentVersion)
	}

	if cache == nil || !cache.HasUpdate || cache.Latest == currentVersion {
		return
	}

	// Set first_notice_at on the first invocation where an update is seen.
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
			"\nhawp %s is available (current: %s). Auto-update in %d min — run `hawp update` to install now.\n",
			cache.Latest, currentVersion, int(remaining.Minutes())+1)

	case elapsed < phase1+phase2:
		remaining := phase1 + phase2 - elapsed
		fmt.Fprintf(os.Stderr,
			"\nhawp %s available — auto-update in %d min. Run `hawp update` to install now.\n",
			cache.Latest, int(remaining.Minutes())+1)

	case elapsed < total:
		remaining := total - elapsed
		fmt.Fprintf(os.Stderr,
			"\nhawp %s available — auto-update in %d min. Last chance: run `hawp update` to control timing.\n",
			cache.Latest, int(remaining.Minutes())+1)

	default:
		fmt.Fprintf(os.Stderr,
			"\nhawp %s is available — auto-updating in 3 seconds. Press Ctrl-C to cancel.\n",
			cache.Latest)
		time.Sleep(autoUpdateLag)

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
		// Clear the update cache so notices stop.
		saveCache(&updateCache{
			CheckedAt: time.Now(),
			Current:   applied,
			Latest:    applied,
			HasUpdate: false,
		})
		fmt.Fprintf(os.Stderr, "Updated to hawp %s. Re-run your command.\n", applied)
	}
}
