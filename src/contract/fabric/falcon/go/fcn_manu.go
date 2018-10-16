package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// 查询生产商产品目录
func (c *Chaincode) queryProducts(stub shim.ChaincodeStubInterface, args []string) []Product {
	prods, _ := stub.GetStateByPartialCompositeKey(TYPE_MANU_PRODUCT, args)
	prodList := make([]Product, 0)
	for prods.HasNext() {
		prod, _ := prods.Next()
		_, keys, _ := stub.SplitCompositeKey(prod.Key)
		price, _ := strconv.ParseUint(keys[3], 10, 16)
		prodList = append(prodList, Product{keys[1], keys[2], uint(price), keys[0]})
	}
	return prodList
}

// 查询生产商库存
func (c *Chaincode) queryStocksOnManu(stub shim.ChaincodeStubInterface, args []string) []Stock {
	stocks, _ := stub.GetStateByPartialCompositeKey(TYPE_MANU_STOCK, args)
	stockList := make([]Stock, 0)
	for stocks.HasNext() {
		stock, _ := stocks.Next()
		_, keys, _ := stub.SplitCompositeKey(stock.Key)
		stockList = append(stockList, Stock{keys[1], keys[0], keys[2], keys[3]})
	}
	return stockList
}

// 查询生产商发货（至经销商）记录
func (c *Chaincode) queryOrdersToDist(stub shim.ChaincodeStubInterface, args []string) []OrderToDist {
	manuDeliveryOrders, _ := stub.GetStateByPartialCompositeKey(TYPE_DELIVERY_MANU_OUT, args)
	manuDeliveryOrderList := make([]OrderToDist, 0)
	for manuDeliveryOrders.HasNext() {
		manuDeliveryOrder, _ := manuDeliveryOrders.Next()
		_, keys, _ := stub.SplitCompositeKey(manuDeliveryOrder.Key)
		OrderToDist := OrderToDist{
			OrderNo:      keys[4],
			DeliveryTime: keys[3],
			ProdCode:     keys[2],
			BarCode:      keys[1],
			DistCode:     keys[0],
		}
		manuDeliveryOrderList = append(manuDeliveryOrderList, OrderToDist)
		//manuDeliveryOrderList = append(manuDeliveryOrderList, OrderToDist{keys[4], keys[3], keys[2], keys[1], keys[0]})
	}
	return manuDeliveryOrderList
}

//生产商新增库存
func (c *Chaincode) createOnManu(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	//get max box Barcode
	bStocks, _ := stub.GetStateByPartialCompositeKey(TYPE_MANU_STOCK, []string{"B"})
	var maxbCode = 0
	for bStocks.HasNext() {
		bStock, _ := bStocks.Next()
		_, bKeys, _ := stub.SplitCompositeKey(bStock.Key)
		bCode, _ := strconv.Atoi(bKeys[2])
		if bCode > maxbCode {
			maxbCode = bCode
		}
	}

	//get max can Barcode
	cStocks, _ := stub.GetStateByPartialCompositeKey(TYPE_MANU_STOCK, []string{"C"})
	var maxcCode = 0
	for cStocks.HasNext() {
		cStock, _ := cStocks.Next()
		_, cKeys, _ := stub.SplitCompositeKey(cStock.Key)
		cCode, _ := strconv.Atoi(cKeys[3])
		if cCode > maxcCode {
			maxcCode = cCode
		}
	}

	var mProdCount = 0
	for _, arg := range args {
		keys := strings.Split(arg, ":")
		ProdCode := keys[0]
		ProdCount, _ := strconv.Atoi(keys[1])
		if ProdCount > mProdCount {
			mProdCount = ProdCount
		}
		for i := 0; i < ProdCount; i++ {
			maxbCode = maxbCode + 1
			skuBKey, _ := stub.CreateCompositeKey(TYPE_MANU_STOCK, []string{"B", ProdCode, strconv.Itoa(maxbCode), ""})
			stub.PutState(skuBKey, []byte{0})
			for j := 0; j < 6; j++ {
				cProdCode := strings.Replace(ProdCode, "B", "C", 1)
				maxcCode = maxcCode + 1
				skuCKey, _ := stub.CreateCompositeKey(TYPE_MANU_STOCK, []string{"C", strconv.Itoa(maxbCode), cProdCode, strconv.Itoa(maxcCode)})
				stub.PutState(skuCKey, []byte{0})
			}
		}
	}

	//返回chaincode执行结果
	createdCount := strconv.Itoa(mProdCount)
	msg := "Create " + createdCount + " Product(s) on manufacture successfully"
	return shim.Success([]byte(msg))
	//return shim.Success(nil)
	
}

// 生产商发货（至经销商）
func (c *Chaincode) deliverToDist(stub shim.ChaincodeStubInterface, args []string) (peer.Response, string) {
	logger.Info("DeliverToDist chaincode args:", args)
	dlTimestamp := time.Now().In(cstSH).Format("2006-01-02 15:04:05")
	DeliveryOrderSeq++
	dlOrderNo := "MDL" + fmt.Sprintf("%06s", strconv.Itoa(DeliveryOrderSeq))
	dlDistCode := ""

	for _, dl := range args {
		keys := strings.Split(dl, ":")
		dlDistCode = keys[0]
		dlProdCode := keys[1]
		dlQuantity, _ := strconv.Atoi(keys[2])
		logger.Info("Checking manu stock for product ["+dlProdCode+"] with quantity", strconv.Itoa(dlQuantity))
		bStocks, _ := stub.GetStateByPartialCompositeKey(TYPE_MANU_STOCK, []string{"B", dlProdCode})
		for dlQuantity > 0 && bStocks.HasNext() {
			bStock, _ := bStocks.Next()
			_, bStockKeys, _ := stub.SplitCompositeKey(bStock.Key)
			bCode := bStockKeys[2]
			logger.Debug("Delivering product:", bCode, " to Distributor:", dlDistCode)
			cStocks, _ := stub.GetStateByPartialCompositeKey(TYPE_MANU_STOCK, []string{"C", bCode})
			for cStocks.HasNext() {
				cStock, _ := cStocks.Next()
				_, cStockKeys, _ := stub.SplitCompositeKey(cStock.Key)
				cCode := cStockKeys[3]
				cProdCode := cStockKeys[2]
				skuKey, _ := stub.CreateCompositeKey(TYPE_DELIVERY_SKU, []string{"B", dlProdCode, bCode, "C", cProdCode, cCode})
				stub.PutState(skuKey, []byte{0})
				err := stub.DelState(cStock.Key)
				if err != nil {
					logger.Error("Fail to remove product from stock - C in manu:" + cCode)
				} else {
					logger.Debug("Removed product from stock - C in manu:" + cCode)
				}
			}
			trailKey, _ := stub.CreateCompositeKey(TYPE_DELIVERY_MANU_OUT, []string{dlDistCode, bCode, dlProdCode, dlTimestamp, dlOrderNo})
			stub.PutState(trailKey, []byte{0})
			err := stub.DelState(bStock.Key)
			if err != nil {
				logger.Error("Fail to remove product from stock - B in manu:" + bCode)
			} else {
				logger.Debug("Removed product from stock - B in manu:" + bCode)
			}
			dlQuantity--
		}

		if dlQuantity > 0 {
			return shim.Error("Not enough stock for product in manu: " + dlProdCode), ""
		}
	}
	msg := "Deliver to distributor [" + dlDistCode + "] successfully - OrderNo: " + dlOrderNo + ", " + dlTimestamp
	return shim.Success([]byte(msg)), dlOrderNo
}
