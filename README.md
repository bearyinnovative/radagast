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
  "pullrequests"
]

[github]
api-token = ""

[bearychat]
rtm-token = "123"

    [bearychat.github-users]
    bcho = "hbc"
    xtang = "tangxm"

[pullrequests]
    [[pullrequests.repos]]
    repo = "bearyinnovative/snitch"
    bearychat-vchannel-id = "=bw52P"

    [[pullrequests.repos]]
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
