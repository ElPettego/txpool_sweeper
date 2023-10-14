package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"
)

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

func FromStringToInt(str string) (*big.Int, bool) {
	return new(big.Int).SetString(str, 16)
}

func ConvertToStringList(input []interface{}) []string {
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

func Contains(list []string, x string) bool {
	for _, item := range list {
		if item == x {
			return true
		}
	}

	return false
}

func GetNow() string {
	return "[" + Purple + time.Now().UTC().Format("2006-01-02|15:04:05.000000") + Reset + "]"
}

func HandleHash(explorer, provider, hash string, id int, addresses []string, client http.Client) {
	var jsonRes map[string]interface{}
	res := BasePost(provider, GetPayload(id, `"eth_getTransactionByHash"`, fmt.Sprintf(`"%s"`, hash)), client)
	err := json.Unmarshal([]byte(res), &jsonRes)
	if err != nil {
		return
		// panic(err)
	}
	if result, ok := jsonRes["result"].(map[string]interface{}); ok {
		to, succ := result["to"].(string)
		if succ {
			if Contains(addresses, to) {
				fmt.Println(fmt.Sprintf("%s <-> %s %s/tx/%s", GetNow(), to, explorer, result["hash"]))
			}
		}
	}
}

func GetPayload(id int, method, params string) string {
	return fmt.Sprintf(`{"id": %d,"jsonrpc":"2.0","method":%s,"params":[%s]}`, id, method, params)
}

func BasePost(url, payload string, client http.Client) string {
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

func FromInterfaceToString(_interface interface{}) string {
	if _, ok := _interface.(string); ok {
		return _interface.(string)
	}
	return "error"
}
