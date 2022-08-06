package sys

type ST string

const (
	Untreated  ST = "untreated"  // 未处理
	Processing ST = "processing" // 处理中
	Failed     ST = "failed"     // 失败
	Success    ST = "success"    // 成功
)

//////////////////////////////////////////////////////////

type OT int

const (
	HF OT = 1 // 话费
	LL OT = 2 // 流量
	KM OT = 3 // 卡密
	ZC OT = 4 // 直充
)

//////////////////////////////////////////////////////////

type PT string

const (
	CDKey          PT = "卡密"
	DirectRecharge PT = "直充"
)

//////////////////////////////////////////////////////////

type SAS string

const (
	OffShelves       SAS = "下架"
	OnShelves        SAS = "上架"
	UnderMaintenance SAS = "维护中"
	StockMaintenance SAS = "库存维护"
)

//////////////////////////////////////////////////////////

type STS string

const (
	OutStock STS = "断货"
	Alert    STS = "警报"
	Adequate STS = "充足"
)
