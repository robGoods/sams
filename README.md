# robSams

[![Build Status](https://app.travis-ci.com/robGoods/sams.svg?branch=master)](https://app.travis-ci.com/robGoods/sams)

[sam's blog](https://robgoods.github.io/sams/)

## 感谢

首先感谢Sam's在上海疫情期间，给我们的帮助，让我们在疫情期间依旧可以买到很多好的东西！

## 使用方式

```shell script
go run main --authToken=xxxxx
```

#### 参数说明

```shell script
$ go run main -h

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
  -trackInfo string
    	可选，HTTP头部track-info
```

### BarkId

![bark.png](https://robgoods.github.io/sams/assets/bark.png)

开始运行后按命令行提示操作即可。

![run.png](https://robgoods.github.io/sams/assets/run.png)

## 关于hack版本

hack版本与master版本基本相同，有兴趣的朋友可以研究下哪些不同。会有惊喜的发现😯😯😯

## 声明
本项目仅供学习交流，严禁用作商业行为，特别禁止黄牛加价代抢等！

因违法违规等不当使用导致的后果与本人无关，如有任何问题可联系本人删除！
