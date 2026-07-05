package git_isolated

import "os"

// Init runs git init without copying hooks from init.templateDir.
func Init(dir, branch string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return Run(dir, "-c", "init.templateDir=", "init", "-b", branch)
}

// InitBare runs git init --bare without copying hooks from init.templateDir.
func InitBare(workDir, branch, path string) error {
	if err := os.MkdirAll(workDir, 0o755); err != nil {
		return err
	}
	return Run(workDir, "-c", "init.templateDir=", "init", "--bare", "-b", branch, path)
}