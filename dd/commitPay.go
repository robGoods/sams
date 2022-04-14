package dd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
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
	Uid                string              `json:"uid"`          //273583094,
	AppId              string              `json:"appId"`        //wx57364320cb03dfba
	AddressId          string              `json:"addressId"`
	DeliveryInfoVO     DeliveryInfoVO      `json:"deliveryInfoVO"`
	Remark             string              `json:"remark"`
	StoreInfo          StoreInfo           `json:"storeInfo"`
	ShortageDesc       string              `json:"shortageDesc"`
	PayMethodId        string              `json:"payMethodId"`
}

type OrderInfo struct {
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

func (s *DingdongSession) GetOrderInfo(result gjson.Result) error {
	s.OrderInfo = OrderInfo{
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
	return nil
}

func (s *DingdongSession) CommitPay() error {
	urlPath := "https://api-sams.walmartmobile.cn/api/v1/sams/trade/settlement/commitPay"

	data := CommitPayPram{
		GoodsList:          s.GoodsList,
		InvoiceInfo:        make(map[int]interface{}),
		DeliveryType:       s.DeliveryType, // 1,急速到达 2,全城配送
		FloorId:            0,              //急速时选1
		Amount:             "13123",        //测试没用但必须有
		PurchaserName:      "",
		SettleDeliveryInfo: s.SettleDeliveryInfo,
		TradeType:          "APP",
		PurchaserId:        "",
		PayType:            0,
		Currency:           "CNY",
		Channel:            s.Channel,
		ShortageId:         1,
		IsSelfPickup:       0,
		OrderType:          0,
		Uid:                "213123", //s.Uid,
		AppId:              fmt.Sprintf("wx51394321bc03adfadf"),
		AddressId:          s.Address.AddressId,
		DeliveryInfoVO:     s.DeliveryInfoVO,
		Remark:             "",
		StoreInfo:          s.FloorInfo.StoreInfo,
		ShortageDesc:       "其他商品继续配送（缺货商品直接退款）",
		PayMethodId:        "1486659732",
	}

	dataStr, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, _ := http.NewRequest("POST", urlPath, bytes.NewReader(dataStr))
	req.Header.Set("Host", "api-sams.walmartmobile.cn")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "*/*")
	//req.Header.Set("auth-token", "xxxxxxxxxxxx")
	req.Header.Set("auth-token", s.AuthToken)
	//req.Header.Set("app-version", "5.0.46.1")
	req.Header.Set("device-type", "ios")
	req.Header.Set("Accept-Language", "zh-Hans-CN;q=1, en-CN;q=0.9, ga-IE;q=0.8")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	//req.Header.Set("apptype", "ios")
	//req.Header.Set("device-name", "iPhone12,8")
	//req.Header.Set("device-os-version", "13.4.1")
	req.Header.Set("User-Agent", "SamClub/5.0.46 (iPhone; iOS 13.4.1; Scale/2.00)")
	req.Header.Set("system-language", "CN")

	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	resp.Body.Close()
	if resp.StatusCode == 200 {
		result := gjson.Parse(string(body))
		switch result.Get("code").Str {
		case "Success":
			if result.Get("data.isSuccess").Bool() {
				return s.GetOrderInfo(result)
			}
			return errors.New(result.Get("data.failReason").Str)
		case "LIMITED":
			return LimitedErr1
		case "CLOSE_ORDER_TIME_EXCEPTION":
			return CloseOrderTimeExceptionErr
		case "DECREASE_CAPACITY_COUNT_ERROR":
			return DecreaseCapacityCountError
		default:
			return errors.New(result.Get("msg").Str)
		}
	} else {
		return errors.New(fmt.Sprintf("[%v] %s", resp.StatusCode, body))
	}
}
