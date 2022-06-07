package dd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
)

type CommitPayPram struct {
	GoodsList          []Goods             `json:"goodsList"`
	InvoiceInfo        map[int]interface{} `json:"invoiceInfo"`
	DeliveryType       int                 `json:"cartDeliveryType"`
	FloorId            int                 `json:"floorId"`
	Amount             string              `json:"amount"`
	PurchaserName      string              `json:"purchaserName"`
	SettleDeliveryInfo SettleDeliveryInfo  `json:"settleDeliveryInfo"`
	TradeType          string              `json:"tradeType"` //"APP"
	PurchaserId        string              `json:"purchaserId"`
	PayType            int                 `json:"payType"`
	Currency           string              `json:"currency"`     // CNY
	Channel            string              `json:"channel"`      // wechat
	ShortageId         int                 `json:"shortageId"`   //1
	IsSelfPickup       int                 `json:"isSelfPickup"` //0
	OrderType          int                 `json:"orderType"`    //0
	CouponList         []CouponInfo        `json:"couponList,omitempty"`
	Uid                string              `json:"uid"`   //273583094,
	AppId              string              `json:"appId"` //wx57364320cb03dfba
	AddressId          string              `json:"addressId"`
	DeliveryInfoVO     DeliveryInfoVO      `json:"deliveryInfoVO"`
	Remark             string              `json:"remark"`
	StoreInfo          Store               `json:"storeInfo"`
	ShortageDesc       string              `json:"shortageDesc"`
	PayMethodId        string              `json:"payMethodId"`
}

type Order struct {
	IsSuccess bool    `json:"isSuccess"`
	OrderNo   string  `json:"orderNo"`
	PayAmount string  `json:"payAmount"`
	Channel   string  `json:"channel"`
	PayInfo   PayInfo `json:"PayInfo"`
}

type PayInfo struct {
	PayInfo    string `json:"PayInfo"`
	OutTradeNo string `json:"OutTradeNo"`
	TotalAmt   int    `json:"TotalAmt"`
}

type SettleDeliveryInfo struct {
	DeliveryType         int    `json:"deliveryType"`         //默认0
	ExpectArrivalTime    string `json:"expectArrivalTime"`    //配送时间: 1649922300000
	ExpectArrivalEndTime string `json:"expectArrivalEndTime"` //配送时间
	ArrivalTimeStr       string `json:"-"`
}

func (s *DingdongSession) GetOrderInfo(result gjson.Result) *Order {
	return &Order{
		IsSuccess: result.Get("data.isSuccess").Bool(),
		OrderNo:   result.Get("data.orderNo").Str,
		PayAmount: result.Get("data.payAmount").Str,
		Channel:   result.Get("data.channel").Str,
		PayInfo: PayInfo{
			PayInfo:    result.Get("data.PayInfo.PayInfo").Str,
			OutTradeNo: result.Get("data.PayInfo.OutTradeNo").Str,
			TotalAmt:   int(result.Get("data.PayInfo.TotalAmt").Num),
		},
	}
}

func (s *DingdongSession) CommitPay(info SettleDeliveryInfo) (*Order, error) {
	urlPath := "https://api-sams.walmartmobile.cn/api/v1/sams/trade/settlement/commitPay"

	data := CommitPayPram{
		GoodsList:          s.GoodsList,
		InvoiceInfo:        make(map[int]interface{}),
		DeliveryType:       s.Conf.DeliveryType, // 1,急速到达 2,全城配送
		FloorId:            s.Conf.FloorId,
		Amount:             s.FloorInfo.Amount,
		PurchaserName:      "",
		SettleDeliveryInfo: info,
		TradeType:          "APP",
		PurchaserId:        "",
		PayType:            0,
		Currency:           "CNY",
		Channel:            "wechat",
		ShortageId:         1,
		IsSelfPickup:       0,
		OrderType:          0,
		CouponList:         make([]CouponInfo, 0),
		Uid:                s.Uid,
		AppId:              fmt.Sprintf("wx51394321bc03adfadf"),
		AddressId:          s.Address.AddressId,
		DeliveryInfoVO: DeliveryInfoVO{
			StoreDeliveryTemplateId: s.StoreList[s.FloorInfo.StoreId].StoreDeliveryTemplateId,
			DeliveryModeId:          s.StoreList[s.FloorInfo.StoreId].DeliveryModeId,
			StoreType:               s.StoreList[s.FloorInfo.StoreId].StoreType,
		},
		Remark:       "",
		StoreInfo:    s.StoreList[s.FloorInfo.StoreId],
		ShortageDesc: "其他商品继续配送（缺货商品直接退款）",
		PayMethodId:  "1486659732",
	}

	if s.Conf.PayMethod == 2 {
		data.Channel = "alipay"
	}

	if len(s.Conf.PromotionId) > 0 {
		for _, id := range s.Conf.PromotionId {
			data.CouponList = append(data.CouponList, CouponInfo{PromotionId: id, StoreId: s.FloorInfo.StoreId})
		}
	}

	dataStr, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req := s.NewRequest("POST", urlPath, dataStr)
	req.Header.Set("track-info", s.Conf.Trackinfo)
	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	if resp.StatusCode == 200 {
		result := gjson.Parse(string(body))
		switch result.Get("code").Str {
		case "Success":
			if result.Get("data.isSuccess").Bool() {
				return s.GetOrderInfo(result), nil
			} else {
				return nil, errors.New(result.Get("data.failReason").Str)
			}
		case "LIMITED":
			return nil, LimitedErr1
		case "GOODS_EXCEED_LIMIT":
			return nil, GoodsExceedLimitErr
		case "CLOSE_ORDER_TIME_EXCEPTION":
			return nil, CloseOrderTimeExceptionErr
		case "DECREASE_CAPACITY_COUNT_ERROR":
			return nil, DecreaseCapacityCountError
		case "OUT_OF_STOCK":
			return nil, OOSErr
		case "NOT_DELIVERY_CAPACITY_ERROR":
			return nil, NotDeliverCapCityErr
		case "STORE_HAS_CLOSED":
			return nil, StoreHasClosedError
		case "PRE_GOOD_NOT_START_SELL":
			return nil, PreGoodNotStartSellErr
		case "CLOUD_GOODS_OVER_WEIGHT":
			return nil, CloudGoodsOverWightErr
		case "CART_GOOD_CHANGE":
			return nil, CartGoodChangeErr
		case "GET_DELIVERY_INFO_ERROR":
			return nil, GetDeliveryInfoErr
		default:
			return nil, errors.New(result.Get("msg").Str)
		}
	} else {
		return nil, errors.New(fmt.Sprintf("[%v] %s", resp.StatusCode, body))
	}
}
