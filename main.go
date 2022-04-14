package main

import (
	"fmt"
	"github.com/robGoods/sams/dd"
	"time"
)

func main() {
	session := dd.DingdongSession{}
	err := session.InitSession("xxxxxxxxxxxx", "xxxxxxxxxx", 1, 2)

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
			fmt.Printf("校验商品失败：%s\n", err)
			switch err {
			case dd.CartGoodChangeErr:
				goto CartLoop
			case dd.LimitedErr:
				goto GoodsLoop
			case dd.NoMatchDeliverMode:
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

		session.SettleDeliveryInfo = dd.SettleDeliveryInfo{}
		for _, caps := range session.Capacity.CapCityResponseList {
			for _, v := range caps.List {
				fmt.Printf("配送时间： %s %s - %s, 是否可用：%v\n", v.CloseDate, v.StartTime, v.EndTime, !v.TimeISFull && !v.Disabled)
				if v.TimeISFull == false && v.Disabled == false && session.SettleDeliveryInfo.ArrivalTimeStr == "" {
					session.SettleDeliveryInfo.ArrivalTimeStr = fmt.Sprintf("%s %s - %s", caps.StrDate, v.StartTime, v.EndTime)
					session.SettleDeliveryInfo.ExpectArrivalTime = v.StartRealTime
					session.SettleDeliveryInfo.ExpectArrivalEndTime = v.EndRealTime
					break
				}
			}
		}

		if session.SettleDeliveryInfo.ArrivalTimeStr != "" {
			fmt.Printf("发现可用的配送时段::%s!\n", session.SettleDeliveryInfo.ArrivalTimeStr)
		} else {
			fmt.Println("当前无可用配送时间段")
			time.Sleep(1 * time.Second)
			goto CapacityLoop
		}
	OrderLoop:
		err = session.CommitPay()
		fmt.Printf("########## 提交订单中【%s】 ###########\n", time.Now().Format("15:04:05"))
		if err == nil {
			fmt.Println("抢购成功，请前往app付款！")
			if session.BarkId != "" {
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
		} else {
			fmt.Printf("下单失败：%s\n", err)
			switch err {
			case dd.LimitedErr1:
				fmt.Println("立即重试...")
				goto OrderLoop
			case dd.CloseOrderTimeExceptionErr:
				goto CapacityLoop
			case dd.DecreaseCapacityCountError:
				goto CapacityLoop
			default:
				goto CapacityLoop
			}
		}
	}
}
