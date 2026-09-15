package distribution

// Network and platform discovery are mocked; downloads, hashes and rename run
// against real files in t.TempDir. The extracted canonical function is unedited.
const binaryDownloadMocks = `
OWNER=fixture
REPO=fixture
SRC="$PWD/source-not-present"
uname() {
  case "$1" in -s) printf 'Linux\n';; -m) printf 'x86_64\n';; *) return 1;; esac
}
curl() {
  local dest="" url=""
  while [ "$#" -gt 0 ]; do
    case "$1" in
      -o) dest="$2"; shift 2;;
      --proto|--proto-redir) shift 2;;
      -*) shift;;
      *) url="$1"; shift;;
    esac
  done
  case "$url" in
    */releases/latest) printf '{"tag_name":"v0.0.24"}\n';;
    */checksums.txt)
      case "$HAWP_TEST_MODE" in
        checksums-fail) return 22;;
        missing-entry) printf '%s other-file\n' "$HAWP_TEST_SHA" > "$dest";;
        duplicate-entry) printf '%s hawp-linux-amd64\n%s hawp-linux-amd64\n' "$HAWP_TEST_SHA" "$HAWP_TEST_SHA" > "$dest";;
        malformed-entry) printf 'not-a-digest hawp-linux-amd64\n' > "$dest";;
        mismatch) printf '%064d hawp-linux-amd64\n' 0 > "$dest";;
        *) printf '%s hawp-linux-amd64\n' "$HAWP_TEST_SHA" > "$dest";;
      esac;;
    */hawp-linux-amd64)
      if [ "$HAWP_TEST_MODE" = binary-fail ]; then printf 'partial' > "$dest"; return 22; fi
      if [ "$HAWP_TEST_MODE" = interrupt ]; then kill -TERM "$BASHPID"; return 143; fi
      cp payload "$dest";;
    *) return 1;;
  esac
}
if [ "$HAWP_TEST_MODE" = no-hash ]; then
  command() {
    if [ "$1" = -v ] && { [ "$2" = sha256sum ] || [ "$2" = shasum ]; }; then return 1; fi
    builtin command "$@"
  }
fi
if [ "$HAWP_TEST_MODE" = hash-fail ]; then
  sha256sum() { return 1; }
fi
`
