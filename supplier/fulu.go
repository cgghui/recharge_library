package supplier

import (
	"github.com/cgghui/recharge_library/supplier/fulu"
	"strconv"
	"time"
)

type Fulu struct {
	app fulu.App
	key []byte
}

func NewFulu(key, secret string) *Fulu {
	return &Fulu{
		app: fulu.NewApp(fulu.OpenApiURL, key, secret),
		key: []byte(secret),
	}
}

func (*Fulu) Name() string {
	return fulu.Name
}

func (f *Fulu) OrderPlaceCard(arg interface{}) (*ResponseOrderPlaceCard, error) {
	ret, err := f.app.OrderCardAdd(arg.(*fulu.OrderCardAddParam))
	if err != nil {
		return nil, err
	}
	ct, _ := time.Parse(fulu.TimeFormat1, ret.CreateTime)
	return &ResponseOrderPlaceCard{
		ResponseBase: ResponseBase{
			Success: ret.ParentInfo.Code == 0,
			Message: ret.ParentInfo.Message,
		},
		ResponseOrderBase: ResponseOrderBase{
			ProductID:            strconv.Itoa(ret.ProductID),
			ProductName:          ret.ProductName,
			OrderNo:              ret.OrderID,
			SystemOrderNo:        ret.CustomerOrderNo,
			BuyNum:               ret.BuyNum,
			OrderType:            ret.OrderType,
			OrderPrice:           ret.OrderPrice,
			OrderState:           ret.OrderState,
			CreateTime:           ct,
			OperatorSerialNumber: ret.OperatorSerialNumber,
		},
	}, nil
}

func (f *Fulu) OrderPlaceDirect(arg interface{}) (*ResponseOrderPlaceDirect, error) {
	ret, err := f.app.OrderDirectAdd(arg.(*fulu.OrderDirectAddParam))
	if err != nil {
		return nil, err
	}
	ct, _ := time.Parse(fulu.TimeFormat1, ret.CreateTime)
	return &ResponseOrderPlaceDirect{
		ResponseBase: ResponseBase{
			Success: ret.ParentInfo.Code == 0,
			Message: ret.ParentInfo.Message,
		},
		ResponseOrderBase: ResponseOrderBase{
			ProductID:            strconv.Itoa(ret.ProductID),
			ProductName:          ret.ProductName,
			OrderNo:              ret.OrderID,
			SystemOrderNo:        ret.CustomerOrderNo,
			BuyNum:               ret.BuyNum,
			OrderType:            ret.OrderType,
			OrderPrice:           ret.OrderPrice,
			OrderState:           ret.OrderState,
			CreateTime:           ct,
			OperatorSerialNumber: ret.OperatorSerialNumber,
		},
		ChargeAccount: ret.ChargeAccount,
		Area:          ret.Area,
		Server:        ret.Server,
	}, nil
}

func (f *Fulu) OrderQuery(arg interface{}) (*ResponseOrderQuery, error) {
	ret, err := f.app.OrderInfoGet(arg.(*fulu.OrderInfoGetParam))
	if err != nil {
		return nil, err
	}
	cards := make([]Card, 0)
	if len(ret.Cards) != 0 {
		for _, card := range ret.Cards {
			n := fulu.DecryptECB([]byte(card.CardNumber), f.key)
			p := fulu.DecryptECB([]byte(card.CardPWD), f.key)
			cards = append(cards, Card{
				Number:   string(n),
				Password: string(p),
				Type:     card.CardType,
				Deadline: card.CardDeadline,
			})
		}
	}
	ct, _ := time.Parse(fulu.TimeFormat1, ret.CreateTime)
	return &ResponseOrderQuery{
		ResponseBase: ResponseBase{
			Success: ret.OrderID != "",
			Message: "",
		},
		ResponseOrderBase: ResponseOrderBase{
			ProductID:            strconv.Itoa(ret.ProductID),
			ProductName:          ret.ProductName,
			OrderNo:              ret.OrderID,
			SystemOrderNo:        ret.CustomerOrderNo,
			BuyNum:               ret.BuyNum,
			OrderType:            ret.OrderType,
			OrderPrice:           ret.OrderPrice,
			OrderState:           ret.OrderState,
			CreateTime:           ct,
			OperatorSerialNumber: ret.OperatorSerialNumber,
		},
		Cards: cards,
	}, nil
}

func (f *Fulu) GoodsQuery(arg interface{}) (*ResponseGoodsQuery, error) {
	ret, err := f.app.GoodsInfoGet(arg.(*fulu.GoodsInfoGetParam))
	if err != nil {
		return nil, err
	}
	return &ResponseGoodsQuery{
		ResponseBase: ResponseBase{
			Success: ret.ProductID != 0,
			Message: "",
		},
		GoodsID:       strconv.Itoa(ret.ProductID),
		GoodsName:     ret.ProductName,
		GoodsType:     ret.ProductType,
		FaceValue:     ret.FaceValue,
		PurchasePrice: ret.PurchasePrice,
		SalesStatus:   ret.SalesStatus,
		StockStatus:   ret.StockStatus,
	}, nil
}

func (f *Fulu) Object() interface{} {
	return f.app
}
