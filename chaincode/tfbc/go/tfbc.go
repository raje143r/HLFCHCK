/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Trade Finance Use Case - WORK IN  PROGRESS
 */

 package main


 import (
	 "bytes"
	 "encoding/json"
	 "fmt"
	 "strconv"
	 "time"
 
	 "github.com/hyperledger/fabric/core/chaincode/shim"
	 sc "github.com/hyperledger/fabric/protos/peer"
 )
 
 const (
	 SOURCE		= "SOURCE"
	 DESTINATION	= "DESTINATION"
 )
 
 
 
 // Define the Smart Contract structure
 type SmartContract struct {
 }
 
 
 // Define the letter of credit
 type LetterOfCredit struct {
	 LCId			string		`json:"lcId"`
	 ExpiryDate		string		`json:"expiryDate"`
	 Buyer    		string   	`json:"buyer"`
	 BuyerBank		string		`json:"buyerbank"`
	 SellerBank		string		`json:"sellerbank"`
	 Seller			string		`json:"seller"`
	 Carrier 		string      `json:"carrier"`
	 Amount			int			`json:"amount,int"`
	 PaidAmount		int			`json:"paidamount,int"`
	 Status			string		`json:"status"`
	 Awb				string      `json:"awb"`
 }
 
 
 func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	 return shim.Success(nil)
 }
 
 func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
 
	 // Retrieve the requested Smart Contract function and arguments
	 function, args := APIstub.GetFunctionAndParameters()
	 // Route to the appropriate handler function to interact with the ledger appropriately
	 if function == "requestLC" {
		 return s.requestLC(APIstub, args)
	 } else if function == "issueLC" {
		 return s.issueLC(APIstub, args)
	 } else if function == "acceptLC" {
		 return s.acceptLC(APIstub, args)
	 }else if function == "getLC" {
		 return s.getLC(APIstub, args)
	 }else if function == "getLCHistory" {
		 return s.getLCHistory(APIstub, args)
	 }else if function == "rejectLC" {
		 return s.rejectLC(APIstub, args)
	 }else if function == "updateCarrier" {
		 return s.updateCarrier(APIstub, args)
	 }else if function == "shipGoods" {
		 return s.shipGoods(APIstub, args)
	 }else if function == "paySeller" {
		 return s.paySeller(APIstub, args)
 
 
	 }
 
	 return shim.Error("Invalid Smart Contract function name.")
 }
 
 
 func getShipmentLocationKey(stub shim.ChaincodeStubInterface, tradeID string) (string, error) {
	 shipmentLocationKey, err := stub.CreateCompositeKey("Shipment", []string{"Location", tradeID})
	 if err != nil {
		 return "", err
	 } else {
		 return shipmentLocationKey, nil
	 }
 }
 
 //
 // This function is initiate by Buyer 
 //
 func (s *SmartContract) requestLC(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
 
	 LC := LetterOfCredit{}
 
	 err  := json.Unmarshal([]byte(args[0]),&LC)
 if err != nil {
		 return shim.Error("Not able to parse args into LC")
	 }
 
 // Verify that the Amount requested is within the limits of Buyer's bank
		if LC.Amount > 1000000 || LC.Amount < 10000 {
			fmt.Printf("The amount requested is not within the limits of Buyer bank to issue a Letter of Credit.\n")
			return shim.Error("The amount requested is not within the limits of Buyer bank to issue a Letter of Credit.\n")
		}
		LC.Status = "New"
	 LCBytes, err := json.Marshal(LC)
	 APIstub.PutState(LC.LCId,LCBytes)
	 fmt.Println("LC Requested -> ", LC)
	 response:=struct{TxID string `json:"txid"`}{TxID: APIstub.GetTxID()}
	 responsebytes,err := json.Marshal(response)
	 return shim.Success(responsebytes)
 }
 //
 // This function is initiate by Buyer.
 //
 func (s *SmartContract) issueLC(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	 lcID := struct {
		 LcID  string `json:"lcID"`
	 }{}
	 err  := json.Unmarshal([]byte(args[0]),&lcID)
	 if err != nil {
		 return shim.Error("Not able to parse args into LCID")
	 }
	 
	 // if err != nil {
	 // 	return shim.Error("No Amount")
	 // }
 
	 LCAsBytes, _ := APIstub.GetState(lcID.LcID)
 
	 var lc LetterOfCredit
 
	 err = json.Unmarshal(LCAsBytes, &lc)
 
	 if err != nil {
		 return shim.Error("Issue with LC json unmarshaling")
	 }
 
    // Verify that the LC is already in Progress before Issueing 
	if (lc.Status == "Accepted" || lc.Status == "Rejected" || lc.Status == "ReadyforShip"  || lc.Status == "Shipped" ||lc.Status ==  "Paid-AWB"){
		fmt.Printf("Letter of Credit request for trade %s is already In Process\n", lcID.LcID)
		return shim.Error("Letter of Credit request#: is already In Process, Operation Not Allowed" )
	}
	 LC := LetterOfCredit{LCId: lc.LCId, ExpiryDate: lc.ExpiryDate, Buyer: lc.Buyer, BuyerBank: lc.BuyerBank,SellerBank: lc.SellerBank, Seller: lc.Seller, Amount: lc.Amount, Status: "Issued"}
	 LCBytes, err := json.Marshal(LC)
 
	 if err != nil {
		 return shim.Error("Issue with LC json marshaling")
	 }
 
	 APIstub.PutState(lc.LCId,LCBytes)
	 fmt.Println("LC Issued -> ", LC)
 
 
	 return shim.Success(nil)
 }
 
 //
 // This function is initiate by Sellerbank 
 //
 func (s *SmartContract) acceptLC(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	 lcID := struct {
		 LcID  string `json:"lcID"`
	 }{}
	 err  := json.Unmarshal([]byte(args[0]),&lcID)
	 if err != nil {
		 return shim.Error("Not able to parse args into LC")
	 }
 
	 LCAsBytes, _ := APIstub.GetState(lcID.LcID)
 
	 var lc LetterOfCredit
 
	 err = json.Unmarshal(LCAsBytes, &lc)
 
	 if err != nil {
		 return shim.Error("Issue with LC json unmarshaling")
	 }
 
 
	 LC := LetterOfCredit{LCId: lc.LCId, ExpiryDate: lc.ExpiryDate, Buyer: lc.Buyer, SellerBank: lc.SellerBank, Seller: lc.Seller, Amount: lc.Amount, Status: "Accepted"}
	 LCBytes, err := json.Marshal(LC)
 
	 if err != nil {
		 return shim.Error("Issue with LC json marshaling")
	 }
 
	 APIstub.PutState(lc.LCId,LCBytes)
	 fmt.Println("LC Accepted -> ", LC)
 
	 return shim.Success(nil)
 }
 
 //
 // This function is initiate by Buyer/Seller 
 //
 func (s *SmartContract) rejectLC(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	 lcID := struct {
			 LcID  string `json:"lcID"`
	 }{}
	 err  := json.Unmarshal([]byte(args[0]),&lcID)
	 if err != nil {
			 return shim.Error("Not able to parse args into LC")
	 }
 
	 LCAsBytes, _ := APIstub.GetState(lcID.LcID)
	 var lc LetterOfCredit
	 err = json.Unmarshal(LCAsBytes, &lc)
 
	 if err != nil {
			 return shim.Error("Issue with LC json unmarshaling")
	 }
 
 
	 //  LC := LetterOfCredit{LCId: lc.LCId, ExpiryDate: lc.ExpiryDate, Buyer: lc.Buyer, SellerBank: lc.SellerBank, Seller: lc.Seller, Amount: lc.Amount, Status: "Rejected", BuyerBank: lc.BuyerBank}
	 LC := LetterOfCredit{LCId: lc.LCId, ExpiryDate: lc.ExpiryDate, Buyer: lc.Buyer, SellerBank: lc.SellerBank, Seller: lc.Seller, Amount: lc.Amount, Status: "Rejcted"}
	 LCBytes, err := json.Marshal(LC)
	 if err != nil {
			 return shim.Error("Issue with LC json marshaling")
	 }
 
	 APIstub.PutState(lc.LCId,LCBytes)
	 fmt.Println("LC Rejected -> ", LC)
 
	 return shim.Success(nil)
 }
 
 
 //
 // This function is initiate by Seller 
 //
 func (s *SmartContract) updateCarrier(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	 lcID := struct {
			 LcID  string `json:"lcID"`	
			 Carrier string 	`json:"carrier"`
	 }{}
 
	 err  := json.Unmarshal([]byte(args[0]),&lcID)
	 if err != nil {
			 return shim.Error("Not able to parse args into LC")
	 }
 
	 LCAsBytes, _ := APIstub.GetState(lcID.LcID)
	 var lc LetterOfCredit
	 err = json.Unmarshal(LCAsBytes, &lc)
	 if err != nil {
			 return shim.Error("Issue with LC json unmarshaling")
	 }
 
	 // Verify that the LC has been agreed to
	 if lc.Status != "Accepted" {
		 fmt.Printf("Letter of Credit request for trade %s accepted\n", lcID.LcID)
		 return shim.Error("LC has not been accepted by the parties")
	 }
 
	 //  LC := LetterOfCredit{LCId: lc.LCId, ExpiryDate: lc.ExpiryDate, Buyer: lc.Buyer, SellerBank: lc.SellerBank, Seller: lc.Seller, Amount: lc.Amount, Status: "Rejected", BuyerBank: lc.BuyerBank}
	 LC := LetterOfCredit{LCId: lc.LCId, ExpiryDate: lc.ExpiryDate, Buyer: lc.Buyer, SellerBank: lc.SellerBank, Seller: lc.Seller, Amount: lc.Amount,Carrier : lcID.Carrier , Status: "ReadyforShip"}
	 LCBytes, err := json.Marshal(LC)
	 if err != nil {
			 return shim.Error("Issue with LC json marshaling")
	 }
 
	 APIstub.PutState(lc.LCId,LCBytes)
	 fmt.Println("LC UpdateCarrier -> ", LC)
 
	 return shim.Success(nil)
 }
 
 
 //
 // This function is initiate by Carrier 
 //
 func (s *SmartContract) shipGoods(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	 lcID := struct {
			 LcID  string `json:"lcID"`	
			 Awb   string `json:"Awb"`
	 }{}
 
	 err  := json.Unmarshal([]byte(args[0]),&lcID)
	 if err != nil {
			 return shim.Error("Not able to parse args into LC")
	 }
 
	 LCAsBytes, _ := APIstub.GetState(lcID.LcID)
	 var lc LetterOfCredit
	 err = json.Unmarshal(LCAsBytes, &lc)
 
	 if err != nil {
			 return shim.Error("Issue with LC json unmarshaling")
	 }
 
 
	 //  LC := LetterOfCredit{LCId: lc.LCId, ExpiryDate: lc.ExpiryDate, Buyer: lc.Buyer, SellerBank: lc.SellerBank, Seller: lc.Seller, Amount: lc.Amount, Status: "Rejected", BuyerBank: lc.BuyerBank}
	 LC := LetterOfCredit{LCId: lc.LCId, ExpiryDate: lc.ExpiryDate, Buyer: lc.Buyer, SellerBank: lc.SellerBank, Seller: lc.Seller, Amount: lc.Amount,Awb : lcID.Awb , Status: "Shipped"}
	 LCBytes, err := json.Marshal(LC)
	 if err != nil {
			 return shim.Error("Issue with LC json marshaling")
	 }
 
	 APIstub.PutState(lc.LCId,LCBytes)
	 fmt.Println("LC ShipStatus -> ", LC)
 
	 return shim.Success(nil)
 }
 
 //
 // This function is initiate by SellerBank 
 //
 func (s *SmartContract) paySeller(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	 lcID := struct {
			 LcID  string `json:"lcID"`	
			 PaidAmount   int `json:"paidamount"`
	 }{}
 
	 err  := json.Unmarshal([]byte(args[0]),&lcID)
	 if err != nil {
			 return shim.Error("Not able to parse args into LC")
	 }
 
	 LCAsBytes, _ := APIstub.GetState(lcID.LcID)
	 var lc LetterOfCredit
	 err = json.Unmarshal(LCAsBytes, &lc)
 
	 if err != nil {
			 return shim.Error("Issue with LC json unmarshaling")
	 }
 
 
	 //  LC := LetterOfCredit{LCId: lc.LCId, ExpiryDate: lc.ExpiryDate, Buyer: lc.Buyer, SellerBank: lc.SellerBank, Seller: lc.Seller, Amount: lc.Amount, Status: "Rejected", BuyerBank: lc.BuyerBank}
	 LC := LetterOfCredit{LCId: lc.LCId, ExpiryDate: lc.ExpiryDate, Buyer: lc.Buyer, SellerBank: lc.SellerBank, Seller: lc.Seller, Amount: lc.Amount,Awb : lc.Awb , PaidAmount: lcID.PaidAmount, Status: "Paid-AWB"}
	 LCBytes, err := json.Marshal(LC)
	 if err != nil {
			 return shim.Error("Issue with LC json marshaling")
	 }
 
	 APIstub.PutState(lc.LCId,LCBytes)
	 fmt.Println("LC PaymentStatus -> ", LC)
 
	 return shim.Success(nil)
 }
 
 
 //
 // This function is initiate by All Orgs
 //
 func (s *SmartContract) getLC(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	 lcId := args[0];
	 
	 // if err != nil {
	 // 	return shim.Error("No Amount")
	 // }
 
	 LCAsBytes, _ := APIstub.GetState(lcId)
	 return shim.Success(LCAsBytes)
 }
 
 //
 // This function is initiate by All Orgs
 //
 func (s *SmartContract) getLCHistory(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
 
	 lcId := args[0];
 
	 resultsIterator, err := APIstub.GetHistoryForKey(lcId)
	 if err != nil {
		 return shim.Error("Error retrieving LC history.")
	 }
	 defer resultsIterator.Close()
 
	 // buffer is a JSON array containing historic values for the marble
	 var buffer bytes.Buffer
	 buffer.WriteString("[")
 
	 bArrayMemberAlreadyWritten := false
	 for resultsIterator.HasNext() {
		 response, err := resultsIterator.Next()
		 if err != nil {
			 return shim.Error("Error retrieving LC history.")
		 }
		 // Add a comma before array members, suppress it for the first array member
		 if bArrayMemberAlreadyWritten == true {
			 buffer.WriteString(",")
		 }
		 buffer.WriteString("{\"TxId\":")
		 buffer.WriteString("\"")
		 buffer.WriteString(response.TxId)
		 buffer.WriteString("\"")
 
		 buffer.WriteString(", \"Value\":")
		 // if it was a delete operation on given key, then we need to set the
		 //corresponding value null. Else, we will write the response.Value
		 //as-is (as the Value itself a JSON marble)
		 if response.IsDelete {
			 buffer.WriteString("null")
		 } else {
			 buffer.WriteString(string(response.Value))
		 }
 
		 buffer.WriteString(", \"Timestamp\":")
		 buffer.WriteString("\"")
		 buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		 buffer.WriteString("\"")
 
		 buffer.WriteString(", \"IsDelete\":")
		 buffer.WriteString("\"")
		 buffer.WriteString(strconv.FormatBool(response.IsDelete))
		 buffer.WriteString("\"")
 
		 buffer.WriteString("}")
		 bArrayMemberAlreadyWritten = true
	 }
	 buffer.WriteString("]")
 
	 fmt.Printf("- getLCHistory returning:\n%s\n", buffer.String())
 
	 return shim.Success(buffer.Bytes())
 }
 
 // The main function is only relevant in unit test mode. Only included here for completeness.
 func main() {
 
	 // Create a new Smart Contract
	 err := shim.Start(new(SmartContract))
	 if err != nil {
		 fmt.Printf("Error creating new Smart Contract: %s", err)
	 }
 } 