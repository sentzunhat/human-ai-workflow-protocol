# nodejs

Reference lane for the existing Node/TypeScript librarian implementation.

Today the checked-in implementation still lives in `librarian/` and remains the
current source of truth for:

- work validation and normalization
- kit validation and normalization
- provider materialization
- distribution sync
- Node CLI and SEA PoCs

This folder is a transition marker for the future split:

- `golang/` grows into the small native librarian product
- `nodejs/` becomes the reference/legacy implementation area

For now, use `librarian/` as the real implementation and treat this folder as a
direction marker only.
