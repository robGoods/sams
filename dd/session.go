package dd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Config struct {
	AuthToken    string
	BarkId       string
	FloorId      int //1,普通商品 2,全球购保税 3,特殊订购自提 4,大件商品 5,厂家直供商品 6,特殊订购商品 7,失效商品
	DeliveryType int //1 急速达，2， 全程配送
	Longitude    string
	Latitude     string
	Deviceid     string
	Trackinfo    string
	PromotionId  []string
	AddressId    string
	PayMethod    int //支付方式
	DeliveryFee  bool
	StoreConf    string
	IsSelected   bool
}

type DingdongSession struct {
	Conf               Config
	Address            Address                    `json:"address"`
	Uid                string                     `json:"uid"`
	Capacity           Capacity                   `json:"capacity"`
	SettleDeliveryInfo map[int]SettleDeliveryInfo `json:"settleDeliveryInfo"`
	GoodsList          []Goods                    `json:"goods"`
	FloorInfo          FloorInfo                  `json:"floorInfo"`
	StoreList          map[string]Store           `json:"store"`
	Client             *http.Client               `json:"client"`
	Cart               Cart                       `json:"cart"`
}

func (s *DingdongSession) InitSession(conf Config) error {
	fmt.Println("########## 初始化 ##########")
	s.Client = &http.Client{Timeout: 60 * time.Second}
	s.Conf = conf

	if len(s.Conf.PromotionId) > 0 {
		fmt.Println("########## 当前选择优惠券 ##########")
		for k, id := range s.Conf.PromotionId {
			fmt.Printf("[%d] %s\n", k, id)
		}
	} else {
		fmt.Println("########## 当前没有选择优惠券 ##########")
	}
	stdin := bufio.NewReader(os.Stdin)

	err, addrList := s.GetAddress()
	if err != nil {
		return err
	}
	if len(addrList) == 0 {
		return errors.New("未查询到有效收货地址，请前往app添加或检查cookie是否正确！")
	}
	if s.Conf.AddressId != "" {
		for _, v := range addrList {
			if v.AddressId == s.Conf.AddressId {
				s.Address = v
				fmt.Printf("收货地址 :  %s %s %s %s %s \n", s.Address.Name, s.Address.DistrictName, s.Address.ReceiverAddress, s.Address.DetailAddress, s.Address.Mobile)
			}
		}
	}
	if s.Address.AddressId == "" {
		fmt.Println("########## 选择收货地址 ##########")
		for i, addr := range addrList {
			fmt.Printf("[%v] Id: %s %s %s %s %s %s \n", i, addr.AddressId, addr.Name, addr.DistrictName, addr.ReceiverAddress, addr.DetailAddress, addr.Mobile)
		}

		var index int
		for true {
			fmt.Println("请输入地址序号（0, 1, 2...)：")
			_, err := fmt.Fscanln(stdin, &index)
			if err != nil {
				fmt.Printf("输入有误：%s!\n", err)
			} else if index >= len(addrList) {
				fmt.Println("输入有误：超过最大序号！")
			} else {
				break
			}
		}
		s.Address = addrList[index]
	}

	fmt.Println("########## 选择支付方式 ##########")
	switch s.Conf.PayMethod {
	case 1:
		fmt.Println("支付方式 : wechat ")
	case 2:
		fmt.Println("支付方式 : alipay ")
	default:
		return errors.New("选择支付方式有误！")
	}

	return nil
}

func (s *DingdongSession) NewRequest(method, url string, dataStr []byte) *http.Request {

	var body io.Reader = nil
	if dataStr != nil {
		body = bytes.NewReader(dataStr)
	}
	req, _ := http.NewRequest(method, url, body)

	req.Header.Set("Host", "api-sams.walmartmobile.cn")
	req.Header.Set("content-type", "application/json;charset=UTF-8")
	//req.Header.Set("accept", "*/*")
	req.Header.Set("auth-token", s.Conf.AuthToken)
	req.Header.Set("longitude", s.Conf.Longitude)
	req.Header.Set("latitude", s.Conf.Latitude)
	req.Header.Set("device-id", s.Conf.Deviceid)
	req.Header.Set("app-version", "5.0.47.0")
	req.Header.Set("device-type", "ios")
	req.Header.Set("Accept-Language", "zh-Hans-CN;q=1")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("apptype", "ios")
	req.Header.Set("device-name", "iPhone14,5")
	req.Header.Set("device-os-version", "15.4.1")
	req.Header.Set("User-Agent", "SamClub/5.0.47 (iPhone; iOS 15.4.1; Scale/3.00)")
	req.Header.Set("system-language", "CN")

	return req
}
