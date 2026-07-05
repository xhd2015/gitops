package git_isolated

// Default test identity injected via -c on every git invocation.
const (
	DefaultUserEmail = "test@test.com"
	DefaultUserName  = "Test"
)

// Env returns environment variables that ignore global/system gitconfig and
// common hook frameworks. Repo-local .git/config still applies.
func Env() []string {
	return []string{
		"GIT_CONFIG_GLOBAL=/dev/null",
		"GIT_CONFIG_SYSTEM=/dev/null",
		"GIT_OPTIONAL_LOCKS=0",
		"HUSKY=0",
		"LEFTHOOK=0",
	}
}