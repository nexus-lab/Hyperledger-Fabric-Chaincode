package main

import (
	"bytes"
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

// SimpleAsset implements a simple chaincode to manage an asset
type SimpleAsset struct {

}

type File struct{
	FileURL url `json:"file-url"`
}
//Able to query, stores all rows calculated in a given file
type Transactions struct{
	RowsCalled []int `json:"rows-called"`
}

//Structure to store transactions (calculations) for the ledger
type Transaction struct{
	RowNum int `json:"row-num"`
	Sum float32	`json:"sum"`
}

//type for saving the URL
type url = string
//index of where to inquire
type inquiryNum = int

//retrieves, parses CSV data from a given URL
func retrieve(fileURL url, s string) ([][]string, error) {
	resp, err := http.Get(fileURL)
	if err != nil {
		log.Fatal(err)

	}

	r := csv.NewReader(resp.Body)
	r.Comma = ','

	rows, err := r.ReadAll()
	if err != nil {
		log.Fatal("Cannot read CSV data: ", err.Error())
	}
	return rows, nil
}


//calculates the functions
//add all numerical values in a user-requested row
func process(rows [][]string, rowNum int) float32 {
	//4 numerical rows are sepal and petal lengths and widths
	//must parse string into float
	sepalLength, err := strconv.ParseFloat(rows[rowNum][0], 32); if err != nil{log.Fatalf(err.Error())}
	sepalWidth, err := strconv.ParseFloat(rows[rowNum][1], 32); if err != nil{log.Fatalf(err.Error())}
	petalLength, err := strconv.ParseFloat(rows[rowNum][2], 32); if err != nil{log.Fatalf(err.Error())}
	petalWidth, err := strconv.ParseFloat(rows[rowNum][3], 32); if err != nil{log.Fatalf(err.Error())}
	var sum float32 = 0
	//add the previous values into a float
	sum = float32(sepalLength + sepalWidth + petalLength + petalWidth)
	return sum
}


// Init is called during chaincode instantiation
func (t *SimpleAsset) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("Init")
	_, args := stub.GetFunctionAndParameters()
	var InitURL string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting a data URL.")
	}

	KeyA = args[0]
	// Initialize the chaincode
	InitURL = args[1]
	if err != nil {
		return shim.Error("Expecting string value for URL")
	}
	fmt.Printf("Successfully accessed: ", InitURL)

	// Write the state to the ledger
	err = stub.PutState([]byte(InitURL))
	if err != nil {
		return shim.Error(err.Error())
	}


	return shim.Success(nil)
}


// {invoke{set{RowNum}}} -> commits rowNum + sum to ledger
// Invoke is called per transaction on the chaincode. Each transaction is
// either a 'get' or a 'set' on the asset created by Init function. The Set
// method may create a new asset by specifying a new key-value pair.
// Get takes in rowNum and prints the sum; Set takes in rowNum, saves rowNum and sum to ledger
func (t *SimpleAsset) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()

	var result string
	var err error
	if fn == "set" {
		result, err = set(stub, args)
	} else { // assume 'get' even if fn is nil
		result, err = get(stub, args)
	}
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return the result as success payload
	return shim.Success([]byte(result))
}


// Set stores the asset (both key and value) on the ledger. If the key exists,
// it will override the value with the new one
func set(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a row index")
	}


	//call the data processing functions
	//retrieve and calculate here, print calculated value
	x, err := retrieve(args[0], args[1])
	if err != nil{
		log.Fatal(err.Error())
	}
	var num int
	num, _ = strconv.Atoi(args[1])
	sum := process(x, num)

	//retrieve rows from CSV URL saved in the ledger init, then call process on the rowNum
	err = stub.PutState(args[0], []byte(float32ToByte(sum)))

	if err != nil {
		return "", fmt.Errorf("Failed to calculate and set row: %s", args[0])
	}

	return args[1], nil
}

func float32ToByte(f float32) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, f)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}

// Get returns the value of the specified asset key
// returns sum of row based on row index
func get(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a row index")
	}

	value, err := stub.GetState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
	}
	if value == nil {
		return "", fmt.Errorf("Asset not found: %s", args[0])
	}

	return string(value), nil
}


// main function starts up the chaincode in the container during instantiate
func main() {
	err := shim.Start(new(SimpleAsset))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}