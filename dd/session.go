package dd

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

type DingdongSession struct {
	AuthToken      string         `json:"auth-token"`
	BarkId         string         `json:"bark_id"`
	FloorId        int            `json:"floorId"` // 1,普通商品 2,全球购保税 3,特殊订购自提 4,大件商品 5,厂家直供商品 6,特殊订购商品 7,失效商品
	Address        Address        `json:"address"`
	Uid            string         `json:"uid"`
	Capacity       Capacity       `json:"capacity"`
	Channel        string         `json:"channel"` //0 => wechat  1 =>alipay
	SettleInfo     SettleInfo     `json:"settleInfo"`
	DeliveryInfoVO DeliveryInfoVO `json:"deliveryInfoVO"`
	GoodsList      []Goods        `json:"goods"`
	FloorInfo      FloorInfo      `json:"floorInfo"`
	StoreList      []Store        `json:"store"`
	OrderInfo      OrderInfo      `json:"orderInfo"`
	Client         *http.Client   `json:"client"`
	Cart           Cart           `json:"cart"`
}

func (s *DingdongSession) InitSession(AuthToken, barkId string, FloorId int) error {
	fmt.Println("########## 初始化 ##########")
	s.Client = &http.Client{Timeout: 60 * time.Second}
	s.AuthToken = AuthToken
	s.BarkId = barkId
	s.FloorId = FloorId //普通商品

	err, addrList := s.GetAddress()
	if err != nil {
		return err
	}
	if len(addrList) == 0 {
		return errors.New("未查询到有效收货地址，请前往app添加或检查cookie是否正确！")
	}
	fmt.Println("########## 选择收货地址 ##########")
	for i, addr := range addrList {
		fmt.Printf("[%v] %s %s %s %s %s \n", i, addr.Name, addr.DistrictName, addr.ReceiverAddress, addr.DetailAddress, addr.Mobile)
	}
	var index int
	for true {
		fmt.Println("请输入地址序号（0, 1, 2...)：")
		stdin := bufio.NewReader(os.Stdin)
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

	fmt.Println("########## 选择支付方式 ##########")
	for true {
		fmt.Println("请输入支付方式序号（0：微信 1：支付宝)：")
		stdin := bufio.NewReader(os.Stdin)
		_, err := fmt.Fscanln(stdin, &index)
		if err != nil {
			fmt.Printf("输入有误：%s!\n", err)
		} else if index == 0 {
			s.Channel = "wechat"
			break
		} else if index == 1 {
			s.Channel = "alipay"
			break
		} else {
			fmt.Println("输入有误：序号无效！")
		}
	}
	return nil
}
