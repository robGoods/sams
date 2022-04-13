# robSams

[sam's blog](https://robgoods.github.io/sams/)

## 感谢

首先感谢Sam's在上海疫情期间，给我们的帮助，让我们在疫情期间依旧可以买到很多好的东西！

## 使用方式
在main.go的main函数中修改该行代码
```go
err := session.InitSession("xxxxxxxxxxxxxxxxx", "xxxxxxxxxxxxxxxxx", 1) // 1,普通商品 2,全球购保税 3,特殊订购自提 4,大件商品 5,厂家直供商品 6,特殊订购商品 7,失效商品
```
其中第一个参数为Sam's登录后的HTTP头部auth-token
第二个参数为通知用的bark id，下载`bark`后从app界面获取, 如果不需要可以填空字符串

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