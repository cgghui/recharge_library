package sys

import (
	"gorm.io/gorm"
)

type Order struct {
	ID                  uint64         `gorm:"primarykey" json:"id"`
	OrderNo             string         `gorm:"primarykey" json:"order_no"`
	OrderStatus         ST             `json:"order_status"`
	OrderType           OT             `json:"order_type"`
	OrderPrice          float64        `json:"order_price"`
	OrderFinishTime     Time           `json:"order_finish_time"`
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
	MerchantOrderNo     string         `json:"merchant_order_no"`
	SupplierName        string         `json:"supplier_name"`
	SupplierGoodsID     uint64         `json:"supplier_goods_id"`
	SupplierGoodsName   string         `json:"supplier_goods_name"`
	SupplierOrderNo     string         `json:"supplier_order_no"`
	SupplierOrderType   string         `json:"supplier_order_type"`
	SupplierOrderPrice  float64        `json:"supplier_order_price"`
	SupplierOrderStatus string         `json:"supplier_order_status"`
}

type OrderDetail struct {
	Order
	MerchantOrderDetail MerchantOrder `gorm:"foreignKey:order_no;references:merchant_order_no" json:"merchant_order_detail"`
	SupplierOrderDetail SupplierOrder `gorm:"foreignKey:order_no;references:supplier_order_no" json:"supplier_order_detail"`
}

type ThirdPartyOrder struct {
	OrderNo    string `json:"order_no" gorm:"primarykey"`
	Name       string `json:"name"`
	OriginData string `json:"origin_data"`
	FinishTime Time   `json:"finish_time"`
	CreatedAt  Time   `json:"created_at"`
}

type MerchantOrder struct {
	ThirdPartyOrder
}

type SupplierOrder struct {
	ThirdPartyOrder
}
