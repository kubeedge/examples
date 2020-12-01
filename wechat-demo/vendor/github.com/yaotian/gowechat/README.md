# GoWechat 微信平台API [![Build Status](https://travis-ci.org/yaotian/gowechat.svg?branch=master)](https://travis-ci.org/yaotian/gowechat)  [![GoDoc](http://godoc.org/github.com/yaotian/gowechat?status.svg)](http://godoc.org/github.com/yaotian/gowechat)  [![Go Report Card](https://goreportcard.com/badge/github.com/yaotian/gowechat)](https://goreportcard.com/report/github.com/yaotian/gowechat)

支持微信公众平台(订阅号，服务号),微信商家平台,微信开放平台,微信企业号;  
与golang的各种框架(beego,gin,net/http)无缝链接

## Quick Start

#### Download and install
    go get github.com/yaotian/gowechat

#### Run examples
    cd ./examples/beego
    go run beego.go

## Features
* mp 微信公众平台API，网页授权，消息发送，菜单等
* mch 微信商户平台API，支付，发红包等
* open 微信开放平台API
* corp 微信企业号API
* mini 微信小程序

## Documentation
* [中文文档](https://yaotian.github.io/gowechat/)

## 附录，目录所对应的平台

目录| 对应 |
----|------|
/mp | 微信公众平台(订阅号，服务号)  |
/corp | 微信企业号  |
/mch | 微信商家平台  |
/open| 微信开放平台|
/mini| 微信小程序|

## Community

## License
MIT