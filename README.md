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

### userinfos

\* requires auth token for the server used

Will query Mastadon server's API for user infos based on their `username`

```bash
export BEARER_TOKEN=<YOUR PERSONAL BEARER TOKEN>
go run main.go userinfos -username dave
```
Expect:

<img src="img/userinfos.png" alt="ctl results for userinfos"/>

### hashtag

Will query Mastadon server's public API for latest post tagged with a specific hashtag

```bash
go run main.go hashtag -name cat
```

Expect:

<img src="img/tablemastodon.png" alt="ctl results for cat"/>

## freely available Mastodon apps

- [mastovue](https://mastovue.glitch.me/#/vis.social/federated/duck)
- [mastoview](http://www.unmung.com/mastoview)
