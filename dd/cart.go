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
	data["deviceType"] = "	"
	data["storeList"] =  s.StoreList
	data["deliveryType"] = s.DeliveryType
	data["homePagelongitude"] = s.Address.Longitude
	data["homePagelatitude"] = s.Address.Latitude

	dataStr, _ := json.Marshal(data)
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
			return s.GetCart(result)
		default:
			return errors.New(result.Get("msg").Str)
		}
	} else {
		return errors.New(fmt.Sprintf("[%v] %s", resp.StatusCode, body))
	}
}
