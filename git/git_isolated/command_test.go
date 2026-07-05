package git_isolated

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnvIgnoresGlobalUserConfig(t *testing.T) {
	globalCfg := filepath.Join(t.TempDir(), "global-gitconfig")
	if err := os.WriteFile(globalCfg, []byte("[user]\n\temail = global@example.com\n\tname = Global\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	repo := filepath.Join(t.TempDir(), "repo")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatal(err)
	}

	cfg := Config{ExtraEnv: []string{"GIT_CONFIG_GLOBAL=" + globalCfg}}
	if err := cfg.Command(repo, "-c", "init.templateDir=", "init", "-b", "main").Run(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "README.md"), []byte("hi\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := cfg.Command(repo, "add", "README.md").Run(); err != nil {
		t.Fatal(err)
	}
	if err := cfg.Command(repo, "commit", "-m", "init").Run(); err != nil {
		t.Fatalf("commit should succeed with default -c user.* despite global config: %v", err)
	}

	author, err := cfg.Output(repo, "log", "-1", "--pretty=%ae")
	if err != nil {
		t.Fatal(err)
	}
	if author != DefaultUserEmail {
		t.Fatalf("author email = %q, want %q", author, DefaultUserEmail)
	}
}

func TestInitSkipsTemplateHooks(t *testing.T) {
	repo := filepath.Join(t.TempDir(), "repo")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := Init(repo, "main"); err != nil {
		t.Fatal(err)
	}
	hooksDir := filepath.Join(repo, ".git", "hooks")
	entries, err := os.ReadDir(hooksDir)
	if os.IsNotExist(err) {
		return
	}
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.Type().IsRegular() && !strings.HasSuffix(e.Name(), ".sample") {
			t.Fatalf("unexpected hook file %q in %s", e.Name(), hooksDir)
		}
	}
}

func TestRepoLocalHooksPathPreserved(t *testing.T) {
	repo := filepath.Join(t.TempDir(), "repo")
	hooksDir := filepath.Join(t.TempDir(), "hooks")
	if err := os.MkdirAll(hooksDir, 0o755); err != nil {
		t.Fatal(err)
	}
	hook := "#!/bin/sh\nexit 2\n"
	if err := os.WriteFile(filepath.Join(hooksDir, "pre-commit"), []byte(hook), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := Init(repo, "main"); err != nil {
		t.Fatal(err)
	}
	if err := Run(repo, "config", "core.hooksPath", hooksDir); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "README.md"), []byte("hi\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Run(repo, "add", "README.md"); err != nil {
		t.Fatal(err)
	}
	if err := Run(repo, "commit", "-m", "should fail"); err == nil {
		t.Fatal("expected repo-local pre-commit hook to fail commit")
	}
}