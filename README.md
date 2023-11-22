# confluence-proxy-image-go
[![toad](https://img.shields.io/badge/yantao-confluenceProxy-FF575D.svg)](https://gitlab.qunhequnhe.com/it/confluence-proxy-image-go/)
[![toad](https://img.shields.io/badge/configuration-Toad-40a9ff.svg)](http://coops.qunhequnhe.com/toad/#/)


## build
```shell
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .
```

## 配置文件介绍
```yaml
urlPatterns:
  - info:
      pattern: pages/viewpage\.action\?pageId=(\d+)&preview=/(\d+)/(\d+)/([^&]+)
      fields:
        type: contentId
        pageId: 1
        contentId: 3
```
- `pattern`: URL 的正则表达式，如果匹配不到则返回空白
- `type`: 根据什么字段来寻找这个附件图片
- `pageId`: 页面 ID 所在链接的正则匹配的第几个子项
- `contengId` 根据 type 主键来获取链接中的第几个子项的值

## 使用
原路径：https://confluence.xxxx.com/pages/viewpage.action?pageId=80733518599&preview=/80733518599/80732910351/image2023-8-9_11-31-51.png
新路径：自己不熟的域名或者IP/pages/viewpage.action?pageId=80733518599&preview=/80733518599/80732910351/image2023-8-9_11-31-51.png