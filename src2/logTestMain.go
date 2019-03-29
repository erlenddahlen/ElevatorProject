package main

import (
	"encoding/json"
	"fmt"

	"./Config"
	"./Governor"
	//"./Peertest"
)

func main() {
	id := "120"
	var GState Config.GlobalState
	GState = Governor.GovernorInit(GState, id)
	GStatejason, _ := json.MarshalIndent(GState, "", " ")
	fmt.Printf("Gstate:\n %s", GStatejason)

}
