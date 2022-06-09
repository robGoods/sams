package dd

import "errors"

var CartGoodChangeErr = errors.New("购物车商品发生变化，请返回购物车页面重新结算")
var CapacityErr = errors.New("获取履约时间异常")
var GetDeliveryInfoErr = errors.New("获取履约配送信息异常")
var LimitedErr = errors.New("服务器正忙,请稍后再试")
var LimitedErr1 = errors.New("当前购物火爆，请稍后再试")
var NoMatchDeliverMode = errors.New("当前区域不支持配送，请重新选择地址")
var GoodsExceedLimitErr = errors.New("商品超过限购数量")
var CloseOrderTimeExceptionErr = errors.New("尊敬的会员，您选择的配送时间已失效，请重新选择")
var NotDeliverCapCityErr = errors.New("当前配送时间段已约满，请重新选择配送时段")
var DecreaseCapacityCountError = errors.New("扣减运力失败")
var StoreHasClosedError = errors.New("门店已打烊")
var PreGoodNotStartSellErr = errors.New("商品还未开始正式售卖，无法购买")
var CloudGoodsOverWightErr = errors.New("出于交通安全考虑，极速达订单限重30公斤，您的订单已超重，请分开下单")

var OOSErr = errors.New("部分商品已缺货")
