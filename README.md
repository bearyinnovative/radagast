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

[airtable]
api-key = ""
base = ""

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

[zhouhui]
repo = "bearyinnovative/snitch"

[zhouhui.issue-template]
labels = "meetup"
title = "[meetup] %s"
body = """
- 时间：%s 17:30 - 18:00
- 参与者：@yuanbohan @bcho @aphawk @unionx @xtang @shonenada @stwind

### 规则
- 会议开始前24小时列出议题，包括议题内容简述或需要解决的问题
- 会后更新各议题讨论结果或解决办法，若没有则延续至下一次会议
- 若有另外需要在这次会议讨论的议题，在回复内提出

### 议题
#### 工作进度同步

> 每人工作目前工作内容及进度

| name | content |
| --- | --- |
| @stwind | |
| @yuanbohan | |
| @aphawk |  |
| @unionx | |
| @bcho | |
| @shonenada | |


##### 发自 laosiji
"""
```

## Build

```
$ make build
```

## Test

```
$ make test
```
