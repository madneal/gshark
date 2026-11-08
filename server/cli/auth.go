package cli

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/madneal/gshark/client"
	"github.com/spf13/cobra"
)

type credentials struct {
	Server   string `json:"server"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

func credentialPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "gshark-cli", "credentials.json"), nil
}

func loadCredentials() (credentials, error) {
	var saved credentials
	path, err := credentialPath()
	if err != nil {
		return saved, err
	}
	body, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return saved, fmt.Errorf("not logged in; run gshark-cli auth login")
	}
	if err != nil {
		return saved, err
	}
	if err = json.Unmarshal(body, &saved); err != nil {
		return saved, fmt.Errorf("invalid credentials file; run auth login again")
	}
	if saved.Token == "" || saved.Server == "" {
		return saved, fmt.Errorf("incomplete credentials; run auth login again")
	}
	return saved, nil
}

func saveCredentials(saved credentials) error {
	path, err := credentialPath()
	if err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	body, err := json.MarshalIndent(saved, "", "  ")
	if err != nil {
		return err
	}
	// Atomic replacement avoids leaving a partial token after an interrupted write.
	f, err := os.CreateTemp(filepath.Dir(path), "credentials-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if _, err = f.Write(body); err != nil {
		f.Close()
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), path)
}

func authCommand(server *string, timeout *time.Duration) *cobra.Command {
	auth := &cobra.Command{Use: "auth", Short: "Log in and manage saved credentials"}
	var webURL string
	var loginTimeout time.Duration
	var noBrowser bool
	login := &cobra.Command{Use: "login", Short: "Log in and authorize the CLI in your browser", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		endpoint := *server
		if !cmd.Flags().Changed("server") && os.Getenv("GSHARK_SERVER") == "" {
			if saved, err := loadCredentials(); err == nil {
				endpoint = saved.Server
			}
		}
		if _, err := client.New(endpoint, "", *timeout); err != nil {
			return err
		}
		page, err := loginPageURL(endpoint, webURL)
		if err != nil {
			return err
		}
		if loginTimeout <= 0 {
			return fmt.Errorf("--login-timeout must be positive")
		}
		ctx, cancel := context.WithTimeout(cmd.Context(), loginTimeout)
		defer cancel()
		saved, err := browserLogin(ctx, page, endpoint, cmd.ErrOrStderr(), noBrowser, openBrowser)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "Logged in as %s on %s\n", saved.Username, saved.Server)
		return err
	}}
	login.Flags().StringVar(&webURL, "web-url", "", "GShark website URL (required for direct backend connections)")
	login.Flags().DurationVar(&loginTimeout, "login-timeout", 5*time.Minute, "Time to complete browser authorization")
	login.Flags().BoolVar(&noBrowser, "no-browser", false, "Print the authorization URL without opening a browser")

	auth.AddCommand(login)
	auth.AddCommand(&cobra.Command{Use: "status", Short: "Show saved login and local token expiry (no server verification)", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		saved, err := loadCredentials()
		if err != nil {
			return err
		}
		expiry := tokenExpiry(saved.Token)
		state := "saved (not verified with server)"
		if !expiry.IsZero() && !expiry.After(time.Now()) {
			state = "expired; run auth login again"
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "Server: %s\nUser: %s\nStatus: %s\n", saved.Server, saved.Username, state)
		if err != nil {
			return err
		}
		if !expiry.IsZero() {
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "Expires: %s\n", expiry.Format(time.RFC3339))
		}
		return err
	}})
	auth.AddCommand(&cobra.Command{Use: "logout", Short: "Remove the saved local login (does not revoke the server token)", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		path, err := credentialPath()
		if err != nil {
			return err
		}
		if err = os.Remove(path); err != nil && !os.IsNotExist(err) {
			return err
		}
		_, err = fmt.Fprintln(cmd.OutOrStdout(), "Saved login removed")
		return err
	}})
	return auth
}

func tokenExpiry(token string) time.Time {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return time.Time{}
	}
	body, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return time.Time{}
	}
	var claims struct {
		Exp int64 `json:"exp"`
	}
	if json.Unmarshal(body, &claims) != nil || claims.Exp == 0 {
		return time.Time{}
	}
	return time.Unix(claims.Exp, 0)
}
