[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

# GRIM

Git Repository Integration Manager

# Installation

To install cli, simply run:

```
$ go get github.com/tamadalab/grim
```

# Usage

```sh
NAME:
  grim - Git Repository Integrated Manager

USAGE:
  grim [GLOBAL OPTIONS] COMMAND [SUB OPTIONS] [ARGUMENTS...]

VERSION
  1.0.0

AUTHOR
  Haruaki Tamada

COMMAND
  list     print managed repository and its group.
  group    print groups.
  prune    remove invalid products, and groups have no entry from database.
  clone    clone from a remote repository.
  add      add given repositories to management database.
  fetch    run `git fetch` command in the all of repositories.

GLOBAL OPTIONS:
  -h, --help     show help
  -v, --version  print the version.
```

# Database

The database for managed repositories is formatted in JSON.
The JSON format is as follows.
The JSON file is placed on `$GRIM_ROOT/.grim.json`.
If `$GRIM_ROOT` was not set, `$HOME` is used as `$GRIM_ROOT`.

```js
{
    last-modified: '2019-01-01T',
    repositories: [
        {
            repository_id: 'repository_id1', // unique key of repository.
            repository_path: 'absolute/path/of/repository',
            repository_url: 'url/of/origin'
        },
        ....
    ]
    groups: [
        {
            group_name: 'group_name',
            group_desc: 'The description of the group.'
            group_items: [ 'repository_id1', repository_id2, ... ]
        },
        ....
    ]
}
```

# Discussion

![Gitter](https://img.shields.io/badge/Gitter-Join_Chat-red.svg)

Join our Gitter channel if you have any problem or suggestions to Grim.

