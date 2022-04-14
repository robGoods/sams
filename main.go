package main

import (
	"fmt"
	"time"

	"github.com/robGoods/sams/dd"
)

func main() {
	conf := dd.Config{
		AuthToken: "xxxxxx", //HTTP头部auth-token
		BarkId:    "xxxxx",  //通知用的bark id，下载bark后从app界面获取, 如果不需要可以填空字符串
		Longitude: "xxxxxx", //HTTP头部longitude
		Latitude:  "xxxxx",  //HTTP头部latitude
		Deviceid:  "xxxxx",  //HTTP头部device-id
		Trackinfo: `xxxxxx`, // HTTP头部track-info
	}
	session := dd.DingdongSession{}
	err := session.InitSession(conf, 1) // 1,普通商品 2,全球购保税 3,特殊订购自提 4,大件商品 5,厂家直供商品 6,特殊订购商品 7,失效商品

	if err != nil {
		fmt.Println(err)
		return
	}

	for true {
	StoreLoop:
		fmt.Println("########## 获取地址附近可用商店 ###########")
		err = session.CheckStore()
		if err != nil {
			fmt.Printf("%s", err)
			goto StoreLoop
		}

		for index, store := range session.StoreList {
			fmt.Printf("[%v] Id：%s 名称：%s, 类型 ：%s\n", index, store.StoreId, store.StoreName, store.StoreType)
		}
	CartLoop:
		fmt.Printf("########## 获取购物车中有效商品【%s】 ###########\n", time.Now().Format("15:04:05"))
		session.CheckCart()
		for _, v := range session.Cart.FloorInfoList {
			if v.FloorId == session.FloorId {
				for index, goods := range v.NormalGoodsList {
					session.GoodsList = append(session.GoodsList, goods.ToGoods())
					fmt.Printf("[%v] %s 数量：%v 总价：%d\n", index, goods.GoodsName, goods.Quantity, goods.Price)
				}
				session.FloorInfo = v
				session.DeliveryInfoVO = dd.DeliveryInfoVO{
					StoreDeliveryTemplateId: v.StoreInfo.StoreDeliveryTemplateId,
					DeliveryModeId:          v.StoreInfo.DeliveryModeId,
					StoreType:               v.StoreInfo.StoreType,
				}
			} else {
				//无效商品
				//for index, goods := range v.NormalGoodsList {
				//	fmt.Printf("----[%v] %s 数量：%v 总价：%d\n", index, goods.SpuId, goods.StoreId, goods.Price)
				//}
			}
		}
		if len(session.GoodsList) == 0 {
			fmt.Println("当前购物车中无有效商品")
			goto CartLoop
		}
	GoodsLoop:
		fmt.Printf("########## 开始校验当前商品【%s】 ###########\n", time.Now().Format("15:04:05"))
		if err = session.CheckGoods(); err != nil {
			fmt.Println(err)
			switch err {
			case dd.OOSErr:
				goto CartLoop
			default:
				goto GoodsLoop
			}
		}
		if err = session.CheckSettleInfo(); err != nil {
			fmt.Println(err)
			switch err {
			case dd.CartGoodChangeErr:
				goto CartLoop
			case dd.LimitedErr:
				goto GoodsLoop
			default:
				goto GoodsLoop
			}
		}
	CapacityLoop:
		fmt.Printf("########## 获取当前可用配送时间【%s】 ###########\n", time.Now().Format("15:04:05"))
		err = session.CheckCapacity()
		if err != nil {
			fmt.Println(err)
			time.Sleep(1 * time.Second)
			continue
		}

		dateISFull := true
		for _, capCityResponse := range session.Capacity.CapCityResponseList {
			if capCityResponse.DateISFull == false && dateISFull {
				dateISFull = false
				fmt.Printf("发现可用的配送时段:%s!\n", capCityResponse.StrDate)
			}
		}

		if dateISFull {
			fmt.Println("当前无可用配送时间段")
			time.Sleep(1 * time.Second)
			goto CapacityLoop
		}
	OrderLoop:
		err = session.CommitPay()
		fmt.Printf("########## 提交订单中【%s】 ###########\n", time.Now().Format("15:04:05"))
		switch err {
		case nil:
			fmt.Println("抢购成功，请前往app付款！")
			if session.Conf.BarkId != "" {
				for true {
					err = session.PushSuccess(fmt.Sprintf("Smas抢单成功，订单号：%s", session.OrderInfo.OrderNo))
					if err == nil {
						break
					} else {
						fmt.Println(err)
					}
					time.Sleep(1 * time.Second)
				}
			}
			return
		case dd.LimitedErr1:
			fmt.Printf("[%s] 立即重试...\n", err)
			goto OrderLoop
		default:
			goto CartLoop
		}
	}
}
