package dd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/tidwall/gjson"
)

type CapCityResponse struct {
	StrDate        string `json:"strDate"`
	DeliveryDesc   string `json:"deliveryDesc"`
	DeliveryDescEn string `json:"deliveryDescEn"`
	DateISFull     bool   `json:"dateISFull"`
}

type Capacity struct {
	Data                      string            `json:"data"`
	CapCityResponseList       []CapCityResponse `json:"capcityResponseList"`
	PortalPerformanceTemplate string            `json:"getPortalPerformanceTemplateResponse"`
}

func parseCapacity(g gjson.Result) (error, CapCityResponse) {
	var sizes []map[string]interface{}
	for _, size := range g.Get("sizes").Array() {
		sizes = append(sizes, size.Value().(map[string]interface{}))
	}
	capacity := CapCityResponse{
		StrDate:        g.Get("strDate").Str,
		DeliveryDesc:   g.Get("deliveryDesc").Str,
		DeliveryDescEn: g.Get("deliveryDescEn").Str,
		DateISFull:     g.Get("dateISFull").Bool(),
	}
	return nil, capacity
}

func (s *DingdongSession) GetCapacity(result gjson.Result) error {
	var capCityResponseList []CapCityResponse
	for _, v := range result.Get("data.capcityResponseList").Array() {
		_, product := parseCapacity(v)
		capCityResponseList = append(capCityResponseList, product)
	}
	s.Capacity = Capacity{
		Data:                      result.String(),
		CapCityResponseList:       capCityResponseList,
		PortalPerformanceTemplate: result.Get("data.getPortalPerformanceTemplateResponse").Str,
	}
	return nil
}

func (s *DingdongSession) CheckCapacity() error {
	urlPath := "https://api-sams.walmartmobile.cn/api/v1/sams/delivery/portal/getCapacityData"

	data := make(map[string]interface{})
	data["perDateList"] = []string{time.Now().Format("2006-01-02"), time.Now().AddDate(0, 0, 1).Format("2006-01-02")}
	data["storeDeliveryTemplateId"] = "552578721878546198"
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
			return s.GetCapacity(result)
		default:
			return errors.New(result.Get("msg").Str)
		}
	} else {
		return errors.New(fmt.Sprintf("[%v] %s", resp.StatusCode, body))
	}
}
