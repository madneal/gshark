package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestGetEnvelopeAndProxyPath(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/rule/getRuleList" || r.URL.Query().Get("name") != "a & b" || r.Header.Get("x-token") != "secret" {
			t.Errorf("unexpected request: %s", r.URL)
		}
		w.Write([]byte(`{"code":0,"data":{"list":[],"total":0},"msg":"ok"}`))
	}))
	defer s.Close()
	c, err := New(s.URL+"/api/", "secret", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	data, err := c.Get(context.Background(), "rule/getRuleList", url.Values{"name": {"a & b"}})
	if err != nil || string(data) != `{"list":[],"total":0}` {
		t.Fatalf("%s %v", data, err)
	}
}

func TestGetFailures(t *testing.T) {
	for _, body := range []string{`{"code":7,"data":{},"msg":"expired"}`, `<html>login</html>`, `{"data":{}}`, `{"code":0}`} {
		t.Run(body, func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(body)) }))
			defer s.Close()
			c, _ := New(s.URL, "secret", time.Second)
			if _, err := c.Get(context.Background(), "test", nil); err == nil {
				t.Fatal("expected failure")
			}
		})
	}
}

func TestRedirectDoesNotForwardToken(t *testing.T) {
	called := false
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true }))
	defer target.Close()
	origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, target.URL, http.StatusFound) }))
	defer origin.Close()
	c, _ := New(origin.URL, "secret", time.Second)
	if _, err := c.Get(context.Background(), "test", nil); err == nil || called {
		t.Fatal("redirect must be rejected")
	}
}

func TestLoginUsesJSONPostWithoutToken(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/api/base/login" || r.Header.Get("Content-Type") != "application/json" || r.Header.Get("x-token") != "" {
			t.Errorf("unexpected login request")
		}
		w.Write([]byte(`{"code":0,"data":{"token":"jwt"}}`))
	}))
	defer s.Close()
	c, err := New(s.URL+"/api", "", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = c.Get(context.Background(), "private", nil); err == nil {
		t.Fatal("unauthenticated get allowed")
	}
	if _, err = c.Post(context.Background(), "base/login", map[string]string{"username": "alice", "password": "password"}); err != nil {
		t.Fatal(err)
	}
}
