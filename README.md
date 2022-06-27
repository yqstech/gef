# goEasyFrame

#### 介绍
goEasyFrame是一款基于Go语言的低代码开发框架

#### 编译插件
    go build -o ../plugins/darwin SmsAm.go
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../plugins/linux SmsAm.go