
<img src="mastodonctl.png" alt="mastodonctl logo" width="96" height="96"/>

# mastodonctl
cli client for mastodon social media platform

## current available commands

### hashtag

Will query Mastadon.social public api for latest post tagged with a specific
hashtag

```bash
go run main.go hashtag -name cat
```

Expect:

```bash
Latest cat pic at this URL: https://files.mastodon.social/cache/media_attachments/files/109/604/600/038/040/501/original/f5b1b999e25b83f9.jpeg
```