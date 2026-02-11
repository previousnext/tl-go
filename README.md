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
tl unsent
tl send
```

```
USAGE

    tl [command] [--flags]

COMMANDS

    add <key> <duration> [description]  Add a time entry
    delete <id>                         Delete a time entry
    edit <id> [--flags]                 Edit a time entry
    issues [--flags]                    List recent issues
    list [--flags]                      List all time entries
    send                                Send time entries to Jira
    setup [--flags]                     Setup tl database and configuration
    show <id>                           Show details of a time entry
    unsent                              List unsent time entries
```

### Acknowledgements

This project was inspired by [larowlan/tl](https://github.com/larowlan/tl) which is written in PHP and does much more than this version.
