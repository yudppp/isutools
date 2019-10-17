# ISUTOOLS

いい感じにSpeedUpするために色々調べるためのツール群です。


## Tools

### profile

pprofの結果をSlackに送る

### tracereporter

datadogのAPMの結果を出力するやつ

### utils/throttle

_.throttle みたいなやつ
Onceと同じInterfaceのfunctionを持つよ

### utils/slackcat

Slackcat wrapper

## TODO tools

- slack log writter

## Setup

```
go get github.com/yudppp/isutools/...
wget https://github.com/bcicen/slackcat/releases/download/v1.5/slackcat-1.5-linux-amd64 -O slackcat
sudo mv slackcat /usr/local/bin/
sudo chmod +x /usr/local/bin/slackcat
slackcat --configure
```


## dependency

- https://github.com/tenntenn/isucontools
- https://github.com/najeira/measure
- http://slackcat.chat/
- http://github.com/rcrowley/go-metrics
- http://gopkg.in/DataDog/dd-trace-go.v1

## Links

- https://gist.github.com/catatsuy/e627aaf118fbe001f2e7c665fda48146
- https://github.com/tohutohu/isucon9/blob/master/Makefile
