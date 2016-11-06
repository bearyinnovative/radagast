# Radagast

Tender Of Beasts

## Usage

```
$ ls
radagast radagast.toml
$ radagast
```

## radagast.toml

```toml
tasks = [
  "monitor-stale-issues"
]

[github]
api-token = ""

[bearychat]
rtm-token = "123"

    [bearychat.users]
    bcho = "hbc"
    xtang = "tangxm"

[monitor-stale-issues]
    [[monitor-stale-issues.repos]]
    repo = "bearyinnovative/snitch"
    bearychat-vchannel-id = "=bw52P"

    [[monitor-stale-issues.repos]]
    repo = "bearyinnovative/pensieve"
    bearychat-vchannel-id = "=bw52Q"
```

## Build

```
$ make build
```

## Test

```
$ make test
```
