# ISUTOOLS

いい感じにSpeedUpするために色々調べるためのツール群です。


## Tools

### profile

pprofの結果をSlackに送る


### utils/throttle

_.throttle みたいなやつ
Onceと同じInterfaceのfunctionを持つよ

### utils/slackcat

Slackcat wrapper

## TODO tools

- access logger
- sql profiler
- slack log writter

## Setup

```
go get github.com/yudppp/isutools/...
wget https://github.com/bcicen/slackcat/releases/download/v1.5/slackcat-1.5-linux-amd64 -O slackcat
sudo mv slackcat /usr/local/bin/
sudo chmod +x /usr/local/bin/slackcat
slackcat --configure
```

## sqlstr

```
go get github.com/najeira/measure
go get github.com/tenntenn/isucontools
cat main.go | sqlstr 
```

## mesuregen

```
go get github.com/tenntenn/isucontools
cat main.go | mesuregen > main.go
```

## dependency

- https://github.com/tenntenn/isucontools
- https://github.com/najeira/measure
- http://slackcat.chat/

## Links

- https://gist.github.com/catatsuy/e627aaf118fbe001f2e7c665fda48146
- https://github.com/tohutohu/isucon9/blob/master/Makefile
