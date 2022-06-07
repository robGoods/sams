package dd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
)

type SettleDelivery struct {
	DeliveryType            int      `json:"deliveryType"` // 1,极速达 2, 全城配 3, 物流配送
	DeliveryName            string   `json:"deliveryName"`
	DeliveryDesc            string   `json:"deliveryDesc"`
	ExpectArrivalTime       string   `json:"expectArrivalTime"`
	ExpectArrivalEndTime    string   `json:"expectArrivalEndTime"`
	StoreDeliveryTemplateId string   `json:"storeDeliveryTemplateId"`
	DeliveryModeIdList      []string `json:"deliveryModeIdList"`
	AreaBlockId             string   `json:"areaBlockId"`
	AreaBlockName           string   `json:"areaBlockName"`
	FirstPeriod             int      `json:"firstPeriod"`
}

func parseSettleDelivery(g gjson.Result) (error, SettleDelivery) {
	r := SettleDelivery{
		DeliveryType:            int(g.Get("deliveryType").Num),
		DeliveryName:            g.Get("deliveryName").Str,
		DeliveryDesc:            g.Get("deliveryDesc").Str,
		ExpectArrivalTime:       g.Get("expectArrivalTime").Str,
		ExpectArrivalEndTime:    g.Get("expectArrivalEndTime").Str,
		StoreDeliveryTemplateId: g.Get("storeDeliveryTemplateId").Str,
		AreaBlockId:             g.Get("areaBlockId").Str,
		AreaBlockName:           g.Get("areaBlockName").Str,
		FirstPeriod:             int(g.Get("firstPeriod").Num),
	}

	for _, v := range g.Get("deliveryModeIdList").Array() {
		r.DeliveryModeIdList = append(r.DeliveryModeIdList, v.Str)
	}
	return nil, r
}

type SettleInfo struct {
	SaasId          string         `json:"saasId"`
	Uid             string         `json:"uid"`
	FloorId         int            `json:"floorId"`
	FloorName       string         `json:"floorName"`
	DeliveryFee     string         `json:"deliveryFee"`
	SettleDelivery  SettleDelivery `json:"settleDelivery"`
	DeliveryAddress Address        `json:"deliveryAddress"`
}

func parseSettleInfo(result gjson.Result) *SettleInfo {
	r := SettleInfo{}

	for _, v := range result.Get("data.settleDelivery").Array() {
		_, settleDelivery := parseSettleDelivery(v)
		r.SettleDelivery = settleDelivery
	}
	r.SaasId = result.Get("data.saasId").Str
	r.Uid = result.Get("data.uid").Str
	r.FloorId = int(result.Get("data.floorId").Num)
	r.FloorName = result.Get("data.floorName").Str
	r.DeliveryFee = result.Get("data.deliveryFee").Str
	address, err := parseAddress(result.Get("data.deliveryAddress"))
	if err == nil {
		r.DeliveryAddress = address
	}

	return &r
}

type DeliveryInfoVO struct {
	StoreDeliveryTemplateId string `json:"storeDeliveryTemplateId"`
	DeliveryModeId          string `json:"deliveryModeId"`
	StoreType               string `json:"storeType"`
}
type CouponInfo struct {
	PromotionId string `json:"promotionId"`
	StoreId     string `json:"storeId"`
}
type SettleParam struct {
	Uid            string         `json:"uid"`
	AddressId      string         `json:"addressId"`
	DeliveryInfoVO DeliveryInfoVO `json:"deliveryInfoVO"`
	DeliveryType   int            `json:"cartDeliveryType"`
	StoreInfo      Store          `json:"storeInfo"`
	CouponList     []CouponInfo   `json:"couponList,omitempty"`
	IsSelfPickup   int            `json:"isSelfPickup"`
	FloorId        int            `json:"floorId"`
	GoodsList      []Goods        `json:"goodsList"`
}

func (s *DingdongSession) CheckSettleInfo() (*SettleInfo, error) {
	urlPath := "https://api-sams.walmartmobile.cn/api/v1/sams/trade/settlement/getSettleInfo"

	data := SettleParam{
		Uid:       s.Uid,
		AddressId: s.Address.AddressId,
		DeliveryInfoVO: DeliveryInfoVO{
			StoreDeliveryTemplateId: s.StoreList[s.FloorInfo.StoreId].StoreDeliveryTemplateId,
			DeliveryModeId:          s.StoreList[s.FloorInfo.StoreId].DeliveryModeId,
			StoreType:               s.StoreList[s.FloorInfo.StoreId].StoreType,
		},
		DeliveryType: s.Conf.DeliveryType,
		StoreInfo:    s.StoreList[s.FloorInfo.StoreId],
		CouponList:   make([]CouponInfo, 0),
		IsSelfPickup: 0,
		FloorId:      s.Conf.FloorId,
		GoodsList:    s.GoodsList,
	}

	if len(s.Conf.PromotionId) > 0 {
		for _, id := range s.Conf.PromotionId {
			data.CouponList = append(data.CouponList, CouponInfo{PromotionId: id, StoreId: s.FloorInfo.StoreId})
		}
	}
	dataStr, _ := json.Marshal(data)
	req := s.NewRequest("POST", urlPath, dataStr)

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
			return parseSettleInfo(result), nil
		case "LIMITED":
			return nil, LimitedErr
		case "NO_MATCH_DELIVERY_MODE":
			return nil, NoMatchDeliverMode
		case "CART_GOOD_CHANGE":
			return nil, CartGoodChangeErr
		default:
			return nil, errors.New(result.Get("msg").Str)
		}
	} else {
		return nil, errors.New(fmt.Sprintf("[%v] %s", resp.StatusCode, body))
	}
}
