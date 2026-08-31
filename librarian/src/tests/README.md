`librarian/src/tests/` holds black-box tests that exercise exported behavior
through stable package boundaries.

Why this lives under `src/` instead of beside it:

- `librarian/src/` is the Go module root today (`go.mod`)
- packages here can still import `internal/...` safely
- a sibling `librarian/tests/` tree would fall outside the module and break
  normal `internal/` import rules

Tests that need unexported symbols or package-internal helpers should remain
co-located beside the code they inspect.
