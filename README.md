
<img src="img/mastodonctl.png" alt="mastodonctl logo" width="196" height="196"/>

# mastodonctl
cli client for mastodon social media platform

## configurations

As an experienced user, you may want to customize your commandline-tool.

This is possible by editing [`conf.json`](conf.json) file

### configurable values

| field               | description                                  |
| ------------------- | -------------------------------------------- |
| ResultsDisplayCount | number of results displayed in your terminal |
| ApiUrl              | URL of targetted mastodon server             |

## current available commands

### hashtag

Will query Mastadon.social public api for latest post tagged with a specific

hashtag

```bash
go run main.go hashtag -name cat
```

Expect:

<img src="img/tablemastodon.png" alt="ctl results for cat"/>