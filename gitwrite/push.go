package gitwrite

import (
	"fmt"

	"github.com/xhd2015/xgo/support/cmd"
)

type WriteOptions struct {
	HTTPSProxy string
	HTTPProxy  string
}

func Push(dir string, url string, sourceRef string, remoteBranch string, opts *WriteOptions) error {
	if dir == "" {
		return fmt.Errorf("requires dir")
	}
	if url == "" {
		return fmt.Errorf("requires url")
	}
	if sourceRef == "" {
		sourceRef = "HEAD"
	}
	if remoteBranch == "" {
		return fmt.Errorf("requires remoteBranch")
	}
	return cmd.Dir(dir).Env(GetProxyEnv(opts)).Run("git", "push", url, fmt.Sprintf("%s:refs/heads/%s", sourceRef, remoteBranch))
}

func PushTag(dir string, url string, tag string, opts *WriteOptions) error {
	if dir == "" {
		return fmt.Errorf("requires dir")
	}
	if url == "" {
		return fmt.Errorf("requires url")
	}
	if tag == "" {
		return fmt.Errorf("requires tag")
	}
	return cmd.Dir(dir).Env(GetProxyEnv(opts)).Run("git", "push", url, tag)
}
func GetProxyEnv(opts *WriteOptions) []string {
	if opts == nil {
		return nil
	}
	var envs []string
	if opts.HTTPProxy != "" {
		envs = append(envs, "http_proxy="+opts.HTTPProxy)
	}
	if opts.HTTPSProxy != "" {
		envs = append(envs, "https_proxy="+opts.HTTPSProxy)
	}
	return envs
}
