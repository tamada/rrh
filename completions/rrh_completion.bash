__rrh_groups(){
    rrh group list --only-groupname
}

__rrh_repositories(){
    rrh list --only-repositoryname
}

__rrh_group_repo_forms(){
    rrh list --group-repository-form
}

__rrh_add() {
    if [[ "${1}" =~ ^\- ]]; then
        COMPREPLY=($(compgen -W "--group -g" -- "$1"))
    elif [ "$2" = "-g" ] || [ "$2" = "--group" ]; then
        groups="$(__rrh_groups)"
        COMPREPLY=($(compgen -W "$groups" -- "$1"))
    else
        compopt -o filenames
        COMPREPLY=($(compgen -d -- "$1"))
    fi
}

__rrh_clone() {
    if [[ "${1}" =~ ^\- ]]; then
        COMPREPLY=($(compgen -W "-g --group -d --dest" -- "$1"))
    elif [ "$2" = "-g" ] || [ "$2" = "--group" ]; then
        groups="$(__rrh_groups)"
        COMPREPLY=($(compgen -W "$groups" -- "$1"))
    elif [ "$2" = "-d" ] || [ "$2" = "--dest" ]; then
        compopt -o filenames
        COMPREPLY=($(compgen -d -- "$1"))
    fi
}

__rrh_config(){
    local rrhenvs="RRH_HOME RRH_CONFIG_PATH RRH_DATABASE_PATH RRH_DEFAULT_GROUP_NAME RRH_CLONE_DESTINATION RRH_ON_ERROR RRH_TIME_FORMAT RRH_AUTO_CREATE_GROUP RRH_AUTO_DELETE_GROUP RRH_SORT_ON_UPDATING"

    if [ "$4" = "$2" ]; then
        COMPREPLY=($(compgen -W "unset set list" -- $1))
    elif [ "$2" = "set" ] || [ "$2" = "unset" ]; then
        COMPREPLY=($(compgen -W rrhenvs $1))
    elif [ "$2" = "RRH_ON_ERROR" ] && [ "${COMP_WORDS[2]}" = "set" ]; then
        COMPREPLY=($(compgen -W "IGNORE WARN FAIL FAIL_IMMEDIATELY" -- $1))
    elif [ "$2" = "RRH_AUTO_CREATE_GROUP" -o "$2" = "RRH_AUTO_DELETE_GROUP" -o "$2" = "RRH_SORT_ON_UPDATING" ] && [ "${COMP_WORDS[2]}" = "set" ]; then
        COMPREPLY=($(compgen -W "true false" -- $1))
    fi
}

__rrh_export() {
    if [[ "$1" =~ ^\- ]]; then
        COMPREPLY=($(compgen -W "--no-indent --no-hide-home" -- "${cur}"))
    fi
}

__rrh_fetch() {
    if [[ "$1" =~ ^\- ]]; then
        COMPREPLY=($(compgen -W "-r --remote" -- "${cur}"))
    elif [ "$2" != "-r" ] && [ "$2" != "--remote" ]; then
        groups="$(__rrh_groups)"
        COMPREPLY=($(compgen -W "$groups" -- "$1"))
    fi
}

__rrh_fetch_all() {
    if [[ "$1" =~ ^\- ]]; then
        COMPREPLY=($(compgen -W "-r --remote" -- "${cur}"))
    fi
}

__rrh_group_add() {
    if [[ "$1" =~ ^\- ]]; then
        COMPREPLY=($(compgen -W "-d --desc -o --omit-list" -- "${cur}"))
    elif [ "$2" = "-o" ] || [ "$2" = "--omit-list" ]; then
        COMPREPLY=($(compgen -W "true false" -- "${cur}"))
    elif [ "$2" != "-d" ] && [ "$2" != "--desc" ]; then
        groups="$(__rrh_groups)"
        COMPREPLY=($(compgen -W "$groups" -- "${cur}"))
    fi
}

__rrh_group_list() {
    if [[ "$1" =~ ^\- ]]; then
        COMPREPLY=($(compgen -W "-d --desc -r --repository -o --only-groupname" -- "${cur}"))
    fi
}

__rrh_group_rm() {
    if [[ "$1" =~ ^\- ]]; then
        COMPREPLY=($(compgen -W "-f --force -i --inquiry -v --verbose" -- "${cur}"))
    elif [ "$2" != "-d" ] && [ "$2" != "--desc" ]; then
        groups="$(__rrh_groups)"
        COMPREPLY=($(compgen -W "$groups" -- "${cur}"))
    fi
}

__rrh_group_update() {
    if [[ "$1" =~ ^\- ]]; then
        COMPREPLY=($(compgen -W "-n --name -d --desc -o --omit-list" -- "${cur}"))
    elif [ "$2" = "-o" ] || [ "$2" = "--omit-list" ]; then
        COMPREPLY=($(compgen -W "true false" -- "${cur}"))
    elif [ "$2" != "-d" ] && [ "$2" != "--desc" ] && [ "$2" != "-n" ] && [ "$2" != "--name" ]; then
        groups="$(__rrh_groups)"
        COMPREPLY=($(compgen -W "$groups" -- "${cur}"))
    fi
}

__rrh_group() {
    if [ "$4" = "$2" ]; then
        COMPREPLY=($(compgen -W "add list rm update" -- "${cur}"))
        return 0
    else
        local subsub="${COMP_WORDS[2]}"
        case "${subsub}" in
            add)
                __rrh_group_add "$1" "$2" "$3" "$4" "$subsub"
                ;;
            list)
                __rrh_group_list "$1" "$2" "$3" "$4" "$subsub"
                ;;
            rm)
                __rrh_group_rm "$1" "$2" "$3" "$4" "$subsub"
                ;;
            update)
                __rrh_group_update "$1" "$2" "$3" "$4" "$subsub"
                ;;
        esac
    fi
}

__rrh_import() {
    if [[ "$1" =~ ^\- ]]; then
        COMPREPLY=($(compgen -W "--auto-clone --overwrite -v --verbose" -- "${cur}"))
    else
        _filedir '@(json)'
    fi
}

__rrh_list() {
    if [[ "$1" =~ ^\- ]]; then
        COMPREPLY=($(compgen -W "-d --desc -p --path -r --remote -A --all-entries -a --all -c --csv" -- "${cur}"))
    else
        groups="$(__rrh_groups)"
        COMPREPLY=($(compgen -W "$groups" -- "${cur}"))
    fi
}

__rrh_mv() {
    if [[ "$1" =~ ^\- ]]; then
        COMPREPLY=($(compgen -W "-v --verbose" -- "${cur}"))
    else
        groups="$(__rrh_groups)"
        gandr="$(__rrh_group_repo_forms)"
        COMPREPLY=($(compgen -W "$groups" -- "${cur}"))
        COMPREPLY+=($(compgen -W "$gandr" -- "${cur}"))
    fi
}

__rrh_path() {
    if [[ "$1" =~ ^\- ]]; then
        COMPREPLY=($(compgen -W "-m --partial-match -p --show-only-path" -- "${cur}"))
    else
        repos="$(__rrh_repositories)"
        COMPREPLY=($(compgen -W "$repos" -- "${cur}"))
    fi
}

__rrh_rm() {
    if [[ "$1" =~ ^\- ]]; then
        COMPREPLY=($(compgen -W "-i --inquiry -r --recursive -v --verbose" -- "${cur}"))
    else
        groups="$(__rrh_groups)"
        repos="$(__rrh_repositories)"
        gandr="$(__rrh_group_repo_forms)"
        COMPREPLY=($(compgen -W "$groups" -- "${cur}"))
        COMPREPLY+=($(compgen -W "$repos" -- "${cur}"))
        COMPREPLY+=($(compgen -W "$gandr" -- "${cur}"))
    fi
}

__rrh_status() {
    if [[ "$1" =~ ^\- ]]; then
        COMPREPLY=($(compgen -W "-b --branches -r --remote -c --csv" -- "${cur}"))
    else
        groups="$(__rrh_groups)"
        repos="$(__rrh_repositories)"
        COMPREPLY=($(compgen -W "$groups" -- "${cur}"))
        COMPREPLY+=($(compgen -W "$repos" -- "${cur}"))
    fi
}

__rrh_completions()
{
    local opts cur prev subcom
    _get_comp_words_by_ref -n : cur prev cword
    subcom="${COMP_WORDS[1]}"
    opts="add clone config export fetch fetch-all group import list mv prune rm status"

    case "${subcom}" in
        add)
            __rrh_add  "$cur" "$prev" "$cword" "$subcom"
            return 0
            ;;
        clone)
            __rrh_clone  "$cur" "$prev" "$cword" "$subcom"
            return 0
            ;;
        config)
            __rrh_config "$cur" "$prev" "$cword" "$subcom"
            return 0
            ;;
        export)
            __rrh_export "$cur" "$prev" "$cword" "$subcom"
            return 0
            ;;
        fetch)
            __rrh_fetch "$cur" "$prev" "$cword" "$subcom"
            return 0
            ;;
        fetch-all)
            __rrh_fetch_all "$cur" "$prev" "$cword" "$subcom"
            return 0
            ;;
        group)
            __rrh_group  "$cur" "$prev" "$cword" "$subcom"
            return 0
            ;;
        import)
            __rrh_import  "$cur" "$prev" "$cword" "$subcom"
            return 0
            ;;
        list)
            __rrh_list  "$cur" "$prev" "$cword" "$subcom"
            return 0
            ;;
        mv)
            __rrh_mv  "$cur" "$prev" "$cword" "$subcom"
            return 0
            ;;
        path)
            __rrh_path  "$cur" "$prev" "$cword" "$subcom"
            return 0
            ;;
        prune)
            return 0
            ;;
        rm)
            __rrh_rm  "$cur" "$prev" "$cword" "$subcom"
            return 0
            ;;
        status)
            __rrh_status  "$cur" "$prev" "$cword" "$subcom"
            return 0
            ;;
    esac
    if [ "$cword" -eq "1" ]; then
        if [[ "$cword" =~ ^\- ]]; then
            COMPREPLY=($(compgen -W "-h --help -v --version" ${cur}))
        else
            COMPREPLY=($(compgen -W "${opts}" -- ${cur}))
        fi
    fi
}

__cdrrh_completions() {
    local opts cur prev subcom
    _get_comp_words_by_ref -n : cur prev cword
    repos="$(__rrh_repositories)"
    COMPREPLY+=($(compgen -W "$repos" -- "${cur}"))
}

complete -F __rrh_completions rrh
complete -F __cdrrh_completions cdrrh
