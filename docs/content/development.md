---
title: ":pencil2: Development"
---

## Database

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
            repository_desc: 'description of the repository.',
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

