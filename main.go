package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/robGoods/sams/dd"
)

var (
	showHelp     = flag.Bool("help", false, "show help")
	authToken    = flag.String("authToken", "", "必选, Sam's App HTTP头部auth-token")
	barkId       = flag.String("barkId", "", "可选，通知用的`bark` id, 可选参数")
	floorId      = flag.Int("floorId", 1, "可选，1,普通商品 2,全球购保税 3,特殊订购自提 4,大件商品 5,厂家直供商品 6,特殊订购商品 7,失效商品")
	deliveryType = flag.Int("deliveryType", 2, "可选，1 急速达，2， 全程配送")
	longitude    = flag.String("longitude", "", "可选，HTTP头部longitude")
	latitude     = flag.String("latitude", "", "可选，HTTP头部latitude")
	deviceId     = flag.String("deviceId", "", "可选，HTTP头部device-id")
	trackInfo    = flag.String("trackInfo", "", "可选，HTTP头部track-info")
	promotionId  = flag.String("promotionId", "", "可选，优惠券id,多个用逗号隔开，山姆app优惠券列表接口中的'ruleId'字段")
	addressId    = flag.String("addressId", "", "可选，地址id")
	payMethod    = flag.Int("payMethod", 1, "可选，1,微信 2,支付宝")
)

func main() {
	flag.Parse()
	if *showHelp {
		flag.Usage()
		return
	}

	if *authToken == "" {
		flag.Usage()
		return
	}

	splitFn := func(c rune) bool {
		return c == ','
	}

	session := dd.DingdongSession{
		SettleDeliveryInfo: map[int]dd.SettleDeliveryInfo{},
		StoreList:          map[string]dd.Store{},
	}
	conf := dd.Config{
		AuthToken:    *authToken,                                //HTTP头部auth-token
		BarkId:       *barkId,                                   //通知用的bark id，下载bark后从app界面获取, 如果不需要可以填空字符串
		FloorId:      *floorId,                                  //1,普通商品 2,全球购保税 3,特殊订购自提 4,大件商品 5,厂家直供商品 6,特殊订购商品 7,失效商品
		DeliveryType: *deliveryType,                             //1 急速达，2， 全程配送
		Longitude:    *longitude,                                //HTTP头部longitude,可选参数
		Latitude:     *latitude,                                 //HTTP头部latitude,可选参数
		Deviceid:     *deviceId,                                 //HTTP头部device-id,可选参数
		Trackinfo:    *trackInfo,                                //HTTP头部track-info,可选参数
		PromotionId:  strings.FieldsFunc(*promotionId, splitFn), //优惠券id
		AddressId:    *addressId,                                //地址
		PayMethod:    *payMethod,                                //支付方式
	}

	err := session.InitSession(conf)

	if err != nil {
		fmt.Println(err)
		return
	}

	for true {
	SaveDeliveryAddress:
		fmt.Println("########## 切换购物车收货地址 ###########")
		err = session.SaveDeliveryAddress()
		if err != nil {
			goto SaveDeliveryAddress
		} else {
			fmt.Println("切换成功!")
			fmt.Printf("%s %s %s %s %s \n", session.Address.Name, session.Address.DistrictName, session.Address.ReceiverAddress, session.Address.DetailAddress, session.Address.Mobile)
		}
	StoreLoop:
		fmt.Println("########## 获取地址附近可用商店 ###########")
		stores, err := session.CheckStore()
		if err != nil {
			fmt.Printf("%s", err)
			goto StoreLoop
		}

		for index, store := range stores {
			if _, ok := session.StoreList[store.StoreId]; !ok {
				session.StoreList[store.StoreId] = store
				fmt.Printf("[%v] Id：%s 名称：%s, 类型 ：%s\n", index, store.StoreId, store.StoreName, store.StoreType)
			}
		}
	CartLoop:
		fmt.Printf("########## 获取购物车中有效商品【%s】 ###########\n", time.Now().Format("15:04:05"))
		err = session.CheckCart()
		for _, v := range session.Cart.FloorInfoList {
			if v.FloorId == session.Conf.FloorId {
				session.GoodsList = make([]dd.Goods, 0)
				for _, goods := range v.NormalGoodsList {
					if goods.StockQuantity > 0 && goods.StockStatus && goods.IsPutOnSale && goods.IsAvailable {
						if goods.StockQuantity <= goods.Quantity {
							goods.Quantity = goods.StockQuantity
						}
						if goods.LimitNum > 0 && goods.Quantity > goods.LimitNum {
							goods.Quantity = goods.LimitNum
						}
						session.GoodsList = append(session.GoodsList, goods.ToGoods())
					}
				}

				for _, goods := range v.ShortageStockGoodsList {
					if goods.StockQuantity > 0 && goods.StockStatus && goods.IsPutOnSale && goods.IsAvailable {
						if goods.StockQuantity <= goods.Quantity {
							goods.Quantity = goods.StockQuantity
						}
						if goods.LimitNum > 0 && goods.Quantity > goods.LimitNum {
							goods.Quantity = goods.LimitNum
						}
						session.GoodsList = append(session.GoodsList, goods.ToGoods())
					}
				}

				//for _, goods := range v.AllOutOfStockGoodsList {
				//	if goods.StockQuantity > 0 && goods.StockStatus && goods.IsPutOnSale && goods.IsAvailable {
				//		if goods.StockQuantity <= goods.Quantity {
				//			goods.Quantity = goods.StockQuantity
				//		}
				//		session.GoodsList = append(session.GoodsList, goods.ToGoods())
				//	}
				//}

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

		for index, goods := range session.GoodsList {
			fmt.Printf("[%v] %s 数量：%v 总价：%d\n", index, goods.GoodsName, goods.Quantity, goods.Price)
		}

		if len(session.GoodsList) == 0 {
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("当前购物车中无有效商品")
			}
			if errors.Is(err, dd.LimitedErr1) {
				time.Sleep(1 * time.Second)
			}
			goto StoreLoop
		}
	GoodsLoop:
		fmt.Printf("########## 开始校验当前商品【%s】 ###########\n", time.Now().Format("15:04:05"))
		if err = session.CheckGoods(); err != nil {
			fmt.Println(err)
			time.Sleep(1 * time.Second)
			switch err {
			case dd.OOSErr:
				goto CartLoop
			default:
				goto GoodsLoop
			}
		}
		if err = session.CheckSettleInfo(); err != nil {
			fmt.Printf("校验商品失败：%s\n", err)
			time.Sleep(1 * time.Second)
			switch err {
			case dd.CartGoodChangeErr:
				goto CartLoop
			case dd.LimitedErr:
				goto GoodsLoop
			case dd.NoMatchDeliverMode:
				goto SaveDeliveryAddress
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
			//刷新可用配送时间， 会出现“服务器正忙,请稍后再试”， 可以忽略。
			goto CapacityLoop
		}

		session.SettleDeliveryInfo = map[int]dd.SettleDeliveryInfo{}
		for _, caps := range session.Capacity.CapCityResponseList {
			for _, v := range caps.List {
				fmt.Printf("配送时间： %s %s - %s, 是否可用：%v\n", caps.StrDate, v.StartTime, v.EndTime, !v.TimeISFull && !v.Disabled)
				if v.TimeISFull == false && v.Disabled == false {
					session.SettleDeliveryInfo[len(session.SettleDeliveryInfo)] = dd.SettleDeliveryInfo{
						ArrivalTimeStr:       fmt.Sprintf("%s %s - %s", caps.StrDate, v.StartTime, v.EndTime),
						ExpectArrivalTime:    v.StartRealTime,
						ExpectArrivalEndTime: v.EndRealTime,
					}
				}
			}
		}

		if len(session.SettleDeliveryInfo) > 0 {
			for _, v := range session.SettleDeliveryInfo {
				fmt.Printf("发现可用的配送时段::%s!\n", v.ArrivalTimeStr)
			}
		} else {
			fmt.Println("当前无可用配送时间段")
			time.Sleep(1 * time.Second)
			goto CapacityLoop
		}
	OrderLoop:
		for len(session.SettleDeliveryInfo) > 0 {
			for k, v := range session.SettleDeliveryInfo {
				fmt.Printf("########## 提交订单中【%s】 ###########\n", time.Now().Format("15:04:05"))
				fmt.Printf("配送时段: %s!\n", v.ArrivalTimeStr)
				err = session.CommitPay(v)
				if err == nil {
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
				} else {
					fmt.Printf("下单失败：%s\n", err)
					switch err {
					case dd.LimitedErr1:
						fmt.Println("立即重试...")
						goto OrderLoop
					case dd.OOSErr, dd.PreGoodNotStartSellErr:
						goto CartLoop
					case dd.StoreHasClosedError:
						goto StoreLoop
					case dd.CloseOrderTimeExceptionErr, dd.DecreaseCapacityCountError, dd.NotDeliverCapCityErr:
						delete(session.SettleDeliveryInfo, k)
					default:
						goto CapacityLoop
					}
				}
			}
		}
		goto CapacityLoop
	}
}
