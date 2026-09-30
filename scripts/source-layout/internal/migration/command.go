package migration

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// Run executes the migration command. No mode is applied implicitly.
func Run(args []string) error {
	flags := flag.NewFlagSet("source-layout", flag.ContinueOnError)
	root := flags.String("root", ".", "repository root")
	planFile := flags.String("plan", "", "reviewed JSON plan; required by -apply")
	writePlan := flags.String("write-plan", "", "write complete source inventory and mapping")
	preview := flags.Bool("preview", false, "show proposed moves without changing source (default)")
	report := flags.String("report", "", "write Markdown preview with target tree and every file decision")
	diff := flags.Bool("diff", false, "print the proposed content diff without changing source")
	doApply := flags.Bool("apply", false, "apply the reviewed plan after candidate compile/vet checks")
	check := flags.Bool("check", false, "compile and vet the proposed source tree in a temporary directory")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return errors.New("unexpected arguments")
	}
	if *doApply && (*preview || *diff || *report != "" || *writePlan != "") {
		return errors.New("-apply cannot be combined with preview or artifact-generation flags")
	}
	absolute, err := filepath.Abs(*root)
	if err != nil {
		return err
	}
	absolute, err = filepath.EvalSymlinks(absolute)
	if err != nil {
		return err
	}
	if *doApply && *planFile == "" {
		return errors.New("-apply requires -plan")
	}
	var reviewed plan
	if *planFile != "" {
		b, err := os.ReadFile(*planFile)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(b, &reviewed); err != nil {
			return err
		}
		state, err := currentState(absolute, reviewed)
		if err != nil {
			return err
		}
		fmt.Println("Plan state:", state)
		if state == "already-applied" {
			return nil
		}
	}
	p, files, err := prepare(absolute)
	if err != nil {
		return err
	}
	if *planFile != "" && !samePlan(reviewed, p) {
		return errors.New("plan differs from reviewed mapping or transformation; regenerate and review")
	}
	moves, rewrites := 0, 0
	for _, f := range files {
		if f.Source != f.Destination {
			moves++
			fmt.Printf("MOVE %s -> %s\n", f.Source, f.Destination)
		}
		if f.Before != f.After {
			rewrites++
		}
	}
	fmt.Printf("Inventory: %d files; %d moves; %d content updates; %d retained paths. No dead-file deletions.\n", len(files), moves, rewrites, len(files)-moves)
	for _, r := range p.Remaining {
		fmt.Println("REMAINING:", r)
	}
	if *writePlan != "" {
		b, _ := json.MarshalIndent(p, "", "  ")
		if err = writeArtifact(absolute, *writePlan, append(b, '\n')); err != nil {
			return err
		}
	}
	if *report != "" {
		if err = writeArtifact(absolute, *report, []byte(previewReport(p))); err != nil {
			return err
		}
	}
	if *diff {
		if err = previewDiff(files); err != nil {
			return err
		}
	}
	if *doApply {
		return apply(absolute, p, files)
	}
	if *check {
		return verify(files)
	}
	fmt.Println("Preview only. Use -check to compile the candidate, -write-plan to save the reviewed map, -plan with -apply to migrate.")
	return nil
}
