package redir

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"golang.org/x/oauth2"

	"github.com/BurntSushi/toml"
	"github.com/google/go-github/github"
)

type entry struct {
	Owner       string `toml:"owner"`
	Repo        string `toml:"repo"`
	RepoURL     string `toml:"repo_url"`
	RedirectURL string `toml:"redirect_url"`
	Private     bool   `toml:"private"`
}

type entries struct {
	Entries []entry `toml:"entry"`
}

const CONFIG_FILE = ".redir.toml"

var state struct {
	File string
	Log  entries
}

// MakeClient creates a new GitHub OAuth2 client with the provided access
// token.
func MakeClient(accessToken string, ctx context.Context) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

// Init creates the log file if it doesn't exist and reads all current
// entries into the program's state.
func Init() error {
	u, err := user.Current()
	if err != nil {
		return err
	}
	state.File = fmt.Sprintf("%s/%s", u.HomeDir, CONFIG_FILE)

	if _, err := os.Stat(state.File); os.IsNotExist(err) {
		// file doesn't exist yet, create it!
		_, err = os.Create(state.File)
		return err
	}

	b, err := ioutil.ReadFile(state.File)
	if err != nil {
		return err
	}

	_, err = toml.Decode(string(b), &state.Log)
	return err
}

// saveLog saves the current log (state.log) to state.file.
func saveLog() error {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(state.Log); err != nil {
		return err
	}
	return ioutil.WriteFile(state.File, buf.Bytes(), 0644)
}
