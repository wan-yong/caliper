/*
 * The sample smart contract
 * Writing the Blockchain Application
 * Author: Wan Yong
 */

 package main

 /* Imports
  * 4 utility libraries for handling bytes, reading and writing JSON, formatting and string manipulation
  * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
  */
 import (
	 "bytes"
	 "encoding/json"
	 "fmt"
	 "strconv"
 
	 "github.com/hyperledger/fabric/core/chaincode/shim"
	 sc "github.com/hyperledger/fabric/protos/peer" //sc代表缩写
 )
 
 // Define the Smart Contract structure
 type SmartContract struct {
 }
 
 // Define the credit structure, with 9 properties.  Structure tags are used by encoding/json library
 type PersonalCredit struct {
	 CreditCardLoan   string `json:"creditcardl"`
	 HouseLoan  string `json:"housel"`
	 OtherLoan string `json:"otherl"`
	 CivilJudgement  string `json:"civilj"`
	 TaxIssue string `json:"tax"`
	 Tele string `json:"tele"`
	 Electricity string `json:"electricity"`
	 Water string `json:"water"`
	 Gas string `json:"gas"`
 }

 /*
  * The Init method is called when the Smart Contract "credit" is instantiated by the blockchain network
  * Best practice is to have any Ledger initialization in separate function
  */
 func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	 return shim.Success(nil)
 }
 
 /*
  * The Invoke method is called as a result of an application request to run the Smart Contract "credit"
  * The calling application program has also specified the particular smart contract function to be called, with arguments
  */
 func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
 
	 // Retrieve function and arguments
	 function, args := APIstub.GetFunctionAndParameters()
	 // Route to the handler function to interact with the ledger
	 if function == "queryCredit" {
		 return s.queryCredit(APIstub, args)
	 } else if function == "initLedger" {
		 return s.initLedger(APIstub)
	 } else if function == "createCredit" {
		 return s.createCredit(APIstub, args)
	 } else if function == "queryAllCredits" {
		 return s.queryAllCredits(APIstub)
	 } else if function == "changeCardLoanCredit" {
		 return s.changeCardLoanCredit(APIstub, args)
	 }
 
	 return shim.Error("Invalid Smart Contract function name.")
 }
 
 func (s *SmartContract) queryCredit(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	 if len(args) != 1 {
		 return shim.Error("Incorrect number of arguments. Expecting 1")
	 }
 
	 creditAsBytes, _ := APIstub.GetState(args[0])
	 return shim.Success(creditAsBytes)
 }
 
 func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	 credits := []PersonalCredit{ //用固定值，做数据铺底
		 PersonalCredit{CreditCardLoan: "-1000000", HouseLoan: "-50000000", OtherLoan: "-3000000", CivilJudgement: "dong cheng fa yuan", TaxIssue: "0", Tele: "-500", Electricity: "0", Water: "-200", Gas: "0"},
		 PersonalCredit{CreditCardLoan: "-500000", HouseLoan: "-70000000", OtherLoan: "-2000000", CivilJudgement: "no", TaxIssue: "0", Tele: "-500", Electricity: "0", Water: "-10", Gas: "0"},
		 PersonalCredit{CreditCardLoan: "-800000", HouseLoan: "-90000000", OtherLoan: "-1000000", CivilJudgement: "chao yang fa yuan", TaxIssue: "-200", Tele: "0", Electricity: "0", Water: "-50", Gas: "0"},
	 }
 
	 i := 0
	 for i < len(credits) {
		 fmt.Println("i is ", i)
		 creditAsBytes, _ := json.Marshal(credits[i])
		 APIstub.PutState("PersonalCredit"+strconv.Itoa(i), creditAsBytes) //第一个参数是"键值"
		 fmt.Println("Added", credits[i])
		 i = i + 1
	 }
 
	 return shim.Success(nil)
 }
 
 func (s *SmartContract) createCredit(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	 if len(args) != 10 {
		 return shim.Error("Incorrect number of arguments. Expecting 5")
	 }
 
	 var credit = PersonalCredit{CreditCardLoan: args[1], HouseLoan: args[2], OtherLoan: args[3], CivilJudgement: args[4], TaxIssue: args[5], Tele: args[6], Electricity: args[7], Water: args[8], Gas: args[9]}
 
	 creditAsBytes, _ := json.Marshal(credit)
	 APIstub.PutState(args[0], creditAsBytes)
 
	 return shim.Success(nil)
 }
 
 func (s *SmartContract) queryAllCredits(APIstub shim.ChaincodeStubInterface) sc.Response {
 
	 startKey := "PersonalCredit0"
	 endKey := "PersonalCredit999"
 
	 resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	 if err != nil {
		 return shim.Error(err.Error())
	 }
	 defer resultsIterator.Close()
 
	 // buffer is a JSON array containing QueryResults
	 var buffer bytes.Buffer
	 buffer.WriteString("[")
 
	 bArrayMemberAlreadyWritten := false
	 for resultsIterator.HasNext() {
		 queryResponse, err := resultsIterator.Next()
		 if err != nil {
			 return shim.Error(err.Error())
		 }
		 // 在数组成员之前添加逗号，忽略第一个数组成员
		 if bArrayMemberAlreadyWritten == true {
			 buffer.WriteString(",")
		 }
		 buffer.WriteString("{\"Key\":")
		 buffer.WriteString("\"")
		 buffer.WriteString(queryResponse.Key)
		 buffer.WriteString("\"")
 
		 buffer.WriteString(", \"Record\":")
		 // Record is a JSON object, so we write as-is
		 buffer.WriteString(string(queryResponse.Value))
		 buffer.WriteString("}")
		 bArrayMemberAlreadyWritten = true
	 }
	 buffer.WriteString("]")
 
	 fmt.Printf("- queryAllCredits:\n%s\n", buffer.String())
 
	 return shim.Success(buffer.Bytes())
 }
 
 func (s *SmartContract) changeCardLoanCredit(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	 if len(args) != 2 {
		 return shim.Error("Incorrect number of arguments. Expecting 2")
	 }
 
	 creditAsBytes, _ := APIstub.GetState(args[0])
	 credit := PersonalCredit{}
 
	 json.Unmarshal(creditAsBytes, &credit)
	 credit.CreditCardLoan = args[1]
 
	 creditAsBytes, _ = json.Marshal(credit)
	 APIstub.PutState(args[0], creditAsBytes)
 
	 return shim.Success(nil)
 }
 
 // The main function is only relevant in unit test mode. Only included here for completeness.
 func main() {
 
	 // Create a new Smart Contract
	 err := shim.Start(new(SmartContract))
	 if err != nil {
		 fmt.Printf("Error creating new Smart Contract: %s", err)
	 }
 }
 