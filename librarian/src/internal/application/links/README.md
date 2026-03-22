# links

`hawp links check` (validate local Markdown links) and `hawp links clean`
(repair the broken ones `check` finds).

## Exports

`Check` / `Render` (validation, `Result` includes both a display-string
`Failures` list and structured `Details` for `Clean` to act on) ·
`Clean` / `RenderClean` (repair).

## `Clean`'s repair order

1. **Relink** — search the whole repo for exactly one file with the same
   base name as the broken link's target, and rewrite the link to point at
   it. Preferred: keeps the reference working rather than just erasing
   evidence it existed.
2. **Neutralize** (fallback) — when no unique match exists (deleted,
   renamed-and-moved, or ambiguous — multiple files share that base name),
   drop the link syntax and keep the visible text:
   `[setup guide](gone.md)` → `setup guide`. A dangling link is worse than
   plain text: it looks like it should work and doesn't, indefinitely.

Dry-run by default (`--apply` to write), matching `work normalize`'s
convention. Never touches the archival directories `Check` already skips
(`.hawp/work/{closed,evidence,notes,status}`) — frozen history is allowed to
reference removed paths by design.

## Quick use

```go
result, err := links.Clean(repoRoot, false) // dry-run
result, err := links.Clean(repoRoot, true)  // apply
```
