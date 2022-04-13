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
			return s.GetStoreList(result)
		default:
			return errors.New(result.Get("msg").Str)
		}
	} else {
		return errors.New(fmt.Sprintf("[%v] %s", resp.StatusCode, body))
	}
}
