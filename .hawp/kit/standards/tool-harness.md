# Tool Harness — CLI Commands and MCP Tools

Every CLI command parser and MCP tool must satisfy this harness before the
handler is wired. Apply the slice harness (`slice-harness.md`) alongside
this document.

## Parser contract (CLI and MCP)

Use a typed `flag.FlagSet` parser — never hand-roll flag loops over `os.Args`
or raw string slices. The parser is a pure function that takes `[]string`
and returns a typed struct plus an error; it performs no filesystem writes
or side effects.

```go
type myArgs struct {
    providers    []string
    noUpdateCheck bool
}

func parseMyArgs(args []string) (myArgs, error) {
    fs := flag.NewFlagSet("my-cmd", flag.ContinueOnError)
    // register flags on fs ...
    if err := fs.Parse(args); err != nil {
        return myArgs{}, err
    }
    if fs.NArg() > 0 {
        return myArgs{}, fmt.Errorf("unexpected positional arguments: %v", fs.Args())
    }
    // validate required/conflicting combinations ...
    return result, nil
}
```

## Required rejections

The parser must return a non-nil error for all of these:

| Case | Example input | Why |
|------|---------------|-----|
| Unknown flag | `["--bogus"]` | `flag.ContinueOnError` propagates this |
| Trailing flag with no value | `["--provider"]` | flag package reports missing value |
| Extra positional argument | `["somepath"]` | checked via `fs.NArg() > 0` |
| Provider flag then extra positional | `["--provider", "claude", "extra"]` | same NArg check |
| Empty required value | `["--provider", ""]` | explicit validation after parse |

Add backend-specific rejection cases where the command has conflicting modes
(e.g., `--semantic` with `--lexical-only`).

## Test table structure

Minimum: 3 valid cases, 4 invalid cases. The patterns below mirror the
existing `init` and `update` parser tests.

```go
func TestParseMyArgsValid(t *testing.T) {
    cases := []struct {
        name      string
        args      []string
        wantProvs []string
    }{
        {name: "no args",         args: nil,                              wantProvs: nil},
        {name: "single provider", args: []string{"--provider", "claude"}, wantProvs: []string{"claude"}},
        {name: "all providers",   args: []string{"--provider", "all"},    wantProvs: []string{"all"}},
    }
    for _, c := range cases {
        t.Run(c.name, func(t *testing.T) {
            got, err := parseMyArgs(c.args)
            if err != nil { t.Fatalf("unexpected error: %v", err) }
            // assert got.providers == c.wantProvs
        })
    }
}

func TestParseMyArgsInvalid(t *testing.T) {
    cases := []struct {
        name string
        args []string
    }{
        {name: "unknown flag",              args: []string{"--bogus"}},
        {name: "provider flag no value",    args: []string{"--provider"}},
        {name: "extra positional",          args: []string{"somepath"}},
        {name: "provider then extra",       args: []string{"--provider", "claude", "extra"}},
    }
    for _, c := range cases {
        t.Run(c.name, func(t *testing.T) {
            if _, err := parseMyArgs(c.args); err == nil {
                t.Errorf("parseMyArgs(%v) expected error, got nil", c.args)
            }
        })
    }
}
```

## No side effects before validation

The handler must not write to the filesystem, open network connections, or
mutate global state until the parser returns without error. Verify this by
inspecting the handler: look for any write/exec call that appears before
`if err := parse...; err != nil`.

## Handler golden-path test

One fixture test proving the happy path produces expected output without
leaking side effects:

```go
func TestMyHandlerHappyPath(t *testing.T) {
    // arrange: stub any I/O dependencies via interfaces
    // act: call the handler with valid parsed args
    // assert: returned output matches expected fixture; no filesystem writes outside temp dir
}
```

Use `t.TempDir()` for any filesystem output the handler must produce.
Assert that no files were written outside `t.TempDir()`.

## MCP tool additional requirements

If the tool is exposed via MCP (`hawp mcp`):

- JSON input schema must be defined in the tool registration (field names,
  types, `required` array).
- Unknown fields in the input JSON must return a typed error, not silent
  ignore.
- Error responses must use JSON-RPC 2.0 error objects; never return an empty
  result for a failed call.

## No-filesystem-write before validation — verification checklist

- [ ] Parser is a pure function (no I/O).
- [ ] Handler calls parser first, returns on error before any write.
- [ ] Table-driven tests: ≥3 valid, ≥4 invalid.
- [ ] Golden-path test uses `t.TempDir()` for any output.
- [ ] MCP schema defined (if applicable).
- [ ] MCP rejects unknown fields (if applicable).
