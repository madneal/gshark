package cli

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/madneal/gshark/client"
	"github.com/spf13/cobra"
	"golang.org/x/term"
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
	var username string
	login := &cobra.Command{Use: "login", Short: "Log in with username, password and image captcha", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		endpoint := *server
		if !cmd.Flags().Changed("server") && os.Getenv("GSHARK_SERVER") == "" {
			if saved, err := loadCredentials(); err == nil {
				endpoint = saved.Server
			}
		}
		input := bufio.NewReader(cmd.InOrStdin())
		prompt := func(label string) (string, error) {
			fmt.Fprint(cmd.ErrOrStderr(), label)
			line, err := input.ReadString('\n')
			if err != nil {
				return "", fmt.Errorf("read %s: %w", label, err)
			}
			return strings.TrimSpace(line), nil
		}
		var err error
		if username == "" {
			username, err = prompt("Username: ")
			if err != nil {
				return err
			}
		}
		if username == "" {
			return fmt.Errorf("username cannot be empty")
		}
		file, ok := cmd.InOrStdin().(*os.File)
		if !ok || !term.IsTerminal(int(file.Fd())) {
			return fmt.Errorf("auth login requires an interactive terminal for hidden password input")
		}
		fmt.Fprint(cmd.ErrOrStderr(), "Password: ")
		password, err := term.ReadPassword(int(file.Fd()))
		fmt.Fprintln(cmd.ErrOrStderr())
		if err != nil {
			return err
		}
		if len(password) == 0 {
			return fmt.Errorf("password cannot be empty")
		}
		c, err := client.New(endpoint, "", *timeout)
		if err != nil {
			return err
		}
		data, err := c.Post(cmd.Context(), "base/captcha", struct{}{})
		if err != nil {
			return err
		}
		var captcha struct {
			ID      string `json:"captchaId"`
			Picture string `json:"picPath"`
		}
		if err = json.Unmarshal(data, &captcha); err != nil {
			return err
		}
		if err = renderCaptcha(cmd.ErrOrStderr(), captcha.Picture); err != nil {
			return err
		}
		answer, err := prompt("Captcha: ")
		if err != nil {
			return err
		}
		data, err = c.Post(cmd.Context(), "base/login", map[string]string{"username": username, "password": string(password), "captchaId": captcha.ID, "captcha": answer})
		if err != nil {
			return fmt.Errorf("login failed (run auth login to retry): %w", err)
		}
		var result struct {
			Token string `json:"token"`
		}
		if err = json.Unmarshal(data, &result); err != nil {
			return err
		}
		if result.Token == "" {
			return fmt.Errorf("login returned no token")
		}
		if err = saveCredentials(credentials{Server: strings.TrimRight(endpoint, "/"), Username: username, Token: result.Token}); err != nil {
			return err
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "Logged in as %s on %s\n", username, endpoint)
		return err
	}}
	login.Flags().StringVarP(&username, "username", "u", "", "Login username")
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

func renderCaptcha(out io.Writer, data string) error {
	_, encoded, ok := strings.Cut(data, ",")
	if !ok {
		return fmt.Errorf("invalid captcha image")
	}
	raw, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return err
	}
	img, _, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		return fmt.Errorf("decode captcha: %w", err)
	}
	bounds := img.Bounds()
	if bounds.Dx() == 0 || bounds.Dy() == 0 {
		return fmt.Errorf("empty captcha")
	}
	width := bounds.Dx()
	if width > 80 {
		width = 80
	}
	height := bounds.Dy() * width / bounds.Dx()
	// A colored half block carries two image rows per terminal character.
	for y := 0; y < height; y += 2 {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(bounds.Min.X+x*bounds.Dx()/width, bounds.Min.Y+y*bounds.Dy()/height).RGBA()
			yy := y + 1
			if yy >= height {
				yy = height - 1
			}
			br, bg, bb, _ := img.At(bounds.Min.X+x*bounds.Dx()/width, bounds.Min.Y+yy*bounds.Dy()/height).RGBA()
			if _, err = fmt.Fprintf(out, "\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm▀", r>>8, g>>8, b>>8, br>>8, bg>>8, bb>>8); err != nil {
				return err
			}
		}
		if _, err = fmt.Fprintln(out, "\x1b[0m"); err != nil {
			return err
		}
	}
	return nil
}
