package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// 查询经销商收货记录
func (c *Chaincode) queryReceiptsOnDist(stub shim.ChaincodeStubInterface, args []string) []ReceiptOnDist {
	receiptsOnDistOrders, _ := stub.GetStateByPartialCompositeKey(TYPE_DELIVERY_DIST_IN, args)
	receiptsOnDistOrderList := make([]ReceiptOnDist, 0)
	for receiptsOnDistOrders.HasNext() {
		receiptsOnDistOrder, _ := receiptsOnDistOrders.Next()
		_, keys, _ := stub.SplitCompositeKey(receiptsOnDistOrder.Key)
		ReceiptOnDist := ReceiptOnDist{
			OrderNo:     keys[4],
			ReceiptTime: keys[3],
			ProdCode:    keys[2],
			BarCode:     keys[1],
			DistCode:    keys[0],
		}
		receiptsOnDistOrderList = append(receiptsOnDistOrderList, ReceiptOnDist)
		//receiptsOnDistOrderList = append(receiptsOnDistOrderList, ReceiptOnDist{keys[3], keys[2], keys[1], keys[0]})
	}
	return receiptsOnDistOrderList
}

// 查询经销商库存
func (c *Chaincode) queryStocksOnDist(stub shim.ChaincodeStubInterface, args []string) []Stock {
	stocks, _ := stub.GetStateByPartialCompositeKey(TYPE_DIST_STOCK, args)
	stockList := make([]Stock, 0)
	for stocks.HasNext() {
		stock, _ := stocks.Next()
		_, keys, _ := stub.SplitCompositeKey(stock.Key)
		stockDist := Stock{
			ProdCode:      keys[2],
			ProdSKU:       keys[1],
			BarCode:       keys[3],
			ParentBarCode: keys[4],
		}
		stockList = append(stockList, stockDist)
	}
	return stockList
}

// 查询经销商发货记录
func (c *Chaincode) queryOrdersToRetl(stub shim.ChaincodeStubInterface, args []string) []OrderToRetl {
	receiptsOnRtlOrders, _ := stub.GetStateByPartialCompositeKey(TYPE_DELIVERY_DIST_OUT, args)
	receiptsOnRtlOrderList := make([]OrderToRetl, 0)
	for receiptsOnRtlOrders.HasNext() {
		distDeliveryOutOrder, _ := receiptsOnRtlOrders.Next()
		_, keys, _ := stub.SplitCompositeKey(distDeliveryOutOrder.Key)
		distDeliverToRetl := OrderToRetl{
			OrderNo:      keys[4],
			DeliveryTime: keys[3],
			ProdCode:     keys[2],
			BarCode:      keys[1],
			RetlCode:     keys[0],
		}
		receiptsOnRtlOrderList = append(receiptsOnRtlOrderList, distDeliverToRetl)
	}
	return receiptsOnRtlOrderList
}

// 经销商收货
func (c *Chaincode) receiveOnDist(stub shim.ChaincodeStubInterface, args []string) (peer.Response, string) {
	logger.Error("ReceiveOnDist chaincode args:", args)
	receiptTime := time.Now().In(cstSH).Format("2006-01-02 15:04:05")
	manuOuts, _ := stub.GetStateByPartialCompositeKey(TYPE_DELIVERY_MANU_OUT, []string{})

	var recOrderNo = ""
	for _, distOrderNo := range args {
		for manuOuts.HasNext() {
			manuOut, _ := manuOuts.Next()
			_, manuOutKeys, _ := stub.SplitCompositeKey(manuOut.Key)
			manuOrderNo := manuOutKeys[4]
			if manuOrderNo == distOrderNo {
				distCode := manuOutKeys[0]
				bCode := manuOutKeys[1]
				recProdCode := manuOutKeys[2]
				recOrderNo = distOrderNo
				receiveOnDistKey, _ := stub.CreateCompositeKey(TYPE_DELIVERY_DIST_IN, []string{distCode, bCode, recProdCode, receiptTime, recOrderNo})
				stub.PutState(receiveOnDistKey, []byte{0})
				//创建Dist库存
				distStockKey, _ := stub.CreateCompositeKey(TYPE_DIST_STOCK, []string{distCode, "B", recProdCode, bCode, ""})
				stub.PutState(distStockKey, []byte{0})
			}
		}
	}
	msg := "Receive delivery order(s) successfully - OrderNo(s): " + recOrderNo
	return shim.Success([]byte(msg)), recOrderNo
}

// 经销商发货（至零售商）
func (c *Chaincode) deliverToRetl(stub shim.ChaincodeStubInterface, args []string) (peer.Response, string) {
	logger.Info("deliverToRetl chaincode args:", args)
	dlTimestamp := time.Now().In(cstSH).Format("2006-01-02 15:04:05")
	DeliveryOrderSeq++
	dlOrderNo := "DDL" + fmt.Sprintf("%06s", strconv.Itoa(DeliveryOrderSeq))
	dlRetlCode := ""

	for _, dl := range args {
		keys := strings.Split(dl, ":")
		dlRetlCode = keys[0]
		dlProdCode := keys[1]
		dlQuantity, _ := strconv.Atoi(keys[2])
		logger.Info("Checking dist stock for product ["+dlProdCode+"] with quantity", strconv.Itoa(dlQuantity))

		bStocks, _ := stub.GetStateByPartialCompositeKey(TYPE_DIST_STOCK, []string{dlRetlCode, "B", dlProdCode})
		for dlQuantity > 0 && bStocks.HasNext() {
			bStock, _ := bStocks.Next()
			_, bStockKeys, _ := stub.SplitCompositeKey(bStock.Key)
			bCode := bStockKeys[3]
			logger.Debug("Delivering product:", bCode, " to Retailer:", dlRetlCode)
			distDlKey, _ := stub.CreateCompositeKey(TYPE_DELIVERY_DIST_OUT, []string{dlRetlCode, bCode, dlProdCode, dlTimestamp, dlOrderNo})
			stub.PutState(distDlKey, []byte{0})
			err := stub.DelState(bStock.Key)
			if err != nil {
				logger.Error("Fail to remove product from stock - B in dist:" + bCode)
			} else {
				logger.Debug("Removed product from stock - B in dist:" + bCode)
			}
			dlQuantity--
		}
		if dlQuantity > 0 {
			return shim.Error("Not enough stock for product in dist: " + dlProdCode), ""
		}
	}
	msg := "Deliver to retailer [" + dlRetlCode + "] successfully - OrderNo: " + dlOrderNo + ", " + dlTimestamp
	return shim.Success([]byte(msg)), dlOrderNo
}
