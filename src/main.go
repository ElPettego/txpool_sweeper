package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
)

var network string

func init() {
	flag.StringVar(&network, "network", "", "EVM network to interact with. must be present in .json.")
}

func get_payload(id int, method, params string) string {
	return fmt.Sprintf(
		`{
		"id"      : %d,
		"jsonrpc" : "2.0",
		"method"  : %s,
		"params"  : [%s]
	}`, id, method, params)

}

func base_post(url, payload string, client http.Client) string {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		panic("error creating request")
	}

	res, err := client.Do(req)
	if err != nil {
		panic("error executing request")
	}

	defer res.Body.Close()

	var resPayload bytes.Buffer

	_, err = resPayload.ReadFrom(res.Body)
	if err != nil {
		panic("error reading request")
	}

	fmt.Println(resPayload.String())

	return resPayload.String()
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
	// var explorer = parsedConfig["explorer"].(string)
	var id = int(parsedConfig["id"].(float64))
	// var addresses = parsedConfig["addresses"].([]interface{})
	// var nodes = parsedConfig["nodes"].([]interface{})

	// fmt.Println(reflect.TypeOf(id))
	// fmt.Println(id)
	// fmt.Println(addresses)
	fmt.Println(get_payload(id, "eth_newPendingTransactionFilter", ""))

	client := &http.Client{}

	var jsonRes map[string]interface{}

	res := base_post(provider, get_payload(id, `"eth_newPendingTransactionFilter"`, ""), *client)
	err = json.Unmarshal([]byte(res), &jsonRes)
	if err != nil {
		panic(err)
	}
	fmt.Println(jsonRes["result"].(string))
	res = base_post(provider, get_payload(id, `"eth_getFilterChanges"`, fmt.Sprintf(`"%s"`, jsonRes["result"].(string))), *client)
	err = json.Unmarshal([]byte(res), &jsonRes)
	if err != nil {
		panic(err)
	}
	fmt.Println(jsonRes["result"].([]interface{}))

}
