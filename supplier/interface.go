package supplier

import (
	"github.com/cgghui/recharge_library/sys"
	"time"
)

type MethodInterface interface {
	Name() string
	OrderPlaceCard(arg interface{}) (*ResponseOrderPlaceCard, error)
	OrderPlaceDirect(arg interface{}) (*ResponseOrderPlaceDirect, error)
	OrderQuery(arg interface{}) (*ResponseOrderQuery, error)
	GoodsQuery(arg interface{}) (*ResponseGoodsQuery, error)
	Object() interface{}
}

type ResponseBase struct {
	Success bool   `json:"success"` // 成功时，这个字段返回true
	Message string `json:"message"`
}

type ResponseOrderBase struct {
	ProductID            string    `json:"product_id"`
	ProductName          string    `json:"product_name"`           // 商品名称
	OrderNo              string    `json:"order_no"`               // 订单编号
	SystemOrderNo        string    `json:"system_order_no"`        // 外部订单号，每次请求必须唯一
	BuyNum               int       `json:"buy_num"`                // 购买数量
	OrderType            sys.OT    `json:"order_type"`             // 订单类型：1-话费 2-流量 3-卡密 4-直充
	OrderPrice           float64   `json:"order_price"`            // 交易单价
	OrderState           sys.ST    `json:"order_state"`            // 订单状态： （success：成功，processing：处理中，failed：失败，untreated：未处理）
	CreateTime           time.Time `json:"create_time"`            // 创建时间 2019-07-27 16:44:30
	OperatorSerialNumber string    `json:"operator_serial_number"` // 运营商流水号
}

type ResponseOrderPlaceCard struct {
	ResponseBase
	ResponseOrderBase
}

type ResponseOrderPlaceDirect struct {
	ResponseBase
	ResponseOrderBase
	ChargeAccount string `json:"charge_account"` // 充值账号
	Area          string `json:"area"`           // 充值区（中文）
	Server        string `json:"server"`         // 充值服（中文）
}

type ResponseOrderQuery struct {
	ResponseBase
	ResponseOrderBase
	Cards []Card `json:"cards"` // 只有当订单的状态是success时，才会返回卡密（券码）信息；
}

type Card struct {
	Number   string    `json:"number"`   // 卡号
	Password string    `json:"password"` // 密码
	Type     sys.CardT `json:"type"`     // 卡类型 0.普通卡密 1.二维码 2.短链
	Deadline string    `json:"deadline"` // 卡密有效期
}

type ResponseGoodsQuery struct {
	ResponseBase
	GoodsID       string  `json:"product_id"`     // 商品Id
	GoodsName     string  `json:"product_name"`   // 商品名称
	GoodsType     sys.PT  `json:"product_type"`   // 库存类型：卡密、直充
	FaceValue     float64 `json:"face_value"`     // 面值
	PurchasePrice float64 `json:"purchase_price"` // 单价（单位：元）
	SalesStatus   sys.SAS `json:"sales_status"`   // 销售状态：下架、上架、维护中、库存维护（本接口只取上架状态的商品）
	StockStatus   sys.STS `json:"stock_status"`   // 库存状态：断货、警报、充足
}

func (r *ResponseGoodsQuery) CheckGoodsType(t sys.PT) bool {
	return t == r.GoodsType
}
