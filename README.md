# go-url-shortener

A simple URL shortener that uses GitHub as a backend. This uses the [`http-equiv` refresh attribute](https://developer.mozilla.org/en/docs/Web/HTML/Element/meta#attr-http-equiv) to do URL redirection.

## Demo

[![asciicast](https://asciinema.org/a/132016.png)](https://asciinema.org/a/132016)

## Usage

You should export your personal GitHub access token as `GO_URL_SHORTENER_ACCESS_TOKEN`. Request one [here](https://github.com/settings/tokens) with the `repo` and `delete_repo` permissions.

Basic usage:

  ```sh
  % git clone https://github.com/kashav/go-url-shortener && cd $_
  % go build ./cmd/shortener/main.go
  % ./main --help
  usage: go-url-shortener [<flags>] <command> [<args> ...]

  Create and manage shortened URLs with GitHub pages.


  Flags:
    --[no-]help     Show context-sensitive help (also try --help-long and --help-man).
    --[no-]verbose  Show detailed output.

  Commands:
  help [<command>...]
      Show help.

  create [<flags>] <url>
      Create a new entry.

  list
      Print a list of active entries.

  remove <repo>...
      Remove one or more entries.


  % ./main create --help
  usage: go-url-shortener create [<flags>] <url>

  Create a new entry.


  Flags:
        --[no-]help     Show context-sensitive help (also try --help-long and --help-man).
        --[no-]verbose  Show detailed output.
    -c, --cname=CNAME   Optional CNAME record for this repository.
    -n, --name=NAME     Endpoint for the shortened URL, chosen randomly if empty.
    -p, --[no-]private  Make this repository private.
    -s, --[no-]subdir   Use a subdirectory in a pre-existing repository, instead of creating a new
                        repository.
    -r, --repo=REPO     A pre-existing repository to be used with the --subdir option (expects `foo/bar`
                        for `https://github.com/foo/bar`). Pushes to the default branch.

  Args:
    <url>  The URL to shorten.

  % ./main list --help
  usage: go-url-shortener list

  Print a list of active entries.


  Flags:
    --[no-]help     Show context-sensitive help (also try --help-long and --help-man).
    --[no-]verbose  Show detailed output.

  % ./main remove --help
  usage: go-url-shortener remove <repo>...

  Remove one or more entries.


  Flags:
    --[no-]help     Show context-sensitive help (also try --help-long and --help-man).
    --[no-]verbose  Show detailed output.

  Args:
    <repo>  List of entries to remove.
  ```

