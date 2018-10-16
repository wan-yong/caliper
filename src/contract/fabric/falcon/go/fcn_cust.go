package main

import (
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

//查询零售商售卖记录
func (c *Chaincode) queryOrdersToCust(stub shim.ChaincodeStubInterface, args []string) []Stock {
	customerBoughts, _ := stub.GetStateByPartialCompositeKey(TYPE_DELIVERY_RETL_OUT, args)
	customerBoughtsList := make([]Stock, 0)
	for customerBoughts.HasNext() {
		customerBought, _ := customerBoughts.Next()
		_, keys, _ := stub.SplitCompositeKey(customerBought.Key)
		customerBoughtStock := Stock{
			ProdCode:      keys[1],
			ProdSKU:       "C",
			BarCode:       keys[0],
			ParentBarCode: keys[2],
		}
		customerBoughtsList = append(customerBoughtsList, customerBoughtStock)
	}
	logger.Error(customerBoughtsList)
	return customerBoughtsList
}

//产品溯源
func (c *Chaincode) queryProductTrail(stub shim.ChaincodeStubInterface, barCode string) peer.Response {
	cBarCode := barCode
	//get retailer trail
	saleRec, _ := stub.GetStateByPartialCompositeKey(TYPE_DELIVERY_RETL_OUT, []string{cBarCode})
	if !saleRec.HasNext() {
		return shim.Error("No sale record found with barcode: " + cBarCode)
	}
	sRec, _ := saleRec.Next()
	_, saleKeys, _ := stub.SplitCompositeKey(sRec.Key)

	bBarCode := saleKeys[2]
	//get distributor trail
	deliveryDistOutFound := false
	distKeys := make([]string, 0)
	for _, retlCodeType := range PRESET_RETL_CODES {
		distDlRec, _ := stub.GetStateByPartialCompositeKey(TYPE_DELIVERY_DIST_OUT, []string{retlCodeType, bBarCode})
		if distDlRec.HasNext() {
			distRec, _ := distDlRec.Next()
			_, distKeys, _ = stub.SplitCompositeKey(distRec.Key)
			deliveryDistOutFound = true
			break
		}
	}
	if !deliveryDistOutFound {
		return shim.Error("No distributor record found with barcode: " + bBarCode)
	}

	//get manufacturer trail
	deliveryManuOutFound := false
	manuKeys := make([]string, 0)
	for _, distCodeType := range PRESET_DIST_CODES {
		manuDlRec, _ := stub.GetStateByPartialCompositeKey(TYPE_DELIVERY_MANU_OUT, []string{distCodeType, bBarCode})
		if manuDlRec.HasNext() {
			manuDlRec, _ := manuDlRec.Next()
			_, manuKeys, _ = stub.SplitCompositeKey(manuDlRec.Key)
			deliveryManuOutFound = true
			break
		}
	}
	if !deliveryManuOutFound {
		return shim.Error("No Manufacturer record found with barcode: " + bBarCode)
	}

	productTrail := ProductTrail{
		BarCode:      saleKeys[0],
		ProdCode:     saleKeys[1],
		RetlSoldTime: saleKeys[3],
		RetlOrderNo:  saleKeys[4],
		RetlCode:     distKeys[0],
		DistBarCode:  distKeys[1],
		DistDlTime:   distKeys[3],
		DistOrderNo:  distKeys[4],
		ManuBarCode:  manuKeys[0],
		ManuDlTime:   manuKeys[2],
		ManuOrderNo:  manuKeys[3],
	}
	jsonBytes, _ := json.Marshal(productTrail)
	return shim.Success(jsonBytes)
}
