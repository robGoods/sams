package dd

import "errors"

var CartGoodChangeErr = errors.New("购物车商品发生变化，请返回购物车页面重新结算")
var LimitedErr = errors.New("服务器正忙,请稍后再试")
var LimitedErr1 = errors.New("当前购物火爆，请稍后再试")
var NoMatchDeliverMode = errors.New("当前区域不支持配送，请重新选择地址")
var CloseOrderTimeExceptionErr = errors.New("尊敬的会员，您选择的配送时间已失效，请重新选择")
var DecreaseCapacityCountError = errors.New("扣减运力失败")
var StoreHasClosedError = errors.New("门店已打烊")

var OOSErr = errors.New("部分商品已缺货")
