#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
config_dir="$repo_root/.codex"
config_path="$config_dir/config.toml"
binary_path="$repo_root/.hawp/bin/hawp"

mkdir -p "$config_dir"

cat >"$config_path" <<EOF
[mcp_servers.hawp]
command = "$binary_path"
args = ["mcp"]
cwd = "$repo_root"
enabled = true
startup_timeout_sec = 30
tool_timeout_sec = 60
enabled_tools = ["hawp_search", "hawp_usage", "hawp_work_new", "hawp_work_validate"]
EOF

echo "Wrote $config_path"
echo "Verify with: codex mcp list && codex mcp get hawp"
