package cli

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func isolateCredentials(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("APPDATA", dir)
	t.Setenv("GSHARK_TOKEN", "")
	t.Setenv("GSHARK_SERVER", "")
}

func TestSavedLoginAndRenewal(t *testing.T) {
	isolateCredentials(t)
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-token") != "original" {
			t.Errorf("wrong token")
		}
		w.Header().Set("new-token", "renewed")
		w.Write([]byte(`{"code":0,"data":{"total":0}}`))
	}))
	defer s.Close()
	if err := saveCredentials(credentials{Server: s.URL, Username: "alice", Token: "original"}); err != nil {
		t.Fatal(err)
	}
	path, _ := credentialPath()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatalf("permissions: %v", info.Mode())
	}
	cmd := NewCommand()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"scans", "overview"})
	if err = cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	saved, err := loadCredentials()
	if err != nil || saved.Token != "renewed" {
		t.Fatalf("renewal: %+v %v", saved, err)
	}
	status := NewCommand()
	out.Reset()
	status.SetOut(&out)
	status.SetArgs([]string{"auth", "status"})
	if err = status.Execute(); err != nil || strings.Contains(out.String(), "renewed") {
		t.Fatalf("status leaked token or failed: %s %v", out.String(), err)
	}
	logout := NewCommand()
	logout.SetOut(&out)
	logout.SetArgs([]string{"auth", "logout"})
	if err = logout.Execute(); err != nil {
		t.Fatal(err)
	}
	if _, err = loadCredentials(); err == nil {
		t.Fatal("logout did not remove credentials")
	}
}

func TestSavedTokenNeverSentToOtherServer(t *testing.T) {
	isolateCredentials(t)
	called := false
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true }))
	defer s.Close()
	if err := saveCredentials(credentials{Server: "https://other.example/api", Token: "secret"}); err != nil {
		t.Fatal(err)
	}
	cmd := NewCommand()
	cmd.SetArgs([]string{"--server", s.URL, "scans", "overview"})
	if err := cmd.Execute(); err == nil || called {
		t.Fatal("saved token crossed server boundary")
	}
}

func TestTokenExpiry(t *testing.T) {
	token := "header." + base64.RawURLEncoding.EncodeToString([]byte(`{"exp":123456789}`)) + ".sig"
	if tokenExpiry(token).Unix() != 123456789 {
		t.Fatal("incorrect expiry")
	}
}
