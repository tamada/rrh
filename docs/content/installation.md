---
title: ":anchor: Installation"
---

## :beer: Homebrew

Install rrh via [Homebrew](https://brew.sh), simply run:

```sh
$ brew tap tamada/brew
$ brew install rrh
```


## Golang

To install cli, simply run:

```
$ go get git@github.com/tamada/rrh.git
```

## :hammer_and_wrench: Build from source codes

```sh
$ git clone https://github.com/tamada/rrh.git
$ cd rrh
$ make
```

## :white_check_mark: Requirements

* Runtime
    * Bash 4.x or after, for completion.
        * [zsh](http://www.zsh.org/)?, and [fish](https://fishshell.com/)?, I do not use them, so I do not know.
        * For macOS user, the default shell of the macOS is bash 3.x, therefore, the completion is not work enough.
             * `rrh` is maybe work on Windows, and Linux. I do not use them.
* Development
    * Go 1.12
    * See `go.mod`
