/*
Copyright 33.cn Corp. 2018 All Rights Reserved.

Chaincode for Suning Corp.
*/

package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var (
	initialCredit = 1e8
	layout        = "2018-01-06 19:01:02"
	date          = "20180106190102"
	loc           *time.Location
)

func init() {
	loc, _ = time.LoadLocation("Asia/Shanghai")
}

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type BlackRecord struct {
	DocType          string `json:"docType"`
	RecordId         string `json:"recordId"`
	ClientId         string `json:"clientId"`
	ClientName       string `json:"clientName"`
	NegativeType     int    `json:"negativeType"`
	NegativeSeverity int    `json:"negativeSeverity"`
	NegativeInfo     string `json:"negativeInfo"`
	OrgAddr          string `json:"orgAddr"`
	Searchable       bool   `json:"searchable"`
	CreateTime       string `json:"createTime"`
	UpdateTime       string `json:"updateTime"`
}

func (blackRecord *BlackRecord) putBlackRecord(stub shim.ChaincodeStubInterface) error {
	brBytes, err := json.Marshal(blackRecord)
	if err != nil {
		fmt.Println("PutBlackRecord Marshal fail:", err.Error())
		return errors.New("PutBlackRecord Marshal fail:" + err.Error())
	}

	err = stub.PutState("BlackRecord:"+blackRecord.RecordId, brBytes)
	if err != nil {
		fmt.Println("PutBlackRecord PutState fail:", err.Error())
		return errors.New("PutBlackRecord PutState fail:" + err.Error())
	}

	return nil
}

func (t *SimpleChaincode) getBlackRecord(stub shim.ChaincodeStubInterface, id string) (*BlackRecord, error) {
	fmt.Println("RecordId:" + id)

	var record BlackRecord
	recordBytes, err := stub.GetState("BlackRecord:" + id)
	if err != nil {
		fmt.Println("getBlackRecord GetState fail:", err.Error())
		return nil, err
	}
	err = json.Unmarshal(recordBytes, &record)
	if err != nil {
		fmt.Println("getBlackRecord Unmarshal fail:", err.Error())
		return nil, err
	}

	fmt.Println(record)
	return &record, nil
}

type Transaction struct {
	DocType    string `json:"docType"`
	TxId       string `json:"txId"`
	From       string `json:"from"`
	To         string `json:"to"`
	Credit     int    `json:"credit"`
	CreateTime string `json:"createTime"`
	UpdateTime string `json:"updateTime"`
}

func (tx *Transaction) putTransaction(stub shim.ChaincodeStubInterface) error {
	txbytes, err := json.Marshal(tx)
	if err != nil {
		fmt.Println("putTransaction Marshal fail:", err.Error())
		return errors.New("putTransaction Marshal fail:" + err.Error())
	}

	err = stub.PutState("Transaction:"+tx.TxId, txbytes)
	if err != nil {
		fmt.Println("putTransaction PutState fail:", err.Error())
		return errors.New("putTransaction PutState fail:" + err.Error())
	}

	return nil
}

func (t *SimpleChaincode) getTransaction(stub shim.ChaincodeStubInterface, id string) (*Transaction, error) {
	fmt.Println("TxId:" + id)

	var tx Transaction
	txBytes, err := stub.GetState("Transaction:" + id)
	if err != nil {
		fmt.Println("getTransaction GetState fail:", err.Error())
		return nil, err
	}
	err = json.Unmarshal(txBytes, &tx)
	if err != nil {
		fmt.Println("getTransaction Unmarshal fail:", err.Error())
		return nil, err
	}

	fmt.Println(tx)
	return &tx, nil
}

type Agency struct {
	Name        string `json:"name"`
	Addr        string `json:"addr"`
	Credit      int    `json:"credit"`
	IssueCredit int    `json:"issueCredit"`
	CreateTime  string `json:"createTime"`
	UpdateTime  string `json:"updateTime"`
}

func (agency *Agency) putAgency(stub shim.ChaincodeStubInterface) error {
	agencyBytes, err := json.Marshal(agency)
	if err != nil {
		fmt.Println("putAgency Marshal fail:", err.Error())
		return errors.New("putAgency Marshal fail:" + err.Error())
	}

	err = stub.PutState("Agency", agencyBytes)
	if err != nil {
		fmt.Println("putAgency PutState fail:", err.Error())
		return errors.New("putAgency PutState fail:" + err.Error())
	}

	return nil
}

func (t *SimpleChaincode) getAgency(stub shim.ChaincodeStubInterface) (*Agency, error) {
	fmt.Println("Get Agency")

	var agency Agency
	agencyBytes, err := stub.GetState("Agency")
	if err != nil {
		fmt.Println("getAgency GetState fail:", err.Error())
		return nil, err
	}
	err = json.Unmarshal(agencyBytes, &agency)
	if err != nil {
		fmt.Println("getAgency Unmarshal fail:", err.Error())
		return nil, err
	}

	fmt.Println(agency)
	return &agency, nil
}

type Org struct {
	OrgId      string `json:"orgId"`
	OrgName    string `json:"orgName"`
	OrgAddr    string `json:"orgAddr"`
	OrgCredit  int    `json:"orgCredit"`
	CreateTime string `json:"createTime"`
	UpdateTime string `json:"updateTime"`
}

func (org *Org) putOrg(stub shim.ChaincodeStubInterface) error {
	orgBytes, err := json.Marshal(org)
	if err != nil {
		fmt.Println("PutOrg Marshal fail:", err.Error())
		return errors.New("PutOrg Marshal fail:" + err.Error())
	}

	err = stub.PutState("Org:"+org.OrgId, orgBytes)
	if err != nil {
		fmt.Println("PutOrg PutState fail:", err.Error())
		return errors.New("PutOrg PutState fail:" + err.Error())
	}

	return nil
}

func (t *SimpleChaincode) getOrg(stub shim.ChaincodeStubInterface, id string) (*Org, error) {
	fmt.Println("OrgId:" + id)

	var org Org
	orgBytes, err := stub.GetState("Org:" + id)
	if err != nil {
		fmt.Println("queryOrg GetState fail:", err.Error())
		return nil, err
	}
	err = json.Unmarshal(orgBytes, &org)
	if err != nil {
		fmt.Println("queryOrg Unmarshal fail:", err.Error())
		return nil, err
	}

	fmt.Println(org)
	return &org, nil
}

func sha1s(s string) string {
	r := sha1.Sum([]byte(s))
	return hex.EncodeToString(r[:])
}

// Create Platform Center Agency , and issue credits.
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	agency := &Agency{
		Name:        "Agency",
		Addr:        sha1s("Agency"),
		Credit:      initialCredit,
		IssueCredit: initialCredit,
		CreateTime:  time.Now().In(loc).Format(layout),
		UpdateTime:  time.Now().In(loc).Format(layout),
	}
	err := agency.putAgency(stub)
	if err != nil {
		return shim.Error("Init putAgency fail:" + err.Error())
	}
	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### BlacklistChain Invoke ###########")
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("function:" + function)
	for _, a := range args {
		fmt.Println("args:" + a)
	}

	if function != "invoke" {
		return shim.Error("Unknown function call:" + function)
	}

	switch args[0] {
	case "createOrg":
		return t.createOrg(stub, args)
	case "submitRecord":
		return t.submitRecord(stub, args)
	case "deleteRecord":
		return t.deleteRecord(stub, args)
	case "queryRecord":
		return t.queryRecord(stub, args)
	case "queryTransaction":
		return t.queryTransaction(stub, args)
	case "issueCredit":
		return t.issueCredit(stub, args[1])
	case "issueCreditToOrg":
		return t.issueCreditToOrg(stub, args)
	case "transfer":
		return t.transfer(stub, args)
	case "queryOrg":
		return t.queryOrg(stub, args[1])
	case "queryAgency":
		return t.queryAgency(stub)

	default:
		return shim.Error("Unknown action, check the first argument:" + args[0])
	}
}

func (t *SimpleChaincode) createOrg(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println(args)
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	org := &Org{
		OrgId:      args[1],
		OrgName:    args[2],
		OrgCredit:  0,
		OrgAddr:    sha1s(args[1]),
		CreateTime: time.Now().In(loc).Format(layout),
		UpdateTime: time.Now().In(loc).Format(layout),
	}
	err := org.putOrg(stub)
	if err != nil {
		fmt.Println("createOrg putOrg fail:", err.Error())
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) submitRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println(args)
	if len(args) != 8 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}

	org := t.getOrg(stub, args[1])
	record := &BlackRecord{
		DocType:          "BlackRecord",
		RecordId:         args[2],
		ClientId:         args[3],
		ClientName:       args[4],
		NegativeType:     args[5],
		NegativeSeverity: args[6],
		NegativeInfo:     args[7],
		OrgAddr:          org.OrgAddr,
		Searchable:       true,
		CreateTime:       time.Now().In(loc).Format(layout),
		UpdateTime:       time.Now().In(loc).Format(layout),
	}
	err := record.putBlackRecord(stub)
	if err != nil {
		fmt.Println("submitRecord putBlackRecord fail:", err.Error())
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) deleteRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println(args)

	var deleteType string

	org := t.getOrg(stub, args[1])
	deleteType = args[2]
	if queryType == "deleteById" {
		if len(args) != 4 {
			return shim.Error("Incorrect number of arguments. Expecting 3")
		}
		arr := strings.Split(args[3], ",")
		var record BlackRecord
		for _, id = range arr {
			record = t.getBlackRecord(stub, id)
			if record.Searchable == false {
				fmt.Println("record-%s does not exist", id)
			} else if record.OrgAddr == org.OrgAddr {
				record.Searchable = false
				record.UpdateTime = time.Now().In(loc).Format(layout)
				record.putBlackRecord(stub)
				fmt.Println("record-%s is deleted ", id)
			} else {
				fmt.Println("record-%s does not belong to org-%s", id, orgId)
			}
		}
	} else if queryType == "deleteAll" {
		//TODO
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) queryRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println(args)
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	var queryType string
	var queryField string
	var queryString string

	queryType = args[1]
	queryField = args[2]
	if queryType == "byClientId" {
		queryString = fmt.Sprintf("{\"selector\":{\"docType\":\"BlackRecord\", \"clientId\":\"%s\", \"searchable\":\"true\"}}", queryField)
	} else if queryType == "byClientName" {
		queryString = fmt.Sprintf("{\"selector\":{\"docType\":\"BlackRecord\", \"clientName\":\"%s\", \"searchable\":\"true\"}}", queryField)
	}
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		fmt.Println("queryRecord getQueryResultForQueryString fail:", err.Error())
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}

func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {
	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		_, value, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString(string(value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

func (t *SimpleChaincode) queryTransaction(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println(args)
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	var queryType string
	var queryField string
	var queryString string

	queryType = args[1]
	queryField = args[2]
	if queryType == "byTxId" {
		queryString = fmt.Sprintf("{\"selector\":{\"docType\":\"Transaction\", \"txId\":\"%s\"}}", queryField)
	} else if queryType == "byToAddr" {
		queryString = fmt.Sprintf("{\"selector\":{\"docType\":\"Transaction\", \"to\":\"%s\"}}", queryField)
	}
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		fmt.Println("queryRecord getQueryResultForQueryString fail:", err.Error())
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}

func (t *SimpleChaincode) issueCredit(stub shim.ChaincodeStubInterface, args string) pb.Response {
	var creditNumber int

	creditNumber, err = strconv.Atoi(args)
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}

	agency, err := t.getAgency(stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	agency.Credit += creditNumber
	agency.IssueCredit += creditNumber
	agency.UpdateTime = time.Now().In(loc).Format(layout)
	err = agency.putAgency(stub)
	if err != nil {
		return shim.Error("write Error: " + err.Error())
	}

	tx := &Transaction{
		DocType:    "Transaction",
		TxId:       time.Now().In(loc).Format(date),
		From:       0,
		To:         agency.Addr,
		Credit:     creditNumber,
		CreateTime: time.Now().In(loc).Format(layout),
		UpdateTime: time.Now().In(loc).Format(layout),
	}
	err = tx.putTransaction(stub)
	if err != nil {
		fmt.Println("issueCredit putTransaction fail:", err.Error())
		return shim.Error("issueCredit putTransaction fail:" + err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) issueCreditToOrg(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println(args)
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	var orgId string
	var creditNumber int

	orgId = args[1]
	creditNumber, err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}

	agency, err := t.getAgency(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	if agency.Credit < creditNumber {
		return shim.Error("Not enough credit for Agency")
	}

	org, err := t.getOrg(stub, orgId)
	if err != nil {
		return shim.Error(err.Error())
	}

	agency.Credit -= creditNumber
	org.OrgCredit += creditNumber

	agency.UpdateTime = time.Now().In(loc).Format(layout)
	err = agency.putAgency(stub)
	if err != nil {
		return shim.Error("write Error: " + err.Error())
	}

	org.UpdateTime = time.Now().In(loc).Format(layout)
	err = org.putOrg(stub)
	if err != nil {
		agency.UpdateTime = time.Now().In(loc).Format(layout)
		agency.Credit += creditNumber
		err = agency.putAgency(stub)
		if err != nil {
			return shim.Error("roll down error")
		}
		return shim.Error("write Error: " + err.Error())
	}

	tx := &Transaction{
		DocType:    "Transaction",
		TxId:       time.Now().In(loc).Format(date),
		From:       agency.Addr,
		To:         org.OrgAddr,
		Credit:     creditNumber,
		CreateTime: time.Now().In(loc).Format(layout),
		UpdateTime: time.Now().In(loc).Format(layout),
	}
	err = tx.putTransaction(stub)
	if err != nil {
		fmt.Println("transfer putTransaction fail:", err.Error())
		return shim.Error("transfer putTransaction fail:" + err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println(args)
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	var fromId string
	var toId string
	var creditNumber int

	fromId = args[1]
	toId = args[2]
	creditNumber, err = strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}

	fromOrg, err := t.getOrg(stub, fromId)
	if err != nil {
		return shim.Error(err.Error())
	}
	if fromOrg.OrgCredit < creditNumber {
		return shim.Error("Not enough credit for Org " + fromOrg.OrgId)
	}

	toOrg, err := t.getOrg(stub, toId)
	if err != nil {
		return shim.Error(err.Error())
	}

	fromOrg.OrgCredit -= creditNumber
	toOrg.OrgCredit += creditNumber

	fromOrg.UpdateTime = time.Now().In(loc).Format(layout)
	err = fromOrg.putOrg(stub)
	if err != nil {
		return shim.Error("write Error: " + err.Error())
	}

	toOrg.UpdateTime = time.Now().In(loc).Format(layout)
	err = toOrg.putOrg(stub)
	if err != nil {
		fromOrg.UpdateTime = time.Now().In(loc).Format(layout)
		fromOrg.OrgCredit += creditNumber
		err = fromOrg.putOrg(stub)
		if err != nil {
			return shim.Error("roll down error")
		}
		return shim.Error("write Error: " + err.Error())
	}

	tx := &Transaction{
		DocType:    "Transaction",
		TxId:       time.Now().In(loc).Format(date),
		From:       fromOrg.OrgAddr,
		To:         toOrg.OrgAddr,
		Credit:     creditNumber,
		CreateTime: time.Now().In(loc).Format(layout),
		UpdateTime: time.Now().In(loc).Format(layout),
	}
	err = tx.putTransaction(stub)
	if err != nil {
		fmt.Println("transfer putTransaction fail:", err.Error())
		return shim.Error("transfer putTransaction fail:" + err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) queryAgency(stub shim.ChaincodeStubInterface) pb.Response {
	agency, err := t.getAgency(stub)
	if err != nil {
		fmt.Println("queryAgency getAgency fail:", err.Error())
		return shim.Error(err.Error())
	}
	agencyBytes, err = json.Marshal(agency)
	if err != nil {
		fmt.Println("queryAgency marshal fail:", err.Error())
		return shim.Error(err.Error())
	}

	return shim.Success(agencyBytes)
}

func (t *SimpleChaincode) queryOrg(stub shim.ChaincodeStubInterface, args string) pb.Response {
	org, err := t.getOrg(stub, args)
	if err != nil {
		fmt.Println("queryOrg getOrg fail:", err.Error())
		return shim.Error(err.Error())
	}
	orgBytes, err = json.Marshal(org)
	if err != nil {
		fmt.Println("queryOrg marshal fail:", err.Error())
		return shim.Error(err.Error())
	}

	return shim.Success(orgBytes)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
