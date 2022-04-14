package dd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/tidwall/gjson"
)

type StoreListParam struct {
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
}

type Store struct {
	StoreId                 string `json:"storeId"`
	StoreName               string `json:"-"`
	StoreType               string `json:"storeType"`
	AreaBlockId             string `json:"areaBlockId"`
	StoreDeliveryTemplateId string `json:"storeDeliveryTemplateId"`
	DeliveryModeId          string `json:"deliveryModeId"`
}

func (s *DingdongSession) GetStoreList(result gjson.Result) error {
	c := make([]Store, 0)

	for _, v := range result.Get("data.storeList").Array() {
		c = append(c, Store{
			StoreId:                 v.Get("storeId").Str,
			StoreName:               v.Get("storeName").Str,
			StoreType:               v.Get("storeType").String(),
			AreaBlockId:             v.Get("storeAreaBlockVerifyData.areaBlockId").Str,
			StoreDeliveryTemplateId: v.Get("storeRecmdDeliveryTemplateData.storeDeliveryTemplateId").Str,
			DeliveryModeId:          v.Get("storeDeliveryModeVerifyData.deliveryModeId").Str,
		})
	}
	s.StoreList = c
	return nil

}

func (s *DingdongSession) CheckStore() error {
	urlPath := "https://api-sams.walmartmobile.cn/api/v1/sams/merchant/storeApi/getRecommendStoreListByLocation"

	data := StoreListParam{
		Longitude: s.Address.Longitude,
		Latitude:  s.Address.Latitude,
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
			return s.GetStoreList(result)
		default:
			return errors.New(result.Get("msg").Str)
		}
	} else {
		return errors.New(fmt.Sprintf("[%v] %s", resp.StatusCode, body))
	}
}
