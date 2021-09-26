package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// SmartContract Define the Smart Contract structure
type SmartContract struct {
}

// Init function
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("Init")
	return shim.Success(nil)
}

// Invoke function
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) peer.Response {
	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()

	// Route to the appropriate handler function to interact with the ledger appropriately
	fmt.Println(function)
	if function == "saveRecord" {
		return saveRecord(APIstub, args)
	} else if function == "queryRecord" {
		return queryRecord(APIstub, args)
	} else if function == "queryRecordByPartial" {
		return queryRecordByPartial(APIstub, args)
	} else if function == "deleteRecord" {
		return deleteRecord(APIstub, args)
	} else if function == "setKeyType" {
		return setKeyType(APIstub, args)
	} else if function == "getKeyType" {
		return getKeyType(APIstub, args)
	}

	fmt.Println("Invalid function name")
	return shim.Error(errorInvalidFunction)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}