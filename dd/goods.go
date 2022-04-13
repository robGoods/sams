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

type Goods struct {
	IsSelected bool   `json:"isSelected"`
	Quantity   int    `json:"quantity"`
	SpuId      string `json:"spuId"`
	StoreId    string `json:"storeId"`
}

type NormalGoods struct {
	StoreId       string `json:"storeId"`
	StoreType     int    `json:"storeType"`
	SpuId         string `json:"spuId"`
	SkuId         string `json:"skuId"`
	BrandId       string `json:"brandId"`
	GoodsName     string `json:"goodsName"`
	Price         int    `json:"price"`
	InvalidReason string `json:"invalidReason"`
	Quantity      int    `json:"quantity"`
}

func (this NormalGoods) ToGoods() Goods {
	return Goods{
		IsSelected: true,
		Quantity:   this.Quantity,
		SpuId:      this.SpuId,
		StoreId:    this.StoreId,
	}
}

func parseNormalGoods(g gjson.Result) (error, NormalGoods) {
	r := NormalGoods{
		StoreId:       g.Get("storeId").Str,
		StoreType:     int(g.Get("storeType").Num),
		SpuId:         g.Get("spuId").Str,
		SkuId:         g.Get("skuId").Str,
		BrandId:       g.Get("brandId").Str,
		GoodsName:     g.Get("goodsName").Str,
		Price:         int(g.Get("price").Int()),
		InvalidReason: g.Get("invalidReason").Str,
		Quantity:      int(g.Get("quantity").Num),
	}
	return nil, r
}

func (s *DingdongSession) CheckGoods() error {
	urlPath := "https://api-sams.walmartmobile.cn/api/v1/sams/trade/settlement/checkGoodsInfo"

	data := make(map[string]interface{})
	data["floorId"] = 1
	data["storeId"] = ""
	goods := make([]Goods, 0)
	for _, v := range s.Cart.FloorInfoList {
		if v.FloorId == s.FloorId {
			for _, v := range v.NormalGoodsList {
				if data["storeId"] == "" {
					data["storeId"] = v.StoreId
				}
				goods = append(goods, Goods{StoreId: v.StoreId, Quantity: v.Quantity, SpuId: v.SpuId})
			}
		}
	}
	data["goodsList"] = goods
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
			if result.Get("data.isHasException").Bool() == false {
				return nil
			} else {
				fmt.Printf("========以下商品已过期====")
				for index, v := range result.Get("data.popUpInfo.goodsList").Array() {
					_, goods := parseNormalGoods(v)
					fmt.Printf("[%v] %s 数量：%v 总价：%d\n", index, goods.SpuId, goods.StoreId, goods.Price)
				}
				return OOSErr
			}
		default:
			return errors.New(result.Get("msg").Str)
		}
	} else {
		return errors.New(fmt.Sprintf("[%v] %s", resp.StatusCode, body))
	}
}
