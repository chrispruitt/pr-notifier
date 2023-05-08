# pr-notifier

## Usage

```text
› go run main.go notify --help
    post open PRs to slack

    Usage:
    pr-notifier notify [flags]

    Flags:
    -a, --authors stringArray   bitbucket usernames, UUIDs, or authors of PRs
    -c, --channel string        slack channel
        --debug                 enable verbose logging
    -h, --help                  help for notify
    -p, --password string       bitbucket app password
    -t, --token string          slack token
    -u, --user string           bitbucket app username
```

### Docker

```sh
docker run -it chrispruitt/pr-notifier:latest \
    notify \
    --author some_bitbucket_author \
    --author some_bitbucket_author \
    --user some_bb_user \
    --channel XXXXXXXX \
    --password XXXXXXXXXXXXXX \
    --token xoxb-XXXXXXX
```

## Roadmap
- [ ] Handle bitbucket response type of "error" properly (when user or project is not found)