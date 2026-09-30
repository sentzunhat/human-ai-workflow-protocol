package mcp

import (
	"fmt"
	"strings"
)

func codexLaunchArgs(raw any, root string) ([]any, error) {
	if raw == nil {
		return []any{"mcp", "--repo-root", root}, nil
	}
	args, ok := raw.([]any)
	if !ok {
		return nil, fmt.Errorf("HAWP args must be a string array")
	}
	args = append([]any(nil), args...)
	for _, arg := range args {
		if _, ok := arg.(string); !ok {
			return nil, fmt.Errorf("HAWP args must be strings")
		}
	}
	if len(args) == 0 {
		return []any{"mcp", "--repo-root", root}, nil
	}
	if args[0] != "mcp" {
		return nil, fmt.Errorf("custom HAWP invocation requires manual merge")
	}
	// Keep the managed repository root before all user-supplied options. This
	// avoids changing the meaning of user arguments that use a separator and
	// makes the generated invocation's ownership boundary explicit.
	if len(args) == 1 {
		return append(args, "--repo-root", root), nil
	}
	found := false
	for i := 1; i < len(args); i++ {
		arg := args[i].(string)
		if !strings.HasPrefix(arg, "-") {
			return nil, fmt.Errorf("custom HAWP subcommand requires manual merge")
		}
		if arg == "--" {
			return nil, fmt.Errorf("HAWP argument separator requires manual merge")
		}
		if arg == "--repo-root" || strings.HasPrefix(arg, "--repo-root=") {
			if found {
				return nil, fmt.Errorf("duplicate HAWP repo-root arguments")
			}
			found = true
			if arg == "--repo-root" {
				if i+1 == len(args) || strings.HasPrefix(args[i+1].(string), "--") {
					return nil, fmt.Errorf("missing HAWP repo-root value")
				}
				i++
				args[i] = root
			} else {
				args[i] = "--repo-root=" + root
			}
		}
	}
	if !found {
		args = append([]any{"mcp", "--repo-root", root}, args[1:]...)
	}
	return args, nil
}
