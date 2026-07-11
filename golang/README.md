# golang

Go scaffold for the future librarian product.

This workspace is the start of the small native HAWP intelligence tool:

- `db init`
- `work/` + `kit/` ingest
- lexical search
- later vector search, local context building, and prompt handoff

Current status:

- Zacatl-shaped layered scaffold only
- binary size proof is under the current target threshold
- no ONNX/model integration yet
- no database implementation yet

Entry point:

```bash
cd golang
go build -trimpath -ldflags='-s -w' -o bin/hawp ./cmd/hawp
./bin/hawp --help
```

Current binary proof on this machine:

- `bin/hawp` = `1,674,450` bytes, about `1.6 MB`
