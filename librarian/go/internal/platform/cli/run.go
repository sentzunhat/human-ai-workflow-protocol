// Package cli routes hawp commands to their application services.
package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	appcheck "github.com/sentzunhat/hawp/librarian/go/internal/application/check"
	appcontext "github.com/sentzunhat/hawp/librarian/go/internal/application/context"
	appdb "github.com/sentzunhat/hawp/librarian/go/internal/application/db"
	appembed "github.com/sentzunhat/hawp/librarian/go/internal/application/embed"
	appindex "github.com/sentzunhat/hawp/librarian/go/internal/application/index"
	appkit "github.com/sentzunhat/hawp/librarian/go/internal/application/kit"
	appkitsync "github.com/sentzunhat/hawp/librarian/go/internal/application/kitsync"
	applinks "github.com/sentzunhat/hawp/librarian/go/internal/application/links"
	appprovision "github.com/sentzunhat/hawp/librarian/go/internal/application/provision"
	appsearch "github.com/sentzunhat/hawp/librarian/go/internal/application/search"
	appupdate "github.com/sentzunhat/hawp/librarian/go/internal/application/update"
	appuuid "github.com/sentzunhat/hawp/librarian/go/internal/application/uuidgen"
	appwork "github.com/sentzunhat/hawp/librarian/go/internal/application/work"
	domainindex "github.com/sentzunhat/hawp/librarian/go/internal/domain/index"
	domainsearch "github.com/sentzunhat/hawp/librarian/go/internal/domain/search"
	domainupdate "github.com/sentzunhat/hawp/librarian/go/internal/domain/update"
	"github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/download"
	"github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/filesystem"
	"github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/githubrelease"
	"github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/repo"
	"github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/sqlite"
)

// ExitError carries a non-zero exit code from a command that already
// reported its findings.
type ExitError struct{ Code int }

func (e ExitError) Error() string { return fmt.Sprintf("exit %d", e.Code) }

func Run(args []string) error {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		fmt.Println(helpText())
		return nil
	}

	command := args[0]
	sub := ""
	if len(args) >= 2 {
		sub = args[1]
	}

	switch {
	case command == "uuid":
		return runUUID(args[1:])

	case command == "links" && sub == "check":
		return runLinksCheck()

	case command == "links" && sub == "clean":
		return runLinksClean(args[2:])

	case command == "kit" && sub == "validate":
		return runKitValidate(args[2:])

	case command == "kit" && sub == "normalize":
		return runKitNormalize(args[2:])

	case command == "work" && sub == "validate":
		return runWorkValidate(args[2:])

	case command == "work" && sub == "normalize":
		return runWorkNormalize(args[2:])

	case command == "work" && sub == "new":
		return runWorkNew(args[2:])

	case command == "check":
		return runCheck()

	case command == "init":
		return runInit()

	case command == "version":
		fmt.Println(domainupdate.Version)
		return nil

	case command == "update" && sub == "latest":
		return runUpdateLatest()

	case command == "update" && sub == "sync":
		return runUpdateSync(args[2:])

	case command == "update" && sub == "verify":
		return runUpdateVerify()

	case command == "update":
		return runUpdateFull(args[1:])

	case command == "commands":
		return runCommands(args[1:])

	case command == "backlog" && sub == "validate":
		return runCheck()

	case command == "backlog" && sub == "upgrade":
		return runWorkNormalize(args[2:])

	case command == "db" && sub == "init":
		service := appdb.NewInitService(filesystem.NewLayoutService())
		result, err := service.Execute()
		if err != nil {
			return err
		}
		fmt.Println(result.String())
		return nil

	case command == "index" && sub == "build":
		return runIndexBuild(args[2:])

	case command == "model" && sub == "pull":
		return runModelPull(args[2:])

	case command == "embed":
		return runEmbed(args[1:])

	case command == "search" && sub == "index":
		return runSearchIndex(args[2:])

	case command == "search" && sub == "embed":
		return runSearchEmbed(args[2:])

	case command == "search" && sub == "benchmark":
		return runSearchBenchmark(args[2:])

	case command == "search":
		return runSearch(args[1:])
	}

	return errors.New("unknown command\n\n" + helpText())
}

func runUUID(args []string) error {
	uuid, err := appuuid.New()
	if err != nil {
		return err
	}
	if len(args) >= 1 && args[0] == "--short" {
		fmt.Println(appuuid.Short(uuid))
	} else {
		fmt.Println(uuid)
	}
	return nil
}

func runLinksCheck() error {
	root, err := repo.FindBacklogRepoRoot(mustGetwd())
	if err != nil {
		return err
	}
	result := applinks.Check(root)
	if code := applinks.Render(os.Stdout, os.Stderr, result); code != 0 {
		return ExitError{Code: code}
	}
	return nil
}

// runLinksClean finds the same broken local Markdown links `hawp links
// check` reports and repairs each one: first it searches the repo for
// exactly one file with the same base name and relinks to it (a doc that
// moved keeps working); only when no unique match exists does it neutralize
// the link — dropping the syntax while keeping the visible text, so
// "[setup guide](removed.md)" becomes plain "setup guide" rather than a
// dead pointer. Dry-run by default (--apply to write), matching work
// normalize's convention. Never touches the archival directories links
// check already skips (.hawp/work/closed, evidence, notes, status) —
// frozen history may reference removed paths by design.
func runLinksClean(args []string) error {
	apply := false
	for _, a := range args {
		if a == "--apply" {
			apply = true
		}
	}

	root, err := repo.FindBacklogRepoRoot(mustGetwd())
	if err != nil {
		return err
	}

	result, err := applinks.Clean(root, apply)
	if err != nil {
		return err
	}
	if code := applinks.RenderClean(os.Stdout, result); code != 0 {
		return ExitError{Code: code}
	}
	return nil
}

func runKitValidate(args []string) error {
	kitPath := ""
	for i := 0; i < len(args); i++ {
		if args[i] == "--kit-path" && i+1 < len(args) {
			kitPath = args[i+1]
			i++
		}
	}
	if kitPath == "" {
		root, err := repo.FindBacklogRepoRoot(mustGetwd())
		if err != nil {
			return err
		}
		kitPath = filepath.Join(root, ".hawp", "kit")
	}
	result := appkit.Validate(kitPath)
	if code := appkit.Render(os.Stdout, os.Stderr, result); code != 0 {
		return ExitError{Code: code}
	}
	return nil
}

func runWorkValidate(args []string) error {
	workDir := ""
	for i := 0; i < len(args); i++ {
		switch {
		case args[i] == "--work-root" && i+1 < len(args):
			workDir = args[i+1]
			i++
		case args[i] == "--hawp-root" && i+1 < len(args):
			workDir = filepath.Join(args[i+1], "work")
			i++
		}
	}
	if workDir == "" {
		root, err := repo.FindBacklogRepoRoot(mustGetwd())
		if err != nil {
			return err
		}
		workDir = filepath.Join(root, ".hawp", "work")
	}
	if !repo.Exists(workDir) {
		return fmt.Errorf("could not resolve .hawp/work directory: %s", workDir)
	}
	report, err := appwork.Validate(workDir)
	if err != nil {
		return err
	}
	if code := appwork.Render(os.Stdout, workDir, report); code != 0 {
		return ExitError{Code: code}
	}
	return nil
}

// runWorkNew scaffolds a work item: generates a UUID, writes an
// investigation plan file, and inserts an "inbox" row into BACKLOG.md. This
// is the mechanical half of HAWP's intake workflow (see
// .hawp/kit/usage/intake-workflow.md) — the investigation and plan content
// itself is a human/AI-agent job, not something this command does.
func runWorkNew(args []string) error {
	if len(args) < 1 {
		return errors.New(`usage: hawp work new "<title>" [--type task|bug|improvement|feature|fix|test|infrastructure|release|decision] [--input "<verbatim request>"] [--hawp-root <path>]`)
	}

	title := ""
	itemType := "task"
	inputText := ""
	hawpRoot := ""

	for i := 0; i < len(args); i++ {
		switch {
		case args[i] == "--type" && i+1 < len(args):
			itemType = args[i+1]
			i++
		case args[i] == "--input" && i+1 < len(args):
			inputText = args[i+1]
			i++
		case args[i] == "--hawp-root" && i+1 < len(args):
			hawpRoot = args[i+1]
			i++
		case !strings.HasPrefix(args[i], "--") && title == "":
			title = args[i]
		}
	}

	if title == "" {
		return errors.New(`usage: hawp work new "<title>" [--type ...] [--input ...] [--hawp-root <path>]`)
	}

	var workDir string
	if hawpRoot != "" {
		workDir = filepath.Join(hawpRoot, "work")
	} else {
		root, err := repo.FindBacklogRepoRoot(mustGetwd())
		if err != nil {
			return err
		}
		workDir = filepath.Join(root, ".hawp", "work")
	}
	if !repo.Exists(workDir) {
		return fmt.Errorf("could not resolve .hawp/work directory: %s", workDir)
	}

	result, err := appwork.NewItem(workDir, itemType, title, inputText)
	if err != nil {
		return err
	}

	fmt.Printf("Created work item %s (%s)\n", appuuid.Short(result.UUID), result.Type)
	fmt.Printf("  Plan file: %s\n", result.PlanFilePath)
	fmt.Printf("  Backlog row added: %s\n", result.BacklogPath)
	fmt.Println()
	fmt.Println("Next: investigate and fill in the plan file (see .hawp/kit/usage/intake-workflow.md Step 2),")
	fmt.Println("then move the backlog status to analyzing/plan-ready as you go.")

	return nil
}

func runInit() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	result := appprovision.Run(download.NewHTTPFetcher(), home, appprovision.DefaultRegistry())
	fmt.Print(result.String())
	if result.Failed() {
		return ExitError{Code: 1}
	}
	return nil
}

// runUpdateVerify checks for an available update and reports it. Exits 1
// if an update is available so scripts can branch on it.
func runUpdateVerify() error {
	client := githubrelease.NewClient()
	status, err := appupdate.Check(client, domainupdate.Repo, domainupdate.Version)
	if err != nil {
		return err
	}
	if status.NoReleases {
		fmt.Println("No published releases found yet.")
		return nil
	}
	fmt.Printf("current: %s\nlatest:  %s\n", status.Current, status.Latest)
	if status.UpdateAvailable {
		fmt.Printf("Update available: run `hawp update` to install %s.\n", status.Latest)
		return ExitError{Code: 1}
	}
	fmt.Println("Already up to date.")
	return nil
}

// runUpdateLatest downloads and installs the latest binary without touching
// .hawp/kit/ or any provider overlays.
func runUpdateLatest() error {
	client := githubrelease.NewClient()
	status, err := appupdate.Check(client, domainupdate.Repo, domainupdate.Version)
	if err != nil {
		return err
	}
	if status.NoReleases {
		fmt.Println("No published releases found yet; nothing to update to.")
		return nil
	}
	fmt.Printf("current: %s\nlatest:  %s\n", status.Current, status.Latest)
	if !status.UpdateAvailable {
		fmt.Println("Already up to date.")
		return nil
	}
	execPath, err := os.Executable()
	if err != nil {
		return err
	}
	applied, err := appupdate.Apply(download.NewHTTPFetcher(), client, domainupdate.Repo, execPath)
	if err != nil {
		return err
	}
	fmt.Printf("Updated binary to %s.\n", applied)
	return nil
}

// runUpdateSync refreshes .hawp/kit/ from the latest release. Pass
// --provider <name> (or --provider all) to also sync provider overlays.
// Skips gracefully when not inside a HAWP repo.
func runUpdateSync(args []string) error {
	providers := parseProviderFlags(args)
	client := githubrelease.NewClient()
	return doKitSync(client, providers)
}

// runUpdateFull is the default for bare `hawp update`: updates the binary
// then syncs kit. Providers are opt-in via --provider.
func runUpdateFull(args []string) error {
	providers := parseProviderFlags(args)

	client := githubrelease.NewClient()
	status, err := appupdate.Check(client, domainupdate.Repo, domainupdate.Version)
	if err != nil {
		return err
	}
	if status.NoReleases {
		fmt.Println("No published releases found yet; nothing to update to.")
		return nil
	}
	fmt.Printf("current: %s\nlatest:  %s\n", status.Current, status.Latest)
	if !status.UpdateAvailable {
		fmt.Println("Already up to date.")
		return nil
	}
	execPath, err := os.Executable()
	if err != nil {
		return err
	}
	applied, err := appupdate.Apply(download.NewHTTPFetcher(), client, domainupdate.Repo, execPath)
	if err != nil {
		return err
	}
	fmt.Printf("Updated binary to %s.\n", applied)
	return doKitSync(client, providers)
}

func parseProviderFlags(args []string) []string {
	var providers []string
	for i := 0; i < len(args); i++ {
		if args[i] == "--provider" && i+1 < len(args) {
			providers = append(providers, args[i+1])
			i++
		}
	}
	return providers
}

func doKitSync(client githubrelease.Client, providers []string) error {
	root, err := repo.FindBacklogRepoRoot(mustGetwd())
	if err != nil {
		return nil // not in a HAWP repo; skip gracefully
	}
	result, err := appkitsync.Sync(download.NewHTTPFetcher(), client, domainupdate.Repo, root, providers)
	if err != nil {
		return err
	}
	if result.NoBundleAsset {
		fmt.Println("This release has no kit bundle; .hawp/kit/ and providers left unchanged.")
		return nil
	}
	fmt.Printf("Kit refreshed: %d file(s).\n", result.KitFilesWritten)
	for _, name := range sortedProviderNames(result.Providers) {
		action := "refreshed"
		if result.ProviderInstalls[name] {
			action = "installed"
		}
		fmt.Printf("Provider %s %s: %d file(s).\n", name, action, result.Providers[name])
	}
	if len(result.Providers) == 0 {
		fmt.Println("No provider overlays synced (use --provider <name> or --provider all to include them).")
	}
	return nil
}

func sortedProviderNames(m map[string]int) []string {
	names := make([]string, 0, len(m))
	for name := range m {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func runIndexBuild(args []string) error {
	scope := domainindex.ScopeAll
	exportPath := ""
	for i := 0; i < len(args); i++ {
		switch {
		case args[i] == "--scope" && i+1 < len(args):
			scope = domainindex.DocumentScope(args[i+1])
			i++
		case args[i] == "--export" && i+1 < len(args):
			exportPath = args[i+1]
			i++
		}
	}
	switch scope {
	case domainindex.ScopeAll, domainindex.ScopeWork, domainindex.ScopeKit:
	default:
		return fmt.Errorf("unknown --scope %q (want all|work|kit)", scope)
	}

	root, err := repo.FindBacklogRepoRoot(mustGetwd())
	if err != nil {
		return err
	}
	service := appindex.NewBuildService(root)
	result, err := service.Execute(scope)
	if err != nil {
		return err
	}
	fmt.Print(result.String())

	if exportPath != "" {
		if err := result.Export(exportPath); err != nil {
			return err
		}
		fmt.Printf("\nExported %d document(s) to %s\n", len(result.Documents), exportPath)
	}
	return nil
}

// modelsRoot resolves ~/.hawp/models, matching internal/domain/provision's
// layout so pulled models live alongside the init-provisioned ones.
func modelsRoot() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filesystem.ResolveHawpHome(home).Models, nil
}

func runModelPull(args []string) error {
	if len(args) < 1 {
		return errors.New("usage: hawp model pull <hf-org/hf-repo> [--onnx-file <path-in-repo>]")
	}
	modelRepo := args[0]
	onnxFile := ""
	for i := 1; i < len(args); i++ {
		if args[i] == "--onnx-file" && i+1 < len(args) {
			onnxFile = args[i+1]
			i++
		}
	}

	dir, err := modelsRoot()
	if err != nil {
		return err
	}
	fmt.Printf("Pulling %s into %s...\n", modelRepo, dir)
	modelPath, err := appembed.PullModel(context.Background(), modelRepo, onnxFile, dir)
	if err != nil {
		return err
	}
	fmt.Printf("Model ready at %s\n", modelPath)
	return nil
}

func runEmbed(args []string) error {
	var texts []string
	modelRepo := ""
	onnxFile := ""
	for i := 0; i < len(args); i++ {
		switch {
		case args[i] == "--model" && i+1 < len(args):
			modelRepo = args[i+1]
			i++
		case args[i] == "--onnx-file" && i+1 < len(args):
			onnxFile = args[i+1]
			i++
		default:
			texts = append(texts, args[i])
		}
	}
	if len(texts) == 0 {
		return errors.New("usage: hawp embed <text> [<text>...] [--model <hf-org/hf-repo>] [--onnx-file <path>]")
	}

	dir, err := modelsRoot()
	if err != nil {
		return err
	}

	ctx := context.Background()
	var modelPath string
	if modelRepo == "" {
		modelPath, err = appembed.PullDefaultModel(ctx, dir)
	} else {
		modelPath, err = appembed.PullModel(ctx, modelRepo, onnxFile, dir)
	}
	if err != nil {
		return err
	}

	vectors, err := appembed.Embed(ctx, modelPath, texts)
	if err != nil {
		return err
	}
	for i, vector := range vectors {
		fmt.Printf("%q: %d dims, first 4: %v\n", texts[i], len(vector), vector[:min(4, len(vector))])
	}
	return nil
}

func runCheck() error {
	root, err := repo.FindBacklogRepoRoot(mustGetwd())
	if err != nil {
		return err
	}
	if code := appcheck.Run(os.Stdout, os.Stderr, root); code != 0 {
		return ExitError{Code: code}
	}
	return nil
}

func runKitNormalize(args []string) error {
	kitPath := ""
	apply := false
	for i := 0; i < len(args); i++ {
		switch {
		case args[i] == "--apply":
			apply = true
		case args[i] == "--dry-run":
			apply = false
		case args[i] == "--kit-path" && i+1 < len(args):
			kitPath = args[i+1]
			i++
		}
	}
	root, err := repo.FindBacklogRepoRoot(mustGetwd())
	if err != nil {
		return err
	}
	if kitPath == "" {
		kitPath = filepath.Join(root, ".hawp", "kit")
	}
	code := appkit.Normalize(os.Stdout, os.Stderr, appkit.NormalizeOptions{
		KitPath: kitPath, RepoRoot: root, Apply: apply,
	})
	if code != 0 {
		return ExitError{Code: code}
	}
	return nil
}

func runWorkNormalize(args []string) error {
	opts := appwork.NormalizeOptions{}
	for i := 0; i < len(args); i++ {
		switch {
		case args[i] == "--apply":
			opts.Apply = true
		case args[i] == "--dry-run":
			opts.Apply = false
		case args[i] == "--validate":
			opts.Validate = true
		case args[i] == "--force-dirty":
			opts.ForceDirty = true
		case args[i] == "--verbose":
			opts.Verbose = true
		case args[i] == "--format" && i+1 < len(args):
			opts.FormatJSON = args[i+1] == "json"
			i++
		case args[i] == "--output" && i+1 < len(args):
			opts.Output = args[i+1]
			i++
		case args[i] == "--export-plan" && i+1 < len(args):
			opts.ExportPlan = args[i+1]
			i++
		case args[i] == "--export-research-queue" && i+1 < len(args):
			opts.ExportResearchQueue = args[i+1]
			i++
		}
	}
	root, err := repo.FindBacklogRepoRoot(mustGetwd())
	if err != nil {
		return err
	}
	opts.RepoRoot = root
	if code := appwork.Normalize(os.Stdout, os.Stderr, opts); code != 0 {
		return ExitError{Code: code}
	}
	return nil
}

func runSearchIndex(args []string) error {
	root, err := repo.FindBacklogRepoRoot(mustGetwd())
	if err != nil {
		fmt.Println("Not in a HAWP repo; nothing to index.")
		return nil
	}

	fmt.Println("Building enriched document index from .hawp/kit/ and .hawp/work/...")

	// Build corpus from actual files
	corpus, err := buildCorpusFromRepo(root)
	if err != nil {
		return fmt.Errorf("build corpus: %w", err)
	}

	// Path to the index DB
	dbPath := filepath.Join(root, ".hawp", "db", "index.sqlite")

	service := appindex.NewIngestService(dbPath)
	result, err := service.Execute(corpus)
	if err != nil {
		return err
	}

	fmt.Print(result.String())
	fmt.Printf("Index ready at: %s\n", dbPath)
	fmt.Println("Try: hawp search vector")
	return nil
}

// buildCorpusFromRepo walks .hawp/kit/ and .hawp/work/ and builds an enriched corpus.
func buildCorpusFromRepo(repoRoot string) (*appindex.EnrichedCorpus, error) {
	corpus := &appindex.EnrichedCorpus{}

	// Index kit files
	kitPath := filepath.Join(repoRoot, ".hawp", "kit")
	if err := walkKitFiles(kitPath, corpus); err != nil {
		return nil, err
	}

	// Index work files (BACKLOG + active plans + evidence + closed)
	workPath := filepath.Join(repoRoot, ".hawp", "work")
	if err := walkWorkFiles(workPath, corpus); err != nil {
		return nil, err
	}

	return corpus, nil
}

func walkKitFiles(kitPath string, corpus *appindex.EnrichedCorpus) error {
	return filepath.Walk(kitPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || filepath.Ext(path) != ".md" {
			return err
		}

		content, _ := os.ReadFile(path)
		rel := strings.TrimPrefix(path, kitPath)
		rel = strings.TrimPrefix(rel, "/")

		corpus.Documents = append(corpus.Documents, appindex.EnrichedDocument{
			Path:       filepath.Join(".hawp/kit", rel),
			Type:       "guide",
			Category:   "kit",
			FolderRole: "kit/" + filepath.Dir(rel),
			Content:    string(content),
			Metadata:   map[string]interface{}{"file": filepath.Base(path)},
		})
		return nil
	})
}

func walkWorkFiles(workPath string, corpus *appindex.EnrichedCorpus) error {
	return filepath.Walk(workPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || filepath.Ext(path) != ".md" {
			return err
		}

		// Skip BACKLOG for now (complex parsing)
		if filepath.Base(path) == "BACKLOG.md" {
			return nil
		}

		content, _ := os.ReadFile(path)
		rel := strings.TrimPrefix(path, workPath)
		rel = strings.TrimPrefix(rel, "/")

		// Determine folder role
		parts := strings.Split(rel, string(filepath.Separator))
		folderRole := "work"
		if len(parts) > 0 {
			folderRole = "work/" + parts[0]
		}

		corpus.Documents = append(corpus.Documents, appindex.EnrichedDocument{
			Path:       filepath.Join(".hawp/work", rel),
			Type:       "plan",
			Category:   "work",
			FolderRole: folderRole,
			Content:    string(content),
			Status:     strPtr("closed"), // simplified; real parse needed
			Metadata:   map[string]interface{}{"file": filepath.Base(path)},
		})
		return nil
	})
}

func strPtr(s string) *string {
	return &s
}

func runSearchEmbed(args []string) error {
	root, err := repo.FindBacklogRepoRoot(mustGetwd())
	if err != nil {
		fmt.Println("Not in a HAWP repo; nothing to embed.")
		return nil
	}

	dbPath := filepath.Join(root, ".hawp", "db", "index.sqlite")

	// Check if DB exists
	db, err := sqlite.Open(dbPath)
	if err != nil {
		fmt.Printf("Index not found at %s. Run `hawp search index` first.\n", dbPath)
		return nil
	}

	// Check how many chunks need embedding
	needEmbed, err := db.ChunksNeedEmbedding()
	if err != nil {
		db.Close()
		return fmt.Errorf("check embeddings: %w", err)
	}
	db.Close()

	if needEmbed == 0 {
		fmt.Println("All chunks already embedded.")
		return nil
	}

	// Parse flags. --backend ollama switches modelID's meaning from an
	// ONNX model pulled via `hawp model pull` to any model available on
	// the local Ollama server (e.g. --model nomic-embed-text).
	backend := "onnx"
	modelID := appindex.DefaultEmbeddingModel
	modelSet := false
	for i := 0; i < len(args); i++ {
		switch {
		case args[i] == "--model" && i+1 < len(args):
			modelID = args[i+1]
			modelSet = true
			i++
		case args[i] == "--backend" && i+1 < len(args):
			backend = args[i+1]
			i++
		}
	}
	if backend == "ollama" && !modelSet {
		modelID = "nomic-embed-text"
	}

	fmt.Printf("Embedding %d chunks with %s (%s)...\n", needEmbed, modelID, backend)
	if backend == "onnx" {
		fmt.Printf("Estimated time: %.0f seconds\n\n", float64(needEmbed)*0.008)
	}

	service := appindex.NewEmbedService(dbPath)
	result, err := service.Execute(context.Background(), backend, modelID)
	if err != nil {
		return err
	}

	fmt.Print(result.String())
	fmt.Println("Vectors ready for hybrid search. Try: hawp search <query>")
	return nil
}

func runSearch(args []string) error {
	if len(args) < 1 {
		return errors.New("usage: hawp search <query> [--limit <n>] [--context] [--llm-reshape] [--format markdown|json] [--max-tokens <n>]")
	}
	query := args[0]
	limit := 10
	wantContext := false
	wantReshape := false
	format := "markdown"
	maxTokens := 2000

	for i := 1; i < len(args); i++ {
		switch {
		case args[i] == "--limit" && i+1 < len(args):
			if n, err := strconv.Atoi(args[i+1]); err == nil {
				limit = n
			}
			i++
		case args[i] == "--context":
			wantContext = true
		case args[i] == "--llm-reshape":
			wantContext = true // reshaping implies --context; the underlying block is always built
			wantReshape = true
		case args[i] == "--format" && i+1 < len(args):
			format = args[i+1]
			i++
		case args[i] == "--max-tokens" && i+1 < len(args):
			if n, err := strconv.Atoi(args[i+1]); err == nil {
				maxTokens = n
			}
			i++
		}
	}

	root, err := repo.FindBacklogRepoRoot(mustGetwd())
	if err != nil {
		fmt.Println("Not in a HAWP repo; no index to search.")
		return nil
	}

	results, err := appsearch.Query(root, query, limit)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	if len(results) == 0 {
		fmt.Printf("No results found for %q\n", query)
		return nil
	}

	if !wantContext {
		// Original output format
		fmt.Printf("Search results for %q (%d found):\n\n", query, len(results))
		for i, result := range results {
			fmt.Printf("[%d] %s | %s (chunk %v)\n",
				i+1,
				result.Source,
				result.Title,
				result.ChunkIndex,
			)
			fmt.Printf("    Context: %s\n", result.FolderContext)
			text := result.Content
			if len(text) > 150 {
				text = text[:150] + "..."
			}
			fmt.Printf("    %q\n\n", text)
		}
		return nil
	}

	// Convert to domain search results
	searchResults := make([]domainsearch.Result, len(results))
	for i, r := range results {
		searchResults[i] = domainsearch.Result{
			ChunkID:       r.ChunkID,
			Source:        r.Source,
			Title:         r.Title,
			Content:       r.Content,
			FolderContext: r.FolderContext,
			ChunkIndex:    r.ChunkIndex,
			Type:          r.Type,
			Category:      r.Category,
			WorkUUID:      r.WorkUUID,
			Status:        r.Status,
			Relevance:     r.Relevance,
			LexicalRank:   r.LexicalRank,
			SemanticScore: r.SemanticScore,
		}
	}

	block := prepareSearchContext(searchResults, query, maxTokens)

	// Optionally reshape via embeddings + LLM (Phase 3). Never silently
	// ignored: any failure to load config, construct backends, or reshape
	// prints an explicit warning to stderr and falls back to the unreshaped
	// context block, so a broken --llm-reshape is always visible to the user.
	var ragOutput *appcontext.RAGPipelineOutput
	if wantReshape {
		ragOutput = tryReshapeViaRAGPipeline(block, maxTokens)
	}

	// Output based on format
	switch format {
	case "json":
		// Convert ContextBlock to JSON-serializable format. "references" is the
		// deduplicated-by-source list (with matched content) — this is the
		// shape a downstream retrieval step consumes as reference docs; use
		// it instead of "results" when you need one entry per unique source.
		jsonBlock := map[string]interface{}{
			"title":        block.Title,
			"query":        block.Query,
			"result_count": block.ResultCount,
			"token_count":  block.TokenCount,
			"results":      block.Results,
			"references":   toJSONReferences(block.References),
			"metadata":     block.Metadata,
		}
		if ragOutput != nil {
			jsonBlock["reshaped_content"] = ragOutput.Content
			jsonBlock["key_concepts"] = ragOutput.KeyConcepts
			jsonBlock["pipeline"] = ragOutput.Pipeline
			// Reshaped references (from the RAG pipeline) supersede the raw
			// ones above once reshaping ran, since they carry the same shape
			// plus pipeline provenance.
			jsonBlock["references"] = toJSONReferences(ragOutput.References)
		}
		out, err := json.MarshalIndent(jsonBlock, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(out))
	default: // markdown
		if ragOutput != nil {
			fmt.Println(renderReshapedWithReferences(block.Title, ragOutput))
		} else {
			fmt.Println(block.String())
		}
	}

	return nil
}

// runSearchBenchmark runs benchmark tests on all 3 search patterns
func runSearchBenchmark(args []string) error {
	root, err := repo.FindBacklogRepoRoot(mustGetwd())
	if err != nil {
		fmt.Println("Not in a HAWP repo; no index to benchmark.")
		return nil
	}

	dbPath := filepath.Join(root, ".hawp", "db", "index.sqlite")
	db, err := sqlite.Open(dbPath)
	if err != nil {
		fmt.Printf("Index not found at %s. Run `hawp search index` first.\n", dbPath)
		return nil
	}
	defer db.Close()

	return RunBenchmark(db)
}

func getStr(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok && v != nil {
		if s, ok := v.(*string); ok && s != nil {
			return *s
		}
		return fmt.Sprintf("%v", v)
	}
	return ""
}

func getInt(m map[string]interface{}, key string) int64 {
	if v, ok := m[key]; ok {
		if i, ok := v.(int64); ok {
			return i
		}
	}
	return 0
}

func getFloat(m map[string]interface{}, key string) float64 {
	if v, ok := m[key]; ok {
		if f, ok := v.(float64); ok {
			return f
		}
	}
	return 0.0
}

func mustGetwd() string {
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return wd
}

func helpText() string {
	return `hawp

Go librarian CLI for HAWP — small native workflow intelligence tool.

USAGE
  hawp <command> [options]

COMMANDS
  uuid [--short]                       generate a work item UUID
  links check                          validate local markdown links (.hawp, docs, README.md)
  links clean [--apply]                relink (or, failing that, neutralize) broken links found by links check
  kit validate [--kit-path <path>]     validate .hawp/kit/ structure
  kit normalize [--apply]              normalize .hawp/kit/ names and links (dry-run default)
  work validate [--work-root <path>]   backlog/plan/evidence integrity checks
  work new "<title>" [--type ...]      scaffold intake: UUID, plan file, inbox backlog row
  work normalize [--apply --validate]  normalize work record drift (dry-run default)
  check                                combined kit + work + links validation
  init                                  provision ~/.hawp (ONNX Runtime + embedding model)
  version                               print the running hawp version
  update                                update binary + kit (no providers by default)
  update latest                         update binary only
  update sync [--provider <name>|all]   sync kit (+ providers if specified)
  update verify                         check whether an update is available (exit 1 = update ready)
  commands [--json]                    list every command; --json is the agent-facing discovery output
  backlog validate                     alias for check
  backlog upgrade                      alias for work normalize
  db init                              plan the ~/.hawp home layout (scaffold)
  index build [--scope all|work|kit] [--export <path>]  enrich kit/work docs with folder context
  search index                                     ingest kit + work documents into SQLite (no vectors yet)
  search embed [--model <name>] [--backend onnx|ollama]  embed all chunks with vectors
  search <query> [--limit <n>] [--context] [--llm-reshape] [--format markdown|json] [--max-tokens <n>]
                                                   lexical + vector hybrid search; --context for LLM-ready format,
                                                   --llm-reshape to additionally restructure via embeddings+LLM
  model pull <hf-org/repo> [--onnx-file <path>]   download any Hugging Face ONNX model into ~/.hawp/models
  embed <text>... [--model <hf-org/repo>]         embed text via a local model (default: all-MiniLM-L6-v2)

WORK NORMALIZE OPTIONS
  --dry-run | --apply    detection only (default) | normalize closed records
  --validate             run workflow validation summary afterwards
  --format text|json     report format (dry-run)
  --output <path>        write report to file
  --export-plan <path>   write plan JSON (dry-run)
  --export-research-queue <path>  write research queue JSON
  --force-dirty          skip the apply-mode dirty-tree guard

SEARCH --CONTEXT OPTIONS (Phase 4)
  --context              output LLM-ready context block (deduped + formatted)
  --llm-reshape           additionally reshape via embeddings+LLM (implies --context;
                          requires backends configured, see librarian/docs/BACKENDS.md;
                          falls back to unreshaped context with a warning if unavailable)
  --format markdown|json output format (default markdown)
  --max-tokens <n>       token budget for context block (default 2000)`
}
