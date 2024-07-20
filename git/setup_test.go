package git

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xhd2015/xgo/support/cmd"
)

func mustGetTmpDir() (string, func()) {
	dir, err := getTmpDir()

	if err != nil {
		if dir != "" {
			os.RemoveAll(dir)
		}
		panic(err)
	}
	return dir, func() {
		os.RemoveAll(dir)
	}
}

func getTmpDir() (string, error) {
	tmpDir, err := os.MkdirTemp("", "gitops-test")
	if err != nil {
		return "", err
	}

	err = cmd.Dir(tmpDir).Run("git", "init")
	if err != nil {
		return tmpDir, err
	}

	err = cmd.Dir(tmpDir).Run("git", "branch", "-M", "master")
	if err != nil {
		return tmpDir, err
	}

	err = os.WriteFile(filepath.Join(tmpDir, "README.md"), []byte("test"), 0755)
	if err != nil {
		return tmpDir, err
	}
	err = cmd.Dir(tmpDir).Run("git", "add", "-A")
	if err != nil {
		return tmpDir, err
	}
	err = commit(tmpDir, "test", "test@test.com", "init")
	if err != nil {
		return tmpDir, err
	}
	return tmpDir, nil
}

// borrowed
func commit(dir string, authorName string, authorMail string, msg string) error {
	if dir == "" {
		return fmt.Errorf("requires dir")
	}
	if authorName == "" {
		return fmt.Errorf("requires authorName")
	}
	if authorMail == "" {
		return fmt.Errorf("requires authorMail")
	}
	if msg == "" {
		return fmt.Errorf("requires msg")
	}
	clean, err := IndexClean(dir)
	if err != nil {
		return err
	}
	if clean {
		return fmt.Errorf("no changes made")
	}

	return cmd.Dir(dir).Run("git", "-c", "user.email="+authorMail, "-c", "user.name="+authorMail, "commit", "--author="+fmt.Sprintf("%s <%s>", authorName, authorMail), "-m", msg)
}
