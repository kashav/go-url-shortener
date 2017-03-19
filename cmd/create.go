package cmd

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"text/template"
	"time"

	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
)

var (
	url       string
	name      string
	cname     string
	isPrivate bool
	files     = []string{"index.html", "README.md"}
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new shortened URL.",
	Long: `Creates a new GitHub Pages repository with the specificed slug which
redirects to the provided URL.

Example:
  $ shorten create --url=https://github.com/golang/go
  $ shorten create --url=https://github.com/golang/go --name=go --private`,
	Run: func(cmd *cobra.Command, args []string) {
		if url == "" {
			log.Fatal("Expected URL flag.")
		}

		if name == "" {
			name = randStr(8)
		}

		// Include the CNAME file iff --cname flag is provided
		if cname != "" {
			files = append(files, "CNAME")
		}

		var tmpDir string = createFiles(map[string]string{
			"url":   url,
			"title": name,
			"cname": cname,
		})

		repo, _, err := client.Repositories.Create(ctx, "", &github.Repository{
			AutoInit:     github.Bool(true),
			MasterBranch: github.String("refs/heads/gh-pages"),
			Name:         github.String(name),
			Private:      github.Bool(isPrivate),
		})
		checkError(err)

		entries := make([]github.TreeEntry, 0, len(files))
		for _, fn := range files {
			hash := sha1.New()

			fp := fmt.Sprintf("/tmp/shorten/%s/%s", tmpDir, fn)

			f, err := os.Open(fp)
			checkError(err)
			defer f.Close()

			st, err := f.Stat()
			checkError(err)

			if _, err := fmt.Fprintf(hash, "blob %d\x00", st.Size()); err != nil {
				log.Fatal(err)
			}

			b, err := ioutil.ReadAll(io.TeeReader(f, hash))
			checkError(err)

			blob, _, err := client.Git.CreateBlob(ctx, *repo.Owner.Login, *repo.Name, &github.Blob{
				Content:  github.String(base64.StdEncoding.EncodeToString(b)),
				Encoding: github.String("base64"),
			})
			checkError(err)

			entries = append(entries, github.TreeEntry{
				Path: github.String(fn),
				Mode: github.String("100644"),
				Type: github.String("blob"),
				SHA:  github.String(*blob.SHA),
			})
		}

		// Commit the changes
		newTree, _, err := client.Git.CreateTree(ctx, *repo.Owner.Login, *repo.Name, "", entries)
		checkError(err)

		newCommit, _, err := client.Git.CreateCommit(ctx, *repo.Owner.Login, *repo.Name, &github.Commit{
			Message: github.String("Initial commit"),
			Tree:    newTree,
		})
		checkError(err)

		// Create the gh-pages branch off of newCommit
		client.Git.CreateRef(ctx, *repo.Owner.Login, *repo.Name, &github.Reference{
			Ref: github.String("refs/heads/gh-pages"),
			Object: &github.GitObject{
				SHA:  newCommit.SHA,
				Type: github.String("commit"),
			},
		})

		// Delete master branch
		client.Git.DeleteRef(ctx, *repo.Owner.Login, *repo.Name, "refs/heads/master")
		checkError(err)

		// Add newly created repo. to config. file
		conf.Entries = append(conf.Entries, entry{
			Owner:       *repo.Owner.Login,
			Repo:        name,
			RepoUrl:     fmt.Sprintf("https://github.com/%s/%s", *repo.Owner.Login, *repo.Name),
			RedirectUrl: url,
		})
		saveConfig()

		fmt.Printf("Success! https://%s.github.io/%s now redirects to %s!\n",
			*repo.Owner.Login, *repo.Name, url)
	},
}

func randStr(n int) string {
	var chars string = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rand.Seed(time.Now().UTC().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

//go:generate go run template
func createFiles(ipt map[string]string) string {
	tmpDir := randStr(24)
	err := os.MkdirAll(fmt.Sprintf("/tmp/shorten/%s", tmpDir), 0777)
	checkError(err)

	for _, fn := range files {
		src := fmt.Sprintf("./template/%s", fn)
		dest := fmt.Sprintf("/tmp/shorten/%s/%s", tmpDir, fn)

		tmpl, err := template.ParseFiles(src)
		checkError(err)

		f, err := os.Create(dest)
		checkError(err)

		err = tmpl.ExecuteTemplate(f, fn, ipt)
		checkError(err)
	}

	return tmpDir
}

func init() {
	RootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&url, "url", "u", "", "Required URL to shorten.")
	createCmd.Flags().StringVarP(&name, "name", "n", "", "Repository name, used as slug of shortened URL. Chosen randomly if empty.")
	createCmd.Flags().StringVarP(&cname, "cname", "c", "", "Optional CNAME record for this repository.")
	createCmd.Flags().BoolVarP(&isPrivate, "private", "p", false, "Make this repository private.")
}
