package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// 查询零售商收货记录
func (c *Chaincode) queryReceiptsOnRetl(stub shim.ChaincodeStubInterface, args []string) []ReceiptOnRetl {
	receipts, _ := stub.GetStateByPartialCompositeKey(TYPE_DELIVERY_RETL_IN, args)
	receiptListOnRetl := make([]ReceiptOnRetl, 0)
	for receipts.HasNext() {
		receipt, _ := receipts.Next()
		_, keys, _ := stub.SplitCompositeKey(receipt.Key)
		receiptOnRetl := ReceiptOnRetl{
			OrderNo:     keys[4],
			ReceiveTime: keys[3],
			ProdCode:    keys[2],
			BarCode:     keys[1],
			RetlCode:    keys[0],
		}
		receiptListOnRetl = append(receiptListOnRetl, receiptOnRetl)
	}
	return receiptListOnRetl
}

// 查询零售商库存
func (c *Chaincode) queryStocksOnRetl(stub shim.ChaincodeStubInterface, args []string) []Stock {
	stocks, _ := stub.GetStateByPartialCompositeKey(TYPE_RETL_STOCK, args)
	stockList := make([]Stock, 0)
	for stocks.HasNext() {
		stock, _ := stocks.Next()
		_, keys, _ := stub.SplitCompositeKey(stock.Key)
		stockOnRetl := Stock{
			ProdCode:      keys[2],
			ProdSKU:       keys[1],
			BarCode:       keys[3],
			ParentBarCode: keys[4],
		}
		stockList = append(stockList, stockOnRetl)
	}
	return stockList
}

// 零售商收货
func (c *Chaincode) receiveOnRetl(stub shim.ChaincodeStubInterface, args []string) (peer.Response, string) {
	logger.Info("receiveOnRetl chaincode args:", args)
	receiptTime := time.Now().In(cstSH).Format("2006-01-02 15:04:05")
	var orderNo = ""
	for idx := range args {
		bStocks, _ := stub.GetStateByPartialCompositeKey(TYPE_DELIVERY_DIST_OUT, []string{})
		notFound := true
		for bStocks.HasNext() {
			bStock, _ := bStocks.Next()
			_, keys, _ := stub.SplitCompositeKey(bStock.Key)
			if orderNo = keys[4]; orderNo == args[idx] {
				notFound = false
				retlCode := keys[0]
				bCode := keys[1]
				recProdCode := keys[2]
				// 创建零售商收货记录
				retlInKey, _ := stub.CreateCompositeKey(TYPE_DELIVERY_RETL_IN, []string{retlCode, bCode, recProdCode, receiptTime, orderNo})
				stub.PutState(retlInKey, []byte{0})
				// 创建零售商库存
				cStocks, _ := stub.GetStateByPartialCompositeKey(TYPE_DELIVERY_SKU, []string{"B", recProdCode, bCode, "C"})
				for cStocks.HasNext() {
					cStock, _ := cStocks.Next()
					_, cStKeys, _ := stub.SplitCompositeKey(cStock.Key)
					retlStockKey, _ := stub.CreateCompositeKey(TYPE_RETL_STOCK, []string{retlCode, "C", cStKeys[4], cStKeys[5], bCode})
					stub.PutState(retlStockKey, []byte{0})
				}
			}
		}

		if notFound {
			return shim.Error("Order [" + args[idx] + "] was not found."), ""
		}
	}

	msg := "Receive delivery order(s) successfully - OrderNo(s): " + strings.Join(args, ", ")
	return shim.Success([]byte(msg)), orderNo
}

// 零售商销售（至消费者）
func (c *Chaincode) sellToCust(stub shim.ChaincodeStubInterface, args []string) (peer.Response, string) {
	logger.Info("SellToCust chaincode args:", args)
	slTimestamp := time.Now().In(cstSH).Format("2006-01-02 15:04:05")
	DeliveryOrderSeq++
	slOrderNo := "RSO" + fmt.Sprintf("%06s", strconv.Itoa(DeliveryOrderSeq))
	for _, sl := range args {
		keys := strings.Split(sl, ":")
		slRetlCode := keys[0]
		slProdCode := keys[1]
		slQuantity, _ := strconv.Atoi(keys[2])
		logger.Info("Checking retl stock for product ["+slProdCode+"] with quantity", strconv.Itoa(slQuantity))
		cStocks, _ := stub.GetStateByPartialCompositeKey(TYPE_RETL_STOCK, []string{slRetlCode, "C", slProdCode})
		for slQuantity > 0 && cStocks.HasNext() {
			cStock, _ := cStocks.Next()
			_, cStockKeys, _ := stub.SplitCompositeKey(cStock.Key)
			cCode := cStockKeys[3]
			bCode := cStockKeys[4]
			logger.Debug("Delivering product:", cCode)
			trailKey, _ := stub.CreateCompositeKey(TYPE_DELIVERY_RETL_OUT, []string{cCode, slProdCode, bCode, slTimestamp, slOrderNo})
			stub.PutState(trailKey, []byte{0})
			err := stub.DelState(cStock.Key)
			if err != nil {
				logger.Error("Fail to remove product from stock - C in retl:" + cCode)
			} else {
				logger.Debug("Removed product from stock - C in retl:" + cCode)
			}
			slQuantity--
		}

		if slQuantity > 0 {
			return shim.Error("Not enough stock for product in retl: " + slProdCode), ""
		}
	}

	msg := "Sell product(s) to customer successfully - OrderNo: " + slOrderNo
	return shim.Success([]byte(msg)), slOrderNo
}
