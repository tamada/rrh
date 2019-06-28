# rrh-new

## Requirements

* [`hub`](https://github.com/github/hub)
    * have to set `oauth_token` for creating repository, because `rrh-new` uses `hub create`.
    * see CONFIGURATION on [https://hub.github.com/hub.1.html](https://hub.github.com/hub.1.html).
* [`rrh`](https://github.com/tamada/rrh)

## Install

* Put the executable to some directory in `PATH` environment.

## Usage

* Run `rrh` with `new` command.

```sh
$ rrh new --help
rrh new [OPTIONS] <[ORGANIZATION/]REPOSITORY>
OPTIONS
    -g, --group <GROUP>         specifies group name.
    -H, --homepage <URL>        specifies homepage url.
    -p, --private               create a private repository.
    -d, --description <DESC>    specifies short description of the repository.
    -P, --parent-path <PATH>    specifies the destination path (default: '.').
    -h, --help                  print this message.
ARGUMENTS
    ORGANIZATION    specifies organization, if needed.
    REPOSITORY      specifies repository name, and it is directory name.
```
