package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	s := make(map[string]interface{})
	g := make(map[string]interface{})
	g["c"] = 3
	g["d"] = "dd"
	s["a"] = 1
	s["b"] = "dfadfafd"
	s["r"] = g

	sb, e := json.Marshal(s)
	fmt.Println(string(sb), e)
}
