# Radagast

Tender Of Beasts

## Usage

```
$ ls
radagast radagast.yml
$ radagast
```

## radagast.yml

```yaml
---
tasks:
  - monitor-stale-issues
  - ...

monitor-stale-issues:
  github:
    token: xxx
  repos:
    - repo: bearyinnovative/snitch
      users:g
        - bcho: hbc
        - xtang: tangxm
      bearychat:
        webhook: https://hook.bearychat.com/incoming/xxx
    - repo: bearyinnovative/pensieve
      users:
        - L42y: l42y
        - xtang: tangxm
      bearychat:
        webhook: https://hook.bearychat.com/incoming/xxx
```

## Build

```
$ make build
```

## Test

```
$ make test
```
