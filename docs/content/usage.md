---
title: ":fork_and_knife: Usage"
---

## Getting Started

`rrh` has various subcommands, however, `list` and `add` subcommand make you happy.

* `rrh list` shows managed repositories.
* `rrh add <REPO>` adds the given repository under the `rrh` management.
* type [`cdrrh`](../utilities/#cdrrh) on Terminal, then type TAB, TAB, TAB!


## Command references

```sh
rrh [GLOBAL OPTIONS] <SUB COMMANDS> [ARGUMENTS]
GLOBAL OPTIONS
    -h, --help                        print this message.
    -v, --version                     print version.
    -c, --config-file <CONFIG_FILE>   specifies the config file path.
AVAILABLE SUB COMMANDS:
    add          add repositories on the local path to rrh.
    clone        run "git clone" and register it to a group.
    config       set/unset and list configuration of rrh.
    export       export rrh database to stdout.
    fetch        run "git fetch" on the given groups.
    fetch-all    run "git fetch" in the all repositories.
    group        add/list/update/remove groups and show groups of the repository.
    help         print this message.
    import       import the given database.
    list         print managed repositories and their groups.
    mv           move the repositories from groups to another group.
    open         open folder or web page of the given repositories.
    prune        prune unnecessary repositories and groups.
    repository   manages repositories.
    rm           remove given repository from database.
    status       show git status of repositories.
    version      show version.
```

If the user specified an unknown subcommand (e.g., `rrh helloworld`), `rrh` treats it as an **external command**.
In that case, `rrh` searches an executable file named `rrh-helloworld` from the PATH environment variable.
If `rrh` found it, `rrh` executes it, if not found, `rrh` prints help and exit.

### Subcommands

#### `rrh add`

Registers the repositories which specified the given paths to the `rrh` database and categorize to the group (Default `no-group`, see [`RRH_DEFAULT_GROUP_NAME`](#rrh_default_group_name)).

```sh
rrh add [OPTIONS] <REPOSITORY_PATHS...>
OPTIONS
    -g, --group <GROUP>        add repository to rrh database.
    -r, --repository-id <ID>   specified repository id of the given repository path.
                               Specifying this option fails with multiple arguments.
ARGUMENTS
    REPOSITORY_PATHS           the local path list of the git repositories.
```

#### `rrh clone`

Runs `git clone` command and registers the cloned repository to `rrh` database.
The following steps identify the id of the repository.

1. If the length of `REMOTE_REPOS` is 1, and `DEST` exists, then the last entry of `REMOTE_REPOS` is repository id by eliminating the suffix `.git`.
3. If the length of `REMOTE_REPOS` is 1, and `DEST` does not exist, then the last entry of `DEST` is repository id.
2. If the length of `REMOTE_REPOS` is greater than 1, then the last entry of each `REMOTE_REPOS` is repository ids by eliminating the suffix `.git`.

```sh
rrh clone [OPTIONS] <REMOTE_REPOS...>
OPTIONS
    -g, --group <GROUP>   print managed repositories categorized in the group.
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
    info      show information of specified groups.
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

##### `rrh group info`

Show information of specified groups.

```sh
rrh group info <GROUPS...>
ARGUMENTS
    GROUPS           group names to show the information.
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

#### `rrh open`

Opens folder or web page of the given repositories.

```sh
rrh open [OPTIONS] <REPOSITORIES...>
OPTIONS
    -f, --folder     open the folder of the specified repository (Default).
    -w, --webpage    open the webpage of the specified repository.
    -h, --help       print this message.
ARGUMENTS
    REPOSITORIES     specifies repository names.
```

#### `rrh prune`

Deletes unnecessary groups and repositories.
The unnecessary groups are no repositories in them.
The unnecessary repositories are to have an invalid path.


```sh
rrh prune [OPTIONS]
OPTIONS
    -d, --dry-run    dry-run mode.
    -v, --verbose    verbose mode.
```

#### `rrh repository`

Prints/Updates the repository.

```sh
rrh repository <SUBCOMMAND>
SUBCOMMAND
    info [OPTIONS] <REPO...>     shows repository information.
    update [OPTIONS] <REPO...>   updates repository information.
    update-remotes [OPTIONS]     update all remote entries.
```

##### `rrh repository info`

prints the repository information.

```sh
rrh repository info [OPTIONS] [REPOSITORIES...]
    -G, --color     prints the results with color.
    -c, --csv       prints the results in the csv format.
ARGUMENTS
    REPOSITORIES    target repositories.  If no repositories are specified,
                    this sub command failed.
```

##### `rrh repository update`

update the information of the repository.

```sh
rrh repository update [OPTIONS] <REPOSITORY>
OPTIONS
    -i, --id <NEWID>     specifies new repository id.
    -d, --desc <DESC>    specifies new description.
    -p, --path <PATH>    specifies new path.
ARGUMENTS
    REPOSITORY           specifies the repository id.
```

##### `rrh repository update-remotes`

update remote entries in the all repositories.

```sh
rrh repository update-remotes [OPTIONS]
OPTIONS
    -d, --dry-run    dry-run mode.
    -v, --verbose    verbose mode.
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
* Default: `"repository:fg=red+group:fg=magenta+label:op=bold+configValue:fg=green"`
* Format: `"repository:fg=<COLOR>;bg=<COLOR>;op=<STYLE>+group:fg=<COLOR>;bg=<COLOR>;op=<STYLE>+label:fg=<COLOR>;bg=<COLOR>;op=<STYLE>+configValue:fg=<COLOR>;bg=<COLOR>;op=<STYLE>"`
    * Available `COLOR`s
        * red, cyan, blue, black, green, white, yellow, magenta.
    * Available `STYLE`s
        * bold, underscore.
    * Delimiter of repository, group and label is `+`, delimiter of type and value is `:`, delimiter of each label is `;`, and delimiter of each value is `,`.
* Examples:
    * `RRH_COLOR: repository:fg=red+group:fg=cyan;op=bold,underscore`
        * Repository: red, Group: cyan in bold with underscore.

#### `RRH_ENABLE_COLORIZED`

* specifies to colorize the output. The colors of output were specified on [`RRH_COLOR`](#rrh_color)
* Default: false
