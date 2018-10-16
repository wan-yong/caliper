package main

type Product struct {
	Code  string
	Name  string
	Price uint
	SKU   string // "C" : Carton; "B" : Box
}

type ProductList struct {
	ProductList []Product
}

type Stock struct {
	ProdCode      string
	ProdSKU       string
	BarCode       string
	ParentBarCode string
}

type StockList struct {
	StockList []Stock
}

// 生产商 发货单
type OrderToDist struct {
	OrderNo      string
	DeliveryTime string
	ProdCode     string
	BarCode      string
	DistCode     string
}

type OrderListToDist struct {
	OrderList []OrderToDist
}

// 经销商 收货单
type ReceiptOnDist struct {
	OrderNo     string
	ReceiptTime string
	ProdCode    string
	BarCode     string
	DistCode    string
}

type ReceiptListOnDist struct {
	ReceiptList []ReceiptOnDist
}

// 经销商 发货单
type OrderToRetl struct {
	OrderNo      string
	DeliveryTime string
	ProdCode     string
	BarCode      string
	RetlCode     string
}

type OrderListToRetl struct {
	OrderList []OrderToRetl
}

// 零售商 收货单
type ReceiptOnRetl struct {
	OrderNo     string
	ReceiveTime string
	ProdCode    string
	BarCode     string
	RetlCode    string
}

type ReceiptListOnRetl struct {
	ReceiptList []ReceiptOnRetl
}

// 零售商 销售记录
type OrderToCust struct {
	OrderNo     string
	Customer    string
	PhoneNumber string
	SoldTime    string
	ProdCode    string
	BarCode     string
}

type OrderListToCust struct {
	OrderList []OrderToCust
}

// 产品溯源 数据结构
type ProductTrail struct {
	ProdCode     string
	BarCode      string
	RetlCode     string
	RetlOrderNo  string
	RetlSoldTime string

	DistBarCode string
	DistOrderNo string
	DistDlTime  string

	ManuBarCode string
	ManuOrderNo string
	ManuDlTime  string
}

type TxLog struct {
	TxID    string
	TxTime  string
	TxType  string
	TxOrder string
}

type TxLogList struct {
	TxLogList []TxLog
}

var PRESET_DIST_CODES = []string{
	"peer",
}

var PRESET_RETL_CODES = []string{
	"peer",
}

var PRESET_PRODUCTS = []Product{
	{"MPC0001", "Milk Powder - Newborn", 300, "C"},
	{"MPB0001", "Milk Powder - Newborn", 1800, "B"},
	{"MPC0002", "Milk Powder - Infant", 250, "C"},
	{"MPB0002", "Milk Powder - Infant", 1500, "B"},
	{"MPC0003", "Milk Powder - Toddler", 200, "C"},
	{"MPB0003", "Milk Powder - Toddler", 1200, "B"},
}

var PRESET_STOCKS = []Stock{
	{"MPB0001", "B", "2004331001", ""},
	{"MPC0001", "C", "3001002001", "2004331001"},
	{"MPC0001", "C", "3001002002", "2004331001"},
	{"MPC0001", "C", "3001002003", "2004331001"},
	{"MPC0001", "C", "3001002004", "2004331001"},
	{"MPC0001", "C", "3001002005", "2004331001"},
	{"MPC0001", "C", "3001002006", "2004331001"},

	{"MPB0001", "B", "2004331002", ""},
	{"MPC0001", "C", "3001002011", "2004331002"},
	{"MPC0001", "C", "3001002012", "2004331002"},
	{"MPC0001", "C", "3001002013", "2004331002"},
	{"MPC0001", "C", "3001002014", "2004331002"},
	{"MPC0001", "C", "3001002015", "2004331002"},
	{"MPC0001", "C", "3001002016", "2004331002"},

	{"MPB0001", "B", "2004331003", ""},
	{"MPC0001", "C", "3001002021", "2004331003"},
	{"MPC0001", "C", "3001002022", "2004331003"},
	{"MPC0001", "C", "3001002023", "2004331003"},
	{"MPC0001", "C", "3001002024", "2004331003"},
	{"MPC0001", "C", "3001002025", "2004331003"},
	{"MPC0001", "C", "3001002026", "2004331003"},

	{"MPB0001", "B", "2004331004", ""},
	{"MPC0001", "C", "3001002031", "2004331004"},
	{"MPC0001", "C", "3001002032", "2004331004"},
	{"MPC0001", "C", "3001002033", "2004331004"},
	{"MPC0001", "C", "3001002034", "2004331004"},
	{"MPC0001", "C", "3001002035", "2004331004"},
	{"MPC0001", "C", "3001002036", "2004331004"},

	{"MPB0001", "B", "2004331005", ""},
	{"MPC0001", "C", "3001002041", "2004331005"},
	{"MPC0001", "C", "3001002042", "2004331005"},
	{"MPC0001", "C", "3001002043", "2004331005"},
	{"MPC0001", "C", "3001002044", "2004331005"},
	{"MPC0001", "C", "3001002045", "2004331005"},
	{"MPC0001", "C", "3001002046", "2004331005"},

	{"MPB0001", "B", "2004331006", ""},
	{"MPC0001", "C", "3001002051", "2004331006"},
	{"MPC0001", "C", "3001002052", "2004331006"},
	{"MPC0001", "C", "3001002053", "2004331006"},
	{"MPC0001", "C", "3001002054", "2004331006"},
	{"MPC0001", "C", "3001002055", "2004331006"},
	{"MPC0001", "C", "3001002056", "2004331006"},

	{"MPB0002", "B", "2004332001", ""},
	{"MPC0002", "C", "3001002061", "2004332001"},
	{"MPC0002", "C", "3001002062", "2004332001"},
	{"MPC0002", "C", "3001002063", "2004332001"},
	{"MPC0002", "C", "3001002064", "2004332001"},
	{"MPC0002", "C", "3001002065", "2004332001"},
	{"MPC0002", "C", "3001002066", "2004332001"},

	{"MPB0002", "B", "2004332002", ""},
	{"MPC0002", "C", "3001002071", "2004332002"},
	{"MPC0002", "C", "3001002072", "2004332002"},
	{"MPC0002", "C", "3001002073", "2004332002"},
	{"MPC0002", "C", "3001002074", "2004332002"},
	{"MPC0002", "C", "3001002075", "2004332002"},
	{"MPC0002", "C", "3001002076", "2004332002"},

	{"MPB0002", "B", "2004332003", ""},
	{"MPC0002", "C", "3001002081", "2004332003"},
	{"MPC0002", "C", "3001002082", "2004332003"},
	{"MPC0002", "C", "3001002083", "2004332003"},
	{"MPC0002", "C", "3001002084", "2004332003"},
	{"MPC0002", "C", "3001002085", "2004332003"},
	{"MPC0002", "C", "3001002086", "2004332003"},

	{"MPB0002", "B", "2004332004", ""},
	{"MPC0002", "C", "3001002091", "2004332004"},
	{"MPC0002", "C", "3001002092", "2004332004"},
	{"MPC0002", "C", "3001002093", "2004332004"},
	{"MPC0002", "C", "3001002094", "2004332004"},
	{"MPC0002", "C", "3001002095", "2004332004"},
	{"MPC0002", "C", "3001002096", "2004332004"},

	{"MPB0003", "B", "2004333001", ""},
	{"MPC0003", "C", "3001002101", "2004333001"},
	{"MPC0003", "C", "3001002102", "2004333001"},
	{"MPC0003", "C", "3001002103", "2004333001"},
	{"MPC0003", "C", "3001002104", "2004333001"},
	{"MPC0003", "C", "3001002105", "2004333001"},
	{"MPC0003", "C", "3001002106", "2004333001"},

	{"MPB0003", "B", "2004333002", ""},
	{"MPC0003", "C", "3001002111", "2004333002"},
	{"MPC0003", "C", "3001002112", "2004333002"},
	{"MPC0003", "C", "3001002113", "2004333002"},
	{"MPC0003", "C", "3001002114", "2004333002"},
	{"MPC0003", "C", "3001002115", "2004333002"},
	{"MPC0003", "C", "3001002116", "2004333002"},

	{"MPB0003", "B", "2004333003", ""},
	{"MPC0003", "C", "3001002121", "2004333002"},
	{"MPC0003", "C", "3001002122", "2004333002"},
	{"MPC0003", "C", "3001002123", "2004333002"},
	{"MPC0003", "C", "3001002124", "2004333002"},
	{"MPC0003", "C", "3001002125", "2004333002"},
	{"MPC0003", "C", "3001002126", "2004333002"},
}
