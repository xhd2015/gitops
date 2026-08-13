package repos

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestReadGitLabFile(t *testing.T) {
	const token = "test-token"
	var gotPath, gotRawQuery, gotToken string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.EscapedPath()
		gotRawQuery = r.URL.RawQuery
		gotToken = r.Header.Get("PRIVATE-TOKEN")
		if got := r.Header.Get("Authorization"); got != "" {
			t.Fatalf("Authorization should be empty for API path, got %q", got)
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("package example\n"))
	}))
	defer server.Close()

	auth, cloneURL := AuthFromToken(server.URL+"/group/repo.git", token)
	ok, content, err := ReadGitLabFile(
		context.Background(),
		cloneURL,
		"release/v1",
		"dir/file name.go",
		FetchOptions{Auth: auth},
	)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected file to exist")
	}
	if content != "package example\n" {
		t.Fatalf("content = %q", content)
	}
	if gotToken != token {
		t.Fatalf("PRIVATE-TOKEN = %q, want %q", gotToken, token)
	}
	wantPath := "/api/v4/projects/" + url.PathEscape("group/repo") +
		"/repository/files/" + url.PathEscape("dir/file name.go") + "/raw"
	if gotPath != wantPath {
		t.Fatalf("path = %q, want %q", gotPath, wantPath)
	}
	if gotRawQuery != "ref="+url.QueryEscape("release/v1") {
		t.Fatalf("query = %q", gotRawQuery)
	}
}

func TestReadGitLabFileNotFound(t *testing.T) {
	server := httptest.NewServer(http.NotFoundHandler())
	defer server.Close()

	ok, content, err := ReadGitLabFile(
		context.Background(),
		server.URL+"/group/repo",
		"deadbeef",
		"missing.go",
		FetchOptions{},
	)
	if err != nil {
		t.Fatal(err)
	}
	if ok || content != "" {
		t.Fatalf("got ok=%v content=%q", ok, content)
	}
}

func TestReadGitLabFileRejectsHTML(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`<!DOCTYPE html><html class="devise-layout-html"><body>Sign in</body></html>`))
	}))
	defer server.Close()

	auth, cloneURL := AuthFromToken(server.URL+"/group/repo.git", "tok")
	ok, content, err := ReadGitLabFile(
		context.Background(),
		cloneURL,
		"abc",
		"file.go",
		FetchOptions{Auth: auth},
	)
	if err == nil {
		t.Fatal("expected error for HTML response")
	}
	if ok || content != "" {
		t.Fatalf("got ok=%v content=%q", ok, content)
	}
	if !strings.Contains(err.Error(), "HTML") {
		t.Fatalf("error = %v", err)
	}
}

func TestReadGitLabFileRejectsRedirect(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/users/sign_in", http.StatusFound)
	}))
	defer server.Close()

	auth, cloneURL := AuthFromToken(server.URL+"/group/repo.git", "tok")
	_, _, err := ReadGitLabFile(
		context.Background(),
		cloneURL,
		"abc",
		"file.go",
		FetchOptions{Auth: auth},
	)
	if err == nil {
		t.Fatal("expected redirect error")
	}
	if !strings.Contains(err.Error(), "redirect") {
		t.Fatalf("error = %v", err)
	}
}

func TestProxyFromGitEnvPrefersHTTPSProxy(t *testing.T) {
	proxy, err := proxyFromGitEnv(map[string]string{
		"https_proxy": "http://credit-squid.example:3128",
		"http_proxy":  "http://wrong.example:3128",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := proxy.String(), "http://credit-squid.example:3128"; got != want {
		t.Fatalf("proxy = %q, want %q", got, want)
	}
}
