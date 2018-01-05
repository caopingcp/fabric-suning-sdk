/*
Copyright 33.cn Corp. 2018 All Rights Reserved.

Chaincode for Suning Corp.
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var (
	layout = "2006-01-02 15:04:05"
	date   = "20060102150405"
	loc    *time.Location
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
		return errors.New("PutBlackRecord PutState Error" + err.Error())
	}

	return nil
}

type Transaction struct {
	DocType    string `json:"docType"`
	TxId       string `json:"txId"`
	From       string `json:"from"`
	To         string `json:"to"`
	Credit     int64  `json:"credit"`
	CreateTime string `json:"createTime"`
	UpdateTime string `json:"updateTime"`
}

func (tx *Transaction) putTransaction(stub shim.ChaincodeStubInterface) error {
	txbytes, err := json.Marshal(tx)
	if err != nil {
		fmt.Println("PutTransaction Marshal fail:", err.Error())
		return errors.New("PutTransaction Marshal fail:" + err.Error())
	}

	err = stub.PutState("Transaction:"+tx.TxId, txbytes)
	if err != nil {
		fmt.Println("PutTransaction PutState fail:", err.Error())
		return errors.New("PutTransaction PutState Error" + err.Error())
	}

	return nil
}

type Agency struct {
	Name        string `json:"name"`
	IssueCredit int    `json:"issueCredit"`
	Credit      int    `json:"credit"`
	CreateTime  string `json:"createTime"`
	UpdateTime  string `json:"updateTime"`
}

func (agency *Agency) putAgency(stub shim.ChaincodeStubInterface) error {
	agencyBytes, err := json.Marshal(agency)
	if err != nil {
		fmt.Println("PutAgency Marshal fail:", err.Error())
		return errors.New("PutAgency Marshal fail:" + err.Error())
	}

	err = stub.PutState("Agency", agencyBytes)
	if err != nil {
		fmt.Println("PutAgency PutState fail:", err.Error())
		return errors.New("PutAgency PutState Error" + err.Error())
	}

	return nil
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
		return errors.New("PutOrg PutState Error" + err.Error())
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

// Create Platform Center Agency , and issue credits.
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	const initialCredit = 1e8
	agency := &Agency{
		Name:        "Agency",
		IssueCredit: initialCredit,
		Credit:      initialCredit,
		CreateTime:  time.Now().In(loc).Format(layout),
		UpdateTime:  time.Now().In(loc).Format(layout),
	}
	err := agency.putAgency(stub)
	if err != nil {
		return shim.Error("Init PutAgency fail:" + err.Error())
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
		return t.queryRecord(stub, args[1])
	case "queryOrg":
		return t.queryOrg(stub, args[1])
	case "queryAgency":
		return t.queryAgency(stub, args[1])
	case "issueCoin":
		return t.issueCoin(stub, args[1])
	case "transfer":
		return t.transfer(stub, args)
	case "queryTransaction":
		return t.queryTransaction(stub, args[1])

	default:
		return shim.Error("Unknown action, check the first argument:" + args[0])
	}
}

func (t *SimpleChaincode) createOrg(stub shim.ChaincodeStubInterface, args string) pb.Response {
	fmt.Println(args)

	org := &Org{
		OrgId:      args[1],
		OrgName:    args[2],
		OrgCredit:  0,
		OrgAddr:    hash(args[1]),
		CreateTime: time.Now().In(loc).Format(layout),
		UpdateTime: time.Now().In(loc).Format(layout),
	}
	err := org.putOrg(stub)
	if err != nil {
		fmt.Println("createOrg PutOrg fail:", err.Error())
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) submitRecord(stub shim.ChaincodeStubInterface, orgId string, args string) pb.Response {
	fmt.Println("Org:" + orgId)
	fmt.Println(args)

	org := t.getOrg(orgId)

	record := &BlackRecord{
		DocType:          "BlackRecord",
		RecordId:         args[1],
		ClientId:         args[2],
		ClientName:       args[3],
		NegativeType:     args[4],
		NegativeSeverity: args[5],
		NegativeInfo:     args[6],
		OrgAddr:          org.OrgAddr,
		CreateTime:       time.Now().In(loc).Format(layout),
		UpdateTime:       time.Now().In(loc).Format(layout),
	}
	err := record.putBlackRecord(stub)
	if err != nil {
		fmt.Println("submitRecord PutBlackRecord fail:", err.Error())
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) queryRecord(stub shim.ChaincodeStubInterface, function string, args string) pb.Response {
	fmt.Println("query is running " + function)

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	var queryString string
	if function == "byClientId" {
		queryString = fmt.Sprintf("{\"selector\":{\"docType\":\"BlackRecord\", \"clientId\":\"%s\"}}", args)
	} else if function == "byClientName" {
		queryString = fmt.Sprintf("{\"selector\":{\"docType\":\"BlackRecord\", \"clientName\":\"%s\"}}", args)
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

func (t *SimpleChaincode) queryOrg(stub shim.ChaincodeStubInterface, args string) pb.Response {
	fmt.Println(args)

	var org Org
	orgBytes, err := stub.GetState("Org:" + args)
	if err != nil {
		fmt.Println("queryOrg GetState fail:", err.Error())
	}
	err = json.Unmarshal(orgBytes, &org)
	if err != nil {
		fmt.Println("queryOrg Unmarshal fail:", err.Error())
	}

	fmt.Println(org)
	return shim.Success(orgBytes)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
