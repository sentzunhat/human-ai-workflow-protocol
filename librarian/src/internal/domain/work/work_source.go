package work

// WorkSource provides infrastructure dependencies to domain/work functions.
// Mirrors the KitSource pattern: domain defines function-typed fields;
// application constructs the source from infrastructure implementations.
type WorkSource struct {
	Exists       func(path string) bool
	ToRepoRelative func(repoRoot, absolutePath string) string
	CollectFiles func(dir string, skipReadme bool) []string
}
