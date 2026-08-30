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

	appcheck "github.com/sentzunhat/hawp/librarian/src/internal/application/check"
	appcontext "github.com/sentzunhat/hawp/librarian/src/internal/application/context"
	appdb "github.com/sentzunhat/hawp/librarian/src/internal/application/db"
	appembed "github.com/sentzunhat/hawp/librarian/src/internal/application/embed"
	appindex "github.com/sentzunhat/hawp/librarian/src/internal/application/index"
	appkit "github.com/sentzunhat/hawp/librarian/src/internal/application/kit"
	appkitsync "github.com/sentzunhat/hawp/librarian/src/internal/application/kitsync"
	applinks "github.com/sentzunhat/hawp/librarian/src/internal/application/links"
	appprovision "github.com/sentzunhat/hawp/librarian/src/internal/application/provision"
	appsearch "github.com/sentzunhat/hawp/librarian/src/internal/application/search"
	appupdate "github.com/sentzunhat/hawp/librarian/src/internal/application/update"
	appuuid "github.com/sentzunhat/hawp/librarian/src/internal/application/uuidgen"
	appwork "github.com/sentzunhat/hawp/librarian/src/internal/application/work"
	domainindex "github.com/sentzunhat/hawp/librarian/src/internal/domain/index"
	domainsearch "github.com/sentzunhat/hawp/librarian/src/internal/domain/search"
	domainupdate "github.com/sentzunhat/hawp/librarian/src/internal/domain/update"
	domainusage "github.com/sentzunhat/hawp/librarian/src/internal/domain/usage"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/download"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/filesystem"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/githubrelease"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repo"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/sqlite"
	appmcp "github.com/sentzunhat/hawp/librarian/src/internal/platform/mcp"
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

	// Emit a non-blocking update notice after the command completes.
	// Skip for mcp (long-running server), update (already handling updates),
	// version (informational), and --no-update-check.
	skipNotify := command == "mcp" || command == "update" || command == "version" ||
		containsArg(args, "--no-update-check")
	defer appupdate.CheckAndNotify(domainupdate.Version, skipNotify)

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
		return runInit(args[1:])

	case command == "mcp":
		return runMCP()

	case command == "version":
		fmt.Println(domainupdate.Version)
		return nil

	case command == "update" && sub == "latest":
		return runUpdateLatest()

	case command == "update" && sub == "sync":
		return runUpdateSync(args[2:])

	case command == "update" && sub == "verify":
		return runUpdateVerify()

	case command == "update" && containsArg(args[1:], "--check"):
		return runUpdateVerify()

	case command == "update" && containsArg(args[1:], "--disable-auto"):
		return runUpdateAutoConfig(false)

	case command == "update" && containsArg(args[1:], "--enable-auto"):
		return runUpdateAutoConfig(true)

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

	case command == "usage" && sub == "enable":
		return runUsageEnable(args[2:])

	case command == "usage" && sub == "disable":
		return runUsageDisable()

	case command == "usage" && sub == "log":
		return runUsageLog()

	case command == "usage" && sub == "report":
		return runUsageReport(args[2:])

	case command == "usage" && sub == "clear":
		return runUsageClear()

	case command == "usage":
		return runUsageTotals()
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

func runInit(args []string) error {
	providers := parseProviderFlags(args)

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	result := appprovision.Run(download.NewHTTPFetcher(), home, appprovision.DefaultRegistry())
	fmt.Print(result.String())
	// Asset failures are non-blocking: kit sync and provider config writes are
	// independent of model/runtime downloads and must proceed regardless.
	// We carry the failure forward and return exit 1 at the very end.

	client := githubrelease.NewClient()
	if err := doKitSync(client, providers); err != nil {
		return err
	}

	if len(providers) > 0 {
		root, rerr := repo.FindBacklogRepoRoot(mustGetwd())
		if rerr != nil {
			fmt.Println("Not in a HAWP repo; skipping MCP config write.")
		} else if err := appmcp.WriteProviderConfigs(root, providers); err != nil {
			return err
		}
	}

	if result.Failed() {
		return ExitError{Code: 1}
	}
	return nil
}

func runMCP() error {
	root, err := repo.FindBacklogRepoRoot(mustGetwd())
	if err != nil {
		root = mustGetwd()
	}
	return appmcp.Serve(root, domainupdate.Version)
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

// runUpdateFull is the default for bare `hawp update`: updates the binary,
// syncs kit, and syncs all provider overlays. Pass --no-providers to skip
// provider sync (kit-only).
func runUpdateFull(args []string) error {
	providers := parseProviderFlags(args)
	if len(providers) == 0 && !containsArg(args, "--no-providers") {
		providers = []string{"all"}
	}

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

// runUpdateAutoConfig writes ~/.hawp/config/update.json to enable or disable
// the Phase-4 auto-install. Notices still print when disabled; only the
// unattended install is suppressed.
func runUpdateAutoConfig(enabled bool) error {
	if err := appupdate.SetAutoUpdate(enabled); err != nil {
		return fmt.Errorf("could not write update config: %w", err)
	}
	if enabled {
		fmt.Println("Auto-update enabled. hawp will self-install after the 21-minute countdown.")
	} else {
		fmt.Println("Auto-update disabled. hawp will still notify you about new versions.")
		fmt.Println("Run `hawp update --enable-auto` to re-enable, or `hawp update` to install manually.")
	}
	return nil
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

func containsArg(args []string, flag string) bool {
	for _, a := range args {
		if a == flag {
			return true
		}
	}
	return false
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
		fmt.Println("No provider overlays synced (use --provider <name>|all or run bare `hawp update`).")
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
		case args[i] == "--migrate-folders":
			opts.MigrateFolders = true
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

	home, _ := os.UserHomeDir()
	hawpHome := ""
	if home != "" {
		hawpHome = filepath.Join(home, ".hawp")
	}
	searchCfg, err := appsearch.LoadSearchConfig(hawpHome, root)
	if err != nil {
		return fmt.Errorf("load search config: %w", err)
	}

	fmt.Printf("Building enriched document index from: %s\n", strings.Join(searchCfg.Index.Paths, ", "))

	// Build corpus from actual files
	corpus, err := buildCorpusFromRepo(root, searchCfg.Index.Paths)
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

// buildCorpusFromRepo walks each configured path and builds an enriched corpus.
// Paths are relative to repoRoot. ".hawp/kit" and ".hawp/work" get their
// enriched walkers; all other paths are walked generically as "custom" corpus.
func buildCorpusFromRepo(repoRoot string, paths []string) (*appindex.EnrichedCorpus, error) {
	corpus := &appindex.EnrichedCorpus{}

	for _, p := range paths {
		abs := filepath.Join(repoRoot, filepath.FromSlash(p))
		switch filepath.ToSlash(p) {
		case ".hawp/kit":
			if err := walkKitFiles(abs, corpus); err != nil {
				return nil, err
			}
		case ".hawp/work":
			if err := walkWorkFiles(abs, corpus); err != nil {
				return nil, err
			}
		default:
			if err := walkCustomPath(abs, p, corpus); err != nil {
				return nil, err
			}
		}
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

// walkCustomPath walks a user-configured path (file or directory) and adds any
// .md files to the corpus under the "custom" category.
func walkCustomPath(abs, configuredPath string, corpus *appindex.EnrichedCorpus) error {
	info, err := os.Stat(abs)
	if os.IsNotExist(err) {
		fmt.Printf("warning: configured index path not found, skipping: %s\n", configuredPath)
		return nil
	}
	if err != nil {
		return err
	}

	if !info.IsDir() {
		// Single file — index it directly if it's a .md file.
		if filepath.Ext(abs) != ".md" {
			return nil
		}
		content, _ := os.ReadFile(abs)
		corpus.Documents = append(corpus.Documents, appindex.EnrichedDocument{
			Path:       configuredPath,
			Type:       "document",
			Category:   "custom",
			FolderRole: "custom/" + filepath.Dir(configuredPath),
			Content:    string(content),
			Metadata:   map[string]interface{}{"file": filepath.Base(abs)},
		})
		return nil
	}

	return filepath.Walk(abs, func(path string, fi os.FileInfo, err error) error {
		if err != nil || fi.IsDir() || filepath.Ext(path) != ".md" {
			return err
		}
		content, _ := os.ReadFile(path)
		rel := strings.TrimPrefix(filepath.ToSlash(path), filepath.ToSlash(abs))
		rel = strings.TrimPrefix(rel, "/")
		docPath := filepath.ToSlash(filepath.Join(configuredPath, rel))
		corpus.Documents = append(corpus.Documents, appindex.EnrichedDocument{
			Path:       docPath,
			Type:       "document",
			Category:   "custom",
			FolderRole: "custom/" + filepath.ToSlash(filepath.Dir(rel)),
			Content:    string(content),
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

	// Parse flags. --backend is required: "onnx" or "ollama".
	// --model overrides the default for that backend.
	backend := ""
	modelID := ""
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
	if backend == "" {
		return errors.New("--backend is required: hawp search embed --backend onnx|ollama [--model <name>]")
	}
	if !modelSet {
		switch backend {
		case "onnx":
			modelID = appindex.DefaultEmbeddingModel
		case "ollama":
			modelID = "nomic-embed-text"
		}
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
		return errors.New("usage: hawp search <query> [--limit <n>] [--semantic] [--context] [--format markdown|json] [--max-tokens <n>] [--verbose|-v] [--hybrid-ratio <f>]")
	}
	query := args[0]
	limit := 10
	wantContext := false
	wantSemantic := false
	format := "markdown"
	maxTokens := 2000
	verbose := false
	hybridRatio := 0.3 // default: 30% lexical, 70% semantic

	for i := 1; i < len(args); i++ {
		switch {
		case args[i] == "--limit" && i+1 < len(args):
			if n, err := strconv.Atoi(args[i+1]); err == nil {
				limit = n
			}
			i++
		case args[i] == "--context":
			wantContext = true
		case args[i] == "--semantic":
			wantSemantic = true
		case args[i] == "--format" && i+1 < len(args):
			format = args[i+1]
			i++
		case args[i] == "--max-tokens" && i+1 < len(args):
			if n, err := strconv.Atoi(args[i+1]); err == nil {
				maxTokens = n
			}
			i++
		case args[i] == "--verbose" || args[i] == "-v":
			verbose = true
		case args[i] == "--hybrid-ratio" && i+1 < len(args):
			f, err := strconv.ParseFloat(args[i+1], 64)
			if err != nil {
				return fmt.Errorf("--hybrid-ratio: %q is not a valid float", args[i+1])
			}
			if f < 0.0 || f > 1.0 {
				return fmt.Errorf("--hybrid-ratio must be in [0.0, 1.0] (got %.4f); 0.0 = pure semantic, 1.0 = pure lexical", f)
			}
			hybridRatio = f
			i++
		}
	}

	root, err := repo.FindBacklogRepoRoot(mustGetwd())
	if err != nil {
		fmt.Println("Not in a HAWP repo; no index to search.")
		return nil
	}

	execution, err := appsearch.DefaultService().Execute(root, appsearch.QueryOptions{
		Query:       query,
		Limit:       limit,
		Semantic:    wantSemantic,
		HybridRatio: float32(hybridRatio),
	})
	if err != nil {
		var missingIndex appsearch.IndexNotFoundError
		if errors.As(err, &missingIndex) {
			fmt.Printf("Index not found at %s. Run `hawp search index` first.\n", missingIndex.Path)
			return nil
		}
		return err
	}

	if wantSemantic && !execution.HasVectors {
		fmt.Println("No vectors found. Run `hawp search embed` first to enable semantic search.")
		return nil
	}

	results := execution.Rows
	if wantSemantic && results == nil {
		fmt.Printf("Semantic search failed for %q — check that your embedding backend is running.\n", query)
		return nil
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
				getStr(result, "path"),
				getStr(result, "folder_role"),
				getInt(result, "chunk_idx"),
			)
			fmt.Printf("    Context: %s\n", getStr(result, "folder_context"))
			text := getStr(result, "text")
			if len(text) > 150 {
				text = text[:150] + "..."
			}
			fmt.Printf("    %q\n\n", text)
		}
		return nil
	}

	// Convert to domain search results
	searchResults := appsearch.RowsToResults(results, execution.HasVectors)

	// Pre-pack content dedup: drop chunks with >70% word-set Jaccard overlap
	// against a higher-ranked chunk. No embeddings needed — fast word-set
	// comparison catches near-duplicate paragraphs from the same document section.
	deduped, droppedByDedup := appcontext.ContentJaccardDedup(searchResults, 0.70)

	// Compute avg chunk token estimate (chars/4) across the full pre-dedup set
	// for the verbose savings report. Done here so the value is accurate even
	// when droppedByDedup is 0.
	avgChunkTokens := 0
	if len(searchResults) > 0 {
		total := 0
		for _, r := range searchResults {
			total += (len(r.Content) + 3) / 4
		}
		avgChunkTokens = total / len(searchResults)
	}

	// Dynamic chunk cap: greedily include deduped chunks until the running token
	// estimate (chars/4) would exceed the budget. Stopping early reduces
	// per-result metadata overhead that FormatAsMarkdown adds, ensuring the
	// wrapper cost doesn't cancel out the dedup savings.
	capped := make([]domainsearch.Result, 0, len(deduped))
	runningTokens := 0
	for _, r := range deduped {
		chunkEst := (len(r.Content) + 3) / 4
		if len(capped) > 0 && runningTokens+chunkEst > maxTokens {
			break
		}
		capped = append(capped, r)
		runningTokens += chunkEst
	}

	// Verbose token accounting: print to stderr so it doesn't pollute stdout
	// context block output (which may be piped to an LLM).
	if verbose {
		savedTokens := droppedByDedup * avgChunkTokens
		fmt.Fprintf(os.Stderr, "context: %d chunks, ~%d tokens (saved ~%d tokens via dedup)\n",
			len(capped), runningTokens, savedTokens)
	}

	// Format as context block
	block := appcontext.FormatAsMarkdown(capped, query, maxTokens)

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
		out, err := json.MarshalIndent(jsonBlock, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(out))
	default: // markdown
		fmt.Println(block.String())
	}

	return nil
}

// toJSONReferences converts DocumentReferences into the JSON shape a
// downstream retrieval step consumes as reference docs: source, title,
// matched content excerpt, relevance, and line range (when known).
func toJSONReferences(refs []appcontext.DocumentReference) []map[string]interface{} {
	out := make([]map[string]interface{}, len(refs))
	for i, r := range refs {
		out[i] = map[string]interface{}{
			"source":     r.Source,
			"title":      r.Title,
			"content":    r.Content,
			"relevance":  r.Relevance,
			"line_start": r.LineStart,
			"line_end":   r.LineEnd,
		}
	}
	return out
}

// runSearchBenchmark runs benchmark tests on all 3 search patterns, or
// the token-savings benchmark when --token is passed.
func runSearchBenchmark(args []string) error {
	tokenMode := false
	exportPath := ""
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--token":
			tokenMode = true
		case "--export":
			if i+1 < len(args) {
				exportPath = args[i+1]
				i++
			}
		}
	}

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

	if !tokenMode && exportPath != "" {
		fmt.Println("Note: --export is only used with --token; ignoring.")
	}
	if tokenMode {
		return RunTokenBenchmark(db, exportPath)
	}
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
  work normalize [--apply --migrate-folders --validate]  normalize work record drift (dry-run default)
  check                                combined kit + work + links validation
  init [--provider <name>|all]           provision ~/.hawp, sync kit, write MCP configs (claude|cursor|codex|continue|all)
  mcp                                   start stdio MCP server (JSON-RPC 2.0) for AI agent tool use
  version                               print the running hawp version
  update                                update binary + kit + all providers (--no-providers for kit-only)
  update --check                        check whether an update is available without installing (exit 1 = update ready)
  update --disable-auto                 disable the 21-min countdown auto-install (notices still print)
  update --enable-auto                  re-enable the auto-install countdown
  update latest                         update binary only
  update sync [--provider <name>|all]   sync kit (+ providers if specified)
  update verify                         check whether an update is available (exit 1 = update ready)
  commands [--json]                    list every command; --json is the agent-facing discovery output
  backlog validate                     alias for check
  backlog upgrade                      alias for work normalize
  db init                              plan the ~/.hawp home layout (scaffold)
  index build [--scope all|work|kit] [--export <path>]  enrich kit/work docs with folder context
  search index                                     ingest configured paths into SQLite (reads .hawp/config/search.json)
  search embed --backend onnx|ollama [--model <name>]  embed all chunks with vectors
  search <query> [--limit <n>] [--semantic] [--context] [--format markdown|json] [--max-tokens <n>] [--hybrid-ratio <f>]
                                                   lexical + vector hybrid search; --semantic for pure-vector mode;
                                                   --context for LLM-ready context block;
                                                   --hybrid-ratio tunes the lexical/semantic blend (default 0.3)
  model pull <hf-org/repo> [--onnx-file <path>]   download any Hugging Face ONNX model into ~/.hawp/models
  embed <text>... [--model <hf-org/repo>]         embed text via a local model (default: all-MiniLM-L6-v2)
  usage                                           show MCP call log totals (opt-in; run hawp usage enable first)
  usage log                                       tail 20 most recent logged calls
  usage report [--export <path>]                  full Markdown report: totals, per-tool breakdown, recent queries
  usage enable [--log-bodies]                     enable call logging (--log-bodies also stores raw input/output)
  usage disable                                   stop recording new calls
  usage clear                                     delete all stored log entries (irreversible)

WORK NORMALIZE OPTIONS
  --dry-run | --apply    detection only (default) | normalize closed records
  --validate             run workflow validation summary afterwards
  --format text|json     report format (dry-run)
  --output <path>        write report to file
  --export-plan <path>   write plan JSON (dry-run)
  --export-research-queue <path>  write research queue JSON
  --force-dirty          skip the apply-mode dirty-tree guard

SEARCH --CONTEXT OPTIONS
  --context              output LLM-ready context block (deduped + formatted)
  --semantic             pure-vector search (requires embed step; skips FTS5)
  --format markdown|json output format (default markdown)
  --max-tokens <n>       token budget for context block (default 2000)
  --verbose | -v         print token accounting summary to stderr (chunks, ~tokens, saved via dedup)
  --hybrid-ratio <f>     lexical fraction for hybrid blend [0.0, 1.0] (default 0.3)`
}

func usageHome() (filesystem.HawpHome, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return filesystem.HawpHome{}, fmt.Errorf("could not resolve home directory: %w", err)
	}
	return filesystem.ResolveHawpHome(home), nil
}

func runUsageTotals() error {
	h, err := usageHome()
	if err != nil {
		return err
	}
	cfg := domainusage.LoadConfig(h.UsageConfigFile)
	if !cfg.Enabled {
		fmt.Println("Usage logging is disabled. Run `hawp usage enable` to start recording calls.")
		return nil
	}
	store, err := domainusage.Open(h.UsageDB)
	if err != nil {
		return fmt.Errorf("usage db: %w", err)
	}
	defer store.Close()
	totals, err := store.GetTotals()
	if err != nil {
		return err
	}
	fmt.Print(domainusage.FormatTotals(totals))
	return nil
}

func runUsageLog() error {
	h, err := usageHome()
	if err != nil {
		return err
	}
	store, err := domainusage.Open(h.UsageDB)
	if err != nil {
		return fmt.Errorf("usage db: %w", err)
	}
	defer store.Close()
	entries, err := store.Recent(20)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		fmt.Println("No entries recorded yet.")
		return nil
	}
	for _, e := range entries {
		fmt.Printf("%s  %-22s  in=%-5d out=%-5d  %s\n",
			e.TS.Format("2006-01-02 15:04:05"),
			e.Tool, e.TokensIn, e.TokensOut, domainusage.EntrySummary(e))
	}
	return nil
}

func runUsageEnable(args []string) error {
	h, err := usageHome()
	if err != nil {
		return err
	}
	cfg := domainusage.LoadConfig(h.UsageConfigFile)
	cfg.Enabled = true
	for _, a := range args {
		if a == "--log-bodies" {
			cfg.LogBodies = true
		}
	}
	if err := domainusage.SaveConfig(h.UsageConfigFile, cfg); err != nil {
		return err
	}
	msg := "Usage logging enabled."
	if cfg.LogBodies {
		msg += " Body capture on (raw input/output stored)."
	}
	fmt.Println(msg)
	return nil
}

func runUsageDisable() error {
	h, err := usageHome()
	if err != nil {
		return err
	}
	cfg := domainusage.LoadConfig(h.UsageConfigFile)
	cfg.Enabled = false
	if err := domainusage.SaveConfig(h.UsageConfigFile, cfg); err != nil {
		return err
	}
	fmt.Println("Usage logging disabled.")
	return nil
}

func runUsageReport(args []string) error {
	exportPath := ""
	for i := 0; i < len(args); i++ {
		if args[i] == "--export" && i+1 < len(args) {
			exportPath = args[i+1]
			i++
		}
	}

	h, err := usageHome()
	if err != nil {
		return err
	}
	store, err := domainusage.Open(h.UsageDB)
	if err != nil {
		return fmt.Errorf("usage db: %w", err)
	}
	defer store.Close()
	rep, err := store.GetReport()
	if err != nil {
		return err
	}
	out := domainusage.FormatReport(rep)
	fmt.Print(out)

	if exportPath != "" {
		if err := os.MkdirAll(filepath.Dir(exportPath), 0o755); err != nil {
			return fmt.Errorf("create export dir: %w", err)
		}
		if err := os.WriteFile(exportPath, []byte(out), 0o644); err != nil {
			return fmt.Errorf("write report: %w", err)
		}
		fmt.Printf("\nReport written to %s\n", exportPath)
	}
	return nil
}

func runUsageClear() error {
	h, err := usageHome()
	if err != nil {
		return err
	}
	store, err := domainusage.Open(h.UsageDB)
	if err != nil {
		return fmt.Errorf("usage db: %w", err)
	}
	defer store.Close()
	fmt.Print("Delete all usage log entries? [y/N] ")
	var answer string
	fmt.Scanln(&answer)
	if strings.ToLower(strings.TrimSpace(answer)) != "y" {
		fmt.Println("Cancelled.")
		return nil
	}
	if err := store.Clear(); err != nil {
		return err
	}
	fmt.Println("Usage log cleared.")
	return nil
}
