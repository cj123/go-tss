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
	tss.DisableMessages()

	ip, err := ipsw.NewIPSWWithIdentifierBuild(client, os.Args[1], os.Args[2])

	if err != nil {
		panic(err)
	}

	raw, err := ip.RawManifest()

	if err != nil {
		panic(err)
	}

	manifest := raw["BuildIdentities"].([]interface{})[0]

	params, err := tss.AddParametersFromManifest(map[string]interface{}{
		"ApECID":           22222222222222,
		"ApProductionMode": true,
	}, manifest)

	if err != nil {
		panic(err)
	}

	tssReq, err := tss.NewRequest(nil)

	if err != nil {
		panic(err)
	}

	defer tssReq.Close()

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
