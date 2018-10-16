package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("falcon")

var cstSH, _ = time.LoadLocation("Asia/Shanghai") //上海时区

type Chaincode struct{}

func (c *Chaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	c.initProdDB(stub)
	c.initStock(stub)
	return shim.Success([]byte("Chaincode initialized successfully."))
}

func (c *Chaincode) initProdDB(stub shim.ChaincodeStubInterface) {
	for _, p := range PRESET_PRODUCTS {
		var key, _ = stub.CreateCompositeKey(TYPE_MANU_PRODUCT, []string{p.SKU, p.Code, p.Name, strconv.Itoa(int(p.Price))})
		stub.PutState(key, []byte{0})
	}
}

func (c *Chaincode) initStock(stub shim.ChaincodeStubInterface) {
	for _, ps := range PRESET_STOCKS {
		var key string
		if ps.ProdSKU == "B" && ps.ParentBarCode == "" {
			key, _ = stub.CreateCompositeKey(TYPE_MANU_STOCK, []string{ps.ProdSKU, ps.ProdCode, ps.BarCode, ps.ParentBarCode})
		} else {
			key, _ = stub.CreateCompositeKey(TYPE_MANU_STOCK, []string{ps.ProdSKU, ps.ParentBarCode, ps.ProdCode, ps.BarCode})
		}
		stub.PutState(key, []byte{0})
	}
}

func (c *Chaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	var res peer.Response
	var orderNo string
	fcn, args := stub.GetFunctionAndParameters()
	switch fcn {
	case "ping":
		return shim.Success([]byte("Chaincode ping successfully."))
	case "query":
		return c.query(stub, args)
	case "createStock":
		return c.createStock(stub, args[0], args[1:])
	case "createOnManu":
		res = c.createOnManu(stub, args)
	case "deliverToDist":
		res, orderNo = c.deliverToDist(stub, args)
	case "receiveOnDist":
		res, orderNo = c.receiveOnDist(stub, args)
	case "deliverToRetl":
		res, orderNo = c.deliverToRetl(stub, args)
	case "receiveOnRetl":
		res, orderNo = c.receiveOnRetl(stub, args)
	case "sellToCust":
		res, orderNo = c.sellToCust(stub, args)
	default:
		return shim.Error("Unsupported chaincode function: " + fcn)
	}

	if res.Status == 200 {
		c.saveTxLog(stub, orderNo, fcn)
	}
	return res
}

func (c *Chaincode) query(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	switch typ := args[0]; strings.ToUpper(typ) {
	case TYPE_MANU_PRODUCT:
		jsonBytes, _ := json.Marshal(ProductList{c.queryProducts(stub, args[1:])})
		return shim.Success(jsonBytes)
	case TYPE_MANU_STOCK:
		jsonBytes, _ := json.Marshal(StockList{c.queryStocksOnManu(stub, args[1:])})
		return shim.Success(jsonBytes)
	case TYPE_DIST_STOCK:
		jsonBytes, _ := json.Marshal(StockList{c.queryStocksOnDist(stub, args[1:])})
		return shim.Success(jsonBytes)
	case TYPE_RETL_STOCK:
		jsonBytes, _ := json.Marshal(StockList{c.queryStocksOnRetl(stub, args[1:])})
		return shim.Success(jsonBytes)
	case TYPE_DELIVERY_MANU_OUT:
		jsonBytes, _ := json.Marshal(OrderListToDist{c.queryOrdersToDist(stub, args[1:])})
		return shim.Success(jsonBytes)
	case TYPE_DELIVERY_DIST_IN:
		jsonBytes, _ := json.Marshal(ReceiptListOnDist{c.queryReceiptsOnDist(stub, args[1:])})
		return shim.Success(jsonBytes)
	case TYPE_DELIVERY_DIST_OUT:
		jsonBytes, _ := json.Marshal(OrderListToRetl{c.queryOrdersToRetl(stub, args[1:])})
		return shim.Success(jsonBytes)
	case TYPE_DELIVERY_RETL_IN:
		jsonBytes, _ := json.Marshal(ReceiptListOnRetl{c.queryReceiptsOnRetl(stub, args[1:])})
		return shim.Success(jsonBytes)
	case TYPE_DELIVERY_RETL_OUT:
		jsonBytes, _ := json.Marshal(StockList{c.queryOrdersToCust(stub, args[1:])})
		return shim.Success(jsonBytes)
	case TYPE_DELIVERY_TRAIL:
		return c.queryProductTrail(stub, args[1])
	case TYPE_TX_LOG:
		jsonBytes, _ := json.Marshal(TxLogList{c.queryTxLog(stub, args[1:])})
		return shim.Success(jsonBytes)
	default:
		return shim.Error("Unknown query type: " + typ)
	}
}

// 发货单号自增器
var DeliveryOrderSeq = 0

func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}

// 增加各环节库存
func (c *Chaincode) createStock(stub shim.ChaincodeStubInterface, stockType string, args []string) peer.Response {
	switch stockType {
	case TYPE_RETL_STOCK:
		for _, arg := range args {
			keys := strings.Split(arg, ":")
			dlProdCode := keys[0]
			dlQuantity, _ := strconv.Atoi(keys[1])
			fmt.Print(dlProdCode, dlQuantity)
			key, _ := stub.CreateCompositeKey(TYPE_RETL_STOCK, []string{})
			stub.PutState(key, []byte{0})
		}
	}
	return shim.Success([]byte("Under construction ..."))
}

// 存储交易记录
func (c *Chaincode) saveTxLog(stub shim.ChaincodeStubInterface, orderNo string, dlTxType string) peer.Response {
	timestamp, _ := stub.GetTxTimestamp()
	txTime := time.Unix(timestamp.Seconds, 0).In(cstSH).Format("2006-01-02 15:04:05")
	if orderNo != "" {
		txKey, _ := stub.CreateCompositeKey(TYPE_TX_LOG, []string{stub.GetTxID(), orderNo, txTime, dlTxType})
		stub.PutState(txKey, []byte{0})
		return shim.Success([]byte("Transaction info was saved ...TxID: " + stub.GetTxID()))
	} else {
		return shim.Error("OrderNo is blank, fail to insert transaction information")
	}
}

// 查询TX记录
func (c *Chaincode) queryTxLog(stub shim.ChaincodeStubInterface, args []string) []TxLog {
	txInfos, _ := stub.GetStateByPartialCompositeKey(TYPE_TX_LOG, args)
	txInfoList := make([]TxLog, 0)
	for txInfos.HasNext() {
		txInfo, _ := txInfos.Next()
		_, keys, _ := stub.SplitCompositeKey(txInfo.Key)
		txInfoList = append(txInfoList, TxLog{keys[0], keys[2], keys[3], keys[1]})
	}
	return txInfoList
}
