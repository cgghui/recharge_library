package fulu

import "github.com/cgghui/recharge_library/sys"

type ParentBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  string `json:"result"`
	Sign    string `json:"sign"`
}

type GoodsListGetResponse struct {
	ProductID     int     `json:"product_id"`     // 商品Id
	ProductName   string  `json:"product_name"`   // 商品名称
	ProductType   sys.PT  `json:"product_type"`   // 库存类型：卡密、直充
	FaceValue     float64 `json:"face_value"`     // 面值
	PurchasePrice float64 `json:"purchase_price"` // 单价（单位：元）
	SalesStatus   sys.SAS `json:"sales_status"`   // 销售状态：下架、上架、维护中、库存维护（本接口只取上架状态的商品）
	StockStatus   sys.STS `json:"stock_status"`   // 库存状态：断货、警报、充足
	TemplateID    string  `json:"template_id"`    // 商品模板Id，可能为空
	Details       string  `json:"details"`        // 商品详情
}

type GoodsStockCheckResponse struct {
	ProductID   int     `json:"product_id"`   // 商品Id
	StockStatus sys.STS `json:"stock_status"` // 库存状态：断货、警报、充足
}

type GoodsInfoGetResponse struct {
	GoodsListGetResponse
	FourCategoryIcon string `json:"four_category_icon"`
	DetailType       int    `json:"detail_type"`
}

type OrderPublic struct {
	OrderID              string  `json:"order_id"`               // 订单编号
	CustomerOrderNo      string  `json:"customer_order_no"`      // 外部订单号，每次请求必须唯一
	ProductName          string  `json:"product_name"`           // 商品名称
	ChargeAccount        string  `json:"charge_account"`         // 充值账号
	BuyNum               int     `json:"buy_num"`                // 购买数量
	OrderType            sys.OT  `json:"order_type"`             // 订单类型：1-话费 2-流量 3-卡密 4-直充
	OrderPrice           float64 `json:"order_price"`            // 交易单价
	OrderState           sys.ST  `json:"order_state"`            // 订单状态： （success：成功，processing：处理中，failed：失败，untreated：未处理）
	CreateTime           string  `json:"create_time"`            // 创建时间 2019-07-27 16:44:30
	FinishTime           string  `json:"finish_time"`            // 订单完成时间，查单接口返回
	Area                 string  `json:"area"`                   // 充值区（中文）
	Server               string  `json:"server"`                 // 充值服（中文）
	Type                 string  `json:"type"`                   // 计费方式（中文）
	OperatorSerialNumber string  `json:"operator_serial_number"` // 运营商流水号
}

// OrderDirectAddResponse 直充下单响应结果
type OrderDirectAddResponse struct {
	ParentInfo *ParentBody `json:"parent_info"`
	OrderPublic
	ProductID int `json:"product_id"` // 商品Idx
}

// OrderInfoGetResponse 订单查询响应结果
type OrderInfoGetResponse struct {
	OrderPublic
	ProductID int    `json:"product_id"` // 商品Id
	Cards     []Card `json:"cards"`      // 卡密（券码）信息，卡密商品或部分直充商品返回（注意：卡密是密文，需要进行解密使用）；只有当订单的状态是success时，才会返回卡密（券码）信息；
}

// RechargeSuccessful 是否充值成功
// 成功返回true 失败返回false
func (o OrderInfoGetResponse) RechargeSuccessful() bool {
	return o.OrderState == sys.Success
}

// RechargeStatusCheck 充值状态检查
// s 可选为： Untreated Failed Processing Success
func (o OrderInfoGetResponse) RechargeStatusCheck(s sys.ST) bool {
	return o.OrderState == s
}

type Card struct {
	CardNumber   string `json:"card_number"`   // 卡号
	CardPWD      string `json:"card_pwd"`      // 密码
	CardType     string `json:"card_type"`     // 卡类型 0.普通卡密 1.二维码 2.短链
	CardDeadline string `json:"card_deadline"` // 卡密有效期
}

// OrderCardAddResponse 卡密
type OrderCardAddResponse struct {
	ParentInfo *ParentBody `json:"parent_info"`
	OrderPublic
	ProductID int `json:"product_id"` // 商品Idx
}
