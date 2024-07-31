package main

import (
	"os"

	"github.com/alecthomas/kingpin/v2"

	shortener "github.com/kashav/go-url-shortener"
)

const accessTokenName = "GO_URL_SHORTENER_ACCESS_TOKEN"

var (
	app        = kingpin.New("go-url-shortener", "Create and manage shortened URLs with GitHub pages.")
	appVerbose = app.Flag("verbose", "Show detailed output.").Bool()

	create        = app.Command("create", "Create a new entry.").Alias("new")
	createCNAME   = create.Flag("cname", "Optional CNAME record for this repository.").Short('c').String()
	createName    = create.Flag("name", "Endpoint for the shortened URL, chosen randomly if empty.").Short('n').String()
	createPrivate = create.Flag("private", "Make this repository private.").Short('p').Bool()
	createSubdir  = create.Flag("subdir", "Use a subdirectory in a pre-existing repository, instead of creating a new repository.").Short('s').Bool()
	createRepo    = create.Flag("repo", "A pre-existing repository to be used with the --subdir option (expects `foo/bar` for `https://github.com/foo/bar`). Pushes to the default branch.").Short('r').String()
	createURL     = create.Arg("url", "The URL to shorten.").Required().String()

	list = app.Command("list", "Print a list of active entries.").Alias("ls")

	remove      = app.Command("remove", "Remove one or more entries.").Alias("rm")
	removeRepos = remove.Arg("repo", "List of entries to remove.").Required().Strings()
)

func main() {
	command := kingpin.MustParse(app.Parse(os.Args[1:]))

	accessToken := os.Getenv(accessTokenName)
	if accessToken == "" {
		app.Fatalf("expected access token, run `export %s=<token>`", accessTokenName)
	}

	var r shortener.Runner

	switch command {
	case create.FullCommand():
		r = &shortener.Creator{
			CNAME:   *createCNAME,
			Name:    *createName,
			Private: *createPrivate,
			Subdir:  *createSubdir,
			Repo:    *createRepo,
			URL:     *createURL,
			Verbose: *appVerbose,
		}

	case list.FullCommand():
		r = &shortener.Lister{}

	case remove.FullCommand():
		r = &shortener.Remover{
			Repos:   *removeRepos,
			Verbose: *appVerbose,
		}
	}

	if err := shortener.Start(r, accessToken); err != nil {
		kingpin.Fatalf(err.Error())
	}
}
