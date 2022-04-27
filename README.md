# robSams

[![Build status](https://ci.appveyor.com/api/projects/status/v5lt859vmjm3v9i5?svg=true)](https://ci.appveyor.com/project/iscod/sams)

[sam's blog](https://robgoods.github.io/sams/)

## 感谢

1. 感谢Sam's在上海疫情期间，给我们的帮助，让我们在疫情期间依旧可以买到很多好的东西！，购买Sam's会员请前往[山姆会员商店](https://www.samsclub.cn/)
1. 感谢各位朋友对该项目的支持和star。
1. 感谢 [gyf19](https://github.com/gyf19), [Matata-lol](https://github.com/Matata-lol), [3096](https://github.com/3096), [Nicolerobinn](https://github.com/Nicolerobinn), [likang7](https://github.com/likang7), [zyr3536](https://github.com/zyr3536)  对本项目的贡献

## 使用方式

```sh
go run main.go --authToken=xxxxx
```

> 如果没有go环境，可以在 [releases](https://github.com/robGoods/sams/releases) 下载编译好的文件，直接运行即可

### 更新说明

1. 增加库存不足商品自动设置订单数量为库存数
1. 配送时间多个可用
1. 支付方式，收货地址均支持`flag`模式选择，而非`Stdin`模式，默认微信支付，地址未指定时依旧会提示选择
1. 优惠券支持多张同时使用，使用前最好确认下订单是否满足使用类型

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
