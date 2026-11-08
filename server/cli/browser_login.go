package cli

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

func loginPageURL(endpoint, website string) (*url.URL, error) {
	if website == "" {
		api, err := url.Parse(endpoint)
		if err != nil {
			return nil, err
		}
		if !strings.HasSuffix(strings.TrimRight(api.Path, "/"), "/api") {
			return nil, fmt.Errorf("direct backend login requires --web-url with the GShark website URL")
		}
		api.Path = strings.TrimSuffix(strings.TrimRight(api.Path, "/"), "api")
		website = api.String()
	}
	page, err := url.Parse(website)
	if err != nil || page == nil || page.Host == "" || (page.Scheme != "http" && page.Scheme != "https") || page.User != nil || page.RawQuery != "" || page.Fragment != "" {
		return nil, fmt.Errorf("--web-url must be an http(s) website URL without credentials, query or fragment")
	}
	page.Host = strings.ToLower(page.Host)
	if (page.Scheme == "https" && page.Port() == "443") || (page.Scheme == "http" && page.Port() == "80") {
		page.Host = page.Hostname()
		if strings.Contains(page.Host, ":") {
			page.Host = "[" + page.Host + "]"
		}
	}
	page.Path = strings.TrimRight(page.Path, "/") + "/"
	return page, nil
}

func openBrowser(address string) error {
	var command string
	var args []string
	switch runtime.GOOS {
	case "darwin":
		command, args = "open", []string{address}
	case "windows":
		command, args = "rundll32", []string{"url.dll,FileProtocolHandler", address}
	default:
		command, args = "xdg-open", []string{address}
	}
	cmd := exec.Command(command, args...)
	if err := cmd.Start(); err != nil {
		return err
	}
	go cmd.Wait()
	return nil
}

func browserLogin(ctx context.Context, page *url.URL, endpoint string, out io.Writer, noBrowser bool, launch func(string) error) (credentials, error) {
	var saved credentials
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		return saved, fmt.Errorf("start browser callback: %w", err)
	}
	defer listener.Close()
	nonce := make([]byte, 32)
	if _, err = rand.Read(nonce); err != nil {
		return saved, err
	}
	state := hex.EncodeToString(nonce)
	callback := "http://" + listener.Addr().String() + "/callback"
	origin := page.Scheme + "://" + page.Host
	type outcome struct {
		credentials credentials
		err         error
	}
	result := make(chan outcome, 1)
	var mu sync.Mutex
	completed := false
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
		if r.Host != listener.Addr().String() || r.Header.Get("Origin") != origin {
			http.Error(w, "Invalid callback origin", http.StatusForbidden)
			return
		}
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", "POST")
			http.Error(w, "POST required", http.StatusMethodNotAllowed)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, 32<<10)
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid callback", http.StatusBadRequest)
			return
		}
		if subtle.ConstantTimeCompare([]byte(r.PostForm.Get("state")), []byte(state)) != 1 {
			http.Error(w, "Invalid authorization state", http.StatusForbidden)
			return
		}
		token, username := r.PostForm.Get("token"), r.PostForm.Get("username")
		if token == "" || username == "" {
			http.Error(w, "Missing login credentials", http.StatusBadRequest)
			return
		}
		mu.Lock()
		defer mu.Unlock()
		if ctx.Err() != nil {
			http.Error(w, "Authorization expired", http.StatusGone)
			return
		}
		if completed {
			http.Error(w, "Authorization already completed", http.StatusConflict)
			return
		}
		saved := credentials{Server: strings.TrimRight(endpoint, "/"), Username: username, Token: token}
		if err := saveCredentials(saved); err != nil {
			http.Error(w, "Could not save CLI credentials", http.StatusInternalServerError)
			completed = true
			result <- outcome{err: fmt.Errorf("save browser login: %w", err)}
			return
		}
		completed = true
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintln(w, "GShark CLI authorized. You can close this tab and return to the terminal.")
		result <- outcome{credentials: saved}
	})
	server := &http.Server{Handler: mux, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second}
	defer server.Close()
	serveErrors := make(chan error, 1)
	go func() { serveErrors <- server.Serve(listener) }()
	query := url.Values{"callback": {callback}, "state": {state}}
	authorization := *page
	authorization.Fragment = "/cli-auth?" + query.Encode()
	if _, err = fmt.Fprintf(out, "Complete login and authorize GShark CLI in your browser:\n%s\n", authorization.String()); err != nil {
		return saved, err
	}
	if !noBrowser {
		if err := launch(authorization.String()); err != nil {
			fmt.Fprintln(out, "Could not open a browser. Open the URL above on this computer.")
		}
	}
	select {
	case next := <-result:
		// Let the confirmation response finish before closing active connections.
		shutdown, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		server.Shutdown(shutdown)
		return next.credentials, next.err
	case err := <-serveErrors:
		return saved, fmt.Errorf("browser callback stopped: %w", err)
	case <-ctx.Done():
		return saved, fmt.Errorf("browser login canceled or timed out: %w", ctx.Err())
	}
}
