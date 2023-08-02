package main

import (
	"fmt"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"net/http"
	"github.com/pelletier/go-toml"
	"os"
	"flag"
)

type Payload struct {
	Id      string `json:"id"`
	JsonRpc string `json:"jsonrpc"`
	Params  string `json:"params"`
	Method  string `json:"method"`
}

var network string

func init() {
	flag.StringVar(&network, "network", "", "evm network to interact with. must be present in .toml.")
}

func main() {
	flag.Parse()
	fmt.Println("yo bro", network)
}
