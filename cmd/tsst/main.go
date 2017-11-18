package main

import (
	"fmt"
	"os"

	"github.com/cj123/go-ipsw"
	"github.com/cj123/go-ipsw/api"
	"github.com/cj123/go-tss"
)

func main() {
	client := api.NewIPSWClientLatest()

	ip, err := ipsw.NewIPSWWithIdentifierBuild(client, os.Args[1], os.Args[2])

	if err != nil {
		panic(err)
	}

	raw, err := ip.RawManifest()

	if err != nil {
		panic(err)
	}

	manifest := raw["BuildIdentities"].([]interface{})[0]

	params := tss.AddParametersFromManifest(map[string]interface{}{
		"ApECID":           22222222222222,
		"ApProductionMode": true,
	}, manifest)

	tssReq := tss.NewRequest(nil)
	tssReq.AddCommonTags(params, nil)

	b, err := tssReq.Bytes()

	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", b)

	resp, err := tssReq.Send()

	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", resp)
}
