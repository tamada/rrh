[![Build Status](https://travis-ci.org/tamada/rrh.svg?branch=master)](https://travis-ci.org/tamada/rrh)
[![Coverage Status](https://coveralls.io/repos/github/tamada/rrh/badge.svg?branch=master)](https://coveralls.io/github/tamada/rrh?branch=master)
[![codebeat badge](https://codebeat.co/badges/15e04551-d448-4ad3-be1d-e98b1e586f1a)](https://codebeat.co/projects/github-com-tamada-rrh-master)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

# RRH

There are too many repositories.
We love programming; however, to manage many repositories is quite hard and bothersome tasks.
Therefore, we built a headquarter for managing the git repositories, named RRH (Repositories, Ready to Hack).
RRH manages repositories by categorizing in groups and execute git command to the groups.

RRH is now growing. Please hack RRH itself.

# Installation

To install cli, simply run:

```
$ go get git@github.com/tamada/rrh.git
```

# Usage

```sh
Usage: rrh [--version] [--help] <command> [<args>]

Available commands are:
    add          add repositories on the local path to RRH
    clone        run "git clone"
    config       set/unset and list configuration of RRH.
    export       export RRH database to stdout.
    fetch        run "git fetch" on the given groups
    fetch-all    run "git fetch" in the all repositories
    group        print groups.
    list         print managed repositories and their groups.
    list-all     print managed repositories and their groups.
    prune        prune unnecessary repositories and groups.
    rm           remove given repository from database.
    status       show git status of repositories.
```

## subcommands

### `rrh add`

```sh
rrh add [OPTIONS] <REPOSITORY_PATHS...>
OPTIONS
    -g, --group <GROUP>    add repository to RRH database.
ARGUMENTS
    REPOSITORY_PATHS       the local path list of the git repositories
```

### `rrh clone`

```sh
rrh clone [OPTIONS] <REMOTE_REPOS...>
OPTIONS
    -g, --group <GROUP>   print managed repositories categoried in the group.
    -d, --dest <DEST>     specify the destination.
ARGUMENTS
    REMOTE_REPOS          repository urls
```

### `rrh config`

```sh
rrh config <COMMAND> [ARGUMENTS]
COMMAND
    set <ENV_NAME> <VALUE>  set ENV_NAME to VALUE
    unset <ENV_NAME>        reset ENV_NAME
    list                    list all of ENVs (default)
```

### `rrh export`

```sh
rrh export [OPTIONS]
OPTiONS
    --no-indent    print result as no indented json (Default indented json)
```

### `rrh fetch`

```sh
rrh fetch [OPTIONS] [GROUPS...]
OPTIONS
    -r, --remote <REMOTE>   specify the remote name. Default is "origin."
ARGUMENTS
    GROUPS                  run "git fetch" command on each repository on the group.
                            if no value is specified, run on the default group.
```

### `rrh fetch-all`

```sh
rrh fetch-all [OPTIONS]
OPTIONS
    -r, --remote <REMOTE>   specify the remote name. Default is "origin."
```

### `rrh group`

```sh
rrh group <SUBCOMMAND>
SUBCOMMAND
    add       add new group.
    list      list groups (default).
    rm        remove group.
    update    update group
```

### `rrh list`

```sh
rrh list [OPTIONS] [GROUPS...]
OPTIONS
    -a, --all       print all (default).
    -d, --desc      print description of group.
    -p, --path      print local paths.
    -r, --remote    print remote urls.
                    if any options of above are specified, '-a' are specified.

    -c, --csv       print result as csv format.
ARGUMENTS
    GROUPS    print managed repositories categoried in the groups.
              if no groups are specified, default groups are printed.
```

### `rrh list-all`

```sh
rrh list-all [OPTIONS]
OPTIONS
    -a, --all       print all (default).
    -d, --desc      print description of group.
    -p, --path      print local paths.
    -r, --remote    print remote urls.
                    if any options of above are specified, '-a' are specified.

    -c, --csv       print result as csv format.
```

### `rrh prune`

```sh
rrh prune
```

### `rrh rm`

```sh
rrh rm [OPTIONS] <REPO_ID|GROUP_ID|REPO_ID/GROUP_ID...>
OPTIONS
    -i, --inquiry       inquiry mode.
    -r, --recursive     recursive mode.
    -v, --verbose       verbose mode.

ARGUMENTS
    REPOY_ID            repository name for removing.
    GROUP_ID            group name. if the group contains repositories,
                        removing will fail without '-r' option.
    GROUP_ID/REPO_ID    remove given REPO_ID from GROUP_ID.
```

### `rrh status`

```sh
rrh status [OPTIONS] [GROUPS||REPOS...]
OPTIONS
    -b, --branches  show the status of the local branches.
	-r, --remote    show the status of the remote branches.
    -c, --csv       print result in csv format.
ARGUMENTS
    GROUPS          target groups.
    REPOS           target repositories.
                    If no arguments were specified, this command
                    shows the result of default group.
```

# Database

The database for managed repositories is formatted in JSON.
The JSON format is as follows.
The JSON file is placed on `$RRH_ROOT/database.json`.
If `$RRH_ROOT` was not set, `$HOME` is used as `$RRH_ROOT`.
Also, configuration file is on `$RRH_ROOT/config.json`

```js
{
    last-modified: '2019-01-01T',
    repositories: [
        {
            repository_id: 'repository_id1', // unique key of repository.
            repository_path: 'absolute/path/of/repository',
            remotes: [
                {
                    Name: "origin",
                    URL: "git@github.com:tamada/rrh.git"
                }
            ]
        },
        ....
    ]
    groups: [
        {
            group_name: 'group_name',
            group_desc: 'The description of the group.'
            group_items: [ 'repository_id1', 'repository_id2', ... ]
        },
        ....
    ]
}
```

# Utilities

## rrcd

list repositories, and filtering them by [`peco`](https://github.com/peco/peco),
then change directory to the filtering result.

```sh
cdrrh(){
  csv=$(rrh list-all --path --csv | peco)
  cd $(echo $csv | awk -F , '{ print $3 }')
  pwd
}
```

# Development Policy

* Separate `foo_cmd.go` and `foo.go` for implementing `foo` command.
    * `foo_cmd.go` includes functions of cli.
    * `foo.go` includes essential functions for `foo`.
* Call `fmt.Print` methods only `foo_cmd.go` file.
* Create test for `foo.go`.

# Candidates of the Product Names

* grim (Git Repository Integrated Manager)
    * However, the means of grim is not good.
* gram (Git Repository Advanced Manager)
    * there are many commands which start with `gr`
* rrh (Repositories, ready to HACK)
    * No red riding hood.
    * RRH is good name for command, since no command starts with `rr`.

# Discussion

[![Gitter](https://img.shields.io/badge/Gitter-Join_Chat-red.svg)](https://gitter.im/rrh_git/community)

Join our Gitter channel if you have any problem or suggestions to Rrh.

[![Gitter misc_ja](https://img.shields.io/badge/Gitter-For_Japanese-red.svg)](https://gitter.im/rrh_git/misc_ja)

For Japanese user, `misc_ja` channel has discussions in Japanese.
The public language of other channels and GitHub pages is English.
