## Utilities

### `cdrrh`

changes directory to the specified repository.

```sh
cdrrh(){
    path=$(rrh repository list --path $1)
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
