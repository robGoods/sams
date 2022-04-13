package dd

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
)

type Address struct {
	AddressId       string `json:"addressId"`
	Mobile          string `json:"mobile"`
	Phone           string `json:"phone"`
	Name            string `json:"name"`
	CountryName     string `json:"countryName"`
	ProvinceName    string `json:"provinceName"`
	CityName        string `json:"cityName"`        //上海市
	DistrictName    string `json:"districtName"`    //浦东新区
	ReceiverAddress string `json:"receiverAddress"` //xx路
	DetailAddress   string `json:"detailAddress"`   //xx楼xx号
	Latitude        string `json:"latitude"`
	Longitude       string `json:"longitude"`
}

func parseAddress(addressMap gjson.Result) (Address, error) {
	address := Address{}
	address.AddressId = addressMap.Get("addressId").Str
	address.Mobile = addressMap.Get("mobile").Str
	address.Name = addressMap.Get("name").Str
	address.CountryName = addressMap.Get("countryName").Str
	address.ProvinceName = addressMap.Get("provinceName").Str
	address.CityName = addressMap.Get("cityName").Str
	address.DistrictName = addressMap.Get("districtName").Str
	address.ReceiverAddress = addressMap.Get("receiverAddress").Str
	address.DetailAddress = addressMap.Get("detailAddress").Str
	address.Latitude = addressMap.Get("latitude").Str
	address.Longitude = addressMap.Get("longitude").Str
	return address, nil
}

func (s *DingdongSession) GetAddress() (error, []Address) {
	urlPath := "https://api-sams.walmartmobile.cn/api/v1/sams/sams-user/receiver_address/address_list"
	req, _ := http.NewRequest("GET", urlPath, nil)
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
		return err, nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, nil
	}
	resp.Body.Close()
	if resp.StatusCode == 200 {

		result := gjson.Parse(string(body))
		switch result.Get("code").Str {
		case "Success":
			var addressList = make([]Address, 0)
			validAddress := result.Get("data.addressList").Array()
			for _, addressMap := range validAddress {
				address, err := parseAddress(addressMap)
				if err != nil {
					return err, nil
				}
				addressList = append(addressList, address)
			}
			return nil, addressList
		case "AUTH_FAIL":
			return errors.New(fmt.Sprintf("%s %s", result.Get("msg").Str, "token过期！！！")), nil
		default:
			return errors.New(result.Get("msg").Str), nil
		}
	} else {
		return errors.New(fmt.Sprintf("[%v] %s", resp.StatusCode, body)), nil
	}
}
