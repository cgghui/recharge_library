package sys

import (
	"gorm.io/gorm"
)

// Order 系统订单
type Order struct {
	ID                  uint64         `gorm:"primarykey" json:"id"`
	OrderNo             string         `gorm:"primarykey" json:"order_no"`
	OrderStatus         ST             `json:"order_status"`
	OrderType           OT             `json:"order_type"`
	OrderPrice          float64        `json:"order_price"`
	OrderFinishTime     Time           `json:"order_finish_time"`
	NotifyRet           string         `json:"-"`
	LoopFinish          string         `json:"loop_finish"`
	ChargeAppName       string         `json:"charge_app_name"`
	ChargeCurrencyName  string         `json:"charge_currency_name"`
	ChargeAccount       string         `json:"charge_account"`
	ChargeParValue      float64        `json:"charge_par_value"`
	ChargeNum           uint           `json:"charge_num"`
	ChargeArea          string         `json:"charge_area"`
	ChargeServer        string         `json:"charge_server"`
	CreatedAt           Time           `json:"created_at"`
	UpdatedAt           Time           `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `json:"deleted_at"`
	MerchantName        string         `json:"merchant_name"`
	MerchantGoodsID     uint64         `json:"merchant_goods_id"`
	MerchantOrderNo     string         `json:"merchant_order_no"`
	SupplierName        string         `json:"supplier_name"`
	SupplierGoodsID     uint64         `json:"supplier_goods_id"`
	SupplierGoodsName   string         `json:"supplier_goods_name"`
	SupplierOrderNo     string         `json:"supplier_order_no"`
	SupplierOrderType   string         `json:"supplier_order_type"`
	SupplierOrderPrice  float64        `json:"supplier_order_price"`
	SupplierOrderStatus string         `json:"supplier_order_status"`
}

// OrderDetail 订单详情
type OrderDetail struct {
	Order
	MerchantOrderDetail MerchantOrder `gorm:"foreignKey:order_no;references:merchant_order_no" json:"merchant_order_detail"`
	SupplierOrderDetail SupplierOrder `gorm:"foreignKey:order_no;references:supplier_order_no" json:"supplier_order_detail"`
}

// OrderDetailForSupplier 查询供应商订单详情
type OrderDetailForSupplier struct {
	Order
	SupplierOrderDetail SupplierOrder `gorm:"foreignKey:order_no;references:supplier_order_no" json:"supplier_order_detail"`
}

func (OrderDetailForSupplier) TableName() string {
	return "order"
}

// OrderDetailForMerchant 查询商户订单详情
type OrderDetailForMerchant struct {
	Order
	MerchantOrderDetail MerchantOrder `gorm:"foreignKey:order_no;references:merchant_order_no" json:"merchant_order_detail"`
}

func (OrderDetailForMerchant) TableName() string {
	return "order"
}

type ThirdPartyOrder struct {
	OrderNo    string `json:"order_no" gorm:"primarykey"`
	Name       string `json:"name"`
	Remark     string `json:"remark"`
	OriginData string `json:"origin_data"`
	FinishTime Time   `json:"finish_time"`
	CreatedAt  Time   `json:"created_at"`
}

type MerchantOrder struct {
	ThirdPartyOrder
}

type MerchantOrderDetail struct {
	ThirdPartyOrder
	OrderDetail []OrderDetailForSupplier `gorm:"foreignKey:merchant_order_no;references:order_no" json:"order_detail"`
}

type SupplierOrder struct {
	ThirdPartyOrder
}

type SupplierOrderDetail struct {
	ThirdPartyOrder
	OrderDetail OrderDetailForMerchant `gorm:"foreignKey:supplier_order_no;references:order_no" json:"order_detail"`
}

// MerchantGoodsThread 单线程多线程的商品
type MerchantGoodsThread struct {
	ID           uint `gorm:"primarykey"`
	GoodsID      string
	SingleThread YN
}

type MerchantGoods struct {
	ID           uint `gorm:"primarykey"`
	MerchantName string
	GoodsID      string
	SkuCode      string
	RelatedCode  string
	Enable       YN
	CreatedAt    Time
	RelatedList  []MerchantGoodsRelated `gorm:"foreignKey:code;references:sku_code"`
}

type MerchantGoodsRelated struct {
	ID              uint   `gorm:"primarykey"`
	Code            string `gorm:"primarykey"`
	SupplierName    string
	SupplierGoodsID string
	Sort            uint
	Buy             uint
	ParValue        float64
	Unit            string
	Enable          YN
	CreatedAt       Time
}
