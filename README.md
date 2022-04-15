# robSams

[![Build Status](https://app.travis-ci.com/robGoods/sams.svg?branch=master)](https://app.travis-ci.com/robGoods/sams)

[sam's blog](https://robgoods.github.io/sams/)

## 感谢

首先感谢Sam's在上海疫情期间，给我们的帮助，让我们在疫情期间依旧可以买到很多好的东西！

## 使用方式
在main.go的main函数中修改配置信息
```go
conf := dd.Config{
    AuthToken:    "xxxxx", //HTTP头部auth-token
    BarkId:       "",      //通知用的bark id，下载bark后从app界面获取, 如果不需要可以填空字符串
    FloorId:      1,       //1,普通商品 2,全球购保税 3,特殊订购自提 4,大件商品 5,厂家直供商品 6,特殊订购商品 7,失效商品
    DeliveryType: 2,       //1,急速达 2,全城配送
    Longitude:    "xxxxx", //HTTP头部longitude
    Latitude:     "xxxxx", //HTTP头部latitude
    Deviceid:     "xxxxx", //HTTP头部device-id
    Trackinfo:    `xxxxx`, // HTTP头部track-info
}
```
#### 参数说明

|参数名|说明|
|----- |-----|
|AuthToken|Sam's登录后的HTTP头部auth-token|
|BarkId|通知用的bark id,如果不需要可以填空字符串|
|FloorId|1,普通商品 2,全球购保税 3,特殊订购自提 4,大件商品 5,厂家直供商品 6,特殊订购商品 7,失效商品|
|DeliveryType|1,急速达 2,全城配送。目前上海地区只开放了全城配送故默认为2|
|Longitude|Sam's登录后的HTTP头部Longitude,可以不填写|
|Latitude|Sam's登录后的HTTP头部Longitude,可以不填写|
|Deviceid|Sam's登录后的HTTP头部Deviceid,可以不填写|
|Trackinfo|Sam's登录后的HTTP头部Trackinfo,可以不填写|

### BarkId

![bark.png](https://robgoods.github.io/sams/assets/bark.png)

## 运行

```sh
go run main.go
```

开始运行后按命令行提示操作即可。

![run.png](https://robgoods.github.io/sams/assets/run.png)

## 声明
本项目仅供学习交流，严禁用作商业行为，特别禁止黄牛加价代抢等！

因违法违规等不当使用导致的后果与本人无关，如有任何问题可联系本人删除！