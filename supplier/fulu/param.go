package fulu

type Param map[string]string

func (p Param) Set(name, value string) {
	p[name] = value
}

func (p Param) Has(name string) bool {
	_, ok := p[name]
	return ok
}

func (p Param) Del(name string) {
	delete(p, name)
}

func (p Param) Get(name string, def ...string) string {
	if _, ok := p[name]; ok {
		return p[name]
	}
	if len(def) == 0 {
		return ""
	}
	return def[0]
}

type GoodsListGetParam struct {
	ProductID        int     `json:"product_id,omitempty"`         // 商品编号
	FirstCategoryID  int     `json:"first_category_id,omitempty"`  // 商品分类Id（一级）（预留参数，暂无商品分类接口提供）
	SecondCategoryID int     `json:"second_category_id,omitempty"` // 商品分类Id（二级）（预留参数，暂无商品分类接口提供）
	ThirdCategoryID  int     `json:"third_category_id,omitempty"`  // 商品分类Id（三级）（预留参数，暂无商品分类接口提供）
	ProductName      string  `json:"product_name,omitempty"`       // 商品名称
	ProductType      PT      `json:"product_type,omitempty"`       // 库存类型：卡密、直充
	FaceValue        float64 `json:"face_value,omitempty"`         // 面值
}

type GoodsStockCheckParam struct {
	ProductID string `json:"product_id"` // 商品编号
	BuyNum    int    `json:"buy_num"`    // 购买数量
}

type GoodsInfoGetParam struct {
	ProductID string `json:"product_id"` // 商品编号
}

// OrderDirectAddParam 直充下单参数
type OrderDirectAddParam struct {
	ProductID        int    `json:"product_id"`                   // 商品编号
	CustomerOrderNo  string `json:"customer_order_no"`            // 外部订单号
	ChargeAccount    string `json:"charge_account"`               // 充值账号
	BuyNum           int    `json:"buy_num"`                      // 购买数量
	ChargeGameName   string `json:"charge_game_name,omitempty"`   // 充值游戏名称
	ChargeGameRegion string `json:"charge_game_region,omitempty"` // 充值游戏区
	ChargeGameSrv    string `json:"charge_game_srv,omitempty"`    // 充值游戏服
	ChargeType       string `json:"charge_type,omitempty"`        // 充值类型
	ChargePassword   string `json:"charge_password,omitempty"`    // 充值密码，部分游戏类要传
	ChargeIP         string `json:"charge_ip,omitempty"`          // 下单真实Ip，区域商品要传
	ContactQQ        string `json:"contact_qq,omitempty"`         // 联系QQ
	ContactTel       string `json:"contact_tel,omitempty"`        // 联系电话
	RemainingNumber  int    `json:"remaining_number,omitempty"`   // 剩余数量
	ChargeGameRole   string `json:"charge_game_role,omitempty"`   // 充值游戏角色
	CustomerPrice    string `json:"customer_price,omitempty"`     // 外部销售价
	ShopType         string `json:"shop_type,omitempty"`          // 店铺类型（PDD、淘宝、天猫、京东、苏宁、其他；非必填字段，可忽略
	ExternalBizID    string `json:"external_biz_id,omitempty"`    // 透传字段
}

func (o *OrderDirectAddParam) Clone(NewProductID int, NewOrderNo string) *OrderDirectAddParam {
	argument := &(*o)
	argument.ProductID = NewProductID
	argument.CustomerOrderNo = NewOrderNo
	return argument
}

// OrderInfoGetParam 订单查询参数
type OrderInfoGetParam struct {
	CustomerOrderNo string `json:"customer_order_no"` // 外部订单号
}
