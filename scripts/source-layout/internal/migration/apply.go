package migration

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"hawp-source-layout/internal/mapping"
)

func samePlan(a, b plan) bool {
	x, _ := json.Marshal(a)
	y, _ := json.Marshal(b)
	return bytes.Equal(x, y)
}

func currentState(root string, p plan) (string, error) {
	if p.MappingRevision != mapping.Revision {
		return "", errors.New("plan belongs to an older mapping; generate and review a new plan after the checkpoint")
	}
	if p.Version != 1 || p.Module != modulePath {
		return "", errors.New("unsupported plan")
	}
	data, modes, err := inventory(root)
	if err != nil {
		return "", err
	}
	if len(data) != len(p.Files) {
		return "", errors.New("source inventory changed: missing or additional files")
	}
	before, after := true, true
	seenSource := map[string]bool{}
	seenDest := map[string]bool{}
	for _, e := range p.Files {
		if !safeRelative(e.Source) || !safeRelative(e.Destination) || seenSource[e.Source] || seenDest[strings.ToLower(e.Destination)] {
			return "", errors.New("unsafe or duplicate mapping")
		}
		seenSource[e.Source] = true
		seenDest[strings.ToLower(e.Destination)] = true
		b, ok := data[e.Source]
		before = before && ok && digest(b) == e.Before && modes[e.Source] == e.Mode
		b, ok = data[e.Destination]
		after = after && ok && digest(b) == e.After && modes[e.Destination] == e.Mode
	}
	if after {
		return "already-applied", nil
	}
	if before {
		if err := checkDestinations(root, p); err != nil {
			return "", err
		}
		return "ready", nil
	}
	return "", errors.New("source changed, partial migration, or destination conflict; regenerate and review plan")
}

func writeFiles(dir string, files []*file) error {
	for _, f := range files {
		p := filepath.Join(dir, filepath.FromSlash(f.Destination))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(p, f.after, fs.FileMode(f.Mode)); err != nil {
			return err
		}
	}
	return nil
}

// Candidate builds run outside the checkout before the first source mutation.
// No tests execute here; both ordinary and opt-in test files must compile.
func verify(files []*file) error {
	dir, err := os.MkdirTemp("", "hawp-source-layout-check-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	// Write files under the module-relative source root so that go.mod sits
	// at modRoot and imports resolve relative to it (mirroring the live tree).
	modRoot := filepath.Join(dir, filepath.FromSlash(sourceRoot))
	if err = writeFiles(modRoot, files); err != nil {
		return err
	}
	for _, args := range [][]string{{"test", "-run", "^$", "./..."}, {"test", "-tags", "integration,benchmark", "-run", "^$", "./..."}, {"vet", "./..."}} {
		cmd := exec.Command("go", args...)
		cmd.Dir = modRoot
		cmd.Env = append(os.Environ(), "GOWORK=off")
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("candidate go %s: %w\n%s", strings.Join(args, " "), err, out)
		}
		fmt.Printf("PASS candidate go %s\n", strings.Join(args, " "))
	}
	return nil
}

// Apply keeps a backup outside src, checks the snapshot again after compilation,
// writes destinations exclusively, and restores originals on an ordinary error.
// An interrupted process leaves its backup for recovery (path printed first).
func apply(root string, p plan, files []*file) (err error) {
	if err = verify(files); err != nil {
		return err
	}
	state, err := currentState(root, p)
	if err != nil {
		return err
	}
	if state != "ready" {
		return errors.New("state changed during verification")
	}
	backup, err := os.MkdirTemp("", "hawp-source-layout-backup-")
	if err != nil {
		return err
	}
	fmt.Println("Recovery backup:", backup)
	changed := []*file{}
	for _, f := range files {
		if f.Source != f.Destination || f.Before != f.After {
			changed = append(changed, f)
			b := filepath.Join(backup, f.Source)
			if err = os.MkdirAll(filepath.Dir(b), 0o755); err != nil {
				return err
			}
			if err = os.WriteFile(b, f.before, fs.FileMode(f.Mode)); err != nil {
				return err
			}
		}
	}
	base := filepath.Join(root, sourceRoot)
	created := []string{}
	touched := []*file{}
	defer func() {
		if err != nil {
			for _, p := range created {
				if e := os.Remove(p); e != nil && !os.IsNotExist(e) {
					fmt.Fprintln(os.Stderr, "rollback remove:", e)
				}
			}
			for _, f := range touched {
				p := filepath.Join(base, f.Source)
				if e := os.WriteFile(p, f.before, fs.FileMode(f.Mode)); e != nil {
					fmt.Fprintln(os.Stderr, "rollback restore:", e)
				}
			}
			fmt.Fprintln(os.Stderr, "Apply failed; originals restored where possible. Backup retained:", backup)
		}
	}()
	for _, f := range changed {
		dest := filepath.Join(base, f.Destination)
		if err = os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return err
		}
		if f.Source != f.Destination {
			var handle *os.File
			handle, err = os.OpenFile(dest, os.O_CREATE|os.O_EXCL|os.O_WRONLY, fs.FileMode(f.Mode))
			if err != nil {
				return err
			}
			created = append(created, dest)
			_, err = handle.Write(f.after)
			closeErr := handle.Close()
			if err != nil {
				return err
			}
			if closeErr != nil {
				return closeErr
			}
			if err = os.Chmod(dest, fs.FileMode(f.Mode)); err != nil {
				return err
			}
		} else {
			touched = append(touched, f)
			if err = os.WriteFile(dest, f.after, fs.FileMode(f.Mode)); err != nil {
				return err
			}
		}
	}
	for _, f := range changed {
		if f.Source != f.Destination {
			touched = append(touched, f)
			if err = os.Remove(filepath.Join(base, f.Source)); err != nil {
				return err
			}
		}
	}
	state, err = currentState(root, p)
	if err != nil {
		return err
	}
	if state != "already-applied" {
		return errors.New("applied tree does not match reviewed plan")
	}
	// Remove only now-empty ancestor directories of successfully moved files.
	for _, f := range changed {
		if f.Source != f.Destination {
			for d := filepath.Dir(filepath.Join(base, f.Source)); d != base; d = filepath.Dir(d) {
				if os.Remove(d) != nil {
					break
				}
			}
		}
	}
	fmt.Println("Applied file moves and imports; backup retained:", backup)
	return nil
}
