[![Build Status](https://travis-ci.org/tamada/rrh.svg?branch=master)](https://travis-ci.org/tamada/rrh)
[![Coverage Status](https://coveralls.io/repos/github/tamada/rrh/badge.svg?branch=master)](https://coveralls.io/github/tamada/rrh?branch=master)
[![codebeat badge](https://codebeat.co/badges/15e04551-d448-4ad3-be1d-e98b1e586f1a)](https://codebeat.co/projects/github-com-tamada-rrh-master)
[![go report](https://goreportcard.com/badge/github.com/tamada/rrh)](https://goreportcard.com/report/github.com/tamada/rrh)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://github.com/tamada/rrh/blob/master/LICENSE)

# RRH

RRH is a simple manager for git repositories.

[Table of contents](#table_of_contents)

## Description

There are too many repositories.
We love programming; however, to manage many repositories is quite hard and bothersome tasks.
Therefore, we built a headquarter for managing the git repositories, named RRH.
RRH manages repositories by categorizing in groups and execute git command to the groups.

I know the tool [ghq](https://github.com/motemen/ghq), manages the git repositories.
However, I cannot use it for the following reasons.

1. there are quite many repositories in my home directory.
    * To start using ghq, we clone the repositories.
      However, I did not accept to clone all of the repositories.
2. The location of repositories is fixed in the config file and is accepted only one location.
    * I decide the directory layout in my home directory.

Additionally, I edit several repositories in a day, when I work hard.
Consequently, the progress of each repository is obscured; I cannot remember a lot of things.
Therefore, it is glad to see the last modified date of branches.

RRH is now growing. Please hack RRH itself.

## Installation

### Homebrew

Install rrh via [Homebrew](https://brew.sh), simply run:

```sh
$ brew tap tamada/rrh
$ brew install rrh
```


### Golang

To install cli, simply run:

```
$ go get git@github.com/tamada/rrh.git
```

### Requirements

* Runtime
    * Bash 4.x or after, for completion.
        * [zsh](http://www.zsh.org/)?, and [fish](https://fishshell.com/)?, I do not use them, so I do not know.
        * For macOS user, the default shell of the macOS is bash 3.x, therefore, the completion is not work enough.
             * `rrh` is maybe work on Windows, and Linux. I do not use them.
* Development
    * [golang/dep](https://github.com/golang/dep)
    * See `Gopkg.toml`


## Usage

```sh
Usage: rrh [--version] [--help] <command> [<args>]

Available commands are:
    add          add repositories on the local path to RRH.
    clone        run "git clone" and register it to a group.
    config       set/unset and list configuration of RRH.
    export       export RRH database to stdout.
    fetch        run "git fetch" on the repositories of the given groups.
    fetch-all    run "git fetch" in the all repositories.
    group        add/list/update/remove groups.
    list         print managed repositories and their groups.
    mv           move the repositories from groups to another group.
    path         print paths of specified repositories.
    prune        prune unnecessary repositories and groups.
    rm           remove given repository from database.
    status       show git status of repositories.
```

### Subcommands

#### `rrh add`

Registers the repositories which specified the given paths to the RRH database and categorize to the group (Default `no-group`, see [`RRH_DEFAULT_GROUP_NAME`](#rrh_default_group_name)).

```sh
rrh add [OPTIONS] <REPOSITORY_PATHS...>
OPTIONS
    -g, --group <GROUP>        add repository to RRH database.
    -r, --repository-id <ID>   specified repository id of the given repository path.
                               Specifying this option fails with multiple arguments.
ARGUMENTS
    REPOSITORY_PATHS           the local path list of the git repositories.
```

#### `rrh clone`

Runs `git clone` command and registers the cloned repository to RRH database.
The following steps identify the id of the repository.

1. If the length of `REMOTE_REPOS` is 1, and `DEST` exists, then the last entry of `REMOTE_REPOS` is repository id by eliminating the suffix `.git`.
3. If the length of `REMOTE_REPOS` is 1, and `DEST` does not exist, then the last entry of `DEST` is repository id.
2. If the length of `REMOTE_REPOS` is greater than 1, then the last entry of each `REMOTE_REPOS` is repository ids by eliminating the suffix `.git`.

```sh
rrh clone [OPTIONS] <REMOTE_REPOS...>
OPTIONS
    -g, --group <GROUP>   print managed repositories categoried in the group.
    -d, --dest <DEST>     specify the destination. Default is the current directory.
ARGUMENTS
    REMOTE_REPOS          repository urls
```

The destination of cloned repository is located based on [`RRH_CLONE_DESTINATION`](#rrh_clone_destination)

#### `rrh config`

Handles the operations of configuration/environment variables.
This subcommand requires sub-sub-command.
If sub-sub-command was not specified, it runs `list` sub-sub-command.

```sh
rrh config <COMMAND> [ARGUMENTS]
COMMAND
    set <ENV_NAME> <VALUE>  set ENV_NAME to VALUE
    unset <ENV_NAME>        reset ENV_NAME
    list                    list all of ENVs (default)
```

#### `rrh export`

Exports the data of RRH database by JSON format.

```sh
rrh export [OPTIONS]
OPTiONS
    --no-indent      print result as no indented json
    --no-hide-home   not replace home directory to '${HOME}' keyword
```

#### `rrh fetch`

Runs `git fetch` command in the repositories of the specified group.

```sh
rrh fetch [OPTIONS] [GROUPS...]
OPTIONS
    -r, --remote <REMOTE>   specify the remote name. Default is "origin."
ARGUMENTS
    GROUPS                  run "git fetch" command on each repository on the group.
                            if no value is specified, run on the default group.
```

#### `rrh fetch-all`

Runs `git fetch` command in all repositories of managing in RRH.
This command may make heavy network traffic; therefore, we do not recommend to run.

```sh
rrh fetch-all [OPTIONS]
OPTIONS
    -r, --remote <REMOTE>   specify the remote name. Default is "origin."
```

#### `rrh group`

Handles the operations of groups of RRH.
This subcommand requires sub-sub-command.
If sub-sub-command was not specified, it runs `list` sub-sub-command.

```sh
rrh group <SUBCOMMAND>
SUBCOMMAND
    add       add new group.
    list      list groups (default).
    of        shows groups of the specified repository.
    rm        remove group.
    update    update group.
```

##### `rrh group add`

Adds new group to the RRH database.

```sh
rrh group add [OPTIONS] <GROUPS...>
OPTIONS
    -d, --desc <DESC>        gives the description of the group.
    -o, --omit-list <FLAG>   gives the omit list flag of the group.
ARGUMENTS
    GROUPS                   gives group names.
```

##### `rrh group list`

Displays group list.

```sh
rrh group list [OPTIONS]
OPTIONS
    -d, --desc             show description.
    -r, --repository       show repositories in the group.
    -o, --only-groupname   show only group name. This option is prioritized.
```

##### `rrh group of`

Displays group of the specified repositories.

```sh
rrh group of <REPOSITORY_ID>
ARGUMENTS
    REPOSITORY_ID     show the groups of the repository.
```

##### `rrh group rm`

Removes groups.

```sh
rrh group rm [OPTIONS] <GROUPS...>
OPTIONS
    -f, --force      force remove.
    -i, --inquery    inquiry mode.
    -v, --verbose    verbose mode.
ARGUMENTS
    GROUPS           target group names.
```

##### `rrh group update`

Update the information of specified group.

```sh
rrh group update [OPTIONS] <GROUP>
OPTIONS
    -n, --name <NAME>        change group name to NAME.
    -d, --desc <DESC>        change description to DESC.
    -o, --omit-list <FLAG>   change omit-list of the group. FLAG must be "true" or "false".
ARGUMENTS
    GROUP                    update target group names.
```

#### `rrh import`

Import the database to the local environment.

```sh
rrh import [OPTIONS] <DATABASE_JSON>
OPTIONS
    --auto-clone    clone the repository, if paths do not exist.
    --overwrite     replace the local RRH database to the given database.
    -v, --verbose   verbose mode.
ARGUMENTS
    DATABASE_JSON   the exported RRH database.
```

#### `rrh list`

Prints the repositories of managing in RRH.

```sh
rrh list [OPTIONS] [GROUPS...]
OPTIONS
    -d, --desc          print description of group.
    -p, --path          print local paths (default).
    -r, --remote        print remote urls.
    -A, --all-entries   print all entries of each repository.

    -a, --all           print all repositories, no omit repositories.
    -c, --csv           print result as csv format.
ARGUMENTS
    GROUPS    print managed repositories categorized in the groups.
              if no groups are specified, all groups are printed.
```

#### `rrh mv`

Move repositories to another group.

```sh
rrh mv [OPTIONS] <FROMS...> <TO>
OPTIONS
    -v, --verbose   verbose mode

ARGUMENTS
    FROMS...        specifies move from, formatted in <GROUP_NAME/REPO_ID>, or <GROUP_NAME>
    TO              specifies move to, formatted in <GROUP_NAME>
```

#### `rrh path`

Prints paths of the specified repositories.

```sh
rrh path [OPTIONS] <REPOSITORIES...>
OPTIONS
    -m, --partial-match        treats the arguments as the patterns.
    -r, --show-repository-id   show repository name.
ARGUMENTS
    REPOSITORIES               repository ids.
```

#### `rrh prune`

Deletes unnecessary groups and repositories.
The unnecessary groups are no repositories in them.
The unnecessary repositories are to have an invalid path.


```sh
rrh prune
```

#### `rrh rm`

Removes the specified groups, repositories, and relations.
If the group has entries is removed by specifying the option `--recursive.`

```sh
rrh rm [OPTIONS] <REPO_ID|GROUP_ID|GROUP_ID/REPO_ID...>
OPTIONS
    -i, --inquiry       inquiry mode.
    -r, --recursive     recursive mode.
    -v, --verbose       verbose mode.

ARGUMENTS
    REPOY_ID            repository name for removing.
    GROUP_ID            group name. if the group contains repositories,
                        removing will fail without '-r' option.
    GROUP_ID/REPO_ID    remove the relation between the given REPO_ID and GROUP_ID.
```

#### `rrh status`

Prints the last modified times of each branch in the repositories of the specified group.

```sh
rrh status [OPTIONS] [GROUPS|REPOS...]
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

### Environment variables

We can see those variables by running `rrh config` sub-command.

#### `RRH_HOME`

* specifies the location of the RRH database and config file.
* Default: `/Users/tamada/.rrh`

#### `RRH_CONFIG_PATH`

* specifies the location of the location path.
    * RRH ignores to specify `RRH_CONFIG_PATH` in the config file.
      This variable availables only environment variable.
* Default: `${RRH_HOME}/config.json`

#### `RRH_DATABASE_PATH`

* specifies the location of the database path.
* Default: `${RRH_HOME}/database.json`

#### `RRH_DEFAULT_GROUP_NAME`

* specifies the default group name.
* Default: `no-group`

#### `RRH_CLONE_DESTINATION`

* specifies the destination by cloning the repository.
* Default: `.`

#### `RRH_ON_ERROR`

* specifies the behaviors of RRH on error.
* Default: `WARN`
* Available values: `FAIL_IMMEDIATELY`, `FAIL`, `WARN`, and `IGNORE`
    * `FAIL_IMMEDIATELY`
        * reports error immediately and quits RRH with a non-zero status.
    * `FAIL`
        * runs through all targets and reports errors if needed, then quits RRH with a non-zero status.
    * `WARN`
        * runs through all targets and reports errors if needed, then quits RRH successfully.
    * `IGNORE`
        * runs all targets and no reports errors.

#### `RRH_TIME_FORMAT`

* specifies the time format for `status` command.
* Default: `relative`
* Available value: `relative` and the time format for Go lang.
    * `relative`
        * shows times by humanized format (e.g., 2 weeks ago)
    * Other strings
        * regard as formatting layout and give to `Format` method of the time.
            * see [Time.Format](https://golang.org/pkg/time/#Time.Format), for more detail.

#### `RRH_AUTO_CREATE_GROUP`

* specifies to create the group when the not existing group was specified, and it needs to create.
* Default: false

#### `RRH_AUTO_DELETE_GROUP`

* specifies to delete the group when some group was no more needed.
* Default: false

#### `RRH_SORT_ON_UPDATING`

* specifies to sort database entries on updating database.
* Default: false

#### `RRH_COLOR`

* specifies the colors of the output.
* Default: `""` (empty string)
* Format: `"repository:fg=<COLOR>;bg=<COLOR>;op=<STYLE>+GROUP:fg=<COLOR>;bg=<COLOR>;op=<STYLE>"`
    * Available `COLOR`s
        * red, cyan, gray, blue, black, green, white, yellow, magenta.
    * Available `STYLE`s
        * bold, underscore.
    * Delimiter of repository and group is `+`, delimiter of type and value is `:`, delimiter of each label is `;`, and delimiter of each value is `,`.
* Examples:
    * `RRH_COLOR: repository:fg=red+group:fg=cyan;op=bold,underscore`
        * Repository: red, Group: cyan in bold with underscore.
* Note
    * The colorized output does not support to arrange the output indentation.

#### `RRH_ENABLE_COLORIZED`

* specifies to colorize the output. The colors of output were specified on [`RRH_COLOR`](#rrh_color)
* Default: true

### Database

The database for managed repositories is formatted in JSON.
The JSON format is as follows.
The JSON file is placed on `$RRH_HOME/database.json`.
If `$RRH_HOME` was not set, `$HOME/.rrh` is used as `$RRH_HOME`.
Also, the configuration file is on `$RRH_HOME/config.json`

```js
{
    last-modified: '2019-01-01T',
    repositories: [
        {
            repository_id: 'rrh', // unique key of repository.
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
            group_name: 'no-group',
            group_desc: 'The description of the group.',
            omit_list: false
        },
        ....
    ],
    relations: [
        {
            repository_id: 'rrh',
            group_name: 'no-group'
        }
    ]
}
```

## Utilities

### `cdrrh`

changes directory to the specified repository.

```sh
cdrrh(){
    path=$(rrh path $1)
    if [ $? -eq 0 ]; then
        cd $path
        pwd
    else
        echo "$1: repository not found"
    fi
}
```

### `rrhpeco`

list repositories, and filtering them by [`peco`](https://github.com/peco/peco),
then change directory to the filtering result.

```sh
rrhpeco(){
  csv=$(rrh list --path --csv | peco)
  cd $(echo $csv | awk -F , '{ print $3 }')
  pwd
}
```

## License

[Apache License version 2.0](https://github.com/tamada/rrh/blob/master/LICENSE)

## Development Policy

* Separate `foo_cmd.go` and `foo.go` for implementing `foo` command.
    * `foo_cmd.go` includes functions of cli.
    * `foo.go` includes essential functions for `foo`.
* Call `fmt.Print` methods only `foo_cmd.go` file.
* Create test for `foo.go`.

## Why the project name RRH

At first, the name of this project was GRIM (Git Repository Integrated Manager).
However, the means of `grim` is not good, and there are many commands which start with `gr`.
Therefore, we changed the project name to RRH.
RRH means "Repositories, Ready to Hack," is not the abbreviation of the Red Riding Hood.

## Discussion

[![Gitter](https://img.shields.io/badge/Gitter-Join_Chat-red.svg)](https://gitter.im/rrh_git/community)

Join our Gitter channel if you have any problem or suggestions to Rrh.

[![Gitter misc_ja](https://img.shields.io/badge/Gitter-For_Japanese-red.svg)](https://gitter.im/rrh_git/misc_ja)

For Japanese user, `misc_ja` channel has discussions in Japanese.
The public language of other channels and GitHub pages are English.

## Table of Contents

* [Description](#description)
* [Installation](#installation)
    * [Homebrew](#homebrew)
    * [Golang](#golang)
    * [Requirements](#requirements)
* [Usage](#usage)
    * [Subcommands](#subcommands)
        * [`rrh add`](#rrh-add)
        * [`rrh clone`](#rrh-clone)
        * [`rrh config`](#rrh-config)
        * [`rrh export`](#rrh-export)
        * [`rrh fetch`](#rrh-fetch)
        * [`rrh fetch-all`](#rrh-fetch-all)
        * [`rrh group`](#rrh-group)
        * [`rrh import`](#rrh-import)
        * [`rrh list`](#rrh-list)
        * [`rrh mv`](#rrh-mv)
        * [`rrh path`](#rrh-path)
        * [`rrh prune`](#rrh-prune)
        * [`rrh rm`](#rrh-rm)
        * [`rrh status`](#rrh-status)
    * [Environment variables](#environment-variables)
    * [Database](#database)
* [Utilities](#utilities)
    * [`cdrrh`](#cdrrh)
    * [`rrhpeco`](#rrhpeco)
* [License](#license)
* [Development Policy](#development-policy)
* [Why the project name RRH?](#why-the-project-name-rrh)
* [Discussion](#discussion)
