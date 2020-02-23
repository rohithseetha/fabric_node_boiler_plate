// peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n kyc -c '{"Args":["signup","tri","cha","ch@","a1","a2","hyd","ap","india","500032","23-2-1990","999","cn"]}'

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------
//   Sl.No           Name                   Date                    Description
//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------
//    01          Rohith Seetha            23-07-2019               Initial Version
//-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type user struct {
	ObjectType  string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	UUID        string `json:"UUID"`    //SN-04
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	UserName    string `json:"userName"`
	AddressOne  string `json:"addressOne"`
	AddressTwo  string `json:"addressTwo"`
	City        string `json:"city"`
	State       string `json:"state"`
	Country     string `json:"country"`
	PinCode     string `json:"pincode"`
	DOB         string `json:"dob"`
	Contact     string `json:"contact"`
	CompanyName string `json:"companyName"`
}

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init initializes chaincode
// ===========================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("entering invoice init successfully")
	return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
// ========================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	//fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "signup" { //create a new User
		return t.signup(stub, args)
		// } else if function == "initUser" { //read a User
		// 	return t.siginitUsernup(stub, args)
	} else if function == "readUser" { //read a User
		return t.readUser(stub, args)
		// } else if function == "readallUsers" { //read a User
		// 	return t.readallUsers(stub)
	} else if function == "deleteUser" {
		return t.deleteUser(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// ============================================================
// initUser - create a new User, store into chaincode state
// ============================================================
func (t *SimpleChaincode) signup(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// if len(args) != 28 {
	// 	return shim.Error("Incorrect number of arguments. Expecting 28")
	// }

	fmt.Println("- start Signup user")

	FirstName := args[0]
	LastName := args[1]
	UserName := args[2]
	AddressOne := args[3]
	AddressTwo := args[4]
	City := args[5]
	State := args[6]
	Country := args[7]
	PinCode := args[8]
	DOB := args[9]
	Contact := args[10]
	CompanyName := args[11]

	objectType := "user"
	user := user{objectType, UserName, FirstName, LastName, UserName, AddressOne, AddressTwo, City, State,
		Country, PinCode, DOB, Contact, CompanyName}
	userJSONasBytes, err3 := json.Marshal(user)
	if err3 != nil {
		return shim.Error(err3.Error())
	}

	err4 := stub.PutState(UserName, userJSONasBytes)
	if err4 != nil {
		return shim.Error(err4.Error())
	}

	fmt.Println("- end init invoice")
	return shim.Success([]byte("Success"))
}

// ===============================================
// readInvoice - read an user from chaincode state
// ===============================================
func (t *SimpleChaincode) readUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var userName, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting username to query")
	}

	userName = args[0]
	valAsbytes, err := stub.GetState(userName) //get the user from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + userName + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Invoice does not exist: " + userName + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success([]byte(valAsbytes))
}

// ================================================================
//          Reading all the User data from the current state ledger
// ================================================================
// func (t *SimpleChaincode) readallUsers(stub shim.ChaincodeStubInterface) pb.Response {

// 	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"user\"}}")
// 	queryResults, err := getQueryResultForQueryString(stub, queryString)
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}
// 	return shim.Success(queryResults)
// }

// =========================================================================================
// getQueryResultForQueryString executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.
// =========================================================================================
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) (string, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)
	var count int

	resultsIterator, err1 := stub.GetQueryResult(queryString)
	if err1 != nil {
		return "", err1
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		count = count + 1
		queryResponse, err2 := resultsIterator.Next()
		if err2 != nil {
			return "", err2
		}
		// Add a comma before array members, suppress it for the first array member
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
	totalCount := strconv.Itoa(count)
	// countAsBytes, err3 := json.Marshal(totalCount)

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", totalCount)
	return (totalCount), nil
}

// ===============================================
// deleteInvoice - delete an invoice request from chaincode state
// ===============================================
func (t *SimpleChaincode) deleteUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var userName string

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting Username to query")
	}

	userName = args[0]
	err := stub.DelState(userName) //delete the user from chaincode state
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to delete user req: %s", args[0]))
	}
	return shim.Success(nil)

}
