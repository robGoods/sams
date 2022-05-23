package dd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/tidwall/gjson"
)

type Goods struct {
	GoodsName  string  `json:"-"`
	Price      int     `json:"-"`
	IsSelected bool    `json:"isSelected"`
	Quantity   int     `json:"quantity"`
	SpuId      string  `json:"spuId"`
	StoreId    string  `json:"storeId"`
	Weight     float64 `json:"-"`
}

type NormalGoods struct {
	StoreId            string  `json:"storeId"`
	StoreType          int     `json:"storeType"`
	SpuId              string  `json:"spuId"`
	SkuId              string  `json:"skuId"`
	BrandId            string  `json:"brandId"`
	GoodsName          string  `json:"goodsName"`
	Price              int     `json:"price"`
	InvalidReason      string  `json:"invalidReason"`
	Quantity           int     `json:"quantity"`
	StockQuantity      int     `json:"stockQuantity"`
	StockStatus        bool    `json:"stockStatus"`
	IsPutOnSale        bool    `json:"isPutOnSale"`
	IsAvailable        bool    `json:"isAvailable"`
	LimitNum           int     `json:"limitNum"`
	ResiduePurchaseNum int     `json:"residuePurchaseNum"`
	IsSelected         bool    `json:"IsSelected"`
	Weight             float64 `json:"-"`
}

func (this NormalGoods) ToGoods() Goods {
	return Goods{
		IsSelected: this.IsSelected,
		GoodsName:  this.GoodsName,
		Price:      this.Price,
		Quantity:   this.Quantity,
		SpuId:      this.SpuId,
		StoreId:    this.StoreId,
		Weight:     this.Weight,
	}
}

func parseNormalGoods(g gjson.Result) NormalGoods {
	return NormalGoods{
		StoreId:            g.Get("storeId").Str,
		StoreType:          int(g.Get("storeType").Num),
		SpuId:              g.Get("spuId").Str,
		SkuId:              g.Get("skuId").Str,
		BrandId:            g.Get("brandId").Str,
		GoodsName:          g.Get("goodsName").Str,
		Price:              int(g.Get("price").Int()),
		InvalidReason:      g.Get("invalidReason").Str,
		Quantity:           int(g.Get("quantity").Num),
		StockQuantity:      int(g.Get("stockQuantity").Num),
		StockStatus:        g.Get("stockStatus").Bool(),
		IsPutOnSale:        g.Get("isPutOnSale").Bool(),
		IsAvailable:        g.Get("isAvailable").Bool(),
		LimitNum:           int(g.Get("purchaseLimitVO.limitNum").Int()),
		ResiduePurchaseNum: int(g.Get("purchaseLimitVO.residuePurchaseNum").Int()),
		IsSelected:         g.Get("isSelected").Bool(),
		Weight:             g.Get("weight").Float(),
	}
}

func (s *DingdongSession) CheckGoods() (map[string]NormalGoods, error) {
	urlPath := "https://api-sams.walmartmobile.cn/api/v1/sams/trade/settlement/checkGoodsInfo"

	data := make(map[string]interface{})
	data["floorId"] = 1
	data["storeId"] = ""
	for _, v := range s.Cart.FloorInfoList {
		if v.FloorId == s.Conf.FloorId {
			data["storeId"] = v.StoreId
		}
	}
	data["goodsList"] = s.GoodsList
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
			if result.Get("data.isHasException").Bool() == false {
				return nil, nil
			} else {
				fmt.Println(result.Get("data.popUpInfo.desc").Str)
				var goods = make(map[string]NormalGoods, 0)
				for _, v := range result.Get("data.popUpInfo.goodsList").Array() {
					g := parseNormalGoods(v)
					goods[g.SpuId] = g
				}
				return goods, OOSErr
			}
		default:
			return nil, errors.New(result.Get("msg").Str)
		}
	} else {
		return nil, errors.New(fmt.Sprintf("[%v] %s", resp.StatusCode, body))
	}
}
