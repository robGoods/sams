package dd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"time"
)

type CapCityResponse struct {
	StrDate        string `json:"strDate"`
	DeliveryDesc   string `json:"deliveryDesc"`
	DeliveryDescEn string `json:"deliveryDescEn"`
	DateISFull     bool   `json:"dateISFull"`
	List           []List `json:"list"`
}

type List struct {
	StartTime     string `json:"startTime"`
	EndTime       string `json:"endTime"`
	TimeISFull    bool   `json:"timeISFull"`
	Disabled      bool   `json:"disabled"`
	CloseDate     string `json:"closeDate"`
	CloseTime     string `json:"closeTime"`
	StartRealTime string `json:"startRealTime"` //1649984400000
	EndRealTime   string `json:"endRealTime"`   //1650016800000
}

type Capacity struct {
	Data                      string            `json:"data"`
	CapCityResponseList       []CapCityResponse `json:"capcityResponseList"`
	PortalPerformanceTemplate string            `json:"getPortalPerformanceTemplateResponse"`
}

func parseCapacityList(g gjson.Result) CapCityResponse {
	var list []List
	for _, v := range g.Get("list").Array() {
		list = append(list, List{
			StartTime:     v.Get("startTime").Str,
			EndTime:       v.Get("endTime").Str,
			TimeISFull:    v.Get("timeISFull").Bool(),
			Disabled:      v.Get("disabled").Bool(),
			CloseDate:     v.Get("closeDate").Str,
			CloseTime:     v.Get("closeTime").Str,
			StartRealTime: v.Get("startRealTime").Str,
			EndRealTime:   v.Get("endRealTime").Str,
		})
	}
	return CapCityResponse{
		StrDate:        g.Get("strDate").Str,
		DeliveryDesc:   g.Get("deliveryDesc").Str,
		DeliveryDescEn: g.Get("deliveryDescEn").Str,
		DateISFull:     g.Get("dateISFull").Bool(),
		List:           list,
	}
}

func parseCapacity(result gjson.Result) *Capacity {
	var capCityResponseList []CapCityResponse
	for _, v := range result.Get("data.capcityResponseList").Array() {
		capCityResponseList = append(capCityResponseList, parseCapacityList(v))
	}
	return &Capacity{
		Data:                      result.String(),
		CapCityResponseList:       capCityResponseList,
		PortalPerformanceTemplate: result.Get("data.getPortalPerformanceTemplateResponse").Str,
	}
}

func (s *DingdongSession) GetCapacity(storeDeliveryTemplateId string) (*Capacity, error) {
	urlPath := "https://api-sams.walmartmobile.cn/api/v1/sams/delivery/portal/getCapacityData"
	data := make(map[string]interface{})
	data["perDateList"] = []string{time.Now().Format("2006-01-02"), time.Now().AddDate(0, 0, 1).Format("2006-01-02")}
	data["storeDeliveryTemplateId"] = storeDeliveryTemplateId
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
			return parseCapacity(result), nil
		case "LIMITED":
			return nil, LimitedErr
		default:
			if result.Get("msg").Str == CapacityErr.Error() {
				return nil, CapacityErr
			}
			return nil, errors.New(result.Get("msg").Str)
		}
	} else {
		return nil, errors.New(fmt.Sprintf("[%v] %s", resp.StatusCode, body))
	}
}
