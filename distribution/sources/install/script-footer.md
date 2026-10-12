```bash
if [ -n "$TMP_DIR" ] && [ -d "$TMP_DIR" ]; then
  rm -rf "$TMP_DIR"
fi

echo "HAWP install complete (provider: ${PROVIDER})."
echo "Refreshed: .hawp/LICENSE, .hawp/kit/**, .hawp/bin/hawp (native executable)"
echo "Seeded if missing: .hawp/.gitignore (bin/ and db/ runtime state)"
echo "Already tracked binaries or databases must be untracked separately; old hawp-bin and hawp-mcp files are not refreshed."
echo "Preserved: .hawp/work/** (no-overwrite)"
echo "Reconciled: Done rows + Active-Work 'done'/'wont-fix' rows moved from .hawp/work/active/ when eligible (see 'reconciled (link):' and 'reconciled (id-fallback):' lines above)"
```
