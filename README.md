## point

> Create and manage shortened URLs with GitHub pages.

point lets you create, view, and manage shortened URLs. All pages are hosted on GitHub Pages and redirection is done with HTML5's [`http-equiv` refresh attribute](https://developer.mozilla.org/en/docs/Web/HTML/Element/meta#attr-http-equiv).

### Demo

- _Coming soon._

### Installation

  - You should have Go [installed](https://golang.org/doc/install) and [configured](https://golang.org/doc/install#testing).

  - Install with Go:

    ```sh
    $ go get -v github.com/kshvmdn/point/...
    $ point --help
    ```

  - Or, install directly via source:

    ```sh
    $ git clone https://github.com/kshvmdn/point.git $GOPATH/src/github.com/kshvmdn/point
    $ cd $_
    $ make install
    $ point --help
    ```

### Usage

  - You should export your personal GitHub access token as `POINT_ACCESS_TOKEN`. You can request one [here](https://github.com/settings/tokens) with the `repo` and `delete_repo` permissions.

  - View the help dialogue with the `--help` flag. View the specific help dialogue for each command by running `point [command] --help`.

    ```console
    $ point --help
    usage: point [<flags>] <command> [<args> ...]

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
    $ point create --help
    usage: point create [<flags>] <url>

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
    $ point list --help
    usage: point list

    Print a list of active entries.

    Flags:
      --help     Show context-sensitive help (also try --help-long and --help-man).
      --verbose  Show detailed output.
      --version  Show application version.

    ```

  - **remove**

    ```console
    $ point remove --help
    usage: point remove <repo>...

    Remove one or more entries.

    Flags:
      --help     Show context-sensitive help (also try --help-long and --help-man).
      --verbose  Show detailed output.
      --version  Show application version.

    Args:
      <repo>  List of entries to remove.

    ```

### Contribute

This project is completely open source. Feel free to [open an issue](https://github.com/kshvmdn/point/issues) or [create a pull request](https://github.com/kshvmdn/point/pulls).

Before submitting code, please ensure that tests are passing and the linter is happy. The following commands may be of use, refer to the [Makefile](./Makefile) to see what they do.

```sh
$ make install \
       get-tools
$ make fmt \
       vet \
       lint
$ make test \
       coverage
$ make bootstrap-dist \
       dist
```

### License

point source code is released under the [MIT license](./LICENSE).
