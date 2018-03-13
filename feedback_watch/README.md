# Feedback Watch

Night gathers, and now my watch begins...

## Build

```
$ go build -o feedback_watch ./cmd/feedback_watch
```

```
$ GOOS=linux go build -o feedback_watch-linux ./cmd/feedback_watch
```

## Usage

### Feedback 值班

```
$ ./feedback_shift feedback_watch.toml
```

```
$ cat feedback_watch.toml

tasks = [
  "feedback_watch"
]

[bearychat]
rtm-token = "63b4c696bdd4a64fdf1da0d2649f1063"
    [bearychat.github-users]
    xtang = "tangxm"

[airtable]
api-key = "<API_KEY>"
base = "appoZUJmA9uVQJv3E"

[feedback_watch.shift]
incharge-table = "值日同学"
misconfig-vchannel-id = "=bw54q"
feedback-vchannel-id = "=bw54q"
shift-table = "排期"
shift-table-view = "viwE1RwWWGOoHzIcE"
template = """今天值日：%s
值日排期：%s [详细](https://airtable.com/shrODz6NSPLykxX03/tblAEMtkfuU84aZSY)
- 每封邮件回复记得抄送 `support@bearyinnovative.com`
- 遇到常见问题可以参照 [用户反馈](http://home.bearychat.com/w/%E7%94%A8%E6%88%B7%E5%8F%8D%E9%A6%88/)"""
```
