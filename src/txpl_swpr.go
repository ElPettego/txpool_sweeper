package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	L "txpool_sweeper/src/lib"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// flags
var network string
var init_db string

func init() {
	flag.StringVar(&network, "network", "", "EVM network to interact with. must be present in .json.")
	flag.StringVar(&init_db, "init_db", "", "if it is set to yes reinits the db")
}

type Record struct {
	test string
	age  int
}

func main() {
	// exp := new(big.Int)
	// exp.Exp(big.NewInt(10), big.NewInt(18), nil)

	// fmt.Println(0xa + 0xa + 0xa001)

	// os.Exit(10)

	db, err := L.ConnectDB("data/db.db")
	defer db.Close()

	// privateKeyFile, err := exec.Command("cat", "../../../.private_keys/crypto/__swg__").Output()
	// if err != nil {
	// 	log.Fatal("error reading private key", err)
	// }
	// privateKey := strings.Split(string(privateKeyFile), "\n")[0]
	// fmt.Println(privateKey)

	// os.Exit(10)

	// w3, err := L.ConnectWeb3("https://bsc.publicnode.com", privateKey)
	// defer w3.Close()

	// fmt.Println(w3.Address)

	// nonce, err := w3.GetNonce(w3.Address.String())

	// fmt.Println(nonce)
	// os.Exit(10)

	flag.Parse()

	if init_db == "yes" {
		db.CreateTable("test", "test STRING PRIMARY KEY, age INTEGER")
		test := make(map[string]interface{})
		test["test"] = "sliem"
		test["age"] = 1001
		db.InsertRecordIntoTable("test", test)
		rows, err := db.SelectFromTable("test", "*")
		if err != nil {
			log.Fatal("failed to retrieve rows")
		}
		var records []Record
		for rows.Next() {
			var record Record
			err := rows.Scan(&record.test, &record.age)
			if err != nil {
				log.Fatal("error selecting records")
			}
			records = append(records, record)

		}

		for i, rec := range records {
			fmt.Printf("row -> %d: test %s, age %d\n", i, rec.test, rec.age)
		}

		os.Exit(10)
	}
	if network == "" {
		panic("must provide network")
	}

	jsonFile, err := os.ReadFile(".json")
	if err != nil {
		panic(err)
	}
	var config map[string]interface{}
	err = json.Unmarshal(jsonFile, &config)
	if err != nil {
		panic(err)
	}
	parsedConfig := config[network].(map[string]interface{})

	var provider = parsedConfig["provider"].(string)
	var explorer = parsedConfig["explorer"].(string)
	var mode = parsedConfig["mode"].(string)
	var id = int(parsedConfig["id"].(float64))
	var addresses = parsedConfig["addresses"].([]interface{})
	var list_addresses = L.ConvertToStringList(addresses)

	router_abi, err := os.ReadFile("build/solidity/IPancakeRouter02.abi")
	if err != nil {
		log.Fatal((err))
	}

	var parsed_abi abi.ABI

	// _err :=
	json.Unmarshal(router_abi, &parsed_abi)

	// fmt.Println(parsed_abi)

	client := &http.Client{}

	var jsonRes map[string]interface{}

	fmt.Println()

	for {
		res := L.BasePost(provider, L.GetPayload(id, `"eth_newPendingTransactionFilter"`, ""), *client)
		err = json.Unmarshal([]byte(res), &jsonRes)
		if err != nil {
			continue
		}

		res = L.BasePost(provider, L.GetPayload(id, `"eth_getFilterChanges"`, fmt.Sprintf(`"%s"`, jsonRes["result"].(string))), *client)
		err = json.Unmarshal([]byte(res), &jsonRes)
		if err != nil {
			continue
		}

		hashes, succ := jsonRes["result"].([]interface{})
		if !succ {
			continue
		}
		for _, hash := range hashes {

			if mode == "hash" {
				go L.HandleHash(explorer, provider, hash.(string), id, list_addresses, *client)
			}
			if mode == "dict" {
				_to, succ := hash.(map[string]interface{})["to"].(string)
				if succ {
					if L.Contains(list_addresses, _to) {
						// convert input to bytes in order to decode method of transacrion on router
						str_input := L.FromInterfaceToString(hash.(map[string]interface{})["input"])
						str_value := L.FromInterfaceToString(hash.(map[string]interface{})["value"])
						str_gasPrice := L.FromInterfaceToString(hash.(map[string]interface{})["gasPrice"])
						inputBytes, err := hex.DecodeString(str_input[2:]) // Remove the "0x" prefix before decoding
						if err != nil {
							fmt.Println("Failed to decode string:", err)
						}
						signature, data := inputBytes[:4], inputBytes[4:] // A byte can represent two hexadecimal characters
						method, err := parsed_abi.MethodById(signature)
						if err != nil {
							log.Fatal(err)
						}
						var args = make(map[string]interface{})
						err = method.Inputs.UnpackIntoMap(args, data)
						if err != nil {
							log.Fatal(err)
						}
						_method := fmt.Sprintf("%s", method)

						// check that the method are the ones needed
						if (strings.Contains(_method, "swapExactETH") || strings.Contains(_method, "swapETH")) && !strings.Contains(_method, "FeeOn") {
							fmt.Printf("%s <-> %s/tx/%s\n", L.GetNow(), explorer, hash.(map[string]interface{})["hash"])
							fmt.Printf("Method: %s\n", method)
							for key, value := range args {
								fmt.Printf("%s: %+v\n", key, value)
							}
							fmt.Println("value ->", str_value)
							valueB, succ := L.FromStringToInt(str_value[2:]) // new(big.Int).SetString(str_value[2:], 16)
							if !succ {
								log.Fatal("failed to convert value")
							}
							fmt.Println("value ->", valueB)
							gasPrice, succ := L.FromStringToInt(str_gasPrice[2:])
							fmt.Println("gasPrice ->", str_gasPrice)
							if succ {
								fmt.Println("gasPrice ->", gasPrice)

							}
							fmt.Println()
							_path := args["path"].([]common.Address)
							fmt.Println(_path[1])

							// _test := "0x60045b3806c973f3aad8497d98f01e582ccb15b1"

							res, err := db.SelectFromTable("tokens", "address", fmt.Sprintf("address = '%s'", _path[1].String()))
							if err != nil {
								log.Fatal("error quering db")
							}
							inDb := res.Next()

							fmt.Println(inDb)
							fmt.Println(L.GetNow())
							fmt.Println()

							if inDb {
								os.Exit(10)
							}

						}
					}
				}
			}
		}
	}
}
