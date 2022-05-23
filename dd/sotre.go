package dd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
)

type StoreListParam struct {
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
}

type Store struct {
	StoreId                 string    `json:"storeId"`
	StoreName               string    `json:"-"`
	StoreType               string    `json:"storeType"`
	AreaBlockId             string    `json:"areaBlockId"`
	StoreDeliveryTemplateId string    `json:"storeDeliveryTemplateId"`
	DeliveryModeId          string    `json:"deliveryModeId"`
	DeliveryType            int       `json:"deliveryType"`
	Capacity                *Capacity `json:"-"`
}

func (s *DingdongSession) GetStoreList(result gjson.Result) []Store {
	c := make([]Store, 0)
	for _, v := range result.Get("data.storeList").Array() {
		c = append(c, Store{
			StoreId:                 v.Get("storeId").Str,
			StoreName:               v.Get("storeName").Str,
			StoreType:               v.Get("storeType").String(),
			AreaBlockId:             v.Get("storeAreaBlockVerifyData.areaBlockId").Str,
			StoreDeliveryTemplateId: v.Get("storeRecmdDeliveryTemplateData.storeDeliveryTemplateId").Str,
			DeliveryModeId:          v.Get("storeDeliveryModeVerifyData.deliveryModeId").Str,
			DeliveryType:            int(v.Get("storeDeliveryModeVerifyData.deliveryType").Int()),
			Capacity:                &Capacity{},
		})
	}

	return c
}

func (s *DingdongSession) CheckStore() ([]Store, error) {
	urlPath := "https://api-sams.walmartmobile.cn/api/v1/sams/merchant/storeApi/getRecommendStoreListByLocation"

	data := StoreListParam{
		Longitude: s.Address.Longitude,
		Latitude:  s.Address.Latitude,
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
			return s.GetStoreList(result), nil
		default:
			return nil, errors.New(result.Get("msg").Str)
		}
	} else {
		return nil, errors.New(fmt.Sprintf("[%v] %s", resp.StatusCode, body))
	}
}
