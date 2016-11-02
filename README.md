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

[monitor-stale-issues]
github-token = "xxx"

[[moniitor-stale-issues.repos]]
repo = bearyinnovative/snitch
users = {
  bcho = "hbc"
  xtang = "tangxm"
}
bearychat-webhook = "https://hook.bearychat.com/incoming/xxx"

[[moniitor-stale-issues.repos]]
repo = bearyinnovative/pensieve
users = {
  L42y = l42y
  xtang = "tangxm"
}
bearychat-webhook = "https://hook.bearychat.com/incoming/xxx"
```

## Build

```
$ make build
```

## Test

```
$ make test
```
