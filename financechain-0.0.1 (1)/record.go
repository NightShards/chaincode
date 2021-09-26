package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

const (
	key_type_simple = "simple"
	key_type_composite ="composite"
)

/*
arg1: object name of record
arg2: key type, simple or composite.
 */
func setKeyType(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error(errorIncorrectParamNum + " Expecting 2")
	}

	fmt.Println("args=", args)
	objectName := args[0]
	if objectName == "" {
		return shim.Error(errorMissParam + " Expect object name")
	}

	keyType := args[1]
	if keyType == "" {
		return shim.Error(errorMissParam + "Expect key type")
	}

	//save the key type
	existType, err := stub.GetState(objectName + "_keytype")
	if err != nil {
		return shim.Error(errorBlockchainError + err.Error())
	}
	if existType == nil {
		err = stub.PutState(objectName + "_keytype", []byte(keyType))
		if err != nil {
			return shim.Error(errorBlockchainError + err.Error())
		}
	} else {
		return shim.Error(errorDataAlreadyExist + err.Error())
	}

	return shim.Success(nil)
}

/*
arg1: object name of record
return: key type, simple or composite.
 */
func getKeyType(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error(errorIncorrectParamNum + " Expecting 1")
	}

	fmt.Println("args=", args)
	objectName := args[0]
	if objectName == "" {
		return shim.Error(errorMissParam + " Expect object name")
	}

	existType, err := stub.GetState(objectName + "_keytype")
	if err != nil {
		return shim.Error(errorBlockchainError + err.Error())
	}
	if existType == nil {
		return shim.Error(errorDataNotExist + err.Error())
	} else {
		return shim.Success([]byte(existType))
	}
}

/*
arg1: object name of record
arg2: list of keywords
arg3: record content that is json format string
 */
func saveRecord(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3 {
		return shim.Error(errorIncorrectParamNum + " Expecting 3")
	}

	fmt.Println("args=", args)
	objectName := args[0]
	if objectName == "" {
		return shim.Error(errorMissParam + " Expect object name")
	}

	keywords := args[1]
	if keywords == "" {
		return shim.Error(errorMissParam + "Expect keywords list")
	}

	//get the key type
	keyType, err := stub.GetState(objectName + "_keytype")
	if err != nil {
		return shim.Error(errorBlockchainError + err.Error())
	}
	if keyType == nil {
		return shim.Error(errorDataNotExist + "Must set key type first")
	}

	var keywordsList []string
	err = json.Unmarshal([]byte(keywords), &keywordsList)
	if err != nil {
		return shim.Error(errorInvalidParam + "Keywords list")
	}

	key, err := createKey(stub, objectName, keywordsList, string(keyType))
	if err != nil {
		return shim.Error(errorBlockchainError + err.Error())
	}

	//save object content
	err = stub.PutState(key, []byte(args[2]))
	if err != nil {
		return shim.Error(errorBlockchainError + err.Error())
	}
	return shim.Success([]byte(key))
}

/*
arg1: object name of record
arg2: list of keywords
 */
func queryRecord(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error(errorIncorrectParamNum + " Expecting 2")
	}

	fmt.Println("args=", args)
	objectName := args[0]
	if objectName == "" {
		return shim.Error(errorMissParam + " Expect object name")
	}

	keywords := args[1]
	if keywords == "" {
		return shim.Error(errorMissParam + "Expect keywords list")
	}

	var keywordsList []string
	err := json.Unmarshal([]byte(keywords), &keywordsList)
	if err != nil {
		return shim.Error(errorInvalidParam + "Keywords list")
	}

	keyType, err := stub.GetState(objectName + "_keytype")
	if err != nil {
		return shim.Error(errorBlockchainError + err.Error())
	}
	if keyType == nil {
		return shim.Error(errorBlockchainError + err.Error())
	}

	key, err := createKey(stub, objectName, keywordsList, string(keyType))
	if err != nil {
		return shim.Error(errorBlockchainError + err.Error())
	}

	ret, err := stub.GetState(key)
	if err != nil {
		return shim.Error(errorBlockchainError + err.Error())
	} else if ret == nil {
		return shim.Error(errorDataNotExist)
	} else {
		return shim.Success(ret)
	}
}

/*
arg1: object name of record
arg2: list of keywords
 */
func queryRecordByPartial(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error(errorIncorrectParamNum + " Expecting 2")
	}

	fmt.Println("args=", args)
	objectName := args[0]
	if objectName == "" {
		return shim.Error(errorMissParam + " Expect object name")
	}

	keywords := args[1]
	if keywords == "" {
		return shim.Error(errorMissParam + "Expect keywords list")
	}

	var keywordsList []string
	err := json.Unmarshal([]byte(keywords), &keywordsList)
	if err != nil {
		return shim.Error(errorInvalidParam + "Keywords list")
	}
	fmt.Println(keywordsList)

	keyType, err := stub.GetState(objectName + "_keytype")
	if err != nil {
		return shim.Error(errorBlockchainError + err.Error())
	}
	if keyType == nil {
		return shim.Error(errorBlockchainError + err.Error())
	}

	if string(keyType) != key_type_composite {
		return shim.Error(errorInvalidFunction + "This object don't support query by partial")
	}

	resultSet := ""
	iterator, err := stub.GetStateByPartialCompositeKey(objectName, keywordsList)
	defer iterator.Close()
	if err != nil {
		return shim.Error(errorBlockchainError + err.Error())
	} else if iterator == nil {
		return shim.Error(errorDataNotExist)
	} else {
		for iterator.HasNext(){
			item, err := iterator.Next()
			if err != nil {
				return shim.Error(errorBlockchainError + err.Error())
			}

			fmt.Println(item.Key)
			ret, err := stub.GetState(item.Key)
			if err != nil {
				return shim.Error(errorBlockchainError + err.Error())
			}
			resultSet += string(ret) + ","
		}

		if len(resultSet) != 0 {
			resultSet = "[" + resultSet[:len(resultSet)-1] + "]"
		} else {
			resultSet = "[]"
		}
		return shim.Success([]byte(resultSet))
	}
}

/*
arg1: object name of record
arg2: list of keywords
 */
func deleteRecord(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error(errorIncorrectParamNum + " Expecting 2")
	}

	fmt.Println("args=", args)
	objectName := args[0]
	if objectName == "" {
		return shim.Error(errorMissParam + " Expect object name")
	}

	keywords := args[1]
	if keywords == "" {
		return shim.Error(errorMissParam + "Expect keywords list")
	}

	var keywordsList []string
	err := json.Unmarshal([]byte(keywords), &keywordsList)
	if err != nil {
		return shim.Error(errorInvalidParam + "Keywords list")
	}

	keyType, err := stub.GetState(objectName + "_keytype")
	if err != nil {
		return shim.Error(errorBlockchainError + err.Error())
	}
	if keyType == nil {
		return shim.Error(errorBlockchainError + err.Error())
	}

	key, err := createKey(stub, objectName, keywordsList, string(keyType))
	if err != nil {
		return shim.Error(errorBlockchainError + err.Error())
	}

	err = stub.DelState(key)
	if err != nil {
		return shim.Error(errorBlockchainError + err.Error())
	} else {
		return shim.Success(nil)
	}
}

func createKey(stub shim.ChaincodeStubInterface, objectName string, keywordsList []string, keyType string) (string, error) {
	if keyType == key_type_simple {
		key := objectName
		for _,element := range keywordsList{
			key += "_" + element
		}
		return key, nil
	} else {
		key, err := stub.CreateCompositeKey(objectName, keywordsList)
		if err != nil {
			return "", err
		} else {
			return key, nil
		}
	}
}
