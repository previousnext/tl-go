# TL

A command-line tool for logging time to Jira Worklogs.

## Setup

You will need your Jira URL, username, and an API token to authenticate with Jira. You can generate an API token from your Atlassian account.

```
tl setup \
  --jira-url https://your-domain.atlassian.net \
  --username yourusername \
  --token yourapitoken
```

## Usage

Basic workflow example:

```
tl add PROJ-123 2h "Worked on feature X"
tl review
tl send
```

Command reference:

```
USAGE

    tl [command] [--flags]

COMMANDS

    add <key> <time> [description] [--flags]  Add a time entry
    alias [command]                           Add, list, and delete command aliases.
    delete <id>                               Delete a time entry
    edit <id> [--flags]                       Edit a time entry
    issues [--flags]                          List recent issues
    list [--flags]                            List all time entries
    migrate                                   Migrate the database schema to the latest version.
    review                                    Review unsent time entries
    send                                      Send time entries to Jira
    setup [--flags]                           Setup tl database and configuration
    show <id>                                 Show details of a time entry
    summary [--flags]                         Show a summary of time spent per project category
```

### Adding bash completion

To enable bash completion for `tl`, you can run the following command on Ubuntu and restart your terminal:

```
tl completion bash > ~/.local/share/bash-completion/completions/tl.bash
```

### Acknowledgements

This project was inspired by [larowlan/tl](https://github.com/larowlan/tl) which is written in PHP and does much more than this version.
