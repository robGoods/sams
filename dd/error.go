package dd

import "errors"

var CartGoodChangeErr = errors.New("购物车商品发生变化，请返回购物车页面重新结算")
var LimitedErr = errors.New("服务器正忙,请稍后再试")
var LimitedErr1 = errors.New("当前购物火爆，请稍后再试")
var OOSErr = errors.New("部分商品已缺货")
