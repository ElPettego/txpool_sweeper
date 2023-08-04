package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

var network string

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
	// _to, succ := jsonRes["result"] //
	// if !succ {
	// 	return
	// }
	// to, succ := _to.(map[string]interface{})["to"]
	// if !succ {
	// 	return
	// }
	// fmt.Println(fmt.Sprintf("%s <-> %s", get_now(), to))
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
	return time.Now().UTC().Format("[2006-01-02|15:04:05.000]")
}

func get_payload(id int, method, params string) string {
	return fmt.Sprintf(`{"id": %d,"jsonrpc":"2.0","method":%s,"params":[%s]}`, id, method, params)
}

func base_post(url, payload string, client http.Client) string {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return ""
		panic("error creating request")
	}
	res, err := client.Do(req)
	if err != nil {
		return ""
		panic("error executing request")
	}
	defer res.Body.Close()
	var resPayload bytes.Buffer
	_, err = resPayload.ReadFrom(res.Body)
	if err != nil {
		return ""
		panic("error reading request")

	}
	return resPayload.String()
}

func init() {
	flag.StringVar(&network, "network", "", "EVM network to interact with. must be present in .json.")
}

func main() {
	flag.Parse()
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
	var id = int(parsedConfig["id"].(float64))
	var addresses = parsedConfig["addresses"].([]interface{})
	var list_addresses = convertToStringList(addresses)
	// var nodes = parsedConfig["nodes"].([]interface{})

	// fmt.Println(addresses)
	// fmt.Printf("%T\n", addresses)

	// fmt.Println(list_addresses)
	// fmt.Printf("%T\n", list_addresses)

	// fmt.Println(Contains(list_addresses, "slime"))
	// fmt.Println(Contains(list_addresses, "0xa5E0829CaCEd8fFDD4De3c43696c57F7D7A678ff"))

	client := &http.Client{}

	var jsonRes map[string]interface{}

	for {
		res := base_post(provider, get_payload(id, `"eth_newPendingTransactionFilter"`, ""), *client)
		err = json.Unmarshal([]byte(res), &jsonRes)
		if err != nil {
			continue
		}
		res = base_post(provider, get_payload(id, `"eth_getFilterChanges"`, fmt.Sprintf(`"%s"`, jsonRes["result"].(string))), *client)
		err = json.Unmarshal([]byte(res), &jsonRes)
		if err != nil {
			continue
		}
		hashes, succ := jsonRes["result"].([]interface{})
		if !succ {
			continue
		}
		for _, hash := range hashes {
			// fmt.Println(fmt.Sprintf("%s <-> %s/tx/%s", get_now(), explorer, hash))
			// conn.WriteToUDP([]byte(hash.(string)), &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: 8888})
			go handle_hash(explorer, provider, hash.(string), id, list_addresses, *client)
		}
	}
}
