package repos

import "os"

// GitEnvFromMap builds git subprocess environment variables from a string map.
func GitEnvFromMap(env map[string]string) []string {
	if len(env) == 0 {
		return nil
	}
	out := make([]string, 0, len(env)+1)
	out = append(out, "GIT_TERMINAL_PROMPT=0")
	for k, v := range env {
		out = append(out, k+"="+v)
	}
	return out
}

// MergeGitEnv merges base process env with git-specific overrides.
func MergeGitEnv(base map[string]string, extra []string) []string {
	merged := append([]string(nil), os.Environ()...)
	if len(base) > 0 {
		merged = append(merged, GitEnvFromMap(base)...)
	}
	if len(extra) > 0 {
		merged = append(merged, extra...)
	}
	return merged
}