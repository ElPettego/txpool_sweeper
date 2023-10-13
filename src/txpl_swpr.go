package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	L "txpool_sweeper/src/lib"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

var network string
var init_db string

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
)

func handle_hash(explorer, provider, hash string, id int, addresses []string, client http.Client) {
	var jsonRes map[string]interface{}
	res := base_post(provider, get_payload(id, `"eth_getTransactionByHash"`, fmt.Sprintf(`"%s"`, hash)), client)
	err := json.Unmarshal([]byte(res), &jsonRes)
	if err != nil {
		return
		// panic(err)
	}
	if result, ok := jsonRes["result"].(map[string]interface{}); ok {
		to, succ := result["to"].(string)
		if succ {
			if Contains(addresses, to) {
				fmt.Println(fmt.Sprintf("%s <-> %s %s/tx/%s", get_now(), to, explorer, result["hash"]))
			}
		}
	}
}

func Contains(list []string, x string) bool {
	for _, item := range list {
		if item == x {
			return true
		}
	}

	return false
}

func convertToStringList(input []interface{}) []string {
	stringList := make([]string, 0, len(input))
	for _, val := range input {
		if strVal, ok := val.(string); ok {
			stringList = append(stringList, strings.ToLower(strVal))
		} else {
			// Handle the case if the interface value is not a string (optional)
			// You can choose to ignore, skip, or handle this differently based on your requirements.
			fmt.Printf("Warning: Value %v is not a string\n", val)
		}
	}
	return stringList
}

func get_now() string {
	return "[" + Purple + time.Now().UTC().Format("2006-01-02|15:04:05.000000") + Reset + "]"
}

func get_payload(id int, method, params string) string {
	return fmt.Sprintf(`{"id": %d,"jsonrpc":"2.0","method":%s,"params":[%s]}`, id, method, params)
}

func base_post(url, payload string, client http.Client) string {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return ""
	}
	res, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer res.Body.Close()
	var resPayload bytes.Buffer
	_, err = resPayload.ReadFrom(res.Body)
	if err != nil {
		return ""
	}
	return resPayload.String()
}

func from_interface_to_string(_interface interface{}) string {
	if _, ok := _interface.(string); ok {
		return _interface.(string)
	}
	return "error"
}

func from_string_to_int(str string) (*big.Int, bool) {
	return new(big.Int).SetString(str, 16)
}

func to_gwei(number *big.Int) {
	// gwei := new(big.Int).Div(number, big.NewInt(1e18))
	_gwei := int(number.Int64())
	div := _gwei / 1e18
	fmt.Println(div)

}

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

	db, err := L.ConnectDB("data/db.db")
	defer db.Close()

	privateKeyFile, err := exec.Command("cat", "../../../.private_keys/crypto/__swg__").Output()
	if err != nil {
		log.Fatal("error reading private key", err)
	}
	privateKey := strings.Split(string(privateKeyFile), "\n")[0]
	// fmt.Println(privateKey)

	// os.Exit(10)

	w3, err := L.ConnectWeb3("https://bsc.publicnode.com", privateKey)
	// defer w3.Close()

	fmt.Println(w3.Address)

	nonce, err := w3.GetNonce(w3.Address.String())

	fmt.Println(nonce)
	// os.Exit(10)

	flag.Parse()

	if init_db == "yes" {
		db.CreateTable("test", "test STRING PRIMARY KEY, age INTEGER")
		test := make(map[string]interface{})
		test["test"] = "sliem"
		test["age"] = 1001
		rows, err := db.SelectFromTable("test", "*")
		if err != nil {
			log.Fatal("failed to retrieve rows")
		}
		db.InsertRecordIntoTable("test", test)
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
	var list_addresses = convertToStringList(addresses)

	// w3_client, err := rpc.Dial(provider)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	router_abi, err := os.ReadFile("build/solidity/IPancakeRouter02.abi")
	if err != nil {
		log.Fatal((err))
	}

	var parsed_abi abi.ABI

	// _err :=
	json.Unmarshal(router_abi, &parsed_abi)

	// fmt.Println(parsed_abi)

	client := &http.Client{}

	// serverAddr, err := net.ResolveUDPAddr("udp", "localhost:8080")
	// if err != nil {
	// 	fmt.Println("Error resolving address:", err)
	// 	return
	// }

	// conn, err := net.DialUDP("udp", nil, serverAddr)
	// if err != nil {
	// 	fmt.Println("Error connecting to server:", err)
	// 	return
	// }
	// defer conn.Close()

	// go exec.Command("nc", "-ul", "8080")

	var jsonRes map[string]interface{}

	fmt.Println()

	for {
		res := base_post(provider, get_payload(id, `"eth_newPendingTransactionFilter"`, ""), *client)
		err = json.Unmarshal([]byte(res), &jsonRes)
		if err != nil {
			continue
		}
		// fmt.Println(res)
		res = base_post(provider, get_payload(id, `"eth_getFilterChanges"`, fmt.Sprintf(`"%s"`, jsonRes["result"].(string))), *client)
		err = json.Unmarshal([]byte(res), &jsonRes)
		if err != nil {
			continue
		}
		// fmt.Println(res)
		hashes, succ := jsonRes["result"].([]interface{})
		if !succ {
			continue
		}
		for _, hash := range hashes {
			// fmt.Println(fmt.Sprintf("%s <-> %s/tx/%s", get_now(), explorer, hash))
			// conn.WriteToUDP([]byte(hash.(string)), &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: 8888})
			if mode == "hash" {
				go handle_hash(explorer, provider, hash.(string), id, list_addresses, *client)
			}
			if mode == "dict" {
				_to, succ := hash.(map[string]interface{})["to"].(string)
				if succ {
					if Contains(list_addresses, _to) {

						// fmt.Println(hash)
						str_input := from_interface_to_string(hash.(map[string]interface{})["input"])
						str_value := from_interface_to_string(hash.(map[string]interface{})["value"])
						str_gasPrice := from_interface_to_string(hash.(map[string]interface{})["gasPrice"])

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

						// fmt.Printf("Method: %s\n", method)
						_method := fmt.Sprintf("%s", method)

						if (strings.Contains(_method, "swapExactETH") || strings.Contains(_method, "swapETH")) && !strings.Contains(_method, "FeeOn") {
							fmt.Printf("%s <-> %s/tx/%s\n", get_now(), explorer, hash.(map[string]interface{})["hash"])
							// fmt.Println(hash)
							fmt.Printf("Method: %s\n", method)
							for key, value := range args {
								fmt.Printf("%s: %+v\n", key, value)
							}

							fmt.Println("value ->", str_value)
							valueB, succ := from_string_to_int(str_value[2:]) // new(big.Int).SetString(str_value[2:], 16)
							if !succ {
								log.Fatal("failed to convert value")
							}

							// res := new(big.Int)
							fmt.Println("value ->", valueB)
							// if str_gasPrice != "error" {
							gasPrice, succ := from_string_to_int(str_gasPrice[2:])
							// div := gasPrice.Int64() / 1e18

							fmt.Println("gasPrice ->", str_gasPrice)
							if succ {
								fmt.Println("gasPrice ->", gasPrice)

							}
							// }

							fmt.Println(get_now())

							fmt.Println()

						}

					}
				}

			}

		}
	}
}
