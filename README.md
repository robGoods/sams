# robSams

[![Build status](https://ci.appveyor.com/api/projects/status/v5lt859vmjm3v9i5?svg=true)](https://ci.appveyor.com/project/iscod/sams)

[sam's blog](https://robgoods.github.io/sams/)

## 感谢

首先感谢Sam's在上海疫情期间，给我们的帮助，让我们在疫情期间依旧可以买到很多好的东西！

## 使用方式

```sh
go run main.go --authToken=xxxxx
```

> 如果没有go环境，可以在 [releases](https://github.com/robGoods/sams/releases) 下载编译好的文件，直接运行即可

#### 参数说明

```sh
$ go run main.go -h

Usage of ./sams:
  -authToken string
    	必选, Sam's App HTTP头部auth-token
  -barkId bark
    	可选，通知用的bark id, 可选参数
  -deliveryType int
    	可选，1 急速达，2， 全程配送 (default 2)
  -deviceId string
    	可选，HTTP头部device-id
  -floorId int
    	可选，1,普通商品 2,全球购保税 3,特殊订购自提 4,大件商品 5,厂家直供商品 6,特殊订购商品 7,失效商品 (default 1)
  -help
    	show help
  -latitude string
    	可选，HTTP头部latitude
  -longitude string
    	可选，HTTP头部longitude
  -payMethod int
    	可选，1,微信 2,支付宝 (default 1)
  -promotionId ruleId
    	可选，优惠券id,多个用逗号隔开，山姆app优惠券列表接口中的'ruleId'字段
  -addressId string
    	可选，地址id
  -trackInfo string
    	可选，HTTP头部track-info
```

### BarkId

![bark.png](https://robgoods.github.io/sams/assets/bark.png)

开始运行后按命令行提示操作即可。

![run.png](https://robgoods.github.io/sams/assets/run.png)

## 声明
本项目仅供学习交流，严禁用作商业行为，特别禁止黄牛加价代抢等！

因违法违规等不当使用导致的后果与本人无关，如有任何问题可联系本人删除！
