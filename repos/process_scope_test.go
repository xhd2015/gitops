//go:build linux || darwin
// +build linux darwin

package repos

import (
 "context"
 "errors"
 "os"
 "os/exec"
 "testing"
 "time"
)

// label: e2e — a process-global single-task scope must be tested in its own
// process, not installed/reset in parallel unit tests in the parent process.
func TestOwnedScopeCancelsLegacyBackgroundCommand(t *testing.T){
 cmd:=exec.Command(os.Args[0],"-test.run=^TestOwnedScopeChild$")
 cmd.Env=append(os.Environ(),"GITOPS_OWNED_SCOPE_CHILD=1")
 output,err:=cmd.CombinedOutput();if err!=nil{t.Fatalf("owned-scope child failed: %v\n%s",err,output)}
}
func TestOwnedScopeChild(t *testing.T){
 if os.Getenv("GITOPS_OWNED_SCOPE_CHILD")!="1"{return}
 root,cancel:=context.WithTimeout(context.Background(),100*time.Millisecond);defer cancel()
 scope,err:=InstallProcessScope(root);if err!=nil{t.Fatal(err)}
 cmd:=exec.Command("/bin/sleep","60")
 if err=runCommandTree(context.Background(),cmd);!errors.Is(err,context.Canceled){t.Fatalf("legacy command did not observe owner cancellation: %v",err)}
 finish,done:=context.WithTimeout(context.Background(),3*time.Second);defer done()
 if !scope.CloseAndWait(finish){t.Fatal("scope did not confirm actual command stop")}
 if err=runCommandTree(context.Background(),exec.Command("/bin/sleep","60"));!errors.Is(err,context.Canceled){t.Fatalf("sealed scope admitted another command: %v",err)}
}
func TestComputeGitPolicyAllowsReadsButNeverNewPreparation(t *testing.T){
 for _,args:=range [][]string{{"git","show","HEAD:a.go"},{"git","-c","http.extraHeader=fixture","diff","a","b"},{"git","--git-dir=/fixture","cat-file","-p","HEAD"}}{if !readOnlyCommand(args){t.Fatal("read rejected")}}
 for _,args:=range [][]string{{"git","fetch","origin"},{"git","clone","fixture"},{"git","worktree","add","fixture"},{"git","config","x.y","value"},{"git","-c","alias.run=!fixture","run"}}{if readOnlyCommand(args){t.Fatal("compute could enter a preparation command")}}
}
