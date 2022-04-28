package dd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/tidwall/gjson"
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
	StockQuantity int    `json:"stockQuantity"`
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
		StockQuantity: int(g.Get("stockQuantity").Num),
	}
	return nil, r
}

func (s *DingdongSession) CheckGoods() error {
	urlPath := "https://api-sams.walmartmobile.cn/api/v1/sams/trade/settlement/checkGoodsInfo"

	data := make(map[string]interface{})
	data["floorId"] = 1
	data["storeId"] = ""
	for _, v := range s.Cart.FloorInfoList {
		if v.FloorId == s.Conf.FloorId {
			data["storeId"] = v.StoreInfo.StoreId
		}
	}
	data["goodsList"] = s.GoodsList
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
			if result.Get("data.isHasException").Bool() == false {
				return nil
			} else {
				fmt.Println(result.Get("data.popUpInfo.desc").Str)
				for index, v := range result.Get("data.popUpInfo.goodsList").Array() {
					_, goods := parseNormalGoods(v)
					fmt.Printf("[%v] %s 数量：%v 总价：%d\n", index, goods.GoodsName, goods.Quantity, goods.Price)
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
