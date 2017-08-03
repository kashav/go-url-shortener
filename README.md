## redir

> Create and manage shortened URLs with GitHub pages.

redir lets you create, view, and manage shortened URLs. All pages are hosted on GitHub Pages and redirection is done with HTML5's [`http-equiv` refresh attribute](https://developer.mozilla.org/en/docs/Web/HTML/Element/meta#attr-http-equiv).

### Demo

- _Coming soon._

### Installation

  - You should have Go [installed](https://golang.org/doc/install) and [configured](https://golang.org/doc/install#testing).

  - Install with Go:

    ```sh
    $ go get -v github.com/kshvmdn/redir/...
    $ redir --help
    ```

  - Or, install directly via source:

    ```sh
    $ git clone https://github.com/kshvmdn/redir.git $GOPATH/src/github.com/kshvmdn/redir
    $ cd $_
    $ make install
    $ redir --help
    ```

### Usage

  - You should export your personal GitHub access token as `REDIR_ACCESS_TOKEN`. You can request one [here](https://github.com/settings/tokens) with the `repo` and `delete_repo` permissions.

  - View the help dialogue with the `--help` flag. View the specific help dialogue for each command by running `redir [command] --help`.

    ```console
    $ redir --help
    usage: redir [<flags>] <command> [<args> ...]

    Create and manage shortened URLs with GitHub pages.

    Flags:
      --help     Show context-sensitive help (also try --help-long and --help-man).
      --verbose  Show detailed output.
      --version  Show application version.

    Commands:
      help [<command>...]
        Show help.

      create [<flags>] <url>
        Create a new entry.

      list
        Print a list of active entries.

      remove <repo>...
        Remove one or more entries.

    ```

  - **create**

    ```console
    $ redir create --help
    usage: redir create [<flags>] <url>

    Create a new entry.

    Flags:
          --help         Show context-sensitive help (also try --help-long and --help-man).
          --verbose      Show detailed output.
          --version      Show application version.
      -c, --cname=CNAME  Optional CNAME record for this repository.
      -n, --name=NAME    Endpoint for the shortened URL, chosen randomly if empty.
      -p, --private      Make this repository private.
      -s, --subdir       Use a subdirectory in a pre-existing repository, instead of creating a new repository.
      -r, --repo=REPO    A pre-existing repository to be used with the --subdir option (expects `foo/bar` for
                         `https://github.com/foo/bar`). Pushes to the default branch.

    Args:
      <url>  The URL to shorten.

    ```

  - **list**

    ```console
    $ redir list --help
    usage: redir list

    Print a list of active entries.

    Flags:
      --help     Show context-sensitive help (also try --help-long and --help-man).
      --verbose  Show detailed output.
      --version  Show application version.

    ```

  - **remove**

    ```console
    $ redir remove --help
    usage: redir remove <repo>...

    Remove one or more entries.

    Flags:
      --help     Show context-sensitive help (also try --help-long and --help-man).
      --verbose  Show detailed output.
      --version  Show application version.

    Args:
      <repo>  List of entries to remove.

    ```

### Contribute

This project is completely open source. Feel free to [open an issue](https://github.com/kshvmdn/redir/issues) with questions / suggestions / requests or [create a pull request](https://github.com/kshvmdn/redir/pulls) to contribute!

### License

redir source code is released under the [MIT license](./LICENSE).
