package main
import (
"fmt"
"github.com/hyperledger/fabric-chaincode-go/shim"
"github.com/hyperledger/fabric/core/peer"
"encoding/json"
"io/ioutil"
)

type TweetChaincode struct {

}

type Tweet struct{

	Created_at Time `json:"created_at"`
	Id int64 'json:"id"'
	Id_str string 'json:"id_str"'
	Text string `json:"text"`
	Source string 'json:"source"'
	Truncated bool 'json:"truncated"'
	In_reply_to_status_id int64 'json"in_reply_to_status_id"'
	In_reply_to_status_id_str string 'json"in_reply_to_status_id_str"'
	In_reply_to_user_id int64 'json:"in_reply_to_user_id"'
	In_reply_to_user_id_str string 'json:"in_reply_to_user_id_str"'
	In_reply_to_screen_name string 'json:"in_reply_to_screen_name"'
	User *User 'json:"User, omitempty'
	Place string `json:"place"`
	Coordinates string `json:"coordinates"`
}
User type User struct{
	Username string 'json:"username"'
	Userid string 'json:"userid"'
	Userlocation string `json:"userlocation"`
}
type CompositeFile struct{
	Filename string 'json:"filename"'
	IncludedFiles string[] 'json:"includedfiles"'
}

//initialize the chaincode
//only call when the chaincode is instantiated
func (t *TweetChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("##### Init Chaincode #####")
	//get args from transaction proposal
	args := stub.GetStringArgs()

	filesSlice := []

	if len(args) != 1 {

		return shim.Error("Incorrect arguments. Expecting output file name.")
	}

	//Store the outfile in the ledger

	err := stub.PutState(args[0]))

	if err != nil {
		return shim.Error(ftm.Sprintf("Failed to create asset: %s", args[0]))
	}
	return shim.Success(nil)
}

func (t *TweetChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.println("##### Invoke Chaincode #####")
	//extract function and args from transaction proposal

	fn, args := stub.GetFunctionAndParameters()

	var result string

	var err error

	if fn  == "filter"{

		result, err = filter(stub, args)

	} else { //assume query
		result, err = query(stub, args)
	}
}

if err != nil { 
	//failed to get function and/or arguments from transaction proposal
	return shim.Error(err.Error())
}

//Return result as success

return shim.Success([]byte(result))
}

func filter(stub shim.ChaincodeStubInterface, args []string)(string, error) {

	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect arguments. Expecting an input and output file name")

	}
	string(inputName) := args[1]
	string(outputName) := args[0]
	filesSlice.append(inputName)

	err := stub.PutState(args[0], filesSlice)

	if err != nil{

		return "", fmt Errorf ("Failed to set asset: %s", args[0])

	}
	
	jsonFile, err := os.Open(inputName)

	byteValue, _ := ioutill.ReadAll(jsonFile)
	var tweets Tweet 
	json.Unmarshal(byteValue, &tweets)
	for i := 0; i < len(tweets.Tweet); i++{
		if tweets.Tweet[13] != null{
			
			writeLine(tweets.Tweet, outFile)

		}

	}
	if err != nil{

		fmt.Println(err)
	}
	defer jsonFile.Close()
	return args[1], nil

}

func query(stub shim.ChaincodeStubInterface, args []string)(string, error) {

	if len(args) != 1 {

		return "", ftm.Errorf("Incorrect arguments. Expecting an output file name")

	}
	value, err := stub.GetState(args[0])
	
		if err != nil{

			return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
		}
		if value == nil{

			return "", fmt.Errorf("Asset not found: %s", args[0])
		}

	return string(value), nil

}

func main(){

	err := shim.Start(new(SampleChaincode))
	if err != nil {
		fmt.Println("could not start TweetChaincode")

	}
	else{

		fmt.Println("TweetChaincode successfully started.")
	}
}
// writeLines writes the lines to the given file.
func writeLine(line string, path string) error {
	file, err := os.Create(path)
	if err != nil {
	  return err
	}
  