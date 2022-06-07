package dd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/tidwall/gjson"
)

type Cart struct {
	FloorInfoList   []FloorInfo `json:"floorInfoList"`
	ParentOrderSign string      `json:"parent_order_sign"`
}

type FloorInfo struct {
	FloorId                int           `json:"floorId"`
	DeliveryType           int           `json:"deliveryType"`
	NormalGoodsList        []NormalGoods `json:"normalGoodsList"`
	ShortageStockGoodsList []NormalGoods `json:"ShortageStockGoodsList"`
	AllOutOfStockGoodsList []NormalGoods `json:"allOutOfStockGoodsList"`
	Amount                 string        `json:"amount"`
	Quantity               int           `json:"quantity"`
	StoreId                string        `json:"storeId"`
}

func parseFloorInfos(g gjson.Result) (error, FloorInfo) {
	r := FloorInfo{
		FloorId:                int(g.Get("floorId").Num),
		DeliveryType:           int(g.Get("deliveryType").Num),
		Amount:                 g.Get("amount").Str,
		Quantity:               int(g.Get("quantity").Num),
		StoreId:                g.Get("storeId").Str,
		NormalGoodsList:        make([]NormalGoods, 0),
		ShortageStockGoodsList: make([]NormalGoods, 0),
		AllOutOfStockGoodsList: make([]NormalGoods, 0),
	}
	for _, normalGoods := range g.Get("normalGoodsList").Array() {
		r.NormalGoodsList = append(r.NormalGoodsList, parseNormalGoods(normalGoods))
	}
	for _, promotionGoodsList := range g.Get("promotionFloorGoodsList").Array() {
		for _, promotionGoods := range promotionGoodsList.Get("promotionGoodsList").Array() {
			r.NormalGoodsList = append(r.NormalGoodsList, parseNormalGoods(promotionGoods))
		}
	}

	for _, shortageStockGoods := range g.Get("shortageStockGoodsList").Array() {
		r.ShortageStockGoodsList = append(r.ShortageStockGoodsList, parseNormalGoods(shortageStockGoods))
	}

	//查询无货商品是否上架
	for _, outOfStockGoods := range g.Get("allOutOfStockGoodsList").Array() {
		r.AllOutOfStockGoodsList = append(r.AllOutOfStockGoodsList, parseNormalGoods(outOfStockGoods))
	}

	return nil, r
}
func (s *DingdongSession) GetCart(result gjson.Result) error {
	c := Cart{
		FloorInfoList: make([]FloorInfo, 0),
	}
	for _, v := range result.Get("data.floorInfoList").Array() {
		_, floor := parseFloorInfos(v)
		c.FloorInfoList = append(c.FloorInfoList, floor)
	}

	s.Cart = c
	return nil
}

type GetCartPram struct {
	Uid               string  `json:"uid"`
	DeviceType        string  `json:"deviceType"`
	StoreList         []Store `json:"storeList"`
	DeliveryType      int     `json:"deliveryType"`
	HomePagelongitude string  `json:"homePagelongitude"`
	HomePagelatitude  string  `json:"homePagelatitude"`
}

func (s *DingdongSession) CheckCart() error {
	urlPath := "https://api-sams.walmartmobile.cn/api/v1/sams/trade/cart/getUserCart"

	data := GetCartPram{
		Uid:               "",
		DeviceType:        "ios",
		StoreList:         make([]Store, 0),
		DeliveryType:      s.Conf.DeliveryType,
		HomePagelongitude: s.Address.Longitude,
		HomePagelatitude:  s.Address.Latitude,
	}

	for _, store := range s.StoreList {
		data.StoreList = append(data.StoreList, store)
	}

	dataStr, _ := json.Marshal(data)
	req := s.NewRequest("POST", urlPath, dataStr)

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
			return s.GetCart(result)
		case "LIMITED":
			return LimitedErr1
		default:
			return errors.New(result.Get("msg").Str)
		}
	} else {
		return errors.New(fmt.Sprintf("[%v] %s", resp.StatusCode, body))
	}
}
