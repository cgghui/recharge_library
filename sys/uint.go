package sys

import "strconv"

type ST string

const (
	Untreated  ST = "untreated"  // 未处理
	Processing ST = "processing" // 处理中
	Failed     ST = "failed"     // 失败
	Success    ST = "success"    // 成功
)

func (s ST) String() string {
	return string(s)
}

//////////////////////////////////////////////////////////

type OT int

const (
	HF OT = 1 // 话费
	LL OT = 2 // 流量
	KM OT = 3 // 卡密
	ZC OT = 4 // 直充
)

func (o OT) String() string {
	return strconv.Itoa(int(o))
}

//////////////////////////////////////////////////////////

type PT string

const (
	CDKey          PT = "卡密"
	DirectRecharge PT = "直充"
)

func (p PT) String() string {
	return string(p)
}

//////////////////////////////////////////////////////////

type SAS string

const (
	OffShelves       SAS = "下架"
	OnShelves        SAS = "上架"
	UnderMaintenance SAS = "维护中"
	StockMaintenance SAS = "库存维护"
)

func (s SAS) String() string {
	return string(s)
}

//////////////////////////////////////////////////////////

type STS string

const (
	OutStock STS = "断货"
	Alert    STS = "警报"
	Adequate STS = "充足"
)

func (s STS) String() string {
	return string(s)
}

//////////////////////////////////////////////////////////

type YN string

const (
	Y YN = "Y"
	N YN = "N"
)

func (s YN) String() string {
	return string(s)
}

///////////////////////////////
// 卡类型 0.普通卡密 1.二维码 2.短链

type CardT int

const (
	Normal CardT = 0
	Qrcode CardT = 1
	Link   CardT = 2
)

func (s CardT) String() string {
	return strconv.Itoa(int(s))
}

///////////////////////////////
