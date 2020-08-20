# com.dk.user
用户平台

# 技术栈
- golang
- go-kit
- nacos

# 文档目录

> [注册逻辑设计](https://github.com/dingkegithub/com.dk.user/blob/master/sidecar/discovery/readme.md)

# 环境安装

```
# install protobuf
$ brew install protobuf

# install grpc
$ export GO111MODULE=on 
$ go get -u google.golang.org/grpc
$ go get -u github.com/golang/protobuf/protoc-gen-go

# set env
$ echo "export PATH="$PATH:$(go env GOPATH)/bin"" >> ~/.zshrc
```
