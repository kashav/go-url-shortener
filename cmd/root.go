package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

type entry struct {
	Owner       string `toml:"owner"`
	Repo        string `toml:"repo"`
	RepoURL     string `toml:"repo_url"`
	RedirectURL string `toml:"redirect_url"`
}

type config struct {
	AccessToken string  `toml:"access_token"`
	Entries     []entry `toml:"entries"`
}

var (
	client   *github.Client
	conf     config
	confFile string
	ctx      context.Context
)

// RootCmd is a cobra command object for wrapping every command that will be executed
var RootCmd = &cobra.Command{
	Use:   "shorten",
	Short: "Shorten URLs via GitHub Pages.",
	Long:  "Create and manage shortened URLs with GitHub Pages.",
}

// Execute adds all child commands to the root command sets flags appropriately.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	var isValid = false
	var fileExists = false

	confFile = fmt.Sprintf("%s/.shorten.toml", usr.HomeDir)

	if _, err = os.Stat(confFile); err == nil {
		// File `confFile` exists. Verify that `access_token` is non-empty.
		fileExists = true

		data, err := ioutil.ReadFile(confFile)
		if err != nil {
			log.Fatal(err)
		}
		toml.Decode(string(data), &conf)

		isValid = conf.AccessToken != ""
	}

	if !fileExists {
		_, err = os.Create(confFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	if !isValid {
		fmt.Print("GitHub access token: ")
		fmt.Scanln(&conf.AccessToken)
		saveConfig()
	}

	ctx = context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: conf.AccessToken})
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)
}

func saveConfig() {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(conf); err != nil {
		log.Fatal(err)
	}

	err := ioutil.WriteFile(confFile, buf.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func nullifyAccessToken() {
	conf.AccessToken = ""
	saveConfig()
}

func checkError(err error) {
	if err == nil {
		return
	}

	if _, ok := err.(*github.RateLimitError); ok {
		log.Fatal("GitHub rate limit exceeded.")
	}

	if strings.Contains(err.Error(), "401 Bad credentials") {
		nullifyAccessToken()
		log.Fatal("Invalid credentials. Create a new access token and try again.")
	}

	log.Fatal(err)
}
