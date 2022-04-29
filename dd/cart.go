package dd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/tidwall/gjson"
)

type Cart struct {
	DeliveryAddress Address     `json:"deliveryAddress"`
	FloorInfoList   []FloorInfo `json:"floorInfoList"`
	ParentOrderSign string      `json:"parent_order_sign"`
}

type FloorInfo struct {
	FloorId         int           `json:"floorId"`
	NormalGoodsList []NormalGoods `json:"normalGoodsList"`
	Amount          string        `json:"amount"`
	Quantity        int           `json:"quantity"`
	StoreInfo       StoreInfo     `json:"storeInfo"`
}

func parseFloorInfos(g gjson.Result) (error, FloorInfo) {
	r := FloorInfo{
		FloorId:  int(g.Get("floorId").Num),
		Amount:   g.Get("amount").Str,
		Quantity: int(g.Get("quantity").Num),
		StoreInfo: StoreInfo{
			StoreId:                 g.Get("storeInfo.storeId").Str,
			StoreType:               fmt.Sprintf("%d", int(g.Get("storeInfo.storeType").Num)),
			AreaBlockId:             g.Get("storeInfo.areaBlockId").Str,
			StoreDeliveryTemplateId: g.Get("storeInfo.storeDeliveryTemplateId").Str,
			DeliveryModeId:          g.Get("storeInfo.deliveryModeId").Str,
		},
	}
	for _, normalGoods := range g.Get("normalGoodsList").Array() {
		_, p := parseNormalGoods(normalGoods)
		r.NormalGoodsList = append(r.NormalGoodsList, p)
	}
	for _, promotionGoodsList := range g.Get("promotionFloorGoodsList").Array() {
		for _, promotionGoods := range promotionGoodsList.Get("promotionGoodsList").Array() {
			_, p := parseNormalGoods(promotionGoods)
			r.NormalGoodsList = append(r.NormalGoodsList, p)
		}
	}

	for _, shortageStockGoods := range g.Get("shortageStockGoodsList").Array() {
		_, p := parseNormalGoods(shortageStockGoods)
		r.NormalGoodsList = append(r.NormalGoodsList, p)
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
	address, err := parseAddress(result.Get("data.deliveryAddress"))
	if err == nil {
		c.DeliveryAddress = address
	}
	s.Cart = c
	return nil
}

func (s *DingdongSession) CheckCart() error {
	urlPath := "https://api-sams.walmartmobile.cn/api/v1/sams/trade/cart/getUserCart"

	data := make(map[string]interface{})
	data["uid"] = ""
	data["deviceType"] = "ios"
	data["storeList"] = s.StoreList
	data["deliveryType"] = s.Conf.DeliveryType
	data["homePagelongitude"] = s.Address.Longitude
	data["homePagelatitude"] = s.Address.Latitude

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
