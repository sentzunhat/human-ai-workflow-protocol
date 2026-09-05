# Review snapshot status

`preview.md` and `plan.json` are the previously checked 39-move snapshot.
They were relocated, not regenerated. The current mapping also proposes context
configuration/deduplication subfolders and carries the later work-normalization
alignment note; the saved plan predates that revision and is deliberately
rejected by the current apply gate.

Commit the cleanup first. Then follow the [preview commands](../README.md) to
replace these artifacts with a freshly generated and checked proposal. Do not
hand-edit source fingerprints or treat earlier test results as new-map evidence.
