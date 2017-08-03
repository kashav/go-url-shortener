package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/kshvmdn/redir"
	"github.com/kshvmdn/redir/version"
)

var (
	app        = kingpin.New("redir", "Create and manage shortened URLs with GitHub pages.")
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
	app.Version(version.VERSION)
	command := kingpin.MustParse(app.Parse(os.Args[1:]))

	accessToken := os.Getenv("REDIR_ACCESS_TOKEN")
	if accessToken == "" {
		app.Fatalf("expected access token, run `export REDIR_ACCESS_TOKEN=<token>`")
	}

	var r redir.Runner

	switch command {
	case create.FullCommand():
		r = &redir.Creator{
			CNAME:   *createCNAME,
			Name:    *createName,
			Private: *createPrivate,
			Subdir:  *createSubdir,
			Repo:    *createRepo,
			URL:     *createURL,
			Verbose: *appVerbose,
		}

	case list.FullCommand():
		r = &redir.Lister{}

	case remove.FullCommand():
		r = &redir.Remover{
			Repos:   *removeRepos,
			Verbose: *appVerbose,
		}
	}

	if err := redir.Start(r, accessToken); err != nil {
		kingpin.Fatalf(err.Error())
	}
}
