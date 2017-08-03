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
    $ make install && make
    $ ./redir --help
    ```

### Usage

  - You should export your personal GitHub access token as `REDIR_ACCESS_TOKEN`. You can request one [here](https://github.com/settings/tokens) with the `repo` and `delete_repo` permissions.

  - View the help dialogue with the `--help` flag. View the specific help dialogue for each command by running `redir [command] --help`.

    ```sh
    $ redir --help
    usage: redir [<flags>] <command> [<args> ...]

    Create and manage shortened URLs with GitHub pages.

    Flags:
      --help     Show context-sensitive help (also try --help-long and --help-man).
      --version  Show application version.

    Commands:
      help [<command>...]
        Show help.

      create [<flags>] <url>
        Create a new entry.

      list
        Print a list of active entries.

      remove <repo>...
        Remove an entry and delete the associated repository.
    ```

### Contribute

This project is completely open source. Feel free to [open an issue](https://github.com/kshvmdn/redir/issues) with questions / suggestions / requests or [create a pull request](https://github.com/kshvmdn/redir/pulls) to contribute!

### License

redir source code is released under the [MIT license](./LICENSE).
